package database

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateBanner(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../../config/connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection()
	d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)

	banner := structures.Banner{Id: 1, Info: "info"}
	new, err := d.CreateBanner(banner)
	require.NoError(t, err)
	require.Equal(t, new.Id, banner.Id)

	new, err = d.CreateBanner(banner)
	require.NoError(t, err)
	require.NotEqual(t, new.Id, banner.Id)
}

func TestGetBanners(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../../config/connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection()

	count := 5
	t.Run("get all banners", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)
		for i := 0; i < count; i++ {
			banner := structures.Banner{Info: "info" + strconv.Itoa(i)}
			d.CreateBanner(banner)
		}

		banners, err := d.GetBanners()
		require.NoError(t, err)
		require.Equal(t, len(banners), count)
	})

	t.Run("get banner by id", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)
		for i := 0; i < count; i++ {
			banner := structures.Banner{Info: "info" + strconv.Itoa(i+1)}
			d.CreateBanner(banner)
		}

		banner, err := d.GetBanner(2)
		require.NoError(t, err)
		require.Equal(t, banner.Id, 2)
		require.Equal(t, banner.Info, "info2")
	})

	t.Run("get from empty", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)
		banners, err := d.GetBanners()
		require.NoError(t, err)
		require.Empty(t, banners)

		banner, err := d.GetBanner(1)
		require.ErrorIs(t, err, ErrNotExist)
		require.Equal(t, banner.Id, invalidId)
		require.Empty(t, banner.Info)
	})
}

func TestUpdateBanner(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../../config/connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection()

	t.Run("update non existed banner", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)

		banner := structures.Banner{Id: 100, Info: "new info"}
		err := d.UpdateBanner(banner)
		require.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("update existed banner", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)
		newBanner, err := d.CreateBanner(structures.Banner{Info: "info"})
		require.NoError(t, err)
		newBanner.Info = "newInfo"
		err = d.UpdateBanner(newBanner)
		require.NoError(t, err)

		updated, _ := d.GetBanner(newBanner.Id)
		require.Equal(t, updated.Info, "newInfo")
	})
}

func TestDeleteBanner(t *testing.T) {
	d := databaseImpl{nil}
	config, _ := configs.GetDBConnectionConfig("../../config/connection_config.yaml")
	closeConnection, _ := d.Connect(config)
	defer closeConnection()

	t.Run("delete non existed banner", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)

		err := d.DeleteBanner(1)
		require.ErrorIs(t, err, ErrNotExist)
	})

	t.Run("delete existed banner", func(t *testing.T) {
		d.db.Exec(`TRUNCATE TABLE "Banners" RESTART IDENTITY CASCADE`)
		newBanner, _ := d.CreateBanner(structures.Banner{Info: "info"})
		newBanner.Info = "newInfo"
		err := d.DeleteBanner(newBanner.Id)
		require.NoError(t, err)

		banners, _ := d.GetBanners()
		require.Empty(t, banners)
	})
}
