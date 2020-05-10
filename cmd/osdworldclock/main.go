package main

//go:generate stringvar -out assets.go -pkg main globeImgFile:globe.png

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"strings"
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

func timein(t time.Time, location string) string {
	l, err := time.LoadLocation(location)
	if err != nil {
		panic(err)
	}
	return t.In(l).Format("3:04pm")
}

func worldclock() error {
	logoImg, _, err := image.Decode(strings.NewReader(globeImgFile))
	if err != nil {
		return errors.WithStack(err)
	}
	logoBounds := logoImg.Bounds()

	const width = 2000
	const fontSize = 48
	sansFont, err := fontloader.LoadCache("Mono")
	if err != nil {
		return errors.WithStack(err)
	}
	face := truetype.NewFace(sansFont, &truetype.Options{Size: fontSize})
	lineHeight := face.Metrics().Height.Round()
	height := (lineHeight * 3) + face.Metrics().Descent.Round() + logoBounds.Dy() + 8
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Over)

	imgBounds := img.Bounds()
	logoR := logoBounds.Add(image.Pt((imgBounds.Dx()-logoBounds.Dx())/2, lineHeight/2))
	draw.DrawMask(img, logoR, image.White, image.Point{}, logoImg, image.Point{}, draw.Over)

	fontDrawer := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
	}

	now := time.Now()
	clocks := [][]string{
		{"Mumbai", timein(now, "Asia/Kolkata")},
		{"Antwerp", timein(now, "Europe/Brussels")},
		{"Dubai", timein(now, "Asia/Dubai")},
		{"London", timein(now, "Europe/London")},
		{"NYC", timein(now, "America/New_York")},
		{"SF", timein(now, "America/Los_Angeles")},
	}

	sectionWidth := width / len(clocks)
	for i, clock := range clocks {
		y := int(lineHeight*2) + logoImg.Bounds().Dx()
		for _, line := range clock {
			lineWidth := fontDrawer.MeasureString(line).Round()
			fontDrawer.Dot = fixed.P(((sectionWidth-lineWidth)/2)+(i*sectionWidth), y)
			fontDrawer.DrawString(line)
			y += lineHeight
		}
	}

	app, err := gtk.ApplicationNew("org.daaku.osdworldclock", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = app.Connect("activate", func() {
		win, err := imagewindow.NewImageWindow(app, img)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}
		win.SetOpacity(0.7)
		win.SetTitle("World Clock")
		win.ShowAll()
	})
	if err != nil {
		return errors.WithStack(err)
	}

	go func() {
		time.Sleep(5 * time.Second)
		_, _ = glib.IdleAdd(func() {
			app.Quit()
		})
	}()

	app.Run(os.Args)
	return nil
}

func main() {
	if err := worldclock(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
