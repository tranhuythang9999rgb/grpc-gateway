package main

import (
	"context"
	"log"
	"net"
	"net/http"

	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/iamrajiv/helloworld-grpc-gateway/proto/api"
	"github.com/iamrajiv/helloworld-grpc-gateway/usecase"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Initialize the tracer (make sure it's properly set up in your actual implementation)
	tracer := opentracing.GlobalTracer()

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object with tracing middleware
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_opentracing.UnaryServerInterceptor(
				grpc_opentracing.WithTracer(tracer),
			),
		),
	)

	// Attach the Greeter service to the server
	api.RegisterGreeterServer(grpcServer, &usecase.Server{})

	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalln("Failed to serve gRPC:", err)
		}
	}()

	// Create a client connection to the gRPC server
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(
				grpc_opentracing.WithTracer(tracer),
			),
		),
	}
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		opts..., // Use variadic parameter syntax
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	// Create a new ServeMux for the gRPC-Gateway
	gwmux := runtime.NewServeMux()

	// Register the Greeter service with the gRPC-Gateway
	err = api.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	// Apply the logging middleware to the ServeMux
	loggedMux := api.LogRequestBody(gwmux)
	// Apply the tracing middleware to the ServeMux
	tracedMux := api.TracingWrapper(loggedMux)

	// Create a new HTTP server for the gRPC-Gateway
	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: tracedMux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	if err := gwServer.ListenAndServe(); err != nil {
		log.Fatalln("Failed to serve gRPC-Gateway:", err)
	}
}
