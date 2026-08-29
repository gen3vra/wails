//go:build linux
// +build linux

package linux

import (
	"strings"

	"github.com/gen3vra/wails/v2/pkg/options"
	"github.com/gen3vra/wails/v2/pkg/options/linux"
)

func (f *Frontend) registerWindow(w *Window) uint {
	f.windowsLock.Lock()
	defer f.windowsLock.Unlock()
	id := f.nextWindowID
	f.nextWindowID++
	f.windows[id] = w
	return id
}

func (f *Frontend) windowForID(id uint) *Window {
	f.windowsLock.Lock()
	defer f.windowsLock.Unlock()
	return f.windows[id]
}

func (f *Frontend) WindowCreate(opts options.Window) uint {
	gpuPolicy := int(linux.WebviewGpuPolicyNever)
	if f.frontendOptions.Linux != nil {
		gpuPolicy = int(f.frontendOptions.Linux.WebviewGpuPolicy)
	}
	url := f.startURL.String() + strings.TrimPrefix(opts.URL, "/")

	var w *Window
	done := make(chan struct{})
	invokeOnMainThread(func() {
		w = NewRuntimeWindow(opts, gpuPolicy)
		w.loadAndShow(url)
		close(done)
	})
	<-done
	return f.registerWindow(w)
}

func (f *Frontend) WindowSetVisibleByID(id uint, visible bool) {
	w := f.windowForID(id)
	if w == nil || w == f.mainWindow {
		return
	}
	invokeOnMainThread(func() { w.setVisibleMainThread(visible) })
}

func (f *Frontend) WindowToggleByID(id uint) {
	w := f.windowForID(id)
	if w == nil || w == f.mainWindow {
		return
	}
	invokeOnMainThread(func() { w.toggleVisibleMainThread() })
}

func (f *Frontend) WindowDestroyByID(id uint) {
	f.windowsLock.Lock()
	w := f.windows[id]
	if w == f.mainWindow {
		f.windowsLock.Unlock()
		return
	}
	delete(f.windows, id)
	f.windowsLock.Unlock()
	if w == nil {
		return
	}
	invokeOnMainThread(func() { w.Destroy() })
}
