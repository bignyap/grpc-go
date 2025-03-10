package router

import (
	"context"
	"log"
	"net/http"
	"time"

	pb "github.com/bignyap/helloworld/service"
	"github.com/bignyap/helloworld/utils"
)

func CallHelloService(
	ctx context.Context,
	c pb.GreeterClient,
	name string) (string, error) {

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	return r.GetMessage(), nil
}

func (grpcClient GrpcClient) GetHelloRequest(w http.ResponseWriter, r *http.Request) {

	c := pb.NewGreeterClient(grpcClient.GrpcConn)
	ctx, cancel := context.WithTimeout(
		context.Background(), time.Second,
	)
	defer cancel()

	name := r.PathValue("name")
	message, err := CallHelloService(ctx, c, name)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	utils.RespondWithJSON(w, http.StatusOK, message)

}
