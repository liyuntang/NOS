package metadata

import "fmt"

func Sha256CodeISOK(sha256Code string) bool {
	// 获取数据库连接
	engine := MetaDataHostInfo.GetEngine(WriteLog)
	// 根据bojectName拼接sql
	sql := fmt.Sprintf("select * from %s.%s where sha256_code='%s';", MetaDataHostInfo.Database, MetaDataHostInfo.Table, sha256Code)
	// 到库里执行sql
	res, err := engine.QueryString(sql)
	if err != nil {
		WriteLog.Println("run sql", sql, "is bad, err is", err)
		return false
	}
	// 说明执行数据库查询成功，下面需要判断ojbect的具体情况
	// object不存在的情况有两种：
	// 1、len(info) == 0
	// 2、is_del == 1
	for _, info := range res {
		fmt.Println(info)
		if len(info) == 0 {
			// 说明object不存在
			return false
		}
		// 说明object存在，返回该map
		return true
	}
	return false
}
