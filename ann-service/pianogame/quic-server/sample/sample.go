package sample

import (
	"fmt"
	"net/http"

	"github.com/lucas-clemente/quic-go/h2quic"
)

func available(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "yes")
}

// Run run quic server
func Run(addr, cert, key string) {
	fmt.Println("Start quic sample server")
	http.HandleFunc("/quic", available)
	h2quic.ListenAndServeQUIC(addr, cert, key, nil)
	//h2quic.ListenAndServe(addr, cert, key, nil)
	//h2quic.ListenAndServeTLS(cert, key)
}
