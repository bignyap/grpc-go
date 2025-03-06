package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/bignyap/helloworld/service"
)

var (
	addr = flag.String("addr", "localhost:50051", "The server address")
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal json response: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 500 error:", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func CallHelloService(name string) (string, error) {
	conn, err := grpc.NewClient(*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Coluld not connect to the grpc server %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(), time.Second,
	)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	return r.GetMessage(), nil
}

func main() {
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /greet/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		message, err := CallHelloService(name)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		respondWithJSON(w, http.StatusOK, message)

	})
	newMux := CorsMiddleware(mux)
	http.ListenAndServe(":8080", newMux)
}
