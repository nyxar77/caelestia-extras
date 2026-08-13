package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/nyxar77/caelestia-extras/internal/config"
	"github.com/nyxar77/caelestia-extras/internal/cursor"
	"github.com/nyxar77/caelestia-extras/internal/integration"
)

const version = "0.1.0"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "caelestia-extras:", err)
		os.Exit(1)
	}
}

func run(arguments []string) error {
	flags := flag.NewFlagSet("caelestia-extras", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	configPath := flags.String("config", config.DefaultPath(), "configuration file")
	showVersion := flags.Bool("version", false, "show version")
	if err := flags.Parse(arguments); err != nil {
		return err
	}
	if *showVersion {
		fmt.Println(version)
		return nil
	}
	args := flags.Args()
	if len(args) == 0 {
		return errors.New("expected cursor, gtk, hyprtoolkit, portal, or pavucontrol")
	}
	configFile, err := config.Load(*configPath)
	if err != nil {
		return err
	}

	switch args[0] {
	case "cursor":
		if configFile.Cursor == nil {
			return errors.New("cursor is not configured")
		}
		if len(args) != 2 {
			return errors.New("expected cursor sync or cursor sync-xcursor")
		}
		if args[1] == "sync" {
			return cursor.Sync(*configFile.Cursor, configFile.Scheme.File)
		}
		if args[1] == "sync-xcursor" {
			return cursor.SyncXCursor(*configFile.Cursor, configFile.Scheme.File)
		}
		return errors.New("expected cursor sync or cursor sync-xcursor")
	case "gtk":
		if configFile.GTK == nil {
			return errors.New("gtk is not configured")
		}
		if len(args) != 2 || args[1] != "sync" {
			return errors.New("expected gtk sync")
		}
		return integration.SyncGTK(*configFile.GTK, configFile.Scheme.File)
	case "hyprtoolkit":
		if configFile.Hyprtoolkit == nil {
			return errors.New("hyprtoolkit is not configured")
		}
		if len(args) != 2 || args[1] != "sync" {
			return errors.New("expected hyprtoolkit sync")
		}
		return integration.SyncHyprtoolkit(*configFile.Hyprtoolkit)
	case "pavucontrol":
		if configFile.Pavucontrol == nil {
			return errors.New("pavucontrol is not configured")
		}
		return integration.LaunchPavucontrol(*configFile.Pavucontrol, args[1:])
	case "portal":
		if configFile.Portal == nil {
			return errors.New("portal is not configured")
		}
		if len(args) != 2 || args[1] != "sync" {
			return errors.New("expected portal sync")
		}
		return integration.SyncPortal(*configFile.Portal)
	default:
		return fmt.Errorf("unknown integration: %s", args[0])
	}
}
