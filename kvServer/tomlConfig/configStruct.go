package tomlConfig

type KVSERVER struct {
	KvServer kvServer
}

type kvServer struct {
	Address string	`toml:"address"`
	Port int		`toml:"port"`
	Kvserver []string	`toml:"kvservers"`
	DataDir string	`toml:"dataDir"`
	LogFile string	`toml:"logFile"`
}
