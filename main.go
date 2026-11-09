package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	
	"project-studio/internal/services"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	authService := services.NewAuthService()
	activityService := services.NewActivityService()
	presenceService := services.NewPresenceService()
	projectService := services.NewProjectService(activityService)
	
	notifService := services.NewNotificationService()
	taskService := services.NewTaskService(notifService, activityService)
	timeService := services.NewTimeTrackingService(activityService)
	workloadService := services.NewWorkloadService()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "PROJECT-STUDIO",
		Width:     1280,
		Height:    800,
		MinWidth:  1024,
		MinHeight: 680,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			authService.Startup(ctx)
			activityService.Startup(ctx)
			presenceService.Startup(ctx)
			projectService.Startup(ctx)
			notifService.Startup(ctx)
			taskService.Startup(ctx)
			timeService.Startup(ctx)
			workloadService.Startup(ctx)
		},
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
			authService,
			activityService,
			presenceService,
			projectService,
			taskService,
			notifService,
			timeService,
			workloadService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
