// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	c, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:10328/login?name=zjw", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	for {
		fmt.Print(": ")
		var msg string
		fmt.Scanf("%s", &msg)
		if msg == "" {
			continue
		}
		err := c.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println("write:", err)
		}
		time.Sleep(800 * time.Millisecond)
	}
	log.Println("over!")
}
