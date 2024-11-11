package sql_server

import (
	"basic"
	"basic/tool/db"
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

type SqlServer struct{}

func GetInstance() *SqlServer {
	return &SqlServer{}
}

func (c *SqlServer) GetName() string {
	return "sql_server"
}

func (r *SqlServer) GetDescribe() string {
	return "sql执行服务"
}

func (r *SqlServer) Register(globalContext *basic.Context) *basic.ComponentMeta {
	command := &basic.ComponentMeta{
		Component: r,
	}
	command.AddParameters(basic.STRING, "-h", "host", true, nil, "")
	command.AddParameters(basic.INT, "-p", "port", true, nil, "")
	command.AddParameters(basic.STRING, "-u", "username", true, nil, "")
	command.AddParameters(basic.STRING, "-w", "password", true, nil, "")
	command.AddParameters(basic.STRING, "-s", "searchPath", true, nil, "")
	command.AddParameters(basic.STRING, "-d", "dbname", true, nil, "")
	command.AddParameters(basic.STRING, "-e", "sql", false, nil, "")
	command.AddParameters(basic.STRING, "-ef", "sqlFile", false, nil, "")
	return command
}

func (r *SqlServer) Start(globalContext *basic.Context) error {
	return nil
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
		return nil
	}
	//未指定sql或者sql文件,进入交互
	reader := bufio.NewReader(os.Stdin)
	for {
		// 读取用户输入
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}
		// 去除输入字符串两侧的空白字符
		input = strings.TrimSpace(input)
		// 检查是否输入了退出命令
		if input == "exit" || input == "q" {
			fmt.Println("Exiting...")
			break
		}
		execSql, err := ExecSql(input, dbBase)
		if err != nil {
			log.Println("sql执行失败:" + err.Error())
		}
		fmt.Println(execSql)
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
			return "", err
		}
		rowsAffected = result.RowsAffected
		log.Printf("Exec Success!,%d rows affected by the operation", rowsAffected)
		return strconv.FormatInt(rowsAffected, 10), nil
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
