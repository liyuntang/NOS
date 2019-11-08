package metadata

import "fmt"

func DeleteObject(objectName, sha256Code string) (isDel bool) {
	// 获取数据库连接
	engine := MetaDataHostInfo.GetEngine(WriteLog)
	// 根据bojectName拼接sql
	sql := fmt.Sprintf("update %s.%s set is_del='1' where object_name='%s' and sha256_code='%s';", MetaDataHostInfo.Database, MetaDataHostInfo.Table, objectName, sha256Code)
	// 到库里执行sql
	_, err := engine.Exec(sql)
	if err != nil {
		WriteLog.Println("run sql", sql, "is bad, err is", err)
		return false
	}
	return true
}
