package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	prosperity_r_place "prosperity-r-place"
	"syscall"
	"time"
)

func main() {
	err := os.MkdirAll(".data", os.ModePerm)
	if err != nil {
		log.Fatal("Failed to make '.data' directory:", err)
	}

	openConf, err := os.Open(".data/config.yml")
	if err != nil {
		log.Fatal("Failed to open '.data/config.yml':", err)
	}
	var conf Config
	err = yaml.NewDecoder(openConf).Decode(&conf)
	if err != nil {
		log.Fatal("Failed to decode config:", err)
	}

	managers := make(map[string]*prosperity_r_place.Manager)
	for i, slot := range conf.Slots {
		manager, err := prosperity_r_place.NewManager(slot.Name, slot.Width, slot.Height)
		if err != nil {
			log.Fatalf("error with slot %d: %v", i, err)
		}
		managers[slot.Name] = manager
	}

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
			go prosperity_r_place.HandleWebsocket(upgrade)
			return
		}
	})))
	server := &http.Server{
		Handler: router,
		Addr:    conf.Listen,
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
