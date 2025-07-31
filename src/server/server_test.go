package main

import (
	"testing"
)

func Test_GetGuidQuery(t *testing.T) {
	cases := map[string]string{
		"cam":  "SELECT Guid FROM Cam WHERE Brand LIKE ? AND Model LIKE ? AND Size = ?",
		"carabiner":  "SELECT Guid FROM Carabiner WHERE Brand LIKE ? AND Model LIKE ?",
		"sling":  "SELECT Guid FROM Sling WHERE Brand LIKE ? AND Model LIKE ? AND LengthInCentimeters = ?",
		"stopper":  "SELECT Guid FROM Stopper WHERE Brand LIKE ? AND Model LIKE ? AND Size = ?",
	}

	for key, expected := range cases {
		actual, err := getGuidQuery(key)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if actual != expected {
			t.Fatalf("Unexpected result for %v. Expected = %v, Actual = %v", key, expected, actual)
		}
	}
}

func Test_GetGuidQuery_Error(t *testing.T) {
	_, err := getGuidQuery("blam")

	if err == nil {
		t.Fatalf("Expected an error")
	}
}

func Test_GetGuidQuery_Not_Equal(t *testing.T) {
	expected := "expected"

	actual, err := getGuidQuery("cam")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual == expected {
		t.Fatalf("Unexpected match")
	}
}

func Test_GetWeightQuery(t *testing.T) {
	cases := map[string]string{
		"cam":  "SELECT WeightInGrams FROM Cam WHERE Guid = ?",
		"carabiner":  "SELECT WeightInGrams FROM Carabiner WHERE Guid = ?",
		"sling":  "SELECT WeightInGrams FROM Sling WHERE Guid = ?",
		"stopper":  "SELECT WeightInGrams FROM Stopper WHERE Guid = ?",
	}

	for key, expected := range cases {
		actual, err := getWeightQuery(key)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if actual != expected {
			t.Fatalf("Unexpected result for %v. Expected = %v, Actual = %v", key, expected, actual)
		}
	}
}

func Test_GetWeightQuery_Error(t *testing.T) {
	_, err := getWeightQuery("rock")

	if err == nil {
		t.Fatalf("Expected an error")
	}
}

func Test_GetWeightQuery_Not_Equal(t *testing.T) {
	expected := "expected"

	actual, err := getWeightQuery("stopper")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if actual == expected {
		t.Fatalf("Unexpected match")
	}
}