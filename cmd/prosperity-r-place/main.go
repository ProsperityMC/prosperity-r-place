package main

import (
	"encoding/json"
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
	err = os.MkdirAll(".data/images", os.ModePerm)
	if err != nil {
		log.Fatal("Failed to make '.data/images' directory:", err)
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
	router.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "text/html")
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("Hello World!\n"))
	})
	router.Handle("/doc/{name}", cors(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		name := vars["name"]
		if manager, ok := managers[name]; ok {
			if websocket.IsWebSocketUpgrade(req) {
				upgrade, err := wsUpgrader.Upgrade(rw, req, rw.Header())
				if err != nil {
					http.Error(rw, "Failed to upgrade to websocket connection", http.StatusServiceUnavailable)
				}
				go prosperity_r_place.HandleWebsocket(upgrade, manager)
				return
			}
			if req.URL.Query().Get("raw") == "image" {
				rw.Header().Set("Content-Type", "image/png")
				rw.WriteHeader(http.StatusOK)
				_, _ = rw.Write(manager.Image())
				return
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(rw).Encode(map[string]int{
				"width":  manager.Width(),
				"height": manager.Height(),
			})
			return
		}
		http.Error(rw, "404 Not Found", http.StatusNotFound)
	})))
	server := &http.Server{
		Handler: router,
		Addr:    conf.Listen,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Println("[Main] Listen and serve error:", err)
		}
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
