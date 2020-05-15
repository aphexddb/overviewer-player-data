package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	mcrcon "github.com/Kelwing/mc-rcon"
)

type rcon struct {
	seconds      int
	conn         *mcrcon.MCConn
	debug        bool
	disconnected chan bool
}

type PlayerPos struct {
	X float64
	Y float64
	Z float64
}

func (c *rcon) start(done <-chan bool) {

	ticker := time.NewTicker(time.Duration(seconds) * time.Second)
	log.Println("Writing to", file, "every", seconds, "seconds")

	for {
		select {
		case <-done:
			return
		case <-ticker.C:

			playerState := make(map[string]PlayerInfo)

			list := c.cmd("list")
			players := playersFromList(list)

			if c.debug {
				log.Printf("%v players online: %s", len(players), strings.Join(players, ","))
			}

			for _, player := range players {
				posData := c.cmd(fmt.Sprintf("data get entity %s Pos", player))
				pos := playerPosFromData(posData)

				dimensionData := c.cmd(fmt.Sprintf("data get entity %s Dimension", player))
				dimension := playerDimensionFromData(dimensionData)

				playerState[player] = PlayerInfo{
					Name:      player,
					X:         pos.X,
					Y:         pos.Y,
					Z:         pos.Z,
					Dimension: dimension,
					LastSeen:  time.Now().UTC(),
					AvatarURL: fmt.Sprintf(avatarURL, player),
				}

			}

			bytes, _ := json.Marshal(playerState)

			if c.debug {
				log.Println(string(bytes))
			}

			writeToFile(file, bytes)
		}
	}

}

func (c *rcon) cmd(cmd string) string {
	resp, err := c.conn.SendCommand(cmd)
	if err != nil {
		if isDisconnected(err.Error()) {
			c.disconnected <- true
		}
		log.Printf("Command '%s' failed %s", cmd, err)
		return ""
	}
	return resp
}

func isDisconnected(msg string) bool {
	if strings.Contains(msg, "broken pipe") || strings.Contains(msg, "closed network connection") {
		return true
	}

	return false
}

func playersFromList(list string) []string {
	players := []string{}

	parts := strings.Split(list, ":")
	if len(parts) != 2 {
		log.Println("Unable to read player list:", list)
		return players
	}

	rawData := strings.TrimSpace(parts[1])
	if len(rawData) == 0 {
		return players
	}

	playerList := strings.Split(rawData, ",")
	for _, player := range playerList {
		players = append(players, strings.TrimSpace(player))
	}

	return players
}

func playerPosFromData(data string) PlayerPos {

	pos := PlayerPos{
		X: float64(0),
		Y: float64(0),
		Z: float64(0),
	}

	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		log.Println("Unable to read player position:", data)
		return pos
	}

	rawData := parts[1]
	rawData = strings.ReplaceAll(rawData, "[", "")
	rawData = strings.ReplaceAll(rawData, "]", "")
	rawData = strings.ReplaceAll(rawData, "d", "")

	positions := strings.Split(rawData, ",")
	if len(positions) != 3 {
		log.Println("Unable to read player positions:", positions)
		return pos
	}

	x, xErr := strconv.ParseFloat(strings.TrimSpace(positions[0]), 64)
	y, yErr := strconv.ParseFloat(strings.TrimSpace(positions[1]), 64)
	z, zErr := strconv.ParseFloat(strings.TrimSpace(positions[2]), 64)

	if xErr != nil {
		log.Println("Unable to read player position X value:", xErr)
		return pos
	}
	if yErr != nil {
		log.Println("Unable to read player position X value:", yErr)
		return pos
	}
	if zErr != nil {
		log.Println("Unable to read player position X value:", zErr)
		return pos
	}

	pos.X = x
	pos.Y = y
	pos.Z = z

	return pos
}

func playerDimensionFromData(data string) int64 {

	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		log.Println("Unable to read player dimension:", data)
		return -1
	}

	dimension, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)

	if err != nil {
		log.Println("Unable to read player dimension value:", err)
		return -1
	}

	return dimension
}

func writeToFile(filePath string, b []byte) {
	err := ioutil.WriteFile(filePath, b, 0644)
	if err != nil {
		log.Println("Error writing to file", err)
	}
}
