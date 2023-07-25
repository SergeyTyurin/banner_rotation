package messagebroker

import (
	"testing"

	"github.com/SergeyTyurin/banner-rotation/configs"
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

	msg, err := m.GetRegisterTransitionEvent()
	require.NoError(t, err)
	require.Equal(t, msg, "test regsiter")

	msg, err = m.GetSelectFromRotationEvent()
	require.NoError(t, err)
	require.Equal(t, msg, "test select")
}
