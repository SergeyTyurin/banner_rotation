package database

import (
	"testing"

	"github.com/SergeyTyurin/banner-rotation/configs"
	"github.com/stretchr/testify/require"
)

const newInfo = "newInfo"

func TestConnectDatabase(t *testing.T) {
	t.Run("valid connection", func(t *testing.T) {
		config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
		closeConnection, err := NewDatabase().DatabaseConnect(config)
		require.NoError(t, err)
		require.NoError(t, closeConnection())
		defer func() {
			_ = closeConnection()
		}()
	})
	t.Run("empty connection", func(t *testing.T) {
		closeConnection, err := NewDatabase().DatabaseConnect(nil)
		require.Error(t, err)
		require.Nil(t, closeConnection)
	})
}

func TestCloseConnection(t *testing.T) {
	config, _ := configs.GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	var d databaseImpl
	for i := 0; i < 10; i++ {
		closeConnection, err := d.DatabaseConnect(config)
		require.NoError(t, err)
		require.NotNil(t, closeConnection)
		require.Equal(t, d.db.Stats().OpenConnections, 1)
		defer func() {
			_ = closeConnection()
		}()
	}
}
