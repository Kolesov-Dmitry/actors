package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/klsvdm/actors/actor"
)

type AddSessionEvent struct {
	Conn net.Conn
}

// Define server as an actor
type Server struct {
	listener net.Listener
}

func NewServer(listenAddr string) (*Server, error) {
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}

	return &Server{
		listener: lis,
	}, nil
}

func (s *Server) Receive(e *actor.Environ, p *actor.Parcel) {
	switch event := p.Message.(type) {
	case actor.StartedEvent:
		log.Printf("server started [%s]\n", s.listener.Addr())
		go s.acceptLoop(e)

	case actor.AboutToStopEvent:
		if err := s.listener.Close(); err != nil {
			log.Printf("failed to stop server: %s\n", err)
			return
		}
		log.Println("server stopped")

	case AddSessionEvent:
		session := &Session{
			conn: event.Conn,
		}

		// spawn session handler as a child actor
		_, err := e.SpawnChild(session, "session", uuid.NewString())
		if err != nil {
			_ = event.Conn.Close()

			log.Printf("failed to spawn session actor: %s\n", err)
			return
		}
	}
}

// acceptLoop accepts incomming connections
// and packs them into AddSessionEvent
func (s *Server) acceptLoop(e *actor.Environ) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}

			log.Printf("accept error: %s\n", err)
			continue
		}

		// send AddSessionEvent so it can be handled thread-safe
		e.Send(e.Self(), AddSessionEvent{Conn: conn})
	}
}

type ReadMessageEvent struct {
	Message []byte
}

// Session is an actor too
type Session struct {
	conn net.Conn
}

func (s *Session) Receive(e *actor.Environ, p *actor.Parcel) {
	switch event := p.Message.(type) {
	case actor.StartedEvent:
		log.Printf("client connected [%s]\n", s.conn.RemoteAddr())
		go s.handleLoop(e)

	case ReadMessageEvent:
		log.Println(string(event.Message))

		if _, err := s.conn.Write(event.Message); err != nil {
			log.Printf("write error: %s\n", err)
		}
	}
}

// handleLoop reads messages from the client and packs them into ReadMessageEvent
func (s *Session) handleLoop(e *actor.Environ) {
	buf := make([]byte, 1024)
	for {
		n, err := s.conn.Read(buf)
		if err != nil {
			if err == io.EOF || errors.Is(err, syscall.ECONNRESET) {
				log.Printf("client [%s] closed connection", s.conn.RemoteAddr())
				return
			}

			log.Printf("read error: %+v\n", err)
			continue
		}

		message := make([]byte, n)
		copy(message, buf)

		e.Send(e.Self(), ReadMessageEvent{Message: message})
	}
}

func main() {
	// 1. create server
	server, err := NewServer(":5000")
	if err != nil {
		log.Fatalln(err)
	}

	// 2. spawn server actor
	engine := actor.NewEngine()
	if _, err = engine.Spawn(server, "server"); err != nil {
		log.Fatalf("unable to spawn server actor: %s\n", err)
	}

	// press Ctrl+C to exit
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	<-exit

	// 3. gracefull shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := engine.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}
}
