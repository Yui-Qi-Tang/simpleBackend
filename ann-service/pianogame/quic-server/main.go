package main

import "simpleBackend/ann-service/pianogame/quic-server/sample"

func main() {
	sample.Run("localhost:4242", "localhost.pem", "localhost.key")
}
