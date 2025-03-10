package router

import (
	"net/http"
)

func HelloHandler(mux *http.ServeMux, grpcClient GrpcClient) {

	mux.HandleFunc("GET /greet/{name}", grpcClient.GetHelloRequest)
}

func RegisterHandlers(httpServer HttpServer, grpcClient GrpcClient) {

	HelloHandler(httpServer.Mux, grpcClient)
}
