// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
)

// Hub maintains the set of active websocks and broadcasts messages to the
// websocks.
type Hub struct {
	// Registered websocks.
	websocks map[int]*Websock

	// Inbound messages from the websocks.
	broadcast chan []byte

	// Register requests from the websocks.
	register chan *Websock

	// Unregister requests from websocks.
	unregister chan *Websock
}

// type Message struct {
// 	Set_Num int `json:"set_num"`
// 	Destination int `json:"destination"`
// 	Content string `json:"content"`
// }

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Websock),
		unregister: make(chan *Websock),
		websocks:   make(map[string]*Websock),
	}
}

func (h *Hub) run() {
	var messageJSON Message
	for {
		select {
		case websock := <-h.register:
			h.websocks[websock.phone] = websock
		case websock := <-h.unregister:
			if _, ok := h.websocks[websock.phone]; ok {
				delete(h.websocks, websock.phone)
				close(websock.send)
			}
		case message := <-h.broadcast:
			fmt.Println(message)
			messageStr := string(message)
			fmt.Println(messageStr)
			err := json.Unmarshal(message, &messageJSON)
			if err != nil{
				panic(err)
			}
			destination := messageJSON.Destination
			fmt.Println(destination)
			content := []byte(messageJSON.Content)

			// for websock := range h.websocks {
			// 	if websock.phone != destination {
			// 		continue
			// 	}
			// 	select {
			// 	case websock.send <- content:
			// 	default:
			// 		// if for whatever reason the hub can't send the message to the websock,
			// 		// then remove the websock from the hash map.
			// 		close(websock.send)
			// 		delete(h.websocks, websock)
			// 	}
			// }

			if websock, ok := h.websocks[destination]; ok {
				select {
					case websock.send <- content:
					default:
						close(websock.send)
						delete(h.websocks, websock.phone)
				}
			}
		}
	}
}
