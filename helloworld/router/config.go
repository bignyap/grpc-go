package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bignyap/helloworld/utils"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	GrpcConn *grpc.ClientConn
}

func NewGrpcClientConn(addr string) GrpcClient {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Coluld not connect to the grpc server %v", err)
	}
	return GrpcClient{
		GrpcConn: conn,
	}
}

type HttpServer interface {
	ListenAndServe()
}

type HttpServerObject struct {
	HttpMux *http.ServeMux
	Mux     *runtime.ServeMux
	Addr    string
}

func NewHttpServer(addr string) HttpServerObject {
	return HttpServerObject{
		Mux:  runtime.NewServeMux(),
		Addr: addr,
	}
}

func (httpServer HttpServerObject) ListenAndServe() {

	fmt.Println("API gateway server is running on " + httpServer.Addr)
	if err := http.ListenAndServe(httpServer.Addr, httpServer.Mux); err != nil {
		log.Fatal("gateway server closed abruptly: ", err)
	}
}

func (httpServer HttpServerObject) AddMiddlewares() {

	httpServer.Mux = utils.CorsMiddleware(httpServer.HttpMux)
}
