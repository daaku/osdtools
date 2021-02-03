package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/daaku/osdtools/internal/bar"
)

func icon(volume int) string {
	if volume == 0 {
		return "\U0001f507"
	}
	return "\U0001f50a"
}

func main() {
	volume := flag.Int("volume", 0, "volume value between 0-100")
	flag.Parse()

	if err := bar.Run(*volume, icon, ".osdvol"); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
