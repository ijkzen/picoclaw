package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGatewayCommand(t *testing.T) {
	cmd := NewGatewayCommand()

	require.NotNil(t, cmd)

	assert.Equal(t, "gateway", cmd.Use)
	assert.Equal(t, "Manage picoclaw gateway", cmd.Short)

	assert.Len(t, cmd.Aliases, 1)
	assert.True(t, cmd.HasAlias("g"))

	assert.Nil(t, cmd.Run)
	assert.Nil(t, cmd.RunE)

	assert.Nil(t, cmd.PersistentPreRun)
	assert.Nil(t, cmd.PersistentPostRun)

	// Should have subcommands: start, stop, status, run
	assert.True(t, cmd.HasSubCommands())

	// Gateway command itself has no flags anymore
	assert.False(t, cmd.HasFlags())
}

func TestNewStartCommand(t *testing.T) {
	cmd := NewStartCommand()

	require.NotNil(t, cmd)

	assert.Equal(t, "start", cmd.Use)
	assert.Equal(t, "Start picoclaw gateway in the background", cmd.Short)

	assert.True(t, cmd.HasFlags())
	assert.NotNil(t, cmd.Flags().Lookup("debug"))
}

func TestNewStopCommand(t *testing.T) {
	cmd := NewStopCommand()

	require.NotNil(t, cmd)

	assert.Equal(t, "stop", cmd.Use)
	assert.Equal(t, "Stop the background picoclaw gateway", cmd.Short)
}

func TestNewStatusCommand(t *testing.T) {
	cmd := NewStatusCommand()

	require.NotNil(t, cmd)

	assert.Equal(t, "status", cmd.Use)
	assert.Equal(t, "Show picoclaw gateway status", cmd.Short)
}

func TestNewRunCommand(t *testing.T) {
	cmd := NewRunCommand()

	require.NotNil(t, cmd)

	assert.Equal(t, "run", cmd.Use)
	assert.Equal(t, "Run picoclaw gateway in foreground (internal use)", cmd.Short)

	assert.True(t, cmd.HasFlags())
	assert.NotNil(t, cmd.Flags().Lookup("debug"))
}
