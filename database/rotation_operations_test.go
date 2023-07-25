package database

import (
	"strconv"
	"testing"

	"github.com/SergeyTyurin/banner-rotation/configs"
	"github.com/SergeyTyurin/banner-rotation/structures"
	"github.com/stretchr/testify/require"
)

func setTestData(d databaseImpl) {
	_, _ = d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)
	_, _ = d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`)
	_, _ = d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)

	for i := 0; i < 10; i++ {
		_, _ = d.DatabaseCreateBanner(structures.Banner{Info: "banner_" + strconv.Itoa(i+1)})
	}
	for i := 0; i < 2; i++ {
		_, _ = d.DatabaseCreateGroup(structures.Group{Info: "group_" + strconv.Itoa(i+1)})
	}
	for i := 0; i < 4; i++ {
		_, _ = d.DatabaseCreateSlot(structures.Slot{Info: "slot_" + strconv.Itoa(i+1)})
	}
}

func TestDatabaseAddToRotation(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()
	setTestData(d)

	t.Run("simple add", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`)
		slotID := 1
		bannerID := 2
		err := d.DatabaseAddToRotation(bannerID, slotID)
		require.NoError(t, err)

		row := d.db.QueryRow(`SELECT count(*) FROM "Statistic"
		WHERE slot_id=$1 AND banner_id=$2`, slotID, bannerID)

		count := 0
		_ = row.Scan(&count)
		// Test Add for each group
		groups, _ := d.DatabaseGetGroups()
		require.Equal(t, count, len(groups))
	})

	t.Run("add non existed intities", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`)

		err := d.DatabaseAddToRotation(20, 1)
		require.ErrorIs(t, err, ErrNotExist)

		err = d.DatabaseAddToRotation(1, 20)
		require.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("add already in rotation", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`)
		err := d.DatabaseAddToRotation(1, 1)
		require.NoError(t, err)

		err = d.DatabaseAddToRotation(1, 1)
		require.ErrorIs(t, err, ErrAlreadyInRotation)
	})
}

func TestDeleteFromRotation(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()
	setTestData(d)

	t.Run("simple delete", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`)
		slotID := 1
		bannerID := 2
		_ = d.DatabaseAddToRotation(bannerID, slotID)
		err := d.DatabaseDeleteFromRotation(bannerID, slotID)
		require.NoError(t, err)

		row := d.db.QueryRow(`SELECT count(*) FROM "Statistic"
	WHERE slot_id=$1 AND banner_id=$2`, slotID, bannerID)
		count := 0
		_ = row.Scan(&count)
		require.Equal(t, count, 0)
	})

	t.Run("delete not in rotation", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`)
		_ = d.DatabaseAddToRotation(1, 1)
		_ = d.DatabaseAddToRotation(2, 1)
		_ = d.DatabaseAddToRotation(1, 2)
		err := d.DatabaseDeleteFromRotation(2, 2)
		require.ErrorIs(t, err, ErrNotInRotation)
	})
}

func TestSelectFromRotation(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()
	setTestData(d)

	t.Run("existing select", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`)
		_ = d.DatabaseAddToRotation(1, 1)
		_ = d.DatabaseAddToRotation(2, 1)
		groups, _ := d.DatabaseGetGroups()
		for _, group := range groups {
			bannerID, err := d.DatabaseSelectFromRotation(1, group.ID)
			require.NoError(t, err)
			require.Equal(t, bannerID, 1)

			bannerID, err = d.DatabaseSelectFromRotation(1, group.ID)
			require.NoError(t, err)
			require.Equal(t, bannerID, 2)
		}
	})

	t.Run("non existing select", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`)
		groups, _ := d.DatabaseGetGroups()
		for _, group := range groups {
			bannerID, notInError := d.DatabaseSelectFromRotation(1, group.ID)
			require.Error(t, notInError)
			require.Equal(t, bannerID, invalidID)
		}
	})
}

func TestRegisterTransition(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()
	setTestData(d)

	t.Run("simple register", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`)
		_ = d.DatabaseAddToRotation(2, 1)

		groups, _ := d.DatabaseGetGroups()
		for _, group := range groups {
			err := d.DatabaseRegisterTransition(1, 2, group.ID)
			require.NoError(t, err)
			row := d.db.QueryRow(`SELECT click_count FROM "Statistic"
	WHERE slot_id=$1 AND banner_id=$2 AND group_id=$3`, 1, 2, group.ID)
			count := 0
			_ = row.Scan(&count)
			require.Equal(t, count, 1)
		}
	})

	t.Run("register for not in rotation", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`)
		_ = d.DatabaseAddToRotation(1, 2)

		groups, _ := d.DatabaseGetGroups()
		for _, group := range groups {
			err := d.DatabaseRegisterTransition(1, 1, group.ID)
			require.ErrorIs(t, err, ErrNotInRotation)
			err = d.DatabaseRegisterTransition(2, 2, group.ID)
			require.ErrorIs(t, err, ErrNotInRotation)
		}

		err := d.DatabaseRegisterTransition(1, 10, 2)
		require.ErrorIs(t, err, ErrNotInRotation)
	})
}
