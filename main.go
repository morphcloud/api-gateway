package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/morphcloud/api-gateway/pb"
)

func run(l *log.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	apiGatewaySockAddr := os.Getenv("API_GATEWAY_SOCK_ADDR")
	advertServiceSockAddr := flag.String("advert-service-sock-addr", os.Getenv("ADVERT_SERVICE_SOCK_ADDR"), "Advert Service TCP Socket Address")

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := pb.RegisterAdvertServiceHandlerFromEndpoint(ctx, mux, *advertServiceSockAddr, opts); err != nil {
		return err
	}

	l.Println("API Gateway is running on port " + apiGatewaySockAddr)
	return http.ListenAndServe(":"+apiGatewaySockAddr, mux)
}

func main() {
	l := log.New(os.Stdout, "api-gateway ", log.LstdFlags)

	flag.Parse()

	if err := godotenv.Load(); err != nil {
		l.Fatalln(err)
	}

	if err := run(l); err != nil {
		l.Fatalln(err)
	}
}
