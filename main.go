package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"path/filepath"
	"time"

	pb ""           // "github.com/grpc-up-and-running/samples/ch08/grpc-gateway/go/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	address  = "localhost:50051"
	hostname = "localhost"
	crtFile  = filepath.Join("06", "mutual-tls-channel", "certs", "client.crt")
	keyFile  = filepath.Join("06", "mutual-tls-channel", "certs", "client.key")
	caFile   = filepath.Join("06", "mutual-tls-channel", "certs", "ca.crt")
)

func main() {
	// Load the client certificate from disk
	certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("could not load client key pair: %s", err)
	}

	// Create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(caFile)
	if err != nil {
		log.Fatalf("could not read ca certificate: %s", err)
	}

	// Append the certificate from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("failed to append ca certs !")
	}

	opts := []grpc.DialOption{
		// transport credentials.
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			ServerName:   hostname, // Note: this is required!
			Certificates: []tls.Certificate{certificate},
			RootCAs:      certPool,
		})),
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	// Contact the server and print out its response.
	name := "Samsung S24"
	description := "Samsung Galaxy S24 is the latest smart phone, launched in October 2024"
	price := float32(1350.00)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})
	if err != nil {
		log.Fatalf("Could not add product: %v", err)
	}
	log.Printf("Product ID: %s added successfully", r.Value)

	product, err := c.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	if err != nil {
		log.Fatalf("Could not get product: %v", err)
	}
	log.Printf("Product: %s", product.String())

}
