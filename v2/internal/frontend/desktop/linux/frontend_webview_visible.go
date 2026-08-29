//go:build linux
// +build linux

package linux

// Wayland never unmaps toplevels on hidden workspaces, so unmapping the widget is the only way WebKit learns it is off-screen and suspends its engine
func (f *Frontend) WindowSetWebviewVisible(visible bool) {
	f.mainWindow.SetWebviewVisible(visible)
}
