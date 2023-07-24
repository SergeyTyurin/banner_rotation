package configs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDBConfig(t *testing.T) {
	conn, err := GetDBConnectionConfig("../config/test/test_connection_config.yaml")
	require.Nil(t, err)
	require.NotNil(t, conn)
}

func TestCreateAppConfig(t *testing.T) {
	conn, err := GetAppSettings("../config/test/test_connection_config.yaml")
	require.Nil(t, err)
	require.NotNil(t, conn)
}

func TestCreateMsgBrokerConfig(t *testing.T) {
	conn, err := GetMessageBrokerConfig("../config/test/test_connection_config.yaml")
	require.Nil(t, err)
	require.NotNil(t, conn)
}
