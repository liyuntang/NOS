package tomlConfig

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_"github.com/go-sql-driver/mysql"
	"log"
)

type NOS struct {
	System system
	Etcd ETCD
	Metadata METADATA
}

type system struct {
	Address string	`toml:"address"`
	Port int	`toml:"port"`
	TmpDir string	`toml:"tmpDir"`
	LogFile string	`toml:"logFile"`
}

type ETCD struct {
	EtcdTimeOut int	`toml:"etcdTimeOut"`
	EtcdDir string	`toml:"etcdDir"`
	EtcdServers []string	`toml:"etcdServers"`
}

type METADATA struct {
	User string	`toml:"user"`
	Passwd string	`toml:"passwd"`
	Address string	`toml:"address"`
	Port int	`toml:"port"`
	Charset string	`toml:"chrset"`
	Database string	`toml:"database"`
	Table string	`toml:"table"`
}

func (meta *METADATA)GetEngine(logger *log.Logger) (engine *xorm.Engine) {
	endPoint := fmt.Sprintf("%s:%d", meta.Address, meta.Port)
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=true", meta.User, meta.Passwd, endPoint, meta.Database, meta.Charset)
	engine, err := xorm.NewEngine("mysql", dataSource)
	if err != nil {
		logger.Println("init db connection", endPoint, "is bad")
		return nil
	}
	return engine
}

