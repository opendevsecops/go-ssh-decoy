package main

import (
	"crypto/rand"
	"crypto/rsa"
	"flag"
	"fmt"
	ssh "golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
)

func main() {
	host := flag.String("host", "0.0.0.0", "Host address to listen to")
	port := flag.Int("port", 2022, "Port number to listen to")
	size := flag.Int("size", 2048, "Private key size")
	priv := flag.String("priv", "", "Private key file")

	flag.Parse()

	config := &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			log.Println(conn.RemoteAddr().String(), conn.User(), string(password))

			return nil, fmt.Errorf("Password rejected for %q", conn.User())
		},
	}

	if *priv == "" {
		key, err := rsa.GenerateKey(rand.Reader, *size)

		if err != nil {
			log.Fatal("Failed to generate ssh key")

			return
		}

		signer, err := ssh.NewSignerFromKey(key)

		if err != nil {
			log.Fatal("Failed to import ssh key")

			return
		}

		config.AddHostKey(signer)
	} else {
		dat, err := ioutil.ReadFile(*priv)

		if err != nil {
			log.Fatal("Failed to load ssh key")

			return
		}

		signer, err := ssh.ParsePrivateKey(dat)

		if err != nil {
			log.Fatal("Failed to parse ssh key")

			return
		}

		config.AddHostKey(signer)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to listen for ssh connections on port %s:%d", *host, *port))

		return
	}

	for {
		nconn, err := listener.Accept()

		if err != nil {
			log.Println("Failed to accept incoming ssh connection:", err)

			continue
		}

		conn, _, _, err := ssh.NewServerConn(nconn, config)

		if err != nil {
			log.Println("Failed to handshake:", err)

			continue
		}

		conn.Close()
	}
}
