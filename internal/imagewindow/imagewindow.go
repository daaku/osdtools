// Package imagewindow provides utilities for managing OSD windows.
package imagewindow

import (
	"bytes"
	"image"
	"image/png"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
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

const styles = `
window {
  background: rgba(5%, 5%, 5%, 0.6);
}
`

func ConfigureDefaultStyles() error {
	cssProvider, err := gtk.CssProviderNew()
	if err != nil {
		return errors.WithStack(err)
	}
	if err := cssProvider.LoadFromData(styles); err != nil {
		return errors.WithStack(err)
	}
	display, err := gdk.DisplayGetDefault()
	if err != nil {
		return errors.WithStack(err)
	}
	screen, err := display.GetDefaultScreen()
	if err != nil {
		return errors.WithStack(err)
	}
	gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	return nil
}
