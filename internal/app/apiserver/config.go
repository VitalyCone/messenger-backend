package apiserver

type Config struct {
	ApiAddr         string `toml:"bind_addr"`
	DatabaseURL     string `toml:"db_addr"`
	TestDatabaseURL string
}

func NewConfig(apiAddr, dbUrl, TestDatabaseURL string) *Config {
	return &Config{
		ApiAddr:         apiAddr,
		DatabaseURL:     dbUrl,
		TestDatabaseURL: TestDatabaseURL,
	}
}
