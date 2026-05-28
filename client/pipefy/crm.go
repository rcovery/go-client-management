package pipefy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

func (r *CRMGateway) CreateCard(ctx context.Context, clientData *client.Client, requestType string) (bool, error) {
	var (
		pipeID = 123
		title  = fmt.Sprintf("%s (%s)", clientData.Name, clientData.Email)
	)

	query := `
		mutation CreateCard($pipeId: ID!, $title: String!, $fields: [FieldValueInput]!) {
			createCard(input: {
				pipe_id: $pipeId,
				title: $title,
				fields_attributes: $fields
			}) {
				card { title }
			}
		}
	`
	variables := map[string]any{
		"pipeId": pipeID,
		"title":  title,
		"fields": []map[string]any{
			{"field_id": "cliente_nome", "field_value": clientData.Name},
			{"field_id": "cliente_email", "field_value": clientData.Email},
			{"field_id": "valor_patrimonio", "field_value": clientData.PortfolioValue},
			{"field_id": "tipo_solicitacao", "field_value": requestType},
		},
	}

	graphqlReqErr := r.runRequest(ctx, query, variables)
	if graphqlReqErr != nil {
		return false, graphqlReqErr
	}

	return true, nil
}

func (r *CRMGateway) UpdateCard(ctx context.Context, cardID string, status client.Status, priority client.Priority) (bool, error) {
	query := `
		mutation UpdateCard($id: ID!, $fields: [FieldValueInput]!) {
			updateCard(input: {
				id: $id,
				fields_attributes: $fields
			}) {
				card { id title }
			}
		}
	`
	variables := map[string]any{
		"id": cardID,
		"fields": []map[string]any{
			{"field_id": "status", "field_value": string(status)},
			{"field_id": "prioridade", "field_value": string(priority)},
		},
	}

	graphqlReqErr := r.runRequest(ctx, query, variables)
	if graphqlReqErr != nil {
		return false, graphqlReqErr
	}

	return true, nil
}

func (r *CRMGateway) runRequest(ctx context.Context, query string, variables map[string]any) error {
	body, _ := json.Marshal(map[string]any{
		"query":     query,
		"variables": variables,
	})

	req, newReqErr := http.NewRequestWithContext(ctx, "POST", "https://api.pipefy.com/graphql", bytes.NewBuffer(body))
	if newReqErr != nil {
		log.Println(newReqErr.Error())
		return fmt.Errorf("could.not.create.crm.request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.pipefyToken)

	_, reqErr := http.DefaultClient.Do(req)
	if reqErr != nil {
		log.Println(reqErr.Error())
		return fmt.Errorf("could.not.do.crm.request")
	}

	return nil
}
