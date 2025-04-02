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

package gorm

import (
	"context"
	"fmt"

	gormModels "github.com/mrsimonemms/opensesame/apps/server/internal/database/gorm/models"
	"github.com/mrsimonemms/opensesame/apps/server/pkg/config"
	"github.com/mrsimonemms/opensesame/apps/server/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SQL struct {
	activeConnection *gorm.DB
	dialector        gorm.Dialector
}

func (s *SQL) Check(ctx context.Context) error {
	db, err := s.activeConnection.DB()
	if err != nil {
		return fmt.Errorf("error retrieving db: %w", err)
	}

	return db.PingContext(ctx)
}

func (s *SQL) Close(ctx context.Context) error {
	return nil
}

func (s *SQL) Connect(ctx context.Context) error {
	db, err := gorm.Open(s.dialector, &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	s.activeConnection = db

	if err := s.Check(ctx); err != nil {
		return fmt.Errorf("failed to check db connection: %w", err)
	}

	models := []any{
		&gormModels.User{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("error migrating databases: %w", err)
	}

	return nil
}

func (s *SQL) DeleteOrganisation(ctx context.Context, orgID, userID string) error {
	panic("delete org unimplemented")
}

func (s *SQL) FindUserByProviderAndUserID(ctx context.Context, providerID, providerUserID string) (user *models.User, err error) {
	panic("find user by provider and user id unimplemented")
}

func (s *SQL) GetOrgByID(ctx context.Context, orgID, userID string) (org *models.Organisation, err error) {
	panic("get org by id unimplemented")
}

func (s *SQL) GetOrgBySlug(ctx context.Context, slug string) (org *models.Organisation, err error) {
	panic("get org by slug unimplemented")
}

func (s *SQL) GetUserByID(ctx context.Context, userID string) (user *models.User, err error) {
	panic("get user by id unimplemented")
}

func (s *SQL) ListOrganisationUsers(
	ctx context.Context,
	offset,
	limit int,
	orgID,
	userID string,
) (users *models.Pagination[*models.OrganisationUser], err error) {
	panic("lsit org users unimplemented")
}

func (s *SQL) ListOrganisations(
	ctx context.Context,
	offset,
	limit int,
	userID string,
) (orgs *models.Pagination[*models.Organisation], err error) {
	panic("list orgs unimplemented")
}

func (s *SQL) SaveOrganisationRecord(ctx context.Context, model *models.Organisation) (user *models.Organisation, err error) {
	panic("save org record unimplemented")
}

func (s *SQL) SaveUserRecord(ctx context.Context, model *models.User) (user *models.User, err error) {
	panic("save user record unimplemented")
}

func (s *SQL) UpdateAllUsers(
	ctx context.Context,
	update func(existing []*models.User) (updated []*models.User, err error),
) (count int64, err error) {
	panic("update all users unimplemented")
}

func New(cfg config.SQL) (*SQL, error) {
	var dialector gorm.Dialector
	switch cfg.Type {
	case "mysql":
		dialector = mysql.Open(cfg.DSN)
	case "postgres":
		dialector = postgres.New(postgres.Config{})
	default:
		return nil, fmt.Errorf("unknown database type: %s", cfg.Type)
	}

	return &SQL{
		dialector: dialector,
	}, nil
}
