package configs

import "github.com/spf13/viper"

type Conf struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
	RabbitMQURL       string `mapstructure:"RABBITMQ_URL"`
}

func LoadConfig(path string) (*Conf, error) {
	viper.SetDefault("DB_DRIVER", "mysql")
	viper.SetDefault("DB_HOST", "mysql")
	viper.SetDefault("DB_PORT", "3306")
	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_PASSWORD", "root")
	viper.SetDefault("DB_NAME", "orders")
	viper.SetDefault("WEB_SERVER_PORT", ":8000")
	viper.SetDefault("GRPC_SERVER_PORT", "50051")
	viper.SetDefault("GRAPHQL_SERVER_PORT", "8080")
	viper.SetDefault("RABBITMQ_URL", "amqp://admin:admin@rabbitmq:5672/")

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	var cfg Conf
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
