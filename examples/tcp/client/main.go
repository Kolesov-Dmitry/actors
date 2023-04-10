package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", ":5000")
	if err != nil {
		log.Fatalf("unable to connecto to the server: %s\n", err)
	}

	response := make([]byte, 1024)
	for i := 1; i <= 3; i++ {
		msg := fmt.Sprintf("message - %d", i)

		log.Println("send message:", msg)
		if _, err := conn.Write([]byte(msg)); err != nil {
			log.Printf("write error: %s\n", err)
			continue
		}

		n, err := conn.Read(response)
		if err != nil {
			log.Printf("read error: %s\n", err)
			continue
		}

		log.Printf("read message: %s\n", string(response[:n]))

		time.Sleep(time.Second)
	}
}
