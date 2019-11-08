package metadata

import (
	"NOS/nosServer/tomlConfig"
	"fmt"
	"log"
)

// 为metadata包声明几个全局变量
var (
	WriteLog *log.Logger
	MetaDataHostInfo tomlConfig.METADATA
)
func GetObjectInfo(objetName string) (objectInfo map[string]string) {
	// 获取数据库连接
	engine := MetaDataHostInfo.GetEngine(WriteLog)
	// 根据bojectName拼接sql
	sql := fmt.Sprintf("select sha256_code, is_del from %s.%s where object_name='%s';", MetaDataHostInfo.Database, MetaDataHostInfo.Table, objetName)
	// 到库里执行sql
	res, err := engine.QueryString(sql)
	if err != nil {
		WriteLog.Println("run sql", sql, "is bad, err is", err)
		return nil
	}

	// 返回结果
	for _, info := range res {
		return info
	}
	return nil
}
