package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "sandbox-grpc/api"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	file = flag.String("file", "/home/seed-data/rack01.json", "input josn file of rack")
)

type Rack struct {
	Type string
	Brand string
	Model string
	Sizes []string
	Quantity int32
	LengthInCentimeters float64
}

func main() {
	flag.Parse()

	log.Printf("Reading input file: %v", file)
	content, err := ioutil.ReadFile(*file)

	if err != nil {
		log.Fatalf("Error opening json file: %v", err)
	}

	var rack []Rack
	err = json.Unmarshal(content, &rack)

	if err != nil {
		log.Fatalf("Error during unmarshal / deserialization: %v", err)
	}

	connection, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed connection: %v", err)
	}

	defer connection.Close()

	client := pb.NewGearStorageClient(connection)
	gear := getGear(client, rack)
	getTotalWeight(client, gear)
}

func buildRequest(rack []Rack) ([]*pb.GearRequest, error) {
	var requests []*pb.GearRequest

	for _, r := range rack {
		if r.Type == "cam" || r.Type == "stopper" {
			for _, size := range r.Sizes {
				requests = append(requests, &pb.GearRequest{Type: r.Type,
					Brand: r.Brand,
					Model: r.Model,
					Size: size,
					Quantity: r.Quantity})
			}
		} else if r.Type == "carabiner" {
			requests = append(requests, &pb.GearRequest{Type: r.Type,
				Brand: r.Brand,
				Model: r.Model,
				Quantity: r.Quantity})
		} else if r.Type == "sling" {
			requests = append(requests, &pb.GearRequest{Type: r.Type,
				Brand: r.Brand,
				Model: r.Model,
				LengthInCentimeters: r.LengthInCentimeters})
		} else {
			return nil, errors.New(fmt.Sprintf("Unknown gear type from input file: %v", r.Type))
		}
	}

	return requests, nil
}

func getGear(client pb.GearStorageClient, rack []Rack) ([]*pb.Gear) {
	requests, err := buildRequest(rack)

	if err != nil {
		log.Fatalf("BuildRequest failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := client.GetGear(ctx);

	if err != nil {
		log.Fatalf("GetGear failed: %v", err)
	}
	
	var gear []*pb.Gear

	waitGetGear := make(chan struct{})
	go func() {
		for {
			response, err := stream.Recv()

			if err == io.EOF {
				close(waitGetGear)
				return
			}

			log.Printf("GetGear response: %v", response.Gear.Guid)

			if err != nil {
				log.Fatalf("GetGear go routine failed: %v", err)
			}

			gear = append(gear, response.Gear)
		}
	}()

	for _, request := range requests {
		j, _ := json.Marshal(request)

		log.Printf("GetGear: %v", string(j))

		// TODO: comment when you're using it for real; have this here
		// to test out the channel / go routine in client
		// time.Sleep(time.Duration(rand.IntN(3)) * time.Second)

		err := stream.Send(request)

		if err != nil {
			log.Fatalf("GetGear: stream.Send(...) failed: %v", err)
		}
	}

	stream.CloseSend()
	<-waitGetGear

	return gear
}

func getTotalWeight(client pb.GearStorageClient, gear []*pb.Gear) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := &pb.WeightRequest{Gear: gear}

	//https://stackoverflow.com/questions/67038598/creating-grpc-client-request-with-repeated-fields
	response, err := client.GetTotalWeight(ctx, request)

	if err != nil {
		log.Fatalf("GetTotalWeight failed: %v", err)
	}

	log.Printf("GetTotalWeight response: %f lbs", response.WeightInGrams / (28*16))
}
