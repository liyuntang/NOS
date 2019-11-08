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
func ObjectISOK(objetName string) (isok bool, objectInfo map[string]string) {
	// 获取数据库连接
	engine := MetaDataHostInfo.GetEngine(WriteLog)
	// 根据bojectName拼接sql
	sql := fmt.Sprintf("select sha256_code, is_del from %s.%s where object_name='%s';", MetaDataHostInfo.Database, MetaDataHostInfo.Table, objetName)
	// 到库里执行sql
	res, err := engine.QueryString(sql)
	if err != nil {
		WriteLog.Println("run sql", sql, "is bad, err is", err)
		return false, nil
	}
	// 说明执行数据库查询成功，下面需要判断ojbect的具体情况
	// object不存在的情况有两种：
	// 1、len(info) == 0
	// 2、is_del == 1
	for _, info := range res {
		if len(info) == 0 {
			// 说明object不存在
			return false, nil
		}
		if info["is_del"] == "1" {
			// 说明该object已经与其数据做了软删除，如果是method是get或delete的话则返回404，如果是put的话应该用新的sha256值替换当前的sha256_code值，并将is_del设置为0
			return false, nil
		}
		// 说明object存在，返回该map
		return true, info
	}
	return false, nil
}
