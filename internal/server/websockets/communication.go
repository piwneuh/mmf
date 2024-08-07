package ws

import (
	"encoding/json"
	"log"
	"mmf/internal/model"

	"github.com/gorilla/websocket"
)

func SendMatchFoundToPlayers(matchId string, matchTickets []model.Ticket) bool {
	mess := GenerateMatchFoundResponse(matchTickets, matchId)
	marshalled, err := json.Marshal(mess)
	if err != nil {
		return false
	}

	for _, ticket := range matchTickets {
		SendMessageToUser(ticket.Member.SteamID, marshalled)
	}
	return true
}

func SendMessageToUser(steamId string, message []byte) {
	userConnectionsMutex.Lock()
	defer userConnectionsMutex.Unlock()

	conn, ok := userConnections[steamId]
	if !ok {
		log.Println("User not connected")
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Println(err)
	}
}

func DisconnectUser(steamId string) {
	userConnectionsMutex.Lock()
	defer userConnectionsMutex.Unlock()

	conn, ok := userConnections[steamId]
	if !ok {
		log.Println("User not connected")
		return
	}

	// Close the connection
	if err := conn.Close(); err != nil {
		log.Println("Error closing connection:", err)
	}

	// Remove the connection from the map
	delete(userConnections, steamId)
}
