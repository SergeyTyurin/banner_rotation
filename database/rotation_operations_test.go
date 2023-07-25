package database

import (
	"strconv"
	"testing"

	"github.com/SergeyTyurin/banner_rotation/configs"
	"github.com/SergeyTyurin/banner_rotation/structures"
	"github.com/stretchr/testify/require"
)

func setTestData(d databaseImpl) {
	d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)  //nolint:all
	d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`)   //nolint:all
	d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`) //nolint:all

	for i := 0; i < 10; i++ {
		d.CreateBanner(structures.Banner{Info: "banner_" + strconv.Itoa(i+1)}) //nolint:all
	}
	for i := 0; i < 2; i++ {
		d.CreateGroup(structures.Group{Info: "group_" + strconv.Itoa(i+1)}) //nolint:all
	}
	for i := 0; i < 4; i++ {
		d.CreateSlot(structures.Slot{Info: "slot_" + strconv.Itoa(i+1)}) //nolint:all
	}
}

func TestAddToRotation(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection() //nolint:all
	setTestData(d)

	t.Run("simple add", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`) //nolint:all
		slotId := 1
		bannerId := 2
		err := d.AddToRotation(bannerId, slotId)
		require.NoError(t, err)

		row := d.db.QueryRow(`SELECT count(*) FROM "Statistic"
		WHERE slot_id=$1 AND banner_id=$2`, slotId, bannerId)

		count := 0
		row.Scan(&count) //nolint:all
		//Test Add for each group
		groups, _ := d.GetGroups()
		require.Equal(t, count, len(groups))
	})

	t.Run("add non existed intities", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`) //nolint:all

		err := d.AddToRotation(20, 1)
		require.ErrorIs(t, err, ErrNotExist)

		err = d.AddToRotation(1, 20)
		require.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("add already in rotation", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`) //nolint:all
		err := d.AddToRotation(1, 1)
		require.NoError(t, err)

		err = d.AddToRotation(1, 1)
		require.ErrorIs(t, err, ErrAlreadyInRotation)
	})
}

func TestDeleteFromRotation(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection() //nolint:all
	setTestData(d)

	t.Run("simple delete", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`) //nolint:all
		slotId := 1
		bannerId := 2
		d.AddToRotation(bannerId, slotId) //nolint:all
		err := d.DeleteFromRotation(bannerId, slotId)
		require.NoError(t, err)

		row := d.db.QueryRow(`SELECT count(*) FROM "Statistic"
	WHERE slot_id=$1 AND banner_id=$2`, slotId, bannerId)
		count := 0
		row.Scan(&count) //nolint:all
		require.Equal(t, count, 0)
	})

	t.Run("delete not in rotation", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`) //nolint:all
		d.AddToRotation(1, 1)                                            //nolint:all
		d.AddToRotation(2, 1)                                            //nolint:all
		d.AddToRotation(1, 2)                                            //nolint:all
		err := d.DeleteFromRotation(2, 2)
		require.ErrorIs(t, err, ErrNotInRotation)
	})
}

func TestSelectFromRotation(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection() //nolint:all
	setTestData(d)

	t.Run("existing select", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`) //nolint:all
		_ = d.AddToRotation(1, 1)
		_ = d.AddToRotation(2, 1)
		groups, _ := d.GetGroups()
		for _, group := range groups {
			banner_id, err := d.SelectFromRotation(1, group.Id)
			require.NoError(t, err)
			require.Equal(t, banner_id, 1)

			banner_id, err = d.SelectFromRotation(1, group.Id)
			require.NoError(t, err)
			require.Equal(t, banner_id, 2)
		}
	})

	t.Run("non existing select", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`) //nolint:all
		groups, _ := d.GetGroups()
		for _, group := range groups {
			banner_id, notInError := d.SelectFromRotation(1, group.Id)
			require.Error(t, notInError)
			require.Equal(t, banner_id, invalidId)
		}
	})
}

func TestRegisterTransition(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection() //nolint:all
	setTestData(d)

	t.Run("simple register", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`) //nolint:all
		_ = d.AddToRotation(2, 1)

		groups, _ := d.GetGroups()
		for _, group := range groups {
			err := d.RegisterTransition(1, 2, group.Id)
			require.NoError(t, err)
			row := d.db.QueryRow(`SELECT click_count FROM "Statistic"
	WHERE slot_id=$1 AND banner_id=$2 AND group_id=$3`, 1, 2, group.Id)
			count := 0
			row.Scan(&count) //nolint:all
			require.Equal(t, count, 1)
		}
	})

	t.Run("register for not in rotation", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Statistic" RESTART IDENTITY CASCADE`) //nolint:all
		_ = d.AddToRotation(1, 2)

		groups, _ := d.GetGroups()
		for _, group := range groups {
			err := d.RegisterTransition(1, 1, group.Id)
			require.ErrorIs(t, err, ErrNotInRotation)
			err = d.RegisterTransition(2, 2, group.Id)
			require.ErrorIs(t, err, ErrNotInRotation)
		}

		err := d.RegisterTransition(1, 10, 2)
		require.ErrorIs(t, err, ErrNotInRotation)
	})
}
