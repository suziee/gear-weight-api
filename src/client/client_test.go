package main

import (
	"encoding/json"
	"errors"
	"testing"

	pb "sandbox-grpc/api"
)

func Test_BuildRequest(t *testing.T) {
	input := []Rack{{
		Type: "cam",
		Brand: "metolius",
		Model: "tcu",
		Sizes: []string{"1", "2", "3"},
		Quantity: 1,
	}}

	expected := []*pb.GearRequest{
		&pb.GearRequest{
			Type: "cam",
			Brand: "metolius",
			Model: "tcu",
			Size: "1",
			Quantity: 1,
		},
		&pb.GearRequest{
			Type: "cam",
			Brand: "metolius",
			Model: "tcu",
			Size: "2",
			Quantity: 1,
		},
		&pb.GearRequest{
			Type: "cam",
			Brand: "metolius",
			Model: "tcu",
			Size: "3",
			Quantity: 1,
		},
	}

	actual, err := buildRequest(input)

	if err != nil {
		t.Fatalf("Unexpected error = %v", err)
	}

	if err := equal_requests(expected, actual); err != nil {
		t.Fatalf("Unexpected equals error = %v.\nExpected = %v,\nActual = %v", err, indent(expected), indent(actual))
	}
}

func Test_BuildRequest_Unknown_Type(t *testing.T) {
	input := []Rack{{
		Type: "blam",
		Brand: "metolius",
		Model: "tcu",
		Sizes: []string{"1", "2", "3"},
		Quantity: 1,
	}}

	_, err := buildRequest(input)

	if err == nil {
		t.Fatal("Expected a buildRequest error")
	}
}

func Test_BuildRequest_Not_Equal(t *testing.T) {
	input := []Rack{{
		Type: "cam",
		Brand: "metolius",
		Model: "tcu",
		Sizes: []string{"1"},
		Quantity: 1,
	}}

	expected := []*pb.GearRequest{
		&pb.GearRequest{
			Type: "cam",
			Brand: "metolius",
			Model: "tcu",
			Size: "3",
			Quantity: 1,
		},
	}

	actual, err := buildRequest(input)

	if err != nil {
		t.Fatalf("Unexpected error = %v", err)
	}

	if err := equal_requests(expected, actual); err == nil {
		t.Fatal("Expected an equals error")
	}
}

func equal_requests(a []*pb.GearRequest, b []*pb.GearRequest) error {
	if len(a) != len(b) {
		return errors.New("length not equal")
	}

	for i := 0; i < len(a); i++ {
		return equal_request(a[i], b[i])
	}

	return nil
}

func equal_request(a *pb.GearRequest, b *pb.GearRequest) error {
	if a.Type != b.Type ||
		a.Brand != b.Brand ||
		a.Model != b.Model ||
		a.Size != b.Size ||
		a.Quantity != b.Quantity {
		return errors.New("objects not equal")
	}

	return nil
}

func indent(obj interface{}) string {
	j, _ := json.MarshalIndent(obj, "", "\t")
	return string(j)
}