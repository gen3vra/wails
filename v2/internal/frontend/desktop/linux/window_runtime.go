//go:build linux
// +build linux

package linux

/*
#cgo linux pkg-config: gtk+-3.0 gtk-layer-shell-0
#cgo !webkit2_41 pkg-config: webkit2gtk-4.0
#cgo webkit2_41 pkg-config: webkit2gtk-4.1

#include <gtk/gtk.h>
#include <gtk-layer-shell/gtk-layer-shell.h>
#include <webkit2/webkit2.h>
#include "window.h"
*/
import "C"
import (
	"unsafe"

	"github.com/gen3vra/wails/v2/pkg/options"
)

// Must run on the GTK main thread
func NewRuntimeWindow(opts options.Window, gpuPolicy int) *Window {
	result := &Window{}

	gtkWindow := C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)
	C.g_object_ref_sink(C.gpointer(gtkWindow))
	result.gtkWindow = unsafe.Pointer(gtkWindow)

	result.contentManager = unsafe.Pointer(C.webkit_user_content_manager_new())
	webview := C.SetupRuntimeWebview(result.contentManager, result.asGTKWindow(), C.int(gpuPolicy))
	result.webview = unsafe.Pointer(webview)
	C.gtk_container_add(result.asGTKContainer(), webview)

	width, height := opts.Width, opts.Height
	if width <= 0 {
		width = 400
	}
	if height <= 0 {
		height = 200
	}
	result.SetKeepAbove(opts.AlwaysOnTop)
	result.SetResizable(opts.Resizable)
	result.SetDefaultSize(width, height)
	// GTK ignores the default size on non resizable windows and uses the webview's natural size instead
	if !opts.Resizable {
		C.gtk_widget_set_size_request(webview, C.gint(width), C.gint(height))
	}
	result.SetDecorated(!opts.Frameless)
	result.SetTitle(opts.Title)
	C.gtk_window_set_skip_taskbar_hint(result.asGTKWindow(), gtkBool(opts.SkipTaskbar))
	if opts.Translucent {
		C.SetWindowTransparency(gtkWindow)
	}
	if opts.LayerShell != nil && opts.LayerShell.Layer != options.LayerNone {
		applyLayerShell(result.asGTKWindow(), opts.Title, opts.LayerShell)
	}
	return result
}

// Must run on the GTK main thread, before the window first maps
func applyLayerShell(win *C.GtkWindow, namespace string, ls *options.LayerShell) {
	C.gtk_layer_init_for_window(win)

	layer := C.GTK_LAYER_SHELL_LAYER_TOP
	switch ls.Layer {
	case options.LayerBackground:
		layer = C.GTK_LAYER_SHELL_LAYER_BACKGROUND
	case options.LayerBottom:
		layer = C.GTK_LAYER_SHELL_LAYER_BOTTOM
	case options.LayerOverlay:
		layer = C.GTK_LAYER_SHELL_LAYER_OVERLAY
	}
	C.gtk_layer_set_layer(win, C.GtkLayerShellLayer(layer))

	if namespace != "" {
		ns := C.CString(namespace)
		C.gtk_layer_set_namespace(win, ns)
		C.free(unsafe.Pointer(ns))
	}

	C.gtk_layer_set_anchor(win, C.GTK_LAYER_SHELL_EDGE_TOP, gtkBool(ls.AnchorTop))
	C.gtk_layer_set_anchor(win, C.GTK_LAYER_SHELL_EDGE_BOTTOM, gtkBool(ls.AnchorBottom))
	C.gtk_layer_set_anchor(win, C.GTK_LAYER_SHELL_EDGE_LEFT, gtkBool(ls.AnchorLeft))
	C.gtk_layer_set_anchor(win, C.GTK_LAYER_SHELL_EDGE_RIGHT, gtkBool(ls.AnchorRight))

	C.gtk_layer_set_margin(win, C.GTK_LAYER_SHELL_EDGE_TOP, C.int(ls.MarginTop))
	C.gtk_layer_set_margin(win, C.GTK_LAYER_SHELL_EDGE_BOTTOM, C.int(ls.MarginBottom))
	C.gtk_layer_set_margin(win, C.GTK_LAYER_SHELL_EDGE_LEFT, C.int(ls.MarginLeft))
	C.gtk_layer_set_margin(win, C.GTK_LAYER_SHELL_EDGE_RIGHT, C.int(ls.MarginRight))

	mode := C.GTK_LAYER_SHELL_KEYBOARD_MODE_NONE
	switch ls.KeyboardMode {
	case options.KeyboardModeExclusive:
		mode = C.GTK_LAYER_SHELL_KEYBOARD_MODE_EXCLUSIVE
	case options.KeyboardModeOnDemand:
		mode = C.GTK_LAYER_SHELL_KEYBOARD_MODE_ON_DEMAND
	}
	C.gtk_layer_set_keyboard_mode(win, C.GtkLayerShellKeyboardMode(mode))
}

// Must run on the GTK main thread
func (w *Window) loadAndShow(url string) {
	_url := C.CString(url)
	C.LoadIndex(w.webview, _url)
	C.free(unsafe.Pointer(_url))
	C.gtk_widget_show_all(w.asGTKWidget())
}

// Must run on the GTK main thread
func (w *Window) setVisibleMainThread(visible bool) {
	if visible {
		C.gtk_widget_show_all(w.asGTKWidget())
	} else {
		C.gtk_widget_hide(w.asGTKWidget())
	}
}

// Must run on the GTK main thread
func (w *Window) toggleVisibleMainThread() {
	w.setVisibleMainThread(C.gtk_widget_is_visible(w.asGTKWidget()) == 0)
}
