package message_broker

import (
	"testing"

	"github.com/SergeyTyurin/banner_rotation/configs"
	"github.com/stretchr/testify/require"
)

func TestSendMessages(t *testing.T) {
	config, _ := configs.GetMessageBrokerConfig("../config/test/test_connection_config.yaml")
	m := messageBrokerImpl{}
	closeFunc, err := m.Connect(config)
	require.NoError(t, err)
	defer closeFunc()

	require.NoError(t, m.SendRegisterTransitionEvent("test regsiter"))
	require.NoError(t, m.SendSelectFromRotationEvent("test select"))
}
