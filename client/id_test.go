package client_test

import (
	"testing"

	"github.com/rcovery/go-client-management/client"
)

func TestNewID(t *testing.T) {
	t.Run("should create a new ShortURL ID", func(t *testing.T) {
		ID, err := client.NewID()
		if err != nil {
			t.Fatalf("NewID() %v", err)
		}

		if ID == "" {
			t.Errorf("Expected a Client ID, received nothing")
		}
	})

	t.Run("should create an unique ID", func(t *testing.T) {
		ID1, err1 := client.NewID()
		if err1 != nil {
			t.Fatalf("NewID() %v", err1)
		}

		stringifiedID1 := ID1
		if stringifiedID1 == "" {
			t.Errorf("ID 1: Expected a Client ID, received nothing")
		}

		ID2, err2 := client.NewID()
		if err2 != nil {
			t.Fatalf("NewID() %v", err2)
		}

		stringifiedID2 := ID2
		if stringifiedID2 == "" {
			t.Errorf("ID 2: Expected a Client ID, received nothing")
		}

		if stringifiedID1 == stringifiedID2 {
			t.Errorf("The IDs are equal! %s / %s", stringifiedID1, stringifiedID2)
		}
	})
}
