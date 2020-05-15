package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mcrcon "github.com/Kelwing/mc-rcon"
)

const playersJS = "players.js"

var (
	seconds   int
	debug     bool
	host      string
	port      int
	password  string
	file      string
	avatarURL string
)

func init() {
	flag.IntVar(&seconds, "seconds", 60, "poll interval in seconds")
	flag.BoolVar(&debug, "debug", false, "log rcon command responses")
	flag.StringVar(&host, "host", "localhost", "rcon server")
	flag.IntVar(&port, "port", 25575, "rcon port")
	flag.StringVar(&password, "password", "", "rcon password")
	flag.StringVar(&file, "file", "players.json", "output file")
	flag.StringVar(&avatarURL, "avatar", "https://minotar.net/avatar/%s/16", "avatar URL where '%s' is replaced by player name")

}

type PlayerInfo struct {
	Name      string    `json:"name"`
	X         float64   `json:"x"`
	Y         float64   `json:"y"`
	Z         float64   `json:"z"`
	Dimension int64     `json:"dimension"`
	LastSeen  time.Time `json:"last_seen"`
	AvatarURL string    `json:"avatar_url"`
}

func connect(conn *mcrcon.MCConn, addr string) {
	dialErr := conn.Open(addr, password)
	if dialErr != nil {
		log.Printf("Connect to %s failed %v", addr, dialErr)
		os.Exit(1)
	}

	authErr := conn.Authenticate()
	if authErr != nil {
		log.Printf("Authentication to %s failed %v", addr, authErr)
		os.Exit(1)
	}
}

func main() {

	flag.Parse()

	conn := new(mcrcon.MCConn)
	defer conn.Close()

	addr := fmt.Sprintf("%s:%v", host, port)
	connect(conn, addr)

	done := make(chan bool)
	disconnected := make(chan bool)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	c := rcon{
		seconds:      seconds,
		conn:         conn,
		debug:        debug,
		disconnected: disconnected,
	}

	go func() {
		for {
			select {
			case <-disconnected:
				log.Println("Potential connection issue, reconnecting")
				connect(conn, addr)
				break
			case sig := <-stop:
				log.Println("Received", sig, "signal")
				done <- true
			}
		}
	}()

	c.start(done)
}
