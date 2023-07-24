package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/SergeyTyurin/banner_rotation/structures"

	"github.com/SergeyTyurin/banner_rotation/int/configs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrNilConfig         = errors.New("config is nil")
	ErrNotExist          = errors.New("entity not exists in database")
	ErrNotInRotation     = errors.New("entities not in rotation")
	ErrAlreadyInRotation = errors.New("entities already in rotation")
)

const invalidId = -1

type Database interface {
	Connect(config configs.DBConnectionConfig) (func() error, error)

	GetBanner(id int) (structures.Banner, error)
	GetSlot(id int) (structures.Slot, error)
	GetGroup(id int) (structures.Group, error)

	DeleteBanner(id int) error
	DeleteSlot(id int) error
	DeleteGroup(id int) error

	CreateBanner(structures.Banner) (structures.Banner, error)
	CreateSlot(structures.Slot) (structures.Slot, error)
	CreateGroup(structures.Group) (structures.Group, error)

	UpdateBanner(structures.Banner) error
	UpdateSlot(structures.Slot) error
	UpdateGroup(structures.Group) error

	AddToRotation(bannerId, slotId int) error
	DeleteFromRotation(bannerId, slotId int) error
	SelectFromRotation(slotId, groupId int) (bannerId int, err error)
	RegisterTransition(slotId, bannerId, groupId int) error
}

type databaseImpl struct {
	db *sql.DB
}

func (di *databaseImpl) Connect(config configs.DBConnectionConfig) (func() error, error) {
	if config == nil {
		return nil, ErrNilConfig
	}
	var connectionError error

	format := `host=%s port=%d user=%s password=%s dbname=%s`
	dsn := fmt.Sprintf(format,
		config.Host(),
		config.Port(),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		config.DatabaseName())

	di.db, connectionError = sql.Open("pgx", dsn)
	if connectionError != nil {
		return nil, connectionError
	}
	closeConnect := func() error {
		return di.db.Close()
	}

	if err := di.db.Ping(); err != nil {
		closeConnect()
		return nil, err
	}
	return closeConnect, nil
}

func NewDatabase() Database {
	return &databaseImpl{db: nil}
}

func checkEntityIsExists(d *databaseImpl, tablename string, id int) error {
	var receivedId int
	switch tablename {
	case "Banners":
		entity, _ := d.GetBanner(id)
		receivedId = entity.Id
	case "Slots":
		entity, _ := d.GetSlot(id)
		receivedId = entity.Id
	case "Groups":
		entity, _ := d.GetGroup(id)
		receivedId = entity.Id
	default:
		return nil
	}
	if receivedId == invalidId {
		return ErrNotExist
	}
	return nil
}
