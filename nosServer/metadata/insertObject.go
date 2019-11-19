package metadata

import "fmt"

func InsertObject(objectName, sha256Code string) bool {
	// 获取数据库连接
	engine := MetaDataHostInfo.GetEngine(WriteLog)
	// 根据bojectName拼接sql
	sql := fmt.Sprintf("insert into %s.%s(object_name, sha256_code) values ('%s', '%s');", MetaDataHostInfo.Database, MetaDataHostInfo.Table, objectName, sha256Code)
	// 到库里执行sql
	fmt.Println("sql is", sql)
	_, err := engine.Exec(sql)
	if err != nil {
		WriteLog.Println("run sql", sql, "is bad, err is", err)
		return false
	}
	return true
}
