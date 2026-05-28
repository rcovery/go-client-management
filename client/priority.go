package client

import (
	"errors"
	"log"
)

type Priority string

const (
	HighPriority   Priority = "prioridade_alta"
	NormalPriority Priority = "prioridade_normal"
)

func ToPriority(rawPriority *string) (*Priority, error) {
	if rawPriority == nil {
		return nil, nil
	}
	p := Priority(*rawPriority)
	if p != HighPriority && p != NormalPriority {
		log.Print("invalid priority: ", *rawPriority)
		return nil, errors.New("invalid.priority.value")
	}

	return &p, nil
}
