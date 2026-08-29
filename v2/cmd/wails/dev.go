package main

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/gen3vra/wails/v2/cmd/wails/flags"
	"github.com/gen3vra/wails/v2/cmd/wails/internal/dev"
	"github.com/gen3vra/wails/v2/internal/colour"
	"github.com/gen3vra/wails/v2/pkg/clilogger"
)

func devApplication(f *flags.Dev) error {
	if f.NoColour {
		pterm.DisableColor()
		colour.ColourEnabled = false
	}

	quiet := f.Verbosity == flags.Quiet

	// Create logger
	logger := clilogger.New(os.Stdout)
	logger.Mute(quiet)

	if quiet {
		pterm.DisableOutput()
	} else {
		app.PrintBanner()
	}

	err := f.Process()
	if err != nil {
		return err
	}

	return dev.Application(f, logger)
}
