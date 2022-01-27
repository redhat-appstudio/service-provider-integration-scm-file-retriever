// Copyright (c) 2021 - 2022 Red Hat, Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package websocket

import (
	"fmt"
	"log"
)

type Pool struct {
	Register    chan *Client
	Unregister  chan *Client
	Clients     map[string]*Client
	SendMessage chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     make(map[string]*Client),
		SendMessage: make(chan Message),
	}
}

func (pool *Pool) IsClientKnown(clientId string) bool {
	_, ok := pool.Clients[clientId]
	return ok
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client.ID] = client
			log.Println("Size of Connection Pool: ", len(pool.Clients))
			log.Println("New User Joined..." + client.ID)
			//for clientId, webClient := range pool.Clients {
			//	fmt.Println(webClient)
			//	webClient.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..." + clientId})
			//}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client.ID)
			log.Println("Size of Connection Pool: ", len(pool.Clients))
			//for client, _ := range pool.Clients {
			//	client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			//}
			break
		case message := <-pool.SendMessage:
			fmt.Println("Sending message to " + message.ClientID + " body " + message.Body)
			client, ok := pool.Clients[message.ClientID]
			if !ok {
				log.Println("Requested message to unknown client " + message.ClientID)
				break
			}
			client.Conn.WriteJSON(Message{Type: message.Type, Body: message.Body})
		}
	}
}
