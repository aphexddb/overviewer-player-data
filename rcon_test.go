package main

import (
	"testing"
)

// func init() {
// 	flag.Parse()
// }

func TestPlayersFromList(t *testing.T) {
	resp := "There are 2 of a max 20 players online: aphexddb, jasonbutler42"

	players := playersFromList(resp)
	count := len(players)

	if len(players) != 2 {
		t.Errorf("players length is %v; want 2", count)
	}

	if players[0] != "aphexddb" {
		t.Errorf("player 1 is %s; want aphexddb", players[0])
	}
	if players[1] != "jasonbutler42" {
		t.Errorf("player 2 is %s; want jasonbutler42", players[1])
	}

	emptyResp := "There are 0 of a max 20 players online:"
	emptyPlayers := playersFromList(emptyResp)

	if len(emptyPlayers) != 0 {
		t.Errorf("players should be empty")
	}
}

func TestPlayerPosFromData(t *testing.T) {
	resp := "aphexddb has the following entity data: [-142.86905620639067d, 72.0d, 145.2172516520357d]"

	pos := playerPosFromData(resp)

	if pos.X != -142.86905620639067 {
		t.Errorf("x is %v; want -142.86905620639067", pos.X)
	}
	if pos.Y != 72.0 {
		t.Errorf("y is %v; want 72.0", pos.Y)
	}
	if pos.Z != 145.2172516520357 {
		t.Errorf("z is %v; want 145.2172516520357", pos.Z)
	}

}

func TestPlayerDimensionFromData(t *testing.T) {
	resp := "aphexddb has the following entity data: 0"

	dimension := playerDimensionFromData(resp)

	if dimension == -1 {
		t.Errorf("dimension is -1; want non-negative value")
	}

	if dimension != 0 {
		t.Errorf("dimension is %v; want 0", dimension)
	}

}
