package sql_server

import (
	"basic"
	"basic/tool/db"
	othertool "basic/tool/other"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"text/tabwriter"
)

type SqlServer struct{}

func GetInstance() *SqlServer {
	return &SqlServer{}
}

func (r *SqlServer) GetOrder() int {
	return math.MaxInt64
}

func (r *SqlServer) Register(globalContext *basic.Context) *basic.ComponentMeta {
	p1 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-h",
		StandardName: "host",
		Required:     true,
		CheckMethod: func(s string) error {
			if !othertool.CheckIp(s) {
				errors.New("ip is not valid")
			}
			return nil
		},
		Describe: "",
	}
	p2 := basic.Parameter{
		ParamType:    basic.INT,
		CommandName:  "-p",
		StandardName: "host",
		Required:     true,
		CheckMethod: func(s string) error {
			if !othertool.CheckPortByString(s) {
				errors.New("port is not valid")
			}
			return nil
		},
		Describe: "",
	}
	p3 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-u",
		StandardName: "username",
		Required:     true,
		Describe:     "",
	}

	p4 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-w",
		StandardName: "password",
		Required:     true,
		Describe:     "",
	}
	p5 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-d",
		StandardName: "dbname",
		Required:     true,
		Describe:     "",
	}

	p6 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-s",
		StandardName: "searchPath",
		Required:     true,
		Describe:     "",
	}

	p7 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-e",
		StandardName: "sql",
		Required:     false,
		Describe:     "",
	}

	p8 := basic.Parameter{
		ParamType:    basic.STRING,
		CommandName:  "-ef",
		StandardName: "sqlFile",
		Required:     false,
		Describe:     "",
	}

	return &basic.ComponentMeta{
		Key:       "sql_server",
		Describe:  "sql执行服务",
		Component: r,
		Params:    []basic.Parameter{p1, p2, p3, p4, p5, p6, p7, p8},
	}
}
func (r *SqlServer) Do(params map[string]any) (resp []byte) {
	config := &dbtool.DbConfig{
		Host:     params["host"].(string),
		Port:     params["port"].(int),
		Username: params["username"].(string),
		Password: params["password"].(string),
		Dbname:   params["dbname"].(string),
	}
	if params["searchPath"] != nil {
		config.SearchPath = params["searchPath"].(string)
	}
	dbBase, err := dbtool.CreateBaseDbByDbConfig(config)
	if err != nil {
		return []byte("connect database fail, " + err.Error())
	}
	if params["sql"] != nil {
		str, err := ExecSql(params["sql"].(string), dbBase)
		if err != nil {
			return []byte("execute sql fail, " + err.Error())
		}
		return []byte(str)
	}
	if params["sqlFile"] != nil {
		err := ExecSqlFile(params["sqlFile"].(string), dbBase)
		if err != nil {
			return []byte("execute sqlFile fail, " + err.Error())
		}
	}
	return nil
}

func ExecSql(sqlStr string, db *dbtool.BaseDb) (string, error) {
	switch {
	case strings.HasPrefix(strings.ToUpper(sqlStr), "SELECT"):
		rows, err := db.Raw(sqlStr).Rows()
		if err != nil {
			// 处理错误，可能是 SQL 语法错误或数据库连接问题
			log.Printf("Error executing query: %v", err)
			return "", err
		}
		defer rows.Close()
		//返回渲染格式化字符串
		return rendering(rows)
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
	return "", nil
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

func rendering(rows *sql.Rows) (string, error) {
	columns, err := rows.Columns()
	if err != nil {
		log.Printf("Error getting columns: %v", err)
		return "", err
	}
	// 创建一个 bytes.Buffer 来捕获输出
	var output bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&output, 0, 8, 1, '\t', 0)

	// 写入列名
	fmt.Fprintln(w, strings.Join(columns, "\t"))
	// 写入分隔线 假设每个列名或值最多占8个字符宽度（这个值可以根据实际情况调整）
	fmt.Fprintln(&output, strings.Repeat("-", len(columns)*8))
	// 为每列创建一个变量，并创建一个切片来保存这些变量的地址
	columnPtrs := make([]interface{}, len(columns))
	columnValues := make([]interface{}, len(columns))
	for i := range columnPtrs {
		columnPtrs[i] = &columnValues[i]
	}
	// 遍历每一行并写入数据
	for rows.Next() {
		// 使用预先创建的变量地址切片来扫描行数据
		err = rows.Scan(columnPtrs...)
		if err != nil {
			return "", err
		}
		// 将扫描到的数据转换为字符串切片
		var valueStrings []string
		for _, val := range columnValues {
			// 这里需要处理不同类型的数据，确保转换为字符串时不会出错
			// 例如，如果 val 是 nil，则可能需要特别处理
			if val == nil {
				valueStrings = append(valueStrings, "NULL")
			} else {
				valueStrings = append(valueStrings, fmt.Sprintf("%v", val))
			}
		}
		// 写入格式化后的数据行
		fmt.Fprintln(w, strings.Join(valueStrings, "\t"))
	}

	// 检查是否有错误发生（在遍历行之后）
	if err = rows.Err(); err != nil {
		return "", err
	}
	// 刷新 tabwriter.Writer 的缓冲区
	w.Flush()
	return output.String(), nil
}
