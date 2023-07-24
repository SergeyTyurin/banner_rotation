package database

import (
	"strconv"
	"testing"

	"github.com/SergeyTyurin/banner_rotation/configs"
	"github.com/SergeyTyurin/banner_rotation/structures"
	"github.com/stretchr/testify/require"
)

func TestCreateGroup(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection()
	d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)

	group := structures.Group{Id: 1, Info: "info"}
	new, err := d.CreateGroup(group)
	require.NoError(t, err)
	require.Equal(t, new.Id, group.Id)

	new, err = d.CreateGroup(group)
	require.NoError(t, err)
	require.NotEqual(t, new.Id, group.Id)
}

func TestGetGroups(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection()

	count := 5
	t.Run("get all groups", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)
		for i := 0; i < count; i++ {
			group := structures.Group{Info: "info" + strconv.Itoa(i)}
			d.CreateGroup(group)
		}

		groups, err := d.GetGroups()
		require.NoError(t, err)
		require.Equal(t, len(groups), count)
	})

	t.Run("get groupby id", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)
		for i := 0; i < count; i++ {
			group := structures.Group{Info: "info" + strconv.Itoa(i+1)}
			d.CreateGroup(group)
		}

		group, err := d.GetGroup(2)
		require.NoError(t, err)
		require.Equal(t, group.Id, 2)
		require.Equal(t, group.Info, "info2")
	})

	t.Run("get from empty", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)
		groups, err := d.GetGroups()
		require.NoError(t, err)
		require.Empty(t, groups)

		group, err := d.GetGroup(1)
		require.ErrorIs(t, err, ErrNotExist)
		require.Equal(t, group.Id, invalidId)
		require.Empty(t, group.Info)
	})
}

func TestUpdateGroup(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection()

	t.Run("update non existed group", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)

		group := structures.Group{Id: 100, Info: "new info"}
		err := d.UpdateGroup(group)
		require.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("update existed group", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)
		newGroup, _ := d.CreateGroup(structures.Group{Info: "info"})
		newGroup.Info = "newInfo"
		err := d.UpdateGroup(newGroup)
		require.NoError(t, err)

		updated, _ := d.GetGroup(newGroup.Id)
		require.Equal(t, updated.Info, "newInfo")
	})
}

func TestDeleteGroup(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection()

	t.Run("delete non existed group", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)

		err := d.DeleteGroup(1)
		require.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("delete existed group", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Groups" RESTART IDENTITY CASCADE`)
		newGroup, _ := d.CreateGroup(structures.Group{Info: "info"})
		newGroup.Info = "newInfo"
		err := d.DeleteGroup(newGroup.Id)
		require.NoError(t, err)

		groups, _ := d.GetGroups()
		require.Empty(t, groups)
	})
}
