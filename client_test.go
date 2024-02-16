package gphoto2cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientListPorts(t *testing.T) {
	c := NewClient()
	ports, err := c.ListPorts(context.Background())
	require.NoError(t, err)
	require.NotEqual(t, len(ports), 0)
}

func TestClientGetConfig(t *testing.T) {
	c := NewClient()
	v, err := c.GetConfig(context.Background(), nil, "serialnumber")
	require.NoError(t, err)
	require.Equal(t, v.Label, "Serial Number")
}

func TestClientListConfig(t *testing.T) {
	c := NewClient()
	cfgs, err := c.ListConfig(context.Background(), nil)
	require.NoError(t, err)

	require.NotEqual(t, len(cfgs), 0)
}

func TestClientCamera(t *testing.T) {
	c := NewClient()
	camera, err := c.GetCamera(context.Background(), "")
	require.NoError(t, err)
	require.NotNil(t, camera)
}

func TestClientCameraWithSerial(t *testing.T) {
	c := NewClient()
	camera, err := c.GetCamera(context.Background(), "6aedfeb1be470fdf85d710953be98415")
	require.NoError(t, err)
	require.NotNil(t, camera)
}