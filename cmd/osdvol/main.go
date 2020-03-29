package main

//go:generate stringvar -out assets.go -pkg main muteImgFile:mute.png volumeImgFile:volume.png

import (
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgbutil"
	"github.com/pkg/errors"

	"github.com/daaku/osdtools/internal/bar"
)

func render(volume int) (*image.RGBA, error) {
	iconFile := volumeImgFile
	if volume == 0 {
		iconFile = muteImgFile
	}

	iconImg, _, err := image.Decode(strings.NewReader(iconFile))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return bar.Render(iconImg, volume)
}

func main() {
	xgb.Logger.SetOutput(ioutil.Discard)
	xgbutil.Logger.SetOutput(ioutil.Discard)

	volume := flag.Int("volume", 0, "volume value between 0-100")
	flag.Parse()

	if err := bar.Run(*volume, render, ".osdvol"); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
