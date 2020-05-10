package bar

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"net"
	"net/rpc"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/daaku/osdtools/internal/imagewindow"
	fontloader "github.com/fxkr/go-freetype-fontloader"
	"github.com/golang/freetype/truetype"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func socketPath(name string) string {
	u, _ := user.Current()
	home := ""
	if u == nil {
		home = os.Getenv("HOME")
	} else {

		home = u.HomeDir
	}
	return filepath.Join(home, name)
}

// Render renders an icon + percent + bar image.
func Render(icon image.Image, pct int) (*image.RGBA, error) {
	if pct < 0 || pct > 100 {
		return nil, errors.New("value must be between 0-100")
	}

	iconBounds := icon.Bounds()

	const fontSize = 24
	sansFont, err := fontloader.LoadCache("Sans")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	face := truetype.NewFace(sansFont, &truetype.Options{Size: fontSize})
	lineHeight := face.Metrics().Height.Round()
	height := (lineHeight * 3) + face.Metrics().Descent.Round() + iconBounds.Dy() + 8
	width := height
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Over)

	imgBounds := img.Bounds()
	iconR := iconBounds.Add(image.Pt((imgBounds.Dx()-iconBounds.Dx())/2, lineHeight/2))
	draw.DrawMask(img, iconR, image.White, image.ZP, icon, image.ZP, draw.Over)

	fontDrawer := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
	}

	barHeight := 30
	barOffset := 20
	gray := image.NewUniform(color.Gray16{0x9999})
	barRect := image.Rect(barOffset, height-barOffset-barHeight, width-barOffset, height-barOffset)
	draw.Draw(img, barRect, gray, image.ZP, draw.Over)
	filledWidth := barOffset + int(barRect.Dx()*pct/100)
	barRect = image.Rect(barOffset, height-barOffset-barHeight, filledWidth, height-barOffset)
	draw.Draw(img, barRect, image.White, image.ZP, draw.Over)

	pctStr := fmt.Sprint(pct)
	pctOffset := 10
	pctWidth := fontDrawer.MeasureString(pctStr).Round()
	fontDrawer.Dot = fixed.P((width-pctWidth)/2, barRect.Min.Y-pctOffset)
	fontDrawer.DrawString(pctStr)

	return img, nil
}

// ImageBarServer to handle RPC calls.
type ImageBarServer struct {
	pct chan int
}

// Set sets the pct.
func (s *ImageBarServer) Set(v int, ok *struct{}) error {
	s.pct <- v
	return nil
}

func sendPct(socketPath string, pct int) error {
	client, err := rpc.Dial("unix", socketPath)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := client.Call("ImageBarServer.Set", pct, nil); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func startServer(socketPath string) (chan int, net.Listener, error) {
	pct := make(chan int)
	rpc.Register(&ImageBarServer{pct: pct})
	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	go rpc.Accept(lis)
	return pct, lis, nil
}

// Run a bar application.
func Run(pct int, render func(int) (*image.RGBA, error), socketName string) error {
	const timeout = 5 * time.Second
	socketPath := socketPath(socketName)

	if err := sendPct(socketPath, pct); err == nil {
		return nil
	}

	pctCh, lis, err := startServer(socketPath)
	if err != nil {
		return err
	}
	defer lis.Close()

	img, err := render(pct)
	if err != nil {
		return err
	}

	appID := fmt.Sprintf("org.daaku%s", socketName)
	app, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return errors.WithStack(err)
	}
	var imgview *gtk.Image
	_, err = app.Connect("activate", func() {
		win, imgview_, err := imagewindow.NewImageWindow(app, img)
		imgview = imgview_
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}
		win.SetOpacity(0.7)
		win.SetTitle("OSD")
		win.ShowAll()
	})
	if err != nil {
		return errors.WithStack(err)
	}

	go func() {
		quitter := time.NewTimer(timeout)
		for {
			select {
			case <-quitter.C:
				_, _ = glib.IdleAdd(app.Quit)
				break
			case pct := <-pctCh:
				quitter.Reset(timeout)
				img, err := render(pct)
				if err != nil {
					panic(err)
				}
				pixbuf, err := imagewindow.ImageToPixbuf(img)
				if err != nil {
					panic(err)
				}
				imgview.SetFromPixbuf(pixbuf)
			}
		}
	}()

	app.Run(nil)
	return nil

}