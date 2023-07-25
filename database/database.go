package database

import (
	"database/sql"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/SergeyTyurin/banner-rotation/configs"
	"github.com/SergeyTyurin/banner-rotation/structures"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrNilConfig         = errors.New("config is nil")
	ErrNotExist          = errors.New("entity not exists in database")
	ErrNotInRotation     = errors.New("entities not in rotation")
	ErrAlreadyInRotation = errors.New("entities already in rotation")
)

const invalidID = -1

type Database interface {
	DatabaseConnect(config configs.DBConnectionConfig) (func() error, error)

	DatabaseGetBanner(id int) (structures.Banner, error)
	DatabaseGetSlot(id int) (structures.Slot, error)
	DatabaseGetGroup(id int) (structures.Group, error)

	DatabaseDeleteBanner(id int) error
	DatabaseDeleteSlot(id int) error
	DatabaseDeleteGroup(id int) error

	DatabaseCreateBanner(structures.Banner) (structures.Banner, error)
	DatabaseCreateSlot(structures.Slot) (structures.Slot, error)
	DatabaseCreateGroup(structures.Group) (structures.Group, error)

	DatabaseUpdateBanner(structures.Banner) error
	DatabaseUpdateSlot(structures.Slot) error
	DatabaseUpdateGroup(structures.Group) error

	DatabaseAddToRotation(bannerID, slotID int) error
	DatabaseDeleteFromRotation(bannerID, slotID int) error
	DatabaseSelectFromRotation(slotID, groupID int) (bannerID int, err error)
	DatabaseRegisterTransition(slotID, bannerID, groupID int) error
}

type databaseImpl struct {
	db *sql.DB
}

func (di *databaseImpl) DatabaseConnect(config configs.DBConnectionConfig) (func() error, error) {
	if config == nil {
		return nil, ErrNilConfig
	}
	var connectionError error
	url := config.URL()
	url = strings.ReplaceAll(url, "{host}", config.Host())
	url = strings.ReplaceAll(url, "{port}", strconv.Itoa(config.Port()))
	url = strings.ReplaceAll(url, "{user}", os.Getenv("DB_USER"))
	url = strings.ReplaceAll(url, "{password}", os.Getenv("DB_PASSWORD"))
	url = strings.ReplaceAll(url, "{dbname}", config.DatabaseName())

	di.db, connectionError = sql.Open("pgx", url)
	if connectionError != nil {
		return nil, connectionError
	}
	closeConnect := func() error {
		return di.db.Close()
	}

	if err := di.db.Ping(); err != nil {
		if err := closeConnect(); err != nil {
			return nil, err
		}
		return nil, err
	}
	return closeConnect, nil
}

func NewDatabase() Database {
	return &databaseImpl{db: nil}
}

func checkEntityIsExists(d *databaseImpl, tablename string, id int) error {
	var receivedID int
	switch tablename {
	case "Banners":
		entity, _ := d.DatabaseGetBanner(id)
		receivedID = entity.ID
	case "Slots":
		entity, _ := d.DatabaseGetSlot(id)
		receivedID = entity.ID
	case "Groups":
		entity, _ := d.DatabaseGetGroup(id)
		receivedID = entity.ID
	default:
		return nil
	}
	if receivedID == invalidID {
		return ErrNotExist
	}
	return nil
}
