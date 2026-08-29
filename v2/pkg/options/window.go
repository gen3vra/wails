package options

type LayerShellLayer string

const (
	LayerNone       LayerShellLayer = ""
	LayerBackground LayerShellLayer = "background"
	LayerBottom     LayerShellLayer = "bottom"
	LayerTop        LayerShellLayer = "top"
	LayerOverlay    LayerShellLayer = "overlay"
)

type LayerShellKeyboardMode string

const (
	KeyboardModeNone      LayerShellKeyboardMode = ""
	KeyboardModeExclusive LayerShellKeyboardMode = "exclusive"
	KeyboardModeOnDemand  LayerShellKeyboardMode = "on-demand"
)

// LayerShell makes the window a wlr-layer-shell surface: the compositor anchors it to output edges above or below normal windows, outside workspace logic
type LayerShell struct {
	Layer        LayerShellLayer
	AnchorTop    bool
	AnchorBottom bool
	AnchorLeft   bool
	AnchorRight  bool
	MarginTop    int
	MarginBottom int
	MarginLeft   int
	MarginRight  int
	KeyboardMode LayerShellKeyboardMode
}

type Window struct {
	Title       string
	Width       int
	Height      int
	URL         string
	Frameless   bool
	AlwaysOnTop bool
	SkipTaskbar bool
	Resizable   bool
	Translucent bool
	LayerShell  *LayerShell
}
