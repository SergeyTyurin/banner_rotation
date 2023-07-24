package database

import (
	"testing"

	"github.com/SergeyTyurin/banner_rotation/internal/configs"

	"github.com/stretchr/testify/require"
)

func TestConnectDatabase(t *testing.T) {
	t.Run("valid connection", func(t *testing.T) {
		config, _ := configs.GetDBConnectionConfig("../../config/connection_config.yaml")
		closeConnection, err := NewDatabase().Connect(config)
		require.NoError(t, err)
		require.NoError(t, closeConnection())
		defer closeConnection()
	})
	t.Run("empty connection", func(t *testing.T) {
		closeConnection, err := NewDatabase().Connect(nil)
		require.Error(t, err)
		require.Nil(t, closeConnection)
	})
}

func TestCloseConnection(t *testing.T) {
	config, _ := configs.GetDBConnectionConfig("../../config/connection_config.yaml")
	var d databaseImpl
	for i := 0; i < 10; i++ {
		closeConnection, err := d.Connect(config)
		require.NoError(t, err)
		require.NotNil(t, closeConnection)
		require.Equal(t, d.db.Stats().OpenConnections, 1)
		defer closeConnection()
	}
}
