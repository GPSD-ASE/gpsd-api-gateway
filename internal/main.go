package main

import (
	"context"
	"gpsd-api-gateway/internal/gateway/pkg/config"
	"gpsd-api-gateway/internal/gateway/pkg/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type App struct {
	server *http.Server
	cc     *config.Config
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize(envs []string) {
	log.Println("Initializing API gateway...")
	a.cc = config.LoadConfig(envs)
	a.server = SetupServer(a.cc)
	log.Printf("Vault running on host %s", a.cc.VaultAddr)
}

func (a *App) Run() error {
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("API Gateway listening on port %s", a.cc.APIGatewayPort)
		serverErrors <- a.server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return err

	case <-shutdown:
		log.Println("Starting shutdown...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := a.server.Shutdown(ctx)
		if err != nil {
			log.Printf("Graceful shutdown did not complete in time: %v", err)

			if err := a.server.Close(); err != nil {
				log.Printf("Error during forced shutdown: %v", err)
			}
		}

		log.Println("Shutdown complete.")
		return nil
	}
}

func (a *App) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return a.server.Shutdown(ctx)
}

func SetupServer(cc *config.Config) *http.Server {
	router := setupRoutes(cc)

	return &http.Server{
		Addr:         ":" + cc.APIGatewayPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func setupRoutes(cc *config.Config) *mux.Router {
	r := mux.NewRouter()
	routes.RegisterRoutes(cc, r)
	return r
}

func Main() error {
	app := NewApp()
	app.Initialize(os.Environ())
	return app.Run()
}

func main() {
	if err := Main(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
