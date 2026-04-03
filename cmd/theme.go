package cmd

import (
	"fmt"
	"strings"

	"github.com/btbytes/bleat/config"
	"github.com/btbytes/bleat/display"
	"github.com/btbytes/bleat/theme"
)

func RunThemeList() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Warning: could not load config: %v\n", err)
	}
	currentTheme := cfg.Theme
	if currentTheme == "" {
		currentTheme = theme.DefaultThemeName()
	}

	names := display.ListThemesSorted()
	display.PrintThemeList(names, currentTheme)
}

func RunThemeSet(name string) {
	if !theme.HasTheme(name) {
		display.PrintThemeNotFound(name)
		return
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	cfg.Theme = name
	if err := cfg.Save(); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		return
	}

	display.PrintThemeSetSuccess(name)
}

func RunThemeShow() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	currentTheme := cfg.Theme
	if currentTheme == "" {
		currentTheme = theme.DefaultThemeName()
	}

	display.PrintThemeShow(currentTheme)
}

func RunTheme(args []string) {
	if len(args) == 0 {
		RunThemeShow()
		return
	}

	sub := strings.ToLower(args[0])
	switch sub {
	case "list", "ls":
		RunThemeList()
	case "set":
		if len(args) < 2 {
			fmt.Println("Usage: bleat theme set <name>")
			return
		}
		RunThemeSet(args[1])
	case "show", "current":
		RunThemeShow()
	default:
		fmt.Printf("Unknown theme command: %s\n", sub)
		fmt.Println("Available commands: list, set, show")
	}
}
