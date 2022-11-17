package store

import (
	"context"
	"fmt"

	"gorm.io/datatypes"
	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/KarolosLykos/json-validation-service/internal/config"
	"github.com/KarolosLykos/json-validation-service/internal/logger"
	"github.com/KarolosLykos/json-validation-service/internal/models/schema"
	"github.com/KarolosLykos/json-validation-service/internal/storage"
	"github.com/KarolosLykos/json-validation-service/internal/utils/exceptions"
)

type store struct {
	db  *gorm.DB
	cfg *config.Config
	log logger.Logger
}

func New(cfg *config.Config, log logger.Logger) storage.Storage {
	return &store{
		cfg: cfg,
		log: log,
	}
}

func (s *store) Connect(ctx context.Context) (storage.Storage, error) {
	s.log.Debug(ctx, "initialize db session")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		s.cfg.Storage.HOST,
		s.cfg.Storage.PORT,
		s.cfg.Storage.User,
		s.cfg.Storage.Name,
		s.cfg.Storage.Password,
	)

	db, err := gorm.Open(driver.Open(dsn), &gorm.Config{})
	if err != nil {
		s.log.Error(ctx, fmt.Errorf("%v:%w", exceptions.ErrConnectingToDatabase, err), "could not initialise db session")

		return nil, fmt.Errorf("%v:%w", exceptions.ErrConnectingToDatabase, err)
	}

	s.db = db

	return s, nil
}

func (s *store) Shutdown(ctx context.Context) error {
	s.log.Debug(ctx, "close database")

	sql, err := s.db.DB()
	if err != nil {
		s.log.Error(ctx, fmt.Errorf("%v:%w", exceptions.ErrGetDB, err), "could not get database handle")

		return fmt.Errorf("%v:%w", exceptions.ErrGetDB, err)
	}

	if err := sql.Close(); err != nil {
		s.log.Error(ctx, fmt.Errorf("%v:%w", exceptions.ErrCloseDB, err), "could not close database")

		return fmt.Errorf("%v:%w", exceptions.ErrCloseDB, err)
	}

	return nil
}

func (s *store) Initialize(ctx context.Context) error {
	s.log.Debug(ctx, "initialize database")

	if err := s.db.AutoMigrate(schema.Schema{}); err != nil {
		s.log.Error(ctx, fmt.Errorf("%v:%w", exceptions.ErrInitializeDatabase, err), "could not automigrate")

		return fmt.Errorf("%v:%w", exceptions.ErrInitializeDatabase, err)
	}

	return nil
}

func (s *store) CreateSchema(ctx context.Context, schemaID, schemaPayload string) error {
	s.log.Debug(ctx, "upload schema")

	schemaJSON := datatypes.JSON(schemaPayload)

	model := &schema.Schema{
		SchemaID: schemaID,
		Schema:   &schemaJSON,
	}

	return s.db.Create(model).Error
}

func (s *store) GetSchema(ctx context.Context, schemaID string) (string, error) {
	s.log.Debug(ctx, "download schema")

	model := &schema.Schema{}

	if err := s.db.Where(&schema.Schema{SchemaID: schemaID}).Take(model).Error; err != nil {
		return "", err
	}

	return model.Schema.String(), nil
}
