package tomlConfig

type KVS struct {
	System SYSTEM
	Etcd ETCD
}

type SYSTEM struct {
	Address string	`toml:"address"`
	Port int	`toml:"port"`
	DataDir string	`toml:"dataDir"`
	LogFile string	`toml:"logFile"`
}

type ETCD struct {
	Lease int64 	`toml:"lease"`
	EtcdTimeOut int	`toml:"etcdTimeOut"`
	EtcdDir string	`toml:"etcdDir"`
	EtcdServers []string	`toml:"etcdServers"`
}
