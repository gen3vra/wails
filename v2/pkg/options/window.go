package options

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
}
