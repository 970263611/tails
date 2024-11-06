package sql

import (
	"basic/tool/db"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

func ExecSql(sqlStr string, db *dbtool.BaseDb) ([]map[string]interface{}, error) {
	switch {
	case strings.HasPrefix(strings.ToUpper(sqlStr), "SELECT"):
		rows, err := db.Raw(sqlStr).Rows()
		if err != nil {
			// 处理错误，可能是 SQL 语法错误或数据库连接问题
			log.Printf("Error executing query: %v", err)
			return nil, err
		}
		defer rows.Close()

		// 使用反射来动态地处理查询结果,所有查询结果都可以映射到一个通用的结构体类型（例如：map[string]interface{}）
		columns, err := rows.Columns()
		if err != nil {
			log.Printf("Error getting columns: %v", err)
			return nil, err
		}
		// 创建一个切片来存储查询结果,每个元素都是一个 map，键是列名，值是列值
		results := make([]map[string]interface{}, 0)
		columnPtrs := make([]interface{}, len(columns))
		columnPtrsValues := make([]interface{}, len(columns))
		for i := range columns {
			columnPtrs[i] = &columnPtrsValues[i]
		}

		for rows.Next() {
			// 将当前行的数据填充到 columnPtrs 中
			if err := rows.Scan(columnPtrs...); err != nil {
				log.Printf("Error scanning row: %v", err)
				continue
			}

			// 创建一个 map 来存储当前行的数据
			rowData := make(map[string]interface{})
			for i, colName := range columns {
				rowData[colName] = columnPtrsValues[i]
			}

			// 将当前行的数据添加到结果切片中
			results = append(results, rowData)
		}

		// 检查是否有任何行被跳过（由于扫描错误）
		if err := rows.Err(); err != nil {
			log.Printf("Error during rows iteration: %v", err)
		}

		// 处理查询结果
		for _, result := range results {
			// 在这里，你可以根据需要处理每个查询结果
			// 例如：打印结果、将其保存到另一个数据结构中等
			fmt.Printf("Query result: %+v\n", result)
		}
		return results, nil
	default:
		// 对于增删改操作，使用 Exec 方法
		var rowsAffected int64
		var err error
		result := db.Exec(sqlStr)
		if err = result.Error; err != nil {
			log.Println("sql执行失败:" + err.Error())
		} else {
			rowsAffected = result.RowsAffected
			log.Printf("Exec Success!,%d rows affected by the operation", rowsAffected)
		}
	}
	return nil, nil
}

func ExecSqlFile(sqlFilePath string, db *dbtool.BaseDb) error {
	_, err := os.Stat(sqlFilePath)
	if os.IsNotExist(err) {
		log.Println("数据库SQL文件不存在:", err)
		return err
	}
	sqls, _ := ioutil.ReadFile(sqlFilePath)
	sqlArr := strings.Split(string(sqls), ";")
	for _, sql := range sqlArr {
		sql = strings.TrimSpace(sql)
		if sql == "" {
			continue
		}
		err := db.Exec(sql).Error
		if err != nil {
			log.Println("sql执行失败:" + err.Error())
			return err
		} else {
			log.Println(sql, "\t success!")
		}
	}
	return nil
}
