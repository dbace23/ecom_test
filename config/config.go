package config

import "os"

type Config struct {
	MongoURI       string
	MongoDBName    string
	ShoppingPort   string
	PaymentPort    string
	PaymentBaseURL string
}

func Load() Config {
	return Config{
		MongoURI:       envOr("MONGO_URI", "mongodb+srv://user:pass@cluster-url/?retryWrites=true&w=majority"),
		MongoDBName:    envOr("MONGO_DB_NAME", "ecom_db"),
		ShoppingPort:   envOr("SHOPPING_PORT", ":9063"),
		PaymentPort:    envOr("PAYMENT_PORT", ":9053"),
		PaymentBaseURL: envOr("PAYMENT_BASE_URL", "http://localhost:9053"),
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
