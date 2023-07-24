package message_broker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendMessages(t *testing.T) {
	config, _ := configs.GetMessageBrokerConfig("../../config/connection_config.yaml")
	m := messageBrokerImpl{}
	closeFunc, err := m.Connect(config)
	require.NoError(t, err)
	defer closeFunc()

	require.NoError(t, m.SendRegisterTransitionEvent("test regsiter"))
	require.NoError(t, m.SendSelectFromRotationEvent("test select"))
}
