package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	"cleanarch/configs"
	"cleanarch/internal/event/handler"
	"cleanarch/internal/infra/graph"
	"cleanarch/internal/infra/grpc/pb"
	"cleanarch/internal/infra/grpc/service"
	"cleanarch/internal/infra/web/webserver"
	"cleanarch/pkg/events"

	graphqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	if err != nil {
		panic(err)
	}

	rabbitMQChannel := getRabbitMQChannel(configs.RabbitMQURL)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := NewListOrdersUseCase(db)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.Router.Post("/order", webOrderHandler.Create)
	webserver.Router.Get("/order", webOrderHandler.List)
	fmt.Println("Web Server started on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("gRPC Server started on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphqlhandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase:   *listOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("GraphQL Server started on port", configs.GraphQLServerPort)
	errhttp := http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
	if errhttp != nil {
		return
	}
}

func getRabbitMQChannel(rabbitURL string) *amqp.Channel {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
