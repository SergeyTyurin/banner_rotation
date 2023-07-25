package database

import (
	"strconv"
	"testing"

	"github.com/SergeyTyurin/banner-rotation/configs"
	"github.com/SergeyTyurin/banner-rotation/structures"
	"github.com/stretchr/testify/require"
)

func TestCreateBanner(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()
	_, _ = d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)

	banner := structures.Banner{ID: 1, Info: "info"}
	newBanner, err := d.DatabaseCreateBanner(banner)
	require.NoError(t, err)
	require.Equal(t, newBanner.ID, banner.ID)

	newBanner, err = d.DatabaseCreateBanner(banner)
	require.NoError(t, err)
	require.NotEqual(t, newBanner.ID, banner.ID)
}

func TestGetBanners(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()

	count := 5
	t.Run("get all banners", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)
		for i := 0; i < count; i++ {
			banner := structures.Banner{Info: "info" + strconv.Itoa(i)}
			_, _ = d.DatabaseCreateBanner(banner)
		}

		banners, err := d.DatabaseGetBanners()
		require.NoError(t, err)
		require.Equal(t, len(banners), count)
	})

	t.Run("get banner by id", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)
		for i := 0; i < count; i++ {
			banner := structures.Banner{Info: "info" + strconv.Itoa(i+1)}
			_, _ = d.DatabaseCreateBanner(banner)
		}

		banner, err := d.DatabaseGetBanner(2)
		require.NoError(t, err)
		require.Equal(t, banner.ID, 2)
		require.Equal(t, banner.Info, "info2")
	})

	t.Run("get from empty", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)
		banners, err := d.DatabaseGetBanners()
		require.NoError(t, err)
		require.Empty(t, banners)

		banner, err := d.DatabaseGetBanner(1)
		require.ErrorIs(t, err, ErrNotExist)
		require.Equal(t, banner.ID, invalidID)
		require.Empty(t, banner.Info)
	})
}

func TestUpdateBanner(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()

	t.Run("update non existed banner", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)

		banner := structures.Banner{ID: 100, Info: "new info"}
		err := d.DatabaseUpdateBanner(banner)
		require.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("update existed banner", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)
		newBanner, err := d.DatabaseCreateBanner(structures.Banner{Info: "info"})
		require.NoError(t, err)
		newBanner.Info = newInfo
		err = d.DatabaseUpdateBanner(newBanner)
		require.NoError(t, err)

		updated, _ := d.DatabaseGetBanner(newBanner.ID)
		require.Equal(t, updated.Info, newInfo)
	})
}

func TestDeleteBanner(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	closeConnection, _ := d.DatabaseConnect(config)
	defer func() {
		_ = closeConnection()
	}()

	t.Run("delete non existed banner", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)

		err := d.DatabaseDeleteBanner(1)
		require.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("delete existed banner", func(t *testing.T) {
		_, _ = d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)
		newBanner, _ := d.DatabaseCreateBanner(structures.Banner{Info: "info"})
		newBanner.Info = newInfo
		err := d.DatabaseDeleteBanner(newBanner.ID)
		require.NoError(t, err)

		banners, _ := d.DatabaseGetBanners()
		require.Empty(t, banners)
	})
}
