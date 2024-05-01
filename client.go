package gphoto2cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	SerialNumberConfig = "serialnumber"
)

type Client struct {
	parser
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) ListPorts(ctx context.Context) ([]*Port, error) {
	output, err := c.runCommand(ctx, nil, "--auto-detect")
	if err != nil {
		return nil, fmt.Errorf("error listing usb ports: %w", err)
	}

	return c.ParseAutoDetect(output), nil
}

func (c *Client) GetConfig(ctx context.Context, p *Port, name string) (*ConfigValue, error) {
	output, err := c.runCommand(ctx, p, fmt.Sprintf("--get-config=%s", name))
	if err != nil {
		return nil, fmt.Errorf("error getting config value: %w", err)
	}

	return c.ParseGetConfig(output), nil
}

func (c *Client) ListConfig(ctx context.Context, p *Port) ([]string, error) {
	output, err := c.runCommand(ctx, p, "--list-config")
	if err != nil {
		return nil, fmt.Errorf("error listing config: %w", err)
	}

	return strings.Split(output, "\n"), nil
}

func (c *Client) Camera(ctx context.Context, serial string) (*Camera, error) {
	ports, err := c.ListPorts(ctx)
	if err != nil {
		return nil, err
	}

	if len(ports) == 0 {
		return nil, fmt.Errorf("no cameras available")
	}

	for _, p := range ports {
		v, err := c.GetConfig(ctx, p, SerialNumberConfig)
		if err != nil {
			return nil, err
		}

		if v.Current == serial {
			return NewCamera(c, p), nil
		}
	}

	return nil, fmt.Errorf("unable to find a camara with serial number: %s", serial)
}

func (c *Client) getBinary() string {
	return "gphoto2"
}

func (c *Client) runCommand(ctx context.Context, p *Port, flags ...string) (string, error) {
	if p != nil {
		p.Lock()
		defer p.Unlock()

		flags = append(flags, fmt.Sprintf("--port=%s", p.Port))
	}

	cmd := exec.CommandContext(ctx, c.getBinary(), flags...)

	debug := os.Getenv("GPHOTO_DEBUG") != ""
	if debug {
		fmt.Println(">", cmd)
	}

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if debug {
		cmd.Stdout = io.MultiWriter(os.Stdout, cmd.Stdout)
		cmd.Stderr = io.MultiWriter(os.Stderr, cmd.Stderr)
	}

	if err := cmd.Run(); err != nil {
		if err := c.ParseErrorMessage(errb.String()); err != nil {
			return "", err
		}

		return "", err
	}

	if err := c.ParseErrorMessage(errb.String()); err != nil {
		return "", err
	}

	return outb.String(), nil
}
