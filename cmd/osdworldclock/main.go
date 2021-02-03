package main

import (
	"fmt"
	"os"
	"time"

	"github.com/daaku/osdtools/internal/imagewindow"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
)

const styles = `
window {
	background: rgba(5%, 5%, 5%, 0.75);
}
.icon {
	font-size: 128px;
	padding-top: 16px;
}
.city {
	padding-right: 24px;
}
.city label {
	font-weight: bold;
	color: rgba(255, 255, 255, 1);
	padding-left: 16px;
	padding-right: 16px;
}
.city .name {
	font-size: 32px;
}
.city .time {
	font-size: 24px;
	padding-bottom: 16px;
}
`

func timein(t time.Time, location string) string {
	l, err := time.LoadLocation(location)
	if err != nil {
		panic(err)
	}
	return t.In(l).Format("3:04pm")
}

func cityTimeBox(city, time string) (*gtk.Box, error) {
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := imagewindow.AddClass(vbox, "city"); err != nil {
		return nil, errors.WithStack(err)
	}

	cityLabel, err := gtk.LabelNew(city)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := imagewindow.AddClass(cityLabel, "name"); err != nil {
		return nil, errors.WithStack(err)
	}
	vbox.Add(cityLabel)

	timeLabel, err := gtk.LabelNew(time)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := imagewindow.AddClass(timeLabel, "time"); err != nil {
		return nil, errors.WithStack(err)
	}
	vbox.Add(timeLabel)
	return vbox, nil
}

func worldclock() error {
	now := time.Now()
	clocks := [][2]string{
		{"Mumbai", timein(now, "Asia/Kolkata")},
		{"Antwerp", timein(now, "Europe/Brussels")},
		{"Dubai", timein(now, "Asia/Dubai")},
		{"London", timein(now, "Europe/London")},
		{"NYC", timein(now, "America/New_York")},
		{"SF", timein(now, "America/Los_Angeles")},
	}

	gtk.Init(nil)

	if err := imagewindow.ConfigureStyles(styles); err != nil {
		return err
	}

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return errors.WithStack(err)
	}
	win.SetTitle("World Clock")
	imagewindow.LayerInit(win)
	imagewindow.LayerSetOverlay(win)

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return errors.WithStack(err)
	}

	emojiLabel, err := gtk.LabelNew("\U0001f557")
	if err != nil {
		return errors.WithStack(err)
	}
	if err := imagewindow.AddClass(emojiLabel, "icon"); err != nil {
		return errors.WithStack(err)
	}
	vbox.Add(emojiLabel)

	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, clock := range clocks {
		ctb, err := cityTimeBox(clock[0], clock[1])
		if err != nil {
			return err
		}
		hbox.Add(ctb)
	}

	vbox.Add(hbox)
	win.Add(vbox)
	win.ShowAll()

	go func() {
		time.Sleep(5 * time.Second)
		_, _ = glib.IdleAdd(gtk.MainQuit)
	}()

	gtk.Main()
	return nil
}

func main() {
	if err := worldclock(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
