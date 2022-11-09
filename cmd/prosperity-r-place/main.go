package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	manager:= NewManager()

	cors := handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Upgrade"}),
		handlers.AllowedMethods([]string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodConnect,
			http.MethodOptions,
			http.MethodTrace,
		}),
	)
	wsUpgrader := &websocket.Upgrader{
		CheckOrigin: func(req *http.Request) bool { return true },
	}

	router := mux.NewRouter()
	router.Handle("/", cors(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if websocket.IsWebSocketUpgrade(req) {
			upgrade, err := wsUpgrader.Upgrade(rw, req, rw.Header())
			if err != nil {
				http.Error(rw, "Failed to upgrade to websocket connection", http.StatusServiceUnavailable)
			}
			go
			return
		}
	})))
	server := &http.Server{
		Handler: router,
		Addr:    os.Getenv("LISTEN"),
	}
	go func() {
		_ = server.ListenAndServe()
	}()

	// Wait for exit signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	fmt.Println()

	// Stop runner
	log.Printf("[Main] Stopping Prosperity r/place")
	n := time.Now()
	_ = server.Close()
	log.Printf("[Main] Took '%s' to shutdown\n", time.Now().Sub(n))
	log.Println("[Main] Goodbye")
}
