package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	gphoto2cli "github.com/mcuadros/go-gphoto2-cli"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	cli := gphoto2cli.NewClient()
	ports, err := cli.ListPorts(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range ports {
		c := gphoto2cli.NewCamera(cli, p)
		r, err := c.CapturePhoto(ctx, "")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("photo captured from %q at port %s and stored in %q\n",
			p.Model, p.Port, r.LocalFile,
		)
	}

}
