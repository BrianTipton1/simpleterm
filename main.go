package main

import (
	"embed"
	"termvim/pkg/server"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	api := server.NewApi()

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "termvim",
		WindowStartState: options.WindowStartState(options.Maximised),
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			api,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
