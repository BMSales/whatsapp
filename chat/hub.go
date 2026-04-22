// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type Message struct {
	destination int `json:"destination"`
	content string `json:"content"`
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	var messageJSON Message
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			fmt.Println(message)
			messageStr := string(message)
			fmt.Println(messageStr)
			err := json.Unmarshal(message, &messageJSON)
			if err != nil{
				panic(err)
			}
			fmt.Println(messageJSON)
			content := []byte(messageJSON.content)
			for client := range h.clients {
				select {
				case client.send <- content:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
