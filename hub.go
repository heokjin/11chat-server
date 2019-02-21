// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
)

type MsgProto struct {
	Sender   string `json:"S,omitempty"` //보내는사람
	Receiver string `json:"R,omitempty"` //받는사람
	Text     string `json:"T,omitempty"` //내용
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	totalChatMap map[string]*Client

	// Registered clients.
	// clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:    make(chan []byte),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		totalChatMap: make(map[string]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			log.Debugf("Register: %s\n", client.id)
			h.totalChatMap[client.id] = client
		case client := <-h.unregister:
			log.Debugf("UnRegister: %s\n", client.id)
			if _, ok := h.totalChatMap[client.id]; ok {
				delete(h.totalChatMap, client.id)
				close(client.send)
			}
		case message := <-h.broadcast:
			msgProto := MsgProto{}
			err := json.Unmarshal(message, &msgProto)
			if err != nil {
				log.Errorf("brodcast err: %s\n", err)
				return
			}

			log.Debug("Broadcast to %s\n", msgProto.Receiver)
			if _, ok := h.totalChatMap[msgProto.Receiver]; ok {
				client := h.totalChatMap[msgProto.Receiver]
				client.send <- []byte(msgProto.Text)
			}
		}
	}
}
