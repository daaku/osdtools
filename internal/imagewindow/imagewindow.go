// Package imagewindow provides utilities for managing OSD windows.
package imagewindow

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
)

func ConfigureDefaultScreenStyles(styles string) error {
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

type StyleContext interface {
	GetStyleContext() (*gtk.StyleContext, error)
}

func AddClass(w StyleContext, class string) error {
	sc, err := w.GetStyleContext()
	if err != nil {
		return errors.WithStack(err)
	}
	sc.AddClass(class)
	return nil
}
