package tomlConfig

type NOS struct {
	System system
}

type system struct {
	Address string	`toml:"address"`
	Port int		`toml:"port"`
	DataDir string	`toml:"dataDir"`
	LogFile string	`toml:"logFile"`
}
