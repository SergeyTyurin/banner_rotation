package database

import (
	"strconv"
	"testing"

	"github.com/SergeyTyurin/banner_rotation/configs"
	"github.com/SergeyTyurin/banner_rotation/structures"
	"github.com/stretchr/testify/require"
)

func TestCreateSlot(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection()                                      //nolint:all
	d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`) //nolint:all

	slot := structures.Slot{Id: 1, Info: "info"}
	new, err := d.CreateSlot(slot)
	require.NoError(t, err)
	require.Equal(t, new.Id, slot.Id)

	new, err = d.CreateSlot(slot)
	require.NoError(t, err)
	require.NotEqual(t, new.Id, slot.Id)
}

func TestGetSlots(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection() //nolint:all

	count := 5
	t.Run("get all slots", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`) //nolint:all // nolint:all
		for i := 0; i < count; i++ {
			slot := structures.Slot{Info: "info" + strconv.Itoa(i)}
			d.CreateSlot(slot) //nolint:all
		}

		slots, err := d.GetSlots()
		require.NoError(t, err)
		require.Equal(t, len(slots), count)
	})

	t.Run("get slotby id", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`) //nolint:all
		for i := 0; i < count; i++ {
			slot := structures.Slot{Info: "info" + strconv.Itoa(i+1)}
			d.CreateSlot(slot) //nolint:alls
		}

		slot, err := d.GetSlot(2)
		require.NoError(t, err)
		require.Equal(t, slot.Id, 2)
		require.Equal(t, slot.Info, "info2")
	})

	t.Run("get from empty", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`) //nolint:all
		slots, err := d.GetSlots()
		require.NoError(t, err)
		require.Empty(t, slots)

		slot, err := d.GetSlot(1)
		require.ErrorIs(t, err, ErrNotExist)
		require.Equal(t, slot.Id, invalidId)
		require.Empty(t, slot.Info)
	})
}

func TestUpdateSlot(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection() //nolint:all

	t.Run("update non existed slot", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`) //nolint:all

		slot := structures.Slot{Id: 100, Info: "new info"}
		err := d.UpdateSlot(slot)
		require.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("update existed slot", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`) //nolint:all
		newSlot, _ := d.CreateSlot(structures.Slot{Info: "info"})
		newSlot.Info = "newInfo"
		err := d.UpdateSlot(newSlot)
		require.NoError(t, err)

		updated, _ := d.GetSlot(newSlot.Id)
		require.Equal(t, updated.Info, "newInfo")
	})
}

func TestDeleteSlot(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection() //nolint:all

	t.Run("delete non existed slot", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`) //nolint:all

		err := d.DeleteSlot(1)
		require.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("delete existed slot", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`) //nolint:all
		newSlot, _ := d.CreateSlot(structures.Slot{Info: "info"})
		newSlot.Info = "newInfo"
		err := d.DeleteSlot(newSlot.Id)
		require.NoError(t, err)

		slots, _ := d.GetSlots()
		require.Empty(t, slots)
	})
}
