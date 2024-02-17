package gphoto2cli

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var SerialNumber string

func init() {
	SerialNumber = os.Getenv("SERIAL_NUMBER")
}

func TestCameraCapturePhoto(t *testing.T) {
	if SerialNumber == "" {
		t.Skip("SERIAL_NUMBER needs to be provided")
	}

	c := NewClient()
	camera, err := c.Camera(context.Background(), SerialNumber)
	require.NoError(t, err)

	folder, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	filename := filepath.Join(folder, "temp.jpg")
	defer func() {
		err = os.Remove(filename)
		require.NoError(t, err)
	}()

	r, err := camera.CapturePhoto(context.Background(), filename)
	require.NoError(t, err)
	require.Equal(t, r.Deleted, true)
	require.True(t, len(r.OnCamera) != 0)
	require.True(t, len(r.LocalFile) != 0)
}

func TestCameraCaptureVideo(t *testing.T) {
	if SerialNumber == "" {
		t.Skip("SERIAL_NUMBER needs to be provided")
	}

	c := NewClient()
	camera, err := c.Camera(context.Background(), SerialNumber)
	require.NoError(t, err)

	folder, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	filename := filepath.Join(folder, "temp.jpg")
	defer func() {
		err = os.Remove(filename)
		require.NoError(t, err)
	}()

	r, err := camera.CaptureVideo(context.Background(), filename, time.Second*5)
	require.NoError(t, err)
	require.Equal(t, r.Deleted, false)
	require.True(t, len(r.OnCamera) == 0)
	require.True(t, len(r.LocalFile) != 0)
}
