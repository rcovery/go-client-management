package id_test

import (
	"testing"

	"github.com/rcovery/go-client-management/packages/id"
)

func TestNewID(t *testing.T) {
	t.Run("should create a new ID", func(t *testing.T) {
		ID, err := id.NewID()
		if err != nil {
			t.Fatalf("NewID() %v", err)
		}

		if ID == nil {
			t.Errorf("Expected an ID, received nothing")
		}
	})

	t.Run("should create an unique ID", func(t *testing.T) {
		ID1, err1 := id.NewID()
		if err1 != nil {
			t.Fatalf("NewID() %v", err1)
		}

		stringifiedID1 := string(*ID1)
		if stringifiedID1 == "" {
			t.Errorf("ID 1: Expected an ID, received nothing")
		}

		ID2, err2 := id.NewID()
		if err2 != nil {
			t.Fatalf("NewID() %v", err2)
		}

		stringifiedID2 := string(*ID2)
		if stringifiedID2 == "" {
			t.Errorf("ID 2: Expected an ID, received nothing")
		}

		if stringifiedID1 == stringifiedID2 {
			t.Errorf("The IDs are equal! %v / %v", stringifiedID1, stringifiedID2)
		}
	})
}
