package pq

import (
	"errors"
	"fmt"

	"github.com/apartapatia/wall_of_comments/internal/config"
	"github.com/apartapatia/wall_of_comments/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ErrMigrateComment = errors.New("failed to migrate comment")
var ErrMigratePost = errors.New("failed to migrate post")

func GetRepo(cfg config.PostgresConfig) (*Repo, error) {
	db, err := newClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Repo{db: db}, nil
}

func newClient(cfg config.PostgresConfig) (*gorm.DB, error) {
	logrus.Info("connecting to postgres")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai", cfg.Host, cfg.User, cfg.Password, cfg.DB, cfg.Port, cfg.SslMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := migrateDB(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&entity.Comment{}); err != nil {
		logrus.Error(ErrMigrateComment)
		return ErrMigrateComment
	}
	if err := db.AutoMigrate(&entity.Post{}); err != nil {
		logrus.Error(ErrMigratePost)
		return ErrMigratePost
	}
	return nil
}
