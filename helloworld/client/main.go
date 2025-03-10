package main

import (
	"flag"
	"net/http"

	"github.com/bignyap/helloworld/router"
)

var (
	addr = flag.String("addr", "localhost:50051", "The server address")
)

func main() {

	flag.Parse()

	// Establish new GrpcClient Object
	gRPCConn := router.NewGrpcClientConn(*addr)
	defer gRPCConn.GrpcConn.Close()

	// Establish a new Http Handler
	httpServer := router.NewHttpServer(":8080")
	httpMux := http.NewServeMux()
	httpMux.Handle("/", httpServer.Mux)
	router.RegisterHandlers(httpServer, gRPCConn)
	// newMux := utils.CorsMiddleware(mux)
	httpServer.ListenAndServe()
}
