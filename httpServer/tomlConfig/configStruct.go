package tomlConfig

type HTTPSERVER struct {
	HttpServer httpServer
}

type httpServer struct {
	Address string	`toml:"address"`
	Port int		`toml:"port"`
	PdServer []string	`toml:"pdServer"`
	LogFile string	`toml:"logFile"`
}
