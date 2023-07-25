package database

import (
	"strconv"
	"testing"

	"github.com/SergeyTyurin/banner-rotation/configs"
	"github.com/SergeyTyurin/banner-rotation/structures"
	"github.com/stretchr/testify/require"
)

func TestCreateGroup(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()
	_, _ = d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)

	group := structures.Group{ID: 1, Info: "info"}
	newGroup, err := d.DatabaseCreateGroup(group)
	require.NoError(t, err)
	require.Equal(t, newGroup.ID, group.ID)

	newGroup, err = d.DatabaseCreateGroup(group)
	require.NoError(t, err)
	require.NotEqual(t, newGroup.ID, group.ID)
}

func TestGetGroups(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()

	count := 5
	t.Run("get all groups", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)
		for i := 0; i < count; i++ {
			group := structures.Group{Info: "info" + strconv.Itoa(i)}
			_, _ = d.DatabaseCreateGroup(group)
		}

		groups, err := d.DatabaseGetGroups()
		require.NoError(t, err)
		require.Equal(t, len(groups), count)
	})

	t.Run("get groupby id", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)
		for i := 0; i < count; i++ {
			group := structures.Group{Info: "info" + strconv.Itoa(i+1)}
			_, _ = d.DatabaseCreateGroup(group)
		}

		group, err := d.DatabaseGetGroup(2)
		require.NoError(t, err)
		require.Equal(t, group.ID, 2)
		require.Equal(t, group.Info, "info2")
	})

	t.Run("get from empty", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)
		groups, err := d.DatabaseGetGroups()
		require.NoError(t, err)
		require.Empty(t, groups)

		group, err := d.DatabaseGetGroup(1)
		require.ErrorIs(t, err, ErrNotExist)
		require.Equal(t, group.ID, invalidID)
		require.Empty(t, group.Info)
	})
}

func TestUpdateGroup(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()

	t.Run("update non existed group", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)

		group := structures.Group{ID: 100, Info: "new info"}
		err := d.DatabaseUpdateGroup(group)
		require.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("update existed group", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)
		newGroup, _ := d.DatabaseCreateGroup(structures.Group{Info: "info"})
		newGroup.Info = newInfo
		err := d.DatabaseUpdateGroup(newGroup)
		require.NoError(t, err)

		updated, _ := d.DatabaseGetGroup(newGroup.ID)
		require.Equal(t, updated.Info, newInfo)
	})
}

func TestDeleteGroup(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()

	t.Run("delete non existed group", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)

		err := d.DatabaseDeleteGroup(1)
		require.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("delete existed group", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)
		newGroup, _ := d.DatabaseCreateGroup(structures.Group{Info: "info"})
		newGroup.Info = newInfo
		err := d.DatabaseDeleteGroup(newGroup.ID)
		require.NoError(t, err)

		groups, _ := d.DatabaseGetGroups()
		require.Empty(t, groups)
	})
}
