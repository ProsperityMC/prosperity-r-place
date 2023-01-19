package main

import (
	bytes "bytes"
	"context"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mrmelon54/mjwt"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	prosperityRPlace "prosperity-r-place"
	"prosperity-r-place/utils"
	"sort"
	"sync"
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

	roleMap := make(map[string]struct{})
	for _, i := range conf.Login.Guild.Roles {
		roleMap[i] = struct{}{}
	}

	var privKey *rsa.PrivateKey
	file, err := os.ReadFile(conf.Auth.Key)
	if os.IsNotExist(err) {
		privKey = generateNewKeys(conf.Auth)
		file = pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privKey),
		})
	} else {
		check("[Main] Failed to read token signing key", err)
		block, _ := pem.Decode(file)
		privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		check("[Main] Failed to parse token signing key", err)
	}
	signer := mjwt.NewMJwtSigner(conf.Auth.Issuer, privKey)

	managers := make(map[string]*prosperityRPlace.Manager)
	for i, slot := range conf.Slots {
		manager, err := prosperityRPlace.NewManager(slot.Name, slot.Width, slot.Height)
		if err != nil {
			log.Fatalf("error with slot %d: %v", i, err)
		}
		managers[slot.Name] = manager
	}

	wsUpgrader := &websocket.Upgrader{
		CheckOrigin: func(req *http.Request) bool { return true },
	}
	oauthConf := &oauth2.Config{
		RedirectURL:  conf.Login.RedirectUrl,
		ClientID:     conf.Login.Id,
		ClientSecret: conf.Login.Token,
		Scopes:       []string{discord.ScopeIdentify, "guilds.members.read"},
		Endpoint:     discord.Endpoint,
	}

	allowedStates := make(map[string]struct{})
	statesLock := &sync.RWMutex{}

	router := mux.NewRouter()
	router.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		http.Error(rw, "Prosperity r/place API endpoint!", http.StatusOK)
	})
	router.HandleFunc("/docs", func(rw http.ResponseWriter, req *http.Request) {
		a := make([]*prosperityRPlace.Manager, len(managers))
		i := 0
		for _, manager := range managers {
			a[i] = manager
			i++
		}
		sort.Slice(a, func(i, j int) bool { return a[i].Name < a[j].Name })
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(rw).Encode(a)
	})
	router.HandleFunc("/doc/{name}", func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		name := vars["name"]
		if manager, ok := managers[name]; ok {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(rw).Encode(manager)
			return
		}
		http.Error(rw, "404 Not Found", http.StatusNotFound)
	})
	router.HandleFunc("/doc/{name}/image", func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		name := vars["name"]
		if manager, ok := managers[name]; ok {
			img, hash := manager.Image()
			if req.Header.Get("If-None-Match") == hash {
				rw.WriteHeader(http.StatusNotModified)
				return
			}
			rw.Header().Set("Content-Type", "image/png")
			rw.Header().Set("ETag", hash)
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write(img)
			return
		}
		http.Error(rw, "404 Not Found", http.StatusNotFound)
	})
	router.HandleFunc("/doc/{name}/live", func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		name := vars["name"]
		if manager, ok := managers[name]; ok {
			if websocket.IsWebSocketUpgrade(req) {
				auth := req.URL.Query().Get("auth")
				upgrade, err := wsUpgrader.Upgrade(rw, req, rw.Header())
				if err != nil {
					http.Error(rw, "Failed to upgrade to websocket connection", http.StatusServiceUnavailable)
				}
				if _, c, err := mjwt.ExtractClaims[utils.DiscordInfo](signer, auth); err == nil {
					go prosperityRPlace.HandleWebsocket(upgrade, manager, c.Claims)
				} else {
					_ = upgrade.WriteMessage(websocket.TextMessage, []byte("no-auth"))
					select {
					case <-time.After(time.Second * 5):
						_ = upgrade.Close()
						break
					}
				}
				return
			}
		}
		http.Error(rw, "404 Not Found", http.StatusNotFound)
	})
	router.HandleFunc("/login", func(rw http.ResponseWriter, req *http.Request) {
		u := uuid.NewString()
		statesLock.Lock()
		allowedStates[u] = struct{}{}
		statesLock.Unlock()
		http.Redirect(rw, req, oauthConf.AuthCodeURL(u), http.StatusTemporaryRedirect)
	})
	router.HandleFunc("/callback", func(rw http.ResponseWriter, req *http.Request) {
		z := req.FormValue("state")
		statesLock.RLock()
		if _, ok := allowedStates[z]; !ok {
			statesLock.RUnlock()
			rw.WriteHeader(http.StatusBadRequest)
			_, _ = rw.Write([]byte("State does not match."))
			return
		}
		statesLock.RUnlock()
		statesLock.Lock()
		delete(allowedStates, z)
		statesLock.Unlock()

		token, err := oauthConf.Exchange(context.Background(), req.FormValue("code"))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte(err.Error()))
			return
		}

		res, err := oauthConf.Client(context.Background(), token).Get("https://discord.com/api/users/@me/guilds/" + conf.Login.Guild.Id + "/member")
		if err != nil {
			http.Error(rw, "Error collecting data from the Discord API", http.StatusInternalServerError)
			return
		}

		// check request status code
		switch res.StatusCode {
		case 200:
			break
		case 404:
			http.Error(rw, "User must be in the Discord guild", http.StatusConflict)
			return
		default:
			http.Error(rw, "Received unexpected response from Discord", http.StatusInternalServerError)
			return
		}

		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)

		var dm utils.DiscordMember
		j := json.NewDecoder(res.Body)
		err = j.Decode(&dm)
		if err != nil {
			http.Error(rw, "Failed to decode Discord API response", http.StatusInternalServerError)
			return
		}

		for _, i := range dm.Roles {
			if _, ok := roleMap[i]; ok {
				goto hasRole
			}
		}

		http.Error(rw, "User is missing a required role in the Discord guild", http.StatusConflict)
		return

	hasRole:
		// no need for the client to get the roles
		dm.Roles = nil

		dcToken, _ := encryptDiscordTokens(&privKey.PublicKey, token)

		u := uuid.NewString()
		h, err := signer.GenerateJwt(u, u, time.Hour*24, utils.DiscordInfo{UserId: dm.User.Id, Name: dm.User.Username, Discord: dcToken})
		if err != nil {
			http.Error(rw, "Failed to generate JWT token", http.StatusInternalServerError)
		}

		_, _ = fmt.Fprintf(rw, "<!DOCTYPE html><html><head><script>window.onload=function(){window.opener.postMessage(")
		encoder := json.NewEncoder(rw)
		_ = encoder.Encode(map[string]any{
			"token":  map[string]string{"access": h},
			"member": dm,
		})
		_, _ = fmt.Fprintf(rw, ",\"%s\");window.close();}</script></head></html>", conf.Login.BaseUrl)
	})
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

func generateNewKeys(auth AuthConfig) *rsa.PrivateKey {
	// Generate key
	fmt.Println("[generateNewKeys()] Generating new RSA private key")
	key, err := rsa.GenerateKey(rand.New(rand.NewSource(time.Now().UnixNano())), 4096)
	check("[generateNewKeys()] Failed to generate new RSA private key", err)

	// Create key files
	createPriv, err := os.Create(auth.Key)
	check("[generateNewKeys()] Failed to open private key file for writing", err)
	createPub, err := os.Create(auth.Public)
	check("[generateNewKeys()] Failed to open public key file for writing", err)

	// Encode and write keys
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	pubBytes := x509.MarshalPKCS1PublicKey(&key.PublicKey)
	err = pem.Encode(createPriv, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes})
	check("[generateNewKeys()] Failed to encode private key to file", err)
	err = pem.Encode(createPub, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubBytes})
	check("[generateNewKeys()] Failed to encode public key to file", err)
	return key
}

func encryptDiscordTokens(key *rsa.PublicKey, token *oauth2.Token) (string, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(map[string]string{
		"access":  token.AccessToken,
		"refresh": token.RefreshToken,
	})
	if err != nil {
		return "", err
	}

	hash := sha512.New()
	data, err := rsa.EncryptOAEP(hash, cryptoRand.Reader, key, b.Bytes(), nil)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(data), nil
}

func decryptDiscordTokens(key *rsa.PrivateKey, data string) (access, refresh string, err error) {
	raw, err := base64.RawStdEncoding.DecodeString(data)
	if err != nil {
		return "", "", err
	}

	hash := sha512.New()
	plain, err := rsa.DecryptOAEP(hash, cryptoRand.Reader, key, raw, nil)
	if err != nil {
		return "", "", err
	}

	var a struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
	}
	b := bytes.NewBuffer(plain)
	err = json.NewDecoder(b).Decode(&a)
	return a.Access, a.Refresh, err
}

func check(msg string, err error) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
