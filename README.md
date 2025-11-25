# Sistema de Pedidos com Clean Architecture (Go)

Visão geral
- Este repositório implementa um sistema de pedidos seguindo os princípios de Clean Architecture em Go.
- Ele expõe três interfaces simultaneamente:
  - API REST (Chi) para criar e listar pedidos
  - API GraphQL (gqlgen) com Playground
  - API gRPC
- A infraestrutura inclui MySQL para persistência e RabbitMQ para publicar eventos OrderCreated.

Stack técnico
- Linguagem: Go 1.25+
- Frameworks/Bibliotecas:
  - Web: github.com/go-chi/chi/v5
  - GraphQL: github.com/99designs/gqlgen
  - gRPC: google.golang.org/grpc
  - Configuração: github.com/spf13/viper
  - Banco de dados: github.com/go-sql-driver/mysql
  - Mensageria: github.com/streadway/amqp (RabbitMQ)
  - Geração de código de DI: github.com/google/wire (wire_gen.go já versionado)
- Gerenciamento de pacotes/build: Go modules (go.mod)
- Serviços em contêiner: docker-compose (MySQL 5.7 e RabbitMQ com management)

Ponto de entrada
- cmd/ordersystem/main.go é o ponto de entrada da aplicação e inicia:
  - Servidor REST (Chi)
  - Servidor gRPC
  - Servidor GraphQL e Playground

Requisitos
- Go 1.25 ou superior
- Docker e Docker Compose

Primeiros passos
1) Clonar e instalar dependências
- Garanta que os módulos estejam em ordem
  - go mod tidy

2) Subir a infraestrutura
- Usando Docker Compose (recomendado para desenvolvimento):
  - docker-compose up -d
  - Serviços:
    - MySQL 5.7 em localhost:3306 (DB: orders, usuário: root, senha: root)
    - RabbitMQ em localhost:5672 com UI de gerenciamento em http://localhost:15672 (guest/guest)

3) Configurar ambiente
- O app lê configurações de um arquivo .env na raiz do repositório (Viper).

Notas sobre portas
- REST: http.ListenAndServe é chamado com o valor de WEB_SERVER_PORT como está. Use um valor como :8000 ou 0.0.0.0:8000.
- gRPC: o servidor escuta em :<GRPC_SERVER_PORT>. Informe apenas o número da porta (ex.: 50051).
- GraphQL: http.ListenAndServe(":"+GRAPHQL_SERVER_PORT). Informe apenas o número da porta (ex.: 8080).

TODOs importantes
- A URL do RabbitMQ está fixa como amqp://guest:guest@localhost:5672/ em cmd/ordersystem/main.go (getRabbitMQChannel). TODO: mover para configuração (.env) para maior flexibilidade.
- O arquivo de licença não está presente. TODO: adicionar um arquivo LICENSE esclarecendo os termos de uso.

Executar a aplicação
- Com Makefile:
  - make start
- Ou diretamente com Go:
  - go run cmd/ordersystem/main.go cmd/ordersystem/wire_gen.go

Serviços disponíveis após iniciar
- API REST (Chi)
  - Endereço base: http://localhost:8000 se usar WEB_SERVER_PORT=:8000
  - Endpoints:
    - POST /order — cria um pedido
    - GET  /order — lista pedidos
- GraphQL
  - Playground: http://localhost:8080/
  - Endpoint:   http://localhost:8080/query
- gRPC
  - Escuta em localhost:50051 por padrão no exemplo do .env.

Exemplos rápidos REST
- Veja o arquivo api.http para exemplos prontos para uso em REST Client (ex.: VS Code / IDEs JetBrains):
  - POST http://localhost:8000/order
  - GET  http://localhost:8000/order

Exemplos de uso gRPC
- Usando grpcurl (instale separadamente):
  - Descobrir serviços: grpcurl -plaintext localhost:50051 list
  - Descrever serviço: grpcurl -plaintext localhost:50051 list cleanarch.internal.infra.grpc.pb.OrderService
  - Invocar (exemplo — ajuste pacote/método se tiver mudado):
    - grpcurl -plaintext -d '{"id":"<uuid>","price":100,"tax":10}' localhost:50051 cleanarch.internal.infra.grpc.pb.OrderService/CreateOrder

Exemplo de uso GraphQL
- Abra o Playground em http://localhost:8080/
- As consultas/mutações dependem do schema em internal/infra/graph/schema.graphqls. Operações comuns:
  - mutation para criar um pedido
  - query para listar pedidos

Variáveis de ambiente
- A aplicação utiliza as seguintes variáveis (via Viper e .env):
  - DB_DRIVER
  - DB_HOST
  - DB_PORT
  - DB_USER
  - DB_PASSWORD
  - DB_NAME
  - WEB_SERVER_PORT
  - GRPC_SERVER_PORT
  - GRAPHQL_SERVER_PORT
- TODO: Adicionar RABBITMQ_URL (e alterar o código para lê-la) em vez de manter fixo no código.

Scripts e ferramentas
- Alvos do Makefile:
  - start: go run cmd/ordersystem/main.go cmd/ordersystem/wire_gen.go
- Serviços do Docker Compose:
  - mysql (5.7)
  - rabbitmq (3-management)
- Ferramentas de geração de código no repositório:
  - wire (google/wire) — DI gerado em cmd/ordersystem/wire_gen.go. Não é necessário executar a geração para rodar, a menos que você altere o wire.go.
  - gqlgen — Artefatos do GraphQL já gerados; gere novamente apenas se modificar schema/resolvers.

Testes
- Executar todos os testes:
  - go test ./...
- Pacotes de teste notáveis:
  - internal/entity
  - internal/infra/database
  - pkg/events

Estrutura do projeto (alto nível)
- cmd/ordersystem/           — ponto de entrada da aplicação, wire (DI)
- configs/                   — carregamento de configuração (Viper)
- internal/entity/           — entidades de domínio e interfaces de repositório
- internal/usecase/          — casos de uso da aplicação (criar/listar pedidos)
- internal/infra/database/   — implementação do repositório MySQL
- internal/infra/web/        — handlers REST e servidor web
- internal/infra/grpc/       — proto gRPC, pb e implementação do serviço
- internal/infra/graph/      — schema GraphQL, resolvers e código gerado
- internal/event/            — eventos de domínio e handlers
- pkg/events/                — abstrações e implementação do despachante de eventos
- docker-compose.yaml        — infraestrutura local (MySQL, RabbitMQ)
- api.http                   — exemplos REST
- Makefile                   — comandos de conveniência

Licença
- TODO: Adicionar arquivo LICENSE.
