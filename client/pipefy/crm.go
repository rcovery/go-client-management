package pipefy

import (
	"context"

	"github.com/rcovery/go-client-management/client"
)

type CRMGateway struct {
	pipefyToken string
}

func NewCRMGateway(pipefyToken string) *CRMGateway {
	return &CRMGateway{
		pipefyToken: pipefyToken,
	}
}

func (r *CRMGateway) CreateCard(ctx context.Context, clientData *client.Client) (bool, error) {
	// _, execErr := r.DB.ExecContext(ctx, query,
	// 	string(*clientData.ID),
	// 	clientData.Name,
	// 	clientData.Email,
	// 	clientData.PortfolioValue,
	// 	string(clientData.Status),
	// )
	// if execErr != nil {
	// 	return false, execErr
	// }

	return true, nil
}
