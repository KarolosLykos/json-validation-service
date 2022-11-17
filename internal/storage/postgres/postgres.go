package postgres

import (
	"context"
	"fmt"

	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/KarolosLykos/json-validation-service/internal/config"
	"github.com/KarolosLykos/json-validation-service/internal/logger"
	"github.com/KarolosLykos/json-validation-service/internal/models/schema"
	"github.com/KarolosLykos/json-validation-service/internal/storage"
	"github.com/KarolosLykos/json-validation-service/internal/utils/exceptions"
)

type postgres struct {
	db  *gorm.DB
	cfg *config.Config
	log logger.Logger
}

func New(cfg *config.Config, log logger.Logger) storage.Storage {
	return &postgres{
		cfg: cfg,
		log: log,
	}
}

func (p *postgres) Connect(ctx context.Context) (storage.Storage, error) {
	p.log.Debug(ctx, "initialize db session")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		p.cfg.Storage.HOST,
		p.cfg.Storage.PORT,
		p.cfg.Storage.User,
		p.cfg.Storage.Name,
		p.cfg.Storage.Password,
	)

	db, err := gorm.Open(driver.Open(dsn), &gorm.Config{})
	if err != nil {
		p.log.Error(ctx, fmt.Errorf("%v:%w", exceptions.ErrConnectingToDatabase, err), "could not initialise db session")

		return nil, fmt.Errorf("%v:%w", exceptions.ErrConnectingToDatabase, err)
	}

	if p.cfg.Debug {
		db = db.Debug()
	}

	p.db = db

	return p, nil
}

func (p *postgres) Shutdown(ctx context.Context) error {
	p.log.Debug(ctx, "close database")

	sql, err := p.db.DB()
	if err != nil {
		p.log.Error(ctx, fmt.Errorf("%v:%w", exceptions.ErrGetDB, err), "could not get database handle")

		return fmt.Errorf("%v:%w", exceptions.ErrGetDB, err)
	}

	if err := sql.Close(); err != nil {
		p.log.Error(ctx, fmt.Errorf("%v:%w", exceptions.ErrCloseDB, err), "could not close database")

		return fmt.Errorf("%v:%w", exceptions.ErrCloseDB, err)
	}

	return nil
}

func (p *postgres) Initialize(ctx context.Context) error {
	p.log.Debug(ctx, "initialize database")

	if err := p.db.AutoMigrate(schema.Schema{}); err != nil {
		p.log.Error(ctx, fmt.Errorf("%v:%w", exceptions.ErrInitializeDatabase, err), "could not automigrate")

		return fmt.Errorf("%v:%w", exceptions.ErrInitializeDatabase, err)
	}

	return nil
}
