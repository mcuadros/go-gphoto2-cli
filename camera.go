package gphoto2cli

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Camera struct {
	cli *Client
	p   *Port
	sync.Mutex
}

func NewCamera(cli *Client, p *Port) *Camera {
	return &Camera{
		cli: cli,
		p:   p,
	}
}

func (c *Camera) Port() *Port {
	return c.p
}

func (c *Camera) CapturePhoto(ctx context.Context, filename string, flags ...string) (*CaptureResult, error) {
	flags = append(flags, "--set-config=capturetarget=0")
	flags = append(flags, "--force-overwrite")
	flags = append(flags, "--capture-image-and-download")

	if filename != "" {
		flags = append(flags, fmt.Sprintf("--filename=%s", filename))

	}
	output, err := c.cli.runCommand(ctx, c.p, flags...)
	if err != nil {
		return nil, fmt.Errorf("error capturing at port %s: %w", c.p.Port, err)
	}

	return c.cli.parser.ParseCapture(output), nil
}

func (c *Camera) CaptureVideo(ctx context.Context, filename string, d time.Duration, flags ...string) (*CaptureResult, error) {
	flags = append(flags, "--set-config=viewfinder=1")
	flags = append(flags, "--set-config=capturetarget=1")
	flags = append(flags, "--set-config=movierecordtarget=0")
	flags = append(flags, "--wait-event", fmt.Sprintf("%ds", int(d.Seconds())))
	flags = append(flags, "--set-config=movierecordtarget=1")
	flags = append(flags, "--wait-event-and-download", "2s")
	flags = append(flags, "--set-config=viewfinder=0")
	flags = append(flags, "--force-overwrite")

	if filename != "" {
		flags = append(flags, fmt.Sprintf("--filename=%s", filename))
	}

	output, err := c.cli.runCommand(ctx, c.p, flags...)
	if err != nil {
		return nil, fmt.Errorf("error capturing at port %s: %w", c.p.Port, err)
	}

	return c.cli.parser.ParseCapture(output), nil
}
