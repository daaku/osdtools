// Package imagewindow provides utilities for managing OSD windows.
package imagewindow

import (
	"image"
	"image/draw"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/pkg/errors"
)

// ImageWindow is a helper to render X windows with images.
type ImageWindow struct {
	X          *xgbutil.XUtil
	Window     *xwindow.Window
	ScreenSize image.Point
}

// Destroy cleans up the associated X windows.
func (w *ImageWindow) Destroy() {
	w.Window.Destroy()
}

// SetOpacity sets the opacity of the window.
func (w *ImageWindow) SetOpacity(opacity float64) {
	ewmh.WmWindowOpacitySet(w.X, w.Window.Id, opacity)
}

// Draw updates the window with a new image
func (w *ImageWindow) Draw(img image.Image) error {
	size := img.Bounds().Size()
	width := size.X
	height := size.Y
	x := (w.ScreenSize.X - width) / 2
	y := (w.ScreenSize.Y - height) / 2

	ewmh.MoveresizeWindow(w.X, w.Window.Id, x, y, width, height)

	ximg := xgraphics.New(w.X, image.Rect(0, 0, width, height))
	draw.Draw(ximg, img.Bounds(), img, image.ZP, draw.Over)
	if err := ximg.CreatePixmap(); err != nil {
		return errors.WithStack(err)
	}
	defer ximg.Destroy()

	ximg.XDraw()
	ximg.XExpPaint(w.Window.Id, 0, 0)
	return nil
}

// New creates a new ImageWindow.
func New() (*ImageWindow, error) {
	X, err := xgbutil.NewConn()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	win, err := xwindow.Generate(X)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	win.Create(X.RootWin(), -1, -1, 1, 1, xproto.CwOverrideRedirect, 1)
	win.Map()

	screen := X.Screen()
	return &ImageWindow{
		X:          X,
		Window:     win,
		ScreenSize: image.Pt(int(screen.WidthInPixels), int(screen.HeightInPixels)),
	}, nil
}
