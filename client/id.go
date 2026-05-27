package client

import "github.com/google/uuid"

type ID string

func NewID() (*ID, error) {
	newuuid, err := uuid.NewV7()
	parseID := ID(newuuid.String())
	return &parseID, err
}
