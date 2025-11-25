# Terceiro desafio pós go expert (clean architecture)
## Requerimentos
  * golang versão 1.22.1 ou superior
  * Docker e docker-compose

## Como executar
 * Certifique-se que todas as dependências estejam baixadas com `go mod tidy`
 * Inicie o banco de dados e o rabbitmq com `docker-compose up`
 * Espere até que os serviços do `docker-compose up` tenham terminado de subir corretamente
 * Inicie o servidor com `make start`, ou se preferir `go run cmd/ordersystem/main.go cmd/ordersystem/wire_gen.go`
