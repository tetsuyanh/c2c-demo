package conf

var c Config

type Config struct {
	Postgres Postgres
}

type Postgres struct {
	DbUser     string
	DbPassword string
	DbHost     string
	DbName     string
}
