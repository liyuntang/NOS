package tomlConfig

type PDSREVER struct {
	PdServer pdServer
}

type pdServer struct {
	Address string	`toml:"address"`
	Port int		`toml:"port"`
	Pdserver []string	`toml:"pdservers"`
	KvServer []string	`toml:"kvServer"`
	LogFile string	`toml:"logFile"`
}
