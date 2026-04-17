//go:build linux
// +build linux

package linux

// WindowStartDrag initiates an interactive move/drag of the main window.
// This bridges into the existing StartDrag implementation on the Window type.
func (f *Frontend) WindowStartDrag() {
    f.startDrag()
}
