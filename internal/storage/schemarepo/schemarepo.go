package schemarepo

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

type schemarepo struct {
	db  *gorm.DB
	cfg *config.Config
	log logger.Logger
}

func New(cfg *config.Config, log logger.Logger) storage.Storage {
	return &schemarepo{
		cfg: cfg,
		log: log,
	}
}

func (r *schemarepo) Connect(ctx context.Context) (storage.Storage, error) {
	r.log.Debug(ctx, "initialize db session")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		r.cfg.Storage.HOST,
		r.cfg.Storage.PORT,
		r.cfg.Storage.User,
		r.cfg.Storage.Name,
		r.cfg.Storage.Password,
	)

	db, err := gorm.Open(driver.Open(dsn), &gorm.Config{})
	if err != nil {
		r.log.Error(ctx, fmt.Errorf("%v:%w", exceptions.ErrConnectingToDatabase, err), "could not initialise db session")

		return nil, fmt.Errorf("%v:%w", exceptions.ErrConnectingToDatabase, err)
	}

	if r.cfg.Debug {
		db = db.Debug()
	}

	r.db = db

	return r, nil
}

func (r *schemarepo) Shutdown(ctx context.Context) error {
	r.log.Debug(ctx, "close database")

	sql, err := r.db.DB()
	if err != nil {
		r.log.Error(ctx, fmt.Errorf("%v:%w", exceptions.ErrGetDB, err), "could not get database handle")

		return fmt.Errorf("%v:%w", exceptions.ErrGetDB, err)
	}

	if err := sql.Close(); err != nil {
		r.log.Error(ctx, fmt.Errorf("%v:%w", exceptions.ErrCloseDB, err), "could not close database")

		return fmt.Errorf("%v:%w", exceptions.ErrCloseDB, err)
	}

	return nil
}

func (r *schemarepo) Initialize(ctx context.Context) error {
	r.log.Debug(ctx, "initialize database")

	if err := r.db.AutoMigrate(schema.Schema{}); err != nil {
		r.log.Error(ctx, fmt.Errorf("%v:%w", exceptions.ErrInitializeDatabase, err), "could not automigrate")

		return fmt.Errorf("%v:%w", exceptions.ErrInitializeDatabase, err)
	}

	return nil
}

func (r *schemarepo) CreateSchema(ctx context.Context, schemaID, schemaPayload string) error {
	r.log.Debug(ctx, "upload schema")

	schemaJSON := datatypes.JSON(schemaPayload)

	s := &schema.Schema{
		SchemaID: schemaID,
		Schema:   &schemaJSON,
	}

	return r.db.Create(s).Error
}

func (r *schemarepo) GetSchema(ctx context.Context, schemaID string) (string, error) {
	r.log.Debug(ctx, "download schema")

	s := &schema.Schema{}

	if err := r.db.Where(&schema.Schema{SchemaID: schemaID}).Take(s).Error; err != nil {
		return "", err
	}

	return s.Schema.String(), nil
}
