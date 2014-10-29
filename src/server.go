package main

import (
    "log"
    "net"
)

type Server struct {
    port int
    ln   *net.TCPListener
    quit chan bool
}

func (s *Server) start() (err error) {
    s.ln, err = net.ListenTCP("tcp4", &net.TCPAddr{Port: s.port})
    if err != nil {
        log.Printf("listen err: %s", err)
        return
    }

    log.Printf("server is listenning on %s", s.ln.Addr())

    s.quit = make(chan bool)
    go s.listenRoutine()

    return
}

func (s *Server) stop() {
    s.ln.Close()
    <-s.quit
}

func (s *Server) listenRoutine() {
    for {
        conn, err := s.ln.AcceptTCP()
        if err != nil {
            log.Printf("listener(%s) AcceptTCP err: %s", s.ln.Addr(), err)
            break
        }
        log.Printf("accept connection from %s", conn.RemoteAddr())
        c := &Session{conn: conn}
        go c.recvRoutine()
    }
    s.quit <- true
}
