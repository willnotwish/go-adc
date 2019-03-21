// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
  // Registered clients.
  clients map[*Client]bool

  // Inbound messages for broadcast to connected clients
  broadcast chan []byte

  // Register requests from the clients.
  register chan *Client

  // Unregister requests from clients.
  unregister chan *Client

  shutdown chan bool
}

func newHub() *Hub {
  return &Hub{
    broadcast:  make(chan []byte),
    register:   make(chan *Client),
    unregister: make(chan *Client),
    clients:    make(map[*Client]bool),
    shutdown:   make(chan bool),
  }
}

func (h *Hub) run() {
  for {
    select {
    case client := <-h.register:
      h.clients[client] = true
    case client := <-h.unregister:
      if _, ok := h.clients[client]; ok {
        delete(h.clients, client)
        close(client.send)
      }
    case json := <-h.broadcast:
      for client := range h.clients {
        select {
        case client.send <- json:
        default:
          close(client.send)
          delete(h.clients, client)
        }
      }
    case <-h.shutdown:
      for client := range h.clients {
        close(client.send)
        delete(h.clients, client)
      }
      return
    }
  }
}
