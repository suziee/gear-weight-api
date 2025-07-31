package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	//https://stackoverflow.com/questions/21220077/what-does-an-underscore-in-front-of-an-import-statement-mean
	_ "github.com/mattn/go-sqlite3"

	"google.golang.org/grpc"

	pb "sandbox-grpc/api"
)

const connection_string = "/home/db/gear.db"

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterGearStorageServer(server, &GearStorageServer{})
	log.Printf("Server is listening at %v", listener.Addr())

	err = server.Serve(listener)

	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

type GearStorageServer struct {
	pb.UnimplementedGearStorageServer
}

func (x *GearStorageServer) GetGear (stream pb.GearStorage_GetGearServer) (error) {
	for {
		request, err := stream.Recv()
		
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		guid, err := getGuid(request)

		if err != nil {
			return err
		}

		// TODO: comment when you're using it for real; have this here
		// to test out the channel / go routine in client
		// time.Sleep(time.Duration(rand.IntN(3)) * time.Second)

		stream.Send(&pb.GearResponse{
			Gear: &pb.Gear{
				Type: request.Type,
				Guid: guid,
				Quantity: request.Quantity}})
	}
}

func (s *GearStorageServer) GetTotalWeight(c context.Context, request *pb.WeightRequest) (*pb.WeightResponse, error) {
	var totalWeight float64

	for _, gear := range request.Gear {
		weight, err := getWeight(gear)

		if err != nil {
			return nil, err
		}

		totalWeight += (weight * float64(gear.Quantity))
	}

	return &pb.WeightResponse{WeightInGrams: totalWeight}, nil
}

func getGuid(request *pb.GearRequest) (string, error) {
	stmt, err := getGuidQuery(request.Type)

	var guid string

	db, err := sql.Open("sqlite3", connection_string)

	if err != nil {
		return guid, err
	}

	defer db.Close()

	statement, err := db.Prepare(stmt)

	if err != nil {
		return guid, err
	}

	if request.Type == "cam" || request.Type == "stopper" {
		err = statement.QueryRow(request.Brand, request.Model, request.Size).Scan(&guid)
	} else if request.Type == "carabiner" {
		err = statement.QueryRow(request.Brand, request.Model).Scan(&guid)
	} else if request.Type == "sling" {
		err = statement.QueryRow(request.Brand, request.Model, request.LengthInCentimeters).Scan(&guid)
	}

	if err != nil {
		j, _ := json.Marshal(request)
		log.Printf("GetGuid: %v, %v", string(j), err)
	}

	return guid, err
}

func getWeight(gear *pb.Gear) (float64, error) {
	stmt, err := getWeightQuery(gear.Type)

	var weight float64

	db, err := sql.Open("sqlite3", connection_string)

	if err != nil {
		return weight, err
	}

	defer db.Close()

	statement, err := db.Prepare(stmt)

	if err != nil {
		return weight, err
	}

	err = statement.QueryRow(gear.Guid).Scan(&weight)

	return weight, err
}

func getGuidQuery(gearType string) (string, error) {
	var stmt string

	//https://go.dev/doc/database/prepared-statements
	//https://go.dev/wiki/Switch implicit break
	switch gearType {
	case "cam": stmt = "SELECT Guid FROM Cam WHERE Brand LIKE ? AND Model LIKE ? AND Size = ?"
	case "carabiner": stmt = "SELECT Guid FROM Carabiner WHERE Brand LIKE ? AND Model LIKE ?"
	case "sling": stmt = "SELECT Guid FROM Sling WHERE Brand LIKE ? AND Model LIKE ? AND LengthInCentimeters = ?"
	case "stopper": stmt = "SELECT Guid FROM Stopper WHERE Brand LIKE ? AND Model LIKE ? AND Size = ?"
	default: return stmt, errors.New(fmt.Sprintf("Invalid gear type: %v", gearType))
	}

	return stmt, nil
}

func getWeightQuery(gearType string) (string, error) {
	var stmt string

	switch gearType {
	case "cam": stmt = "SELECT WeightInGrams FROM Cam WHERE Guid = ?"
	case "carabiner": stmt = "SELECT WeightInGrams FROM Carabiner WHERE Guid = ?"
	case "sling": stmt = "SELECT WeightInGrams FROM Sling WHERE Guid = ?"
	case "stopper": stmt = "SELECT WeightInGrams FROM Stopper WHERE Guid = ?"
	default: return stmt, errors.New(fmt.Sprintf("Invalid gear type: %v", gearType))
	}

	return stmt, nil
}