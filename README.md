## Instruções de execução local

Para rodar o projeto localmente:
```sh
docker compose up --build
```
- O script irá subir o contêiner, preparar as variáveis de ambiente, instalar dependências, rodar migrations e rodar a aplicação

---

Para rodar os testes:
```sh
go test -v ./...
```
- Este comando ira executar os testes de integraçao (com testcontainers) e testes unitarios (com comportamento mockado)

## Exemplos de requisição

Criação de cliente
```sh
curl --location 'localhost:9000/clientes' \
--header 'Content-Type: application/json' \
--data-raw '{
    "cliente_nome": "João Silva",
    "cliente_email": "joao.silva@example.com",
    "tipo_solicitacao": "Atualização cadastral",
    "valor_patrimonio": 250000
}'
```

Retorno em caso de sucesso:
```json
{
    "message": "created.successfully",
    "success": true
}
```

---
Simulação de webhook do Pipefy
```sh
curl --location 'localhost:9000/webhooks/pipefy/card-updated' \
--header 'Content-Type: application/json' \
--data-raw '{
    "event_id": "evt_123",
    "card_id": "card_456",
    "cliente_email": "joao.silva@example.com",
    "timestamp": "2026-05-18T12:00:00Z"
}'
```

Retorno em caso de sucesso:
```json
{
    "message": "processed.successfully",
    "success": true
}
```

## Como este projeto escalaria na AWS?

### Banco de dados

Neste projeto utilizaríamos uma instância do Amazon RDS, e nela podemos escalar verticalmente ou 
horizontalmente dependendo do problema que queremos resolver:

- Verticalmente: aumentaríamos os recursos do servidor (CPU, memória, storage)
- Horizontalmente: criaríamos uma réplica read-only do banco (geralmente é usado para fazer leituras mais pesadas. Ex: relatórios), a própria AWS gerencia isto para nós. Neste caso, teríamos que fazer algumas alterações no código da aplicação para que seja possível nos conectarmos nas réplicas.

Há também a opção de fazermos sharding (separar o banco em partições distribuídas), seria uma opção um pouco mais avançada e talvez teríamos que utilizar DynamoDB ao invés de um banco relacional (a AWS também já gerencia os shardings para nós).

### Aplicação

Esta API foi desenvolvida da maneira "padrão", ou seja, um servidor Golang irá subir e expor 2 endpoints. Poderíamos também subir esta aplicação no modelo serverless (com Lambda).

- API Gateway + Lambda: A escalabilidade é gerenciada pela própria AWS (horizontal), um dos pontos que precisamos nos atentar é o limite de requisições simultâneas (caso seja uma API com high throughput), e também a quantidade de memória alocada em cada função Lambda (geralmente o padrão nos atende bem). Um detalhe é que não conseguimos definir a CPU usada nas funções, ela escala proporcionalmente com a memória. Para usar Lambdas neste projeto, precisaríamos fazer algumas alterações no código para carregar a biblioteca da AWS.
- AWS Elastic Beanstalk: Por baixo dos panos ele sobe um load balancer e um EC2 que é responsável por rodar a aplicação, então conseguiríamos escalar horizontalmente. Também é possível aumentarmos os recursos das máquinas se quisermos escalar verticalmente.

Obs: também precisaríamos conectar a aplicação com o banco de dados, para isto teríamos que usar uma rede privada (VPC) em ambos os serviços, assim a conexão entre eles é feita e o banco não fica exposto para a internet.