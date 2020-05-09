// This code was taken from kahlys' proxy server which is available here (https://github.com/kahlys/proxy/blob/master/proxy.go) under the MIT license

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// Server is a TCP server that takes an incoming request and sends it to another
// server, proxying the response back to the client.
type Server struct {
	// TCP address to listen on
	Addr string

	// TCP address of target server
	Target string

	// The path to the directory in which the save file should be written
	// Note: must end in '/', e.g. "/path/to/bucket/"
	PathPrefix string

	// The name of the file that the proxy server is currently writing network traffic to
	SaveFile string

	// The path to the file that should be loaded and replayed upon the initial call to ListenAndServe
	// Should be "" (the empty string) if the user does not want to replay any network traffic
	// example path: "/path/to/the/file"
	ReplayPath string

	mu sync.Mutex
}

// ListenAndServe listens on the TCP network address laddr and then handle packets
// on incoming connections.
func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	if s.ReplayPath != "" {
		// the user wants to replay network traffic
		s.mu.Lock()
		time.Sleep(5 * time.Second)
		log.Println("Replaying traffic from:", s.ReplayPath)
		replay_b, err := ioutil.ReadFile(s.ReplayPath)
		s.ReplayPath = "" // prevent future connections from replaying the same traffic
		if err != nil {
			log.Println(err)
		} else {
			log.Println(" >-> replaying")
			str := string(replay_b)
			requests := strings.Split(str, "\r\n\r\n")
			for i, req := range requests {
				replay_conn, err := net.Dial("tcp", s.Target)
				if err != nil {
					log.Println("failed to form replay connection")
					continue
				}
				new_b := []byte(req + "\r\n\r\n")
				log.Println(fmt.Sprintf("sending replay message: %d", i))
				log.Println("Contents:" + req)
				_, err = replay_conn.Write(new_b)
				if err != nil {
					log.Println("replay err:")
					log.Println(err)
				}
				replay_conn.Close()
			}
		}

		log.Println("Replay Complete")
		s.mu.Unlock()
	}
	return s.serve(listener)
}

func (s *Server) serve(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go s.handleConn(conn)
	}
}

// These functions are how the checkpointer stops the proxy from relaying traffic when it
// is actively checkpointing and resumes when it is finished. Note that StopProxy automatically
// creates a new SaveFile and returns its name so that the checkpointer can set ReplayPath
// accordingly in the event of a crash
// ------------------------------------------------------------------------------------------
func (s *Server) StopProxy(version int) string {
	s.mu.Lock()

	//update SaveFile
	s.SaveFile = fmt.Sprintf("network-%d", version)
	return s.SaveFile
}

func (s *Server) ResumeProxy() {
	s.mu.Unlock() // allows proxy to continue functioning as normal
}

// ------------------------------------------------------------------------------------------

func (s *Server) handleConn(conn net.Conn) {
	// connects to target server
	rconn, err := net.Dial("tcp", s.Target)
	if err != nil {
		return
	}

	// write to dst what it reads from src
	var pipe = func(src, dst net.Conn) {
		defer func() {
			conn.Close()
			rconn.Close()
		}()

		buff := make([]byte, 65535)

		for {
			n, err := src.Read(buff)
			if err != nil {
				log.Println(err)
				return
			}
			b := buff[:n]

			_, err = dst.Write(b)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

	// write to dst what it reads from src
	// while saving all network traffic to s.SaveFile
	var save_pipe = func(src, dst net.Conn, s *Server) {
		defer func() {
			conn.Close()
			rconn.Close()
		}()

		buff := make([]byte, 65535)

		for {
			n, err := src.Read(buff)
			if err != nil {
				return
			}

			s.mu.Lock()
			b := buff[:n]

			// append this traffic to the current SaveFile
			filePath := s.PathPrefix + s.SaveFile
			log.Println("saving network traffic to:", filePath)
			f, err := os.OpenFile(filePath,
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println("Error opening file:")
				log.Println(err)
			}

			if _, err := f.Write(b); err != nil {
				log.Println("Write error to file:")
				log.Println(err)
				f.Close()
				s.mu.Unlock()
				return
			}
			f.Sync()

			f.Close()

			// send the data along the connection now that it has been saved
			_, err = dst.Write(b)
			if err != nil {
				log.Println("Write error to dest:")
				log.Println(err)
				s.mu.Unlock()
				return
			}

			s.mu.Unlock()
		}
	}

	go save_pipe(conn, rconn, s)
	go pipe(rconn, conn)
}
