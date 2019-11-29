package tomlConfig

type PDS struct {
	System system
	Etcd ETCD
}

type system struct {
	Address string	`toml:"address"`
	Port int	`toml:"port"`
	LogFile string	`toml:"logFile"`
	Max_replicas int
	MvOjbectCount int	`toml:"mvOjbectCount"`
}

type ETCD struct {
	Lease int	`toml:"lease"`
	EtcdTimeOut int	`toml:"etcdTimeOut"`
	EtcdDir string	`toml:"etcdDir"`
	EtcdServers []string	`toml:"etcdServers"`
}

