package database

import (
	"strconv"
	"testing"

	"github.com/SergeyTyurin/banner-rotation/configs"
	"github.com/SergeyTyurin/banner-rotation/structures"
	"github.com/stretchr/testify/require"
)

func TestCreateSlot(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()
	_, _ = d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`)

	slot := structures.Slot{ID: 1, Info: "info"}
	newSlot, err := d.DatabaseCreateSlot(slot)
	require.NoError(t, err)
	require.Equal(t, newSlot.ID, slot.ID)

	newSlot, err = d.DatabaseCreateSlot(slot)
	require.NoError(t, err)
	require.NotEqual(t, newSlot.ID, slot.ID)
}

func TestGetSlots(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()

	count := 5
	t.Run("get all slots", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`)
		for i := 0; i < count; i++ {
			slot := structures.Slot{Info: "info" + strconv.Itoa(i)}
			_, _ = d.DatabaseCreateSlot(slot)
		}

		slots, err := d.DatabaseGetSlots()
		require.NoError(t, err)
		require.Equal(t, len(slots), count)
	})

	t.Run("get slotby ID", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`)
		for i := 0; i < count; i++ {
			slot := structures.Slot{Info: "info" + strconv.Itoa(i+1)}
			_, _ = d.DatabaseCreateSlot(slot)
		}

		slot, err := d.DatabaseGetSlot(2)
		require.NoError(t, err)
		require.Equal(t, slot.ID, 2)
		require.Equal(t, slot.Info, "info2")
	})

	t.Run("get from empty", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`)
		slots, err := d.DatabaseGetSlots()
		require.NoError(t, err)
		require.Empty(t, slots)

		slot, err := d.DatabaseGetSlot(1)
		require.ErrorIs(t, err, ErrNotExist)
		require.Equal(t, slot.ID, invalidID)
		require.Empty(t, slot.Info)
	})
}

func TestUpdateSlot(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()

	t.Run("update non existed slot", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`)

		slot := structures.Slot{ID: 100, Info: "new info"}
		err := d.DatabaseUpdateSlot(slot)
		require.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("update existed slot", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`)
		newSlot, _ := d.DatabaseCreateSlot(structures.Slot{Info: "info"})
		newSlot.Info = newInfo
		err := d.DatabaseUpdateSlot(newSlot)
		require.NoError(t, err)

		updated, _ := d.DatabaseGetSlot(newSlot.ID)
		require.Equal(t, updated.Info, newInfo)
	})
}

func TestDeleteSlot(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()

	t.Run("delete non existed slot", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`)

		err := d.DatabaseDeleteSlot(1)
		require.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("delete existed slot", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Slots" RESTART IDENTITY CASCADE`)
		newSlot, _ := d.DatabaseCreateSlot(structures.Slot{Info: "info"})
		newSlot.Info = newInfo
		err := d.DatabaseDeleteSlot(newSlot.ID)
		require.NoError(t, err)

		slots, _ := d.DatabaseGetSlots()
		require.Empty(t, slots)
	})
}
