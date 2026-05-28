package webhook

import "github.com/rcovery/go-client-management/packages/id"

type ID id.ID

func NewID() (*ID, error) {
	i, err := id.NewID()
	return (*ID)(i), err
}
