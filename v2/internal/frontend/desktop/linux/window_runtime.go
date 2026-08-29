//go:build linux
// +build linux

package linux

/*
#cgo linux pkg-config: gtk+-3.0
#cgo !webkit2_41 pkg-config: webkit2gtk-4.0
#cgo webkit2_41 pkg-config: webkit2gtk-4.1

#include <gtk/gtk.h>
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
	return result
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
