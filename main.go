package main

import (
	"embed"
	"flag"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

func main() {
	port := flag.String("port", "2525", "SMTP server port")
	flag.Parse()

	app := NewApp(*port)

	err := wails.Run(&options.App{
		Title:     "MailTap",
		Width:     1200,
		Height:    800,
		MinWidth:  900,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			About: &mac.AboutInfo{
				Title:   "MailTap v0.1.0",
				Message: "Local email capture and inspection tool.\n\n\u00a9 2026 Mike Wall.",
				Icon:    appIcon,
			},
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
