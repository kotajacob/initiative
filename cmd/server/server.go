package main

import (
	"log"
	"net"
	"strings"
)

func server(messages chan<- message) {
	log.Println("listening on :6666")
	listener, err := net.Listen("tcp", ":6666")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err) // Connection aborted.
			continue
		}
		go handleConn(conn, messages)
	}
}

type message struct {
	cmd     string
	options []string
}

func handleConn(c net.Conn, messages chan<- message) {
	log.Println("client connected")
	defer c.Close()
	for {
		buf := make([]byte, 1024)
		reqLen, err := c.Read(buf)
		if err != nil {
			log.Println(err)
			break
		}
		req := string(buf[:reqLen])

		parts := strings.Split(strings.TrimSpace(req), ",")
		if len(parts) == 0 {
			continue
		}
		msg := message{
			cmd:     parts[0],
			options: parts[1:],
		}
		log.Println(msg.cmd, msg.options)
		messages <- msg
	}
}
