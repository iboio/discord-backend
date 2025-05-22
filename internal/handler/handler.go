package handler

import (
	"discord/config"
	"discord/context"
	"discord/internal/cron"
	"discord/internal/worker"
	echo2 "discord/pkg/server/echo"
	"fmt"
)

func Handler() {
	config.LoadConfig()
	appCtx, err := context.AppContextInit()
	if err != nil {
		fmt.Println("Error initializing app context:", err)
	}

	err = echo2.StartServer(appCtx)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	newCron := cron.NewCron(appCtx) // Initialize the cron job with the batcher
	newCron.StartCron()             // Start the cron job
	fmt.Println("Cron job started successfully")

	newWorker := worker.NewWorker(appCtx) // Initialize the worker with the app context
	newWorker.StartWorkers()              // Start the workers
	fmt.Println("Workers started successfully")
	select {} // Keeps the application running
}
