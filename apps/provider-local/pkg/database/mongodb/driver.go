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
	"context"
	"fmt"
	"time"

	"github.com/mrsimonemms/opensesame/apps/provider-local/pkg/config"
	mongoModels "github.com/mrsimonemms/opensesame/apps/provider-local/pkg/database/mongodb/models"
	"github.com/mrsimonemms/opensesame/packages/go-sdk/provider/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoConnection struct {
	client *mongo.Client
	db     *mongo.Database
}

type MongoDB struct {
	activeConnection mongoConnection
	connectionURI    string
	collection       string
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
	if err := db.applyIndices(ctx); err != nil {
		return fmt.Errorf("error applying indices: %w", err)
	}

	return nil
}

// FindUserByEmailAddress implements database.Driver.
func (db *MongoDB) FindUserByEmailAddress(ctx context.Context, emailAddress string) (*models.User, error) {
	filter := bson.D{
		{
			Key:   "email_address",
			Value: emailAddress,
		},
	}

	var result mongoModels.User
	err := db.activeConnection.db.Collection(db.collection).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, fmt.Errorf("error retrieving user email address in mongodb: %w", err)
	}

	return result.ToModel(), nil
}

// FindUserByUsername implements database.Driver.
func (db *MongoDB) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	filter := bson.D{
		{
			Key:   "username",
			Value: username,
		},
	}

	var result mongoModels.User
	err := db.activeConnection.db.Collection(db.collection).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, fmt.Errorf("error retrieving user email address in mongodb: %w", err)
	}

	return result.ToModel(), nil
}

// Save implements database.Driver.
func (db *MongoDB) Save(ctx context.Context, model *models.User) (*models.User, error) {
	record, err := mongoModels.UserToModel(model)
	if err != nil {
		return nil, err
	}

	collection := db.activeConnection.db.Collection(db.collection)

	if record.ID.IsZero() {
		// Insert new record
		result, err := collection.InsertOne(ctx, record)
		if err != nil {
			return nil, fmt.Errorf("error insert new record: %w", err)
		}
		record.ID = result.InsertedID.(bson.ObjectID)
	} else {
		// Update existing record
		if _, err := collection.UpdateByID(ctx, record.ID, bson.M{"$set": record}); err != nil {
			return nil, fmt.Errorf("error updating record: %w", err)
		}
	}

	return record.ToModel(), nil
}

func (db *MongoDB) applyIndices(ctx context.Context) error {
	indices := map[string][]mongo.IndexModel{
		db.collection: {
			{
				Keys: bson.D{
					{Key: "email_address", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bson.D{
					{Key: "username", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
	}

	for collection, indexModels := range indices {
		for _, index := range indexModels {
			if _, err := db.activeConnection.db.Collection(collection).Indexes().CreateOne(ctx, index); err != nil {
				return fmt.Errorf("error creating index for %s collection: %w", collection, err)
			}
		}
	}

	return nil
}

func New(cfg config.MongoDB) *MongoDB {
	return &MongoDB{
		connectionURI: cfg.ConnectionURI,
		collection:    cfg.Collection,
		database:      cfg.Database,
	}
}
