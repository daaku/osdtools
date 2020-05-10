// Package imagewindow provides utilities for managing OSD windows.
package imagewindow

import (
	"bytes"
	"image"
	"image/png"

	"github.com/gotk3/gotk3/gdk"
	"github.com/pkg/errors"
)

var pngEncoder = png.Encoder{
	CompressionLevel: png.BestSpeed,
}

func ImageToPixbuf(img image.Image) (*gdk.Pixbuf, error) {
	var buf bytes.Buffer
	if err := pngEncoder.Encode(&buf, img); err != nil {
		return nil, errors.WithStack(err)
	}
	loader, err := gdk.PixbufLoaderNewWithType("png")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	pixbuf, err := loader.WriteAndReturnPixbuf(buf.Bytes())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return pixbuf, nil
}
