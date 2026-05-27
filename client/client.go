package client

type Client struct {
	ID             *ID
	IdempotencyKey *IdempotencyKey
	Name           string
	Email          string
	RequestType    string
	Status         Status
	PortfolioValue int
}

type PostClientBody struct {
	Name           string `json:"cliente_nome"`
	Email          string `json:"cliente_email"`
	RequestType    string `json:"tipo_solicitacao"`
	PortfolioValue int    `json:"valor_patrimonio"`
}
