package request

import (
	"github.com/rcovery/go-client-management/packages/id"
)

type ID id.ID

func NewID() (*ID, error) {
	id, err := id.NewID()
	return (*ID)(id), err
}
