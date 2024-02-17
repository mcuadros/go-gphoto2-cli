# go-gphoto2-cli [![GoDoc](http://godoc.org/github.com/mcuadros/go-gphoto2-cli?status.svg)](http://godoc.org/github.com/mcuadros/go-gphoto2-cli)

Golang wrapper arround the gphoto2 cli.

[gphoto2](http://www.gphoto.org/), is a set of software applications and libraries
for use in digital photography. gPhoto supports not just retrieving of images from
camera devices, but also upload and remote controlled configuration and capture,
depending on whether the camera supports those features.

This library has a goal to cover basic features as capturing video and photos from
several cameras at the same time.

## Installation

The recommended way to install go-tsunami

```
go get github.com/mcuadros/go-gphoto2-cli
```

## Example

```go
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
```

## License

MIT, see [LICENSE](LICENSE)
