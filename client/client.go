// Package client - Defines what a client is, and what data it stores
package client

type Client struct {
	ID             *ID
	Name           string
	Email          string
	Status         Status
	Priority       *Priority
	PortfolioValue int
}

type PostClientBody struct {
	Name           string `json:"cliente_nome"`
	Email          string `json:"cliente_email"`
	RequestType    string `json:"tipo_solicitacao"`
	PortfolioValue int    `json:"valor_patrimonio"`
}
