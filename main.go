package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const configFile = "/etc/gotify/cli.json"

type cliConfig struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

type Message struct {
	ID       int    `json:"id"`
	AppID    int    `json:"appid"`
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
	Date     string `json:"date"`
}

func loadConfig() (cliConfig, error) {
	path := configFile
	f, err := os.Open(path)
	if err != nil {
		return cliConfig{}, err
	}
	defer f.Close()
	var cfg cliConfig
	return cfg, json.NewDecoder(f).Decode(&cfg)
}

func wsURL(serverURL string) string {
	u := strings.TrimRight(serverURL, "/")
	u = strings.Replace(u, "https://", "wss://", 1)
	u = strings.Replace(u, "http://", "ws://", 1)
	return u + "/stream"
}

func notify(title, message string) {
	if err := exec.Command("notify-send", "--", title, message).Run(); err != nil {
		log.Printf("notify-send: %v", err)
	}
}

func run(cfg cliConfig) {
	url := wsURL(cfg.URL) + "?token=" + cfg.Token
	backoff := 5 * time.Second
	for {
		log.Printf("connecting to wss://push.wirecrop.net/stream")
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Printf("dial: %v - retry in %s", err, backoff)
			time.Sleep(backoff)
			continue
		}
		log.Println("connected")
		backoff = 5 * time.Second
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read: %v - reconnecting", err)
				conn.Close()
				break
			}
			var msg Message
			if err := json.Unmarshal(data, &msg); err != nil {
				log.Printf("parse: %v (raw: %s)", err, data)
				continue
			}
			notify(msg.Title, msg.Message)
		}
		time.Sleep(backoff)
	}
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to read %s: %v", configFile, err)
	}
	if cfg.Token == "" {
		log.Fatalf("token is empty in %s", configFile)
	}
	if cfg.URL == "" {
		log.Fatalf("url is empty in %s", configFile)
	}
	run(cfg)
}
