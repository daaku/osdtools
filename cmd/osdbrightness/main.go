package main

//go:generate stringvar -out assets.go -pkg main brightnessImgFile:brightness.png

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

func render(brightness int) (*image.RGBA, error) {
	iconImg, _, err := image.Decode(strings.NewReader(brightnessImgFile))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return bar.Render(iconImg, brightness)
}

func main() {
	xgb.Logger.SetOutput(ioutil.Discard)
	xgbutil.Logger.SetOutput(ioutil.Discard)

	brightness := flag.Int("brightness", 0, "brightness value between 0-100")
	flag.Parse()

	if err := bar.Run(*brightness, render, ".osdbrightness"); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
