package imagewindow

/*
#cgo pkg-config: gtk-layer-shell-0
#include <gtk-layer-shell.h>
*/
import "C"

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

func LayerInit(window *gtk.Window) {
	p := unsafe.Pointer(window.GObject)
	C.gtk_layer_init_for_window((*C.GtkWindow)(p))
}

func LayerSetOverlay(window *gtk.Window) {
	p := unsafe.Pointer(window.GObject)
	C.gtk_layer_set_layer((*C.GtkWindow)(p), C.GTK_LAYER_SHELL_LAYER_OVERLAY)
}
