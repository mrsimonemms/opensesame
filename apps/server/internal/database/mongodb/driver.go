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
	"math"
	"time"

	"github.com/mrsimonemms/cloud-native-auth/apps/server/internal/common"
	mongoModels "github.com/mrsimonemms/cloud-native-auth/apps/server/internal/database/mongodb/models"
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
	if err := db.applyIndices(ctx); err != nil {
		return fmt.Errorf("error applying indices: %w", err)
	}

	return nil
}

func (db *MongoDB) DeleteOrganisation(ctx context.Context, orgID, userID string) error {
	id, err := bson.ObjectIDFromHex(orgID)
	if err != nil {
		return fmt.Errorf("error converting org id to bson object id: %w", err)
	}

	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "users.userId", Value: userID},
	}

	result, err := db.activeConnection.db.Collection(OrgsCollection).DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting org: %w", err)
	}

	if result.DeletedCount == 0 {
		return common.ErrNotDeleted
	}

	return nil
}

func (db *MongoDB) FindUserByProviderAndUserID(ctx context.Context, providerID, providerUserID string) (*models.User, error) {
	filter := bson.D{
		{
			Key:   fmt.Sprintf("accounts.%s.providerUserId", providerID),
			Value: providerUserID,
		},
	}

	var result mongoModels.User
	err := db.activeConnection.db.Collection(UsersCollection).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, fmt.Errorf("error retrieving user by provider in mongodb: %w", err)
	}

	return result.ToModel(), nil
}

func (db *MongoDB) GetOrgByID(ctx context.Context, orgID, userID string) (*models.Organisation, error) {
	id, err := bson.ObjectIDFromHex(orgID)
	if err != nil {
		return nil, fmt.Errorf("error converting org id to bson object id: %w", err)
	}

	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "users.userId", Value: userID},
	}

	var result mongoModels.Organisation
	err = db.activeConnection.db.Collection(OrgsCollection).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, fmt.Errorf("error getting org by id: %w", err)
	}

	return result.ToModel(), nil
}

func (db *MongoDB) GetOrgBySlug(ctx context.Context, slug string) (*models.Organisation, error) {
	filter := bson.D{
		{Key: "slug", Value: slug},
	}

	var result mongoModels.Organisation
	if err := db.activeConnection.db.Collection(OrgsCollection).FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, fmt.Errorf("error getting org by slug: %w", err)
	}

	return result.ToModel(), nil
}

func (db *MongoDB) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	id, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("error converting user id to bson object id: %w", err)
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: id,
		},
	}

	var result mongoModels.User
	err = db.activeConnection.db.Collection(UsersCollection).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, fmt.Errorf("error getting user by id: %w", err)
	}

	return result.ToModel(), nil
}

func (db *MongoDB) ListOrganisations(
	ctx context.Context,
	offset,
	limit int,
	userID string,
) (*models.Pagination[*models.Organisation], error) {
	col := db.activeConnection.db.Collection(OrgsCollection)
	filter := bson.D{
		{Key: "users.userId", Value: userID},
	}
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset)).SetSort(bson.D{
		{Key: "name", Value: 1},
	})

	totalDocs, err := col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error counting organisations: %w", err)
	}

	cursor, err := col.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding organisations: %w", err)
	}

	var mongodbOrgs []*mongoModels.Organisation
	if err := cursor.All(ctx, &mongodbOrgs); err != nil {
		return nil, fmt.Errorf("error getting all organisation records in cursor: %w", err)
	}

	orgs := make([]*models.Organisation, 0)
	for _, o := range mongodbOrgs {
		orgs = append(orgs, o.ToModel())
	}

	totalPages := int(math.Ceil(float64(totalDocs) / float64(limit)))
	page := int(math.Ceil(float64(offset-1)/float64(limit)) + 1)
	if page < 0 {
		page = 0
	} else if page > totalPages {
		page = totalPages
	}

	return &models.Pagination[*models.Organisation]{
		Data:       orgs,
		Count:      len(orgs),
		Page:       page,
		PerPage:    limit,
		TotalPages: totalPages,
		Total:      totalDocs,
	}, nil
}

func (db *MongoDB) ListOrganisationUsers(
	ctx context.Context,
	offset,
	limit int,
	orgID,
	userID string,
) (*models.Pagination[*models.OrganisationUser], error) {
	org, err := db.GetOrgByID(ctx, orgID, userID)
	if err != nil {
		return nil, fmt.Errorf("errot getting org by id: %w", err)
	}
	if org == nil {
		return nil, nil
	}

	totalDocs := len(org.Users)
	start := max(offset, 0)
	end := min(offset+limit, totalDocs)

	// Prune the org's users
	org.Users = org.Users[start:end]

	totalPages := int(math.Ceil(float64(totalDocs) / float64(limit)))
	page := int(math.Ceil(float64(offset-1)/float64(limit)) + 1)
	if page < 0 {
		page = 0
	} else if page > totalPages {
		page = totalPages
	}

	userIDs, err := mongoModels.PaginateUniqueUsers(org)
	if err != nil {
		return nil, fmt.Errorf("error extracting unique users from org %s: %w", org.ID, err)
	}

	if len(userIDs) > 0 {
		users, err := db.findMultipleUsers(ctx, userIDs)
		if err != nil {
			return nil, err
		}

		for key, value := range org.Users {
			org.Users[key].Name = users[value.UserID].Name
		}
	}

	return &models.Pagination[*models.OrganisationUser]{
		Data:       org.Users,
		Count:      len(org.Users),
		Page:       page,
		PerPage:    limit,
		TotalPages: totalPages,
		Total:      int64(totalDocs),
	}, nil
}

func (db *MongoDB) SaveOrganisationRecord(ctx context.Context, model *models.Organisation) (*models.Organisation, error) {
	mongoModel, err := mongoModels.OrganisationToMongo(model)
	if err != nil {
		return nil, fmt.Errorf("error converting before saving org record: %w", err)
	}

	col := db.activeConnection.db.Collection(OrgsCollection)

	recordID, err := saveGenericRecord(ctx, col, mongoModel.ID, mongoModel)
	if err != nil {
		return nil, err
	}

	mongoModel.ID = recordID

	fmt.Printf("%+v\n", mongoModel)

	return mongoModel.ToModel(), nil
}

func (db *MongoDB) SaveUserRecord(ctx context.Context, model *models.User) (*models.User, error) {
	mongoModel, err := mongoModels.UserToMongo(model)
	if err != nil {
		return nil, fmt.Errorf("error converting before saving user record: %w", err)
	}

	col := db.activeConnection.db.Collection(UsersCollection)

	recordID, err := saveGenericRecord(ctx, col, mongoModel.ID, mongoModel)
	if err != nil {
		return nil, err
	}

	mongoModel.ID = recordID

	return mongoModel.ToModel(), nil
}

func (db *MongoDB) UpdateAllUsers(
	ctx context.Context,
	update func(existing []*models.User) (updated []*models.User, err error),
) (int64, error) {
	collection := db.activeConnection.db.Collection(UsersCollection)
	filter := bson.D{}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("error retrieving all user records: %w", err)
	}

	var mongodbUsers []*mongoModels.User
	if err := cursor.All(ctx, &mongodbUsers); err != nil {
		return 0, fmt.Errorf("error getting all user records in cursor: %w", err)
	}

	users := make([]*models.User, 0)
	for _, u := range mongodbUsers {
		users = append(users, u.ToModel())
	}

	updatedRecords, err := update(users)
	if err != nil {
		return 0, fmt.Errorf("error generating updated user records: %w", err)
	}

	models := []mongo.WriteModel{}
	for _, model := range updatedRecords {
		s, err := mongoModels.UserToMongo(model)
		if err != nil {
			return 0, fmt.Errorf("error converting user to mongo model: %w", err)
		}

		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.D{
				{
					Key:   "_id",
					Value: s.ID,
				},
			}).SetUpdate(bson.M{"$set": s}),
		)
	}

	result, err := collection.BulkWrite(ctx, models)
	if err != nil {
		return 0, fmt.Errorf("error updating user records: %w", err)
	}

	return result.ModifiedCount, nil
}

func (db *MongoDB) applyIndices(ctx context.Context) error {
	indices := map[string][]mongo.IndexModel{
		OrgsCollection: {
			{
				Keys: bson.D{
					{Key: "_id", Value: 1},
					{Key: "users.userID", Value: 1},
				},
			},
			{
				Keys: bson.D{
					{Key: "slug", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
		UsersCollection: {
			{
				Keys: bson.D{
					{Key: "accounts.providerId", Value: 1},
					{Key: "accounts.providerUserId", Value: 1},
				},
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

func (db *MongoDB) findMultipleUsers(ctx context.Context, userIDs []bson.M) (map[string]*mongoModels.User, error) {
	filter := bson.M{
		"$or": userIDs,
	}

	cursor, err := db.activeConnection.db.Collection(UsersCollection).Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error getting multiple users: %w", err)
	}

	var users []*mongoModels.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("error populating multiple users from cursor: %w", err)
	}

	result := map[string]*mongoModels.User{}
	for _, u := range users {
		result[u.ID.Hex()] = u
	}

	return result, nil
}

func New(cfg config.MongoDB) *MongoDB {
	return &MongoDB{
		connectionURI: cfg.ConnectionURI,
		database:      cfg.Database,
	}
}

func saveGenericRecord(ctx context.Context, collection *mongo.Collection, recordID bson.ObjectID, record any) (bson.ObjectID, error) {
	if recordID.IsZero() {
		// No ID - create record
		result, err := collection.InsertOne(ctx, record)
		if err != nil {
			return recordID, fmt.Errorf("error inserting org record: %w", err)
		}

		recordID = result.InsertedID.(bson.ObjectID)
	}
	// ID exists - update record
	if _, err := collection.UpdateByID(ctx, recordID, bson.M{"$set": record}); err != nil {
		return recordID, fmt.Errorf("error updating org record: %w", err)
	}

	return recordID, nil
}
