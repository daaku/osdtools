package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/daaku/osdtools/internal/bar"
)

func icon(brightness int) string {
	return "\U0001f506"
}

func main() {
	brightness := flag.Int("brightness", 0, "brightness value between 0-100")
	flag.Parse()

	if err := bar.Run(*brightness, icon, ".osdbrightness"); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
