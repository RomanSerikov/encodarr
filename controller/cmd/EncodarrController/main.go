package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BrenekH/encodarr/controller"
	"github.com/BrenekH/encodarr/controller/globals"
	"github.com/BrenekH/encodarr/controller/job_health"
	"github.com/BrenekH/encodarr/controller/settings"
	"github.com/BrenekH/encodarr/controller/sqlite"
	"github.com/BrenekH/logange"
)

func main() {
	configDir := "."

	// Setup main logger
	mainLogger := logange.NewLogger("main")

	formatter := logange.StandardFormatter{FormatString: "${datetime}|${name}|${lineno}|${levelname}|${message}\n"}

	// Setup the root logger to print info
	rootStdoutHandler := logange.NewStdoutHandler()
	rootStdoutHandler.SetFormatter(formatter)
	rootStdoutHandler.SetLevel(logange.LevelInfo)

	logange.RootLogger.AddHandler(&rootStdoutHandler)

	// Root logging to a file
	rootFileHandler, err := logange.NewFileHandler(fmt.Sprintf("%v/controller.log", configDir))
	if err != nil {
		log.Printf("Error creating rootFileHandler: %v", err)
		os.Exit(10)
		return
	}
	rootFileHandler.SetFormatter(formatter)
	rootFileHandler.SetLevel(logange.LevelInfo)

	logange.RootLogger.AddHandler(&rootFileHandler)

	mainLogger.Info("Starting Encodarr Controller version %v\n", globals.Version)
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		mainLogger.Info("Received stop signal: %v", sig)
		cancel()
	}()

	sqliteDatabase, err := sqlite.NewSQLiteDatabase(configDir)
	if err != nil {
		mainLogger.Critical("%v", err)
	}

	settingsStore, err := settings.NewSettingsStore(configDir)
	if err != nil {
		mainLogger.Critical("NewSettingsStore Error: %v", err)
	}

	hcDBAdapter := sqlite.NewHealthCheckerAdapater(&sqliteDatabase)
	healthChecker := job_health.NewChecker(&hcDBAdapter, &settingsStore)

	// TODO: Replace mocks with actual implemented structs
	mockLibraryManager := controller.MockLibraryManager{}
	mockRunnerCommunicator := controller.MockRunnerCommunicator{}
	mockUserInterfacer := controller.MockUserInterfacer{}

	controller.Run(&ctx, &healthChecker, &mockLibraryManager, &mockRunnerCommunicator, &mockUserInterfacer, false)
}
