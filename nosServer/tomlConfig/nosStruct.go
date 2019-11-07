package tomlConfig

type NOS struct {
	System SYSTEM
	Etcd ETCD
}

type SYSTEM struct {
	Address string	`toml:"address"`
	Port int	`toml:"port"`
	LogFile string	`toml:"logFile"`
}

type ETCD struct {
	EtcdTimeOut int	`toml:"etcdTimeOut"`
	EtcdDir string	`toml:"etcdDir"`
	EtcdServers []string	`toml:"etcdServers"`
}
