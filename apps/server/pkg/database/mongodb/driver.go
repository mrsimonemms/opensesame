/*
 * Copyright 2025 Simon Emms <simon@simonemms.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mongodb

import (
	"fmt"
	"time"

	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/config"
	"github.com/mrsimonemms/cloud-native-auth/apps/server/pkg/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/net/context"
)

type mongoConnection struct {
	client *mongo.Client
	db     *mongo.Database
}

type MongoDB struct {
	activeConnection mongoConnection
	connectionURI    string
	database         string
}

func (db *MongoDB) Check(ctx context.Context) error {
	return db.activeConnection.client.Ping(ctx, nil)
}

func (db *MongoDB) Close(ctx context.Context) error {
	return db.activeConnection.client.Disconnect(ctx)
}

func (db *MongoDB) Connect(ctx context.Context) error {
	opts := options.Client().
		ApplyURI(db.connectionURI).
		SetTimeout(time.Second).
		SetConnectTimeout(time.Second)

	client, err := mongo.Connect(opts)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	// Store the Mongo connection
	db.activeConnection = mongoConnection{
		client: client,
		db:     client.Database(db.database),
	}

	// Check the connection works
	if err := db.Check(ctx); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	// Apply the indices
	userModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "accounts.providerId", Value: 1},
			{Key: "accounts.providerUserId", Value: 1},
		},
	}

	if _, err := db.activeConnection.db.Collection(UsersCollection).Indexes().CreateOne(context.TODO(), userModel); err != nil {
		return fmt.Errorf("error creating index for user collection: %w", err)
	}

	return nil
}

func (db *MongoDB) FindUserByProviderAndUserID(ctx context.Context, providerID, providerUserID string) (*models.User, error) {
	filter := bson.D{
		{
			Key:   "accounts.providerId",
			Value: providerID,
		},
		{
			Key:   "accounts.providerUserId",
			Value: providerUserID,
		},
	}

	var result user
	err := db.activeConnection.db.Collection(UsersCollection).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, fmt.Errorf("error retrieving user by provider in mongodb: %w", err)
	}

	return result.toModel(), nil
}

func (db *MongoDB) SaveUserRecord(ctx context.Context, model *models.User) (*models.User, error) {
	mongoModel, err := userToMongo(model)
	if err != nil {
		return nil, fmt.Errorf("error converting before saving user record: %w", err)
	}

	col := db.activeConnection.db.Collection(UsersCollection)
	if mongoModel.ID.IsZero() {
		// No ID - create record
		result, err := col.InsertOne(ctx, mongoModel)
		if err != nil {
			return nil, fmt.Errorf("error inserting user record: %w", err)
		}

		mongoModel.ID = result.InsertedID.(bson.ObjectID)
	} else {
		// ID exists - update record
		if _, err := col.UpdateByID(ctx, mongoModel.ID, bson.M{"$set": mongoModel}); err != nil {
			return nil, fmt.Errorf("error updating user record: %w", err)
		}
	}

	return mongoModel.toModel(), nil
}

func New(cfg config.MongoDB) *MongoDB {
	return &MongoDB{
		connectionURI: cfg.ConnectionURI,
		database:      cfg.Database,
	}
}
