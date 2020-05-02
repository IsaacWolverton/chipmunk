package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	localPort  int = 42069
	remotePort int
)

func init() {
	// get the application image name from the environment
	var err error
	remotePort, err = strconv.Atoi(os.Getenv("APPLICATION_PORT"))
	if err != nil {
		log.Println("failed to parse port number", err)
		remotePort = 8080
	}
}

func main() {
	// TODO: UDP as well
	server, err := net.Listen("tcp", fmt.Sprintf(":%d", localPort))
	if err != nil {
		panic(err)
	}
	log.Println("chipmunk proxy running")

	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		log.Printf("handling request from '%v'\n", conn.RemoteAddr())

		destination := fmt.Sprintf("%s:%d", strings.Split(conn.LocalAddr().String(), ":")[0], remotePort)
		appConn, err := net.Dial("tcp", destination)
		if err != nil {
			panic(err)
		}
		defer appConn.Close()

		log.Println("proxied connect, sending data")
		go func() {
			io.Copy(appConn, conn)
			// log.Println("incoming: ", data)

		}()
		go func() {
			io.Copy(conn, appConn)
			// log.Println("outcoming: ", data)
		}()
	}
}
