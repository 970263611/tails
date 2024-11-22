package db_tool

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/db"
	"bufio"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type SqlServer struct{}

func GetInstance(globalContext iface.Context) iface.Component {
	return &SqlServer{}
}

func (c *SqlServer) GetName() string {
	return "db_tool"
}

func (r *SqlServer) GetDescribe() string {
	return "sql执行服务 例: db_tool -h 127.0.0.1 -p 5432 -u postgres -w postgre -d postgres -s public -e \"select * from test \" "
}

func (r *SqlServer) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, cons.LOWER_H, "ip", "host", true, nil, "数据库Ip")
	cm.AddParameters(cons.INT, cons.LOWER_P, "port", "port", true, nil, "数据库port")
	cm.AddParameters(cons.STRING, cons.LOWER_U, "username", "username", true, nil, "数据库登录用户名")
	cm.AddParameters(cons.STRING, cons.LOWER_W, "password", "password", true, nil, "数据库登录密码")
	cm.AddParameters(cons.STRING, cons.LOWER_D, "dbname", "dbname", true, nil, "数据库名称")
	cm.AddParameters(cons.STRING, cons.LOWER_S, "searchPath", "searchPath", true, nil, "数据库结构")
	cm.AddParameters(cons.STRING, cons.LOWER_C, "sql", "sql", false, nil, "执行sql语句")
	cm.AddParameters(cons.STRING, cons.LOWER_O, "outPutFile", "outPutFile", false, nil, "执行sql语句后,将查询结果导出到指定文件")
	cm.AddParameters(cons.STRING, cons.LOWER_N, "sqlFile", "sqlFile", false, nil, "执行sql文件")
	cm.AddParameters(cons.NO_VALUE, cons.LOWER_F, "force", "force", false, nil, "强制执行")
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
	//执行单个sql,并输出到指定文件
	if params["sql"] != nil {
		str, err := ExecSql(params, dbBase)
		if err != nil {
			log.Error(err)
			return []byte(err.Error())
		}
		outPutFile, ok := params["outPutFile"].(string)
		if !ok {
			//命令行没有指定输出到文件,直接将查询结果返回
			return []byte(str)
		} else {
			// 检查 outPutFile 是否是一个完整的文件路径
			isFullPath := strings.ContainsAny(outPutFile, string(os.PathSeparator))
			// 如果不是完整路径，则添加当前目录
			if !isFullPath {
				wd, err := os.Getwd()
				if err != nil {
					log.Error("获取当前目录出错: " + err.Error())
					return []byte("获取当前目录出错: " + err.Error())
				}
				outPutFile = filepath.Join(wd, outPutFile)
			}
			err := os.WriteFile(outPutFile, []byte(str), 0644)
			if err != nil {
				log.Error("写入文件错误: " + err.Error())
				return []byte("写入文件错误: " + err.Error())
			}
			return []byte("输出到指定文件成功!")
		}
	}
	//执行sql文件
	if params["sqlFile"] != nil {
		err := ExecSqlFile(params["sqlFile"].(string), dbBase)
		if err != nil {
			return []byte(err.Error())
		}
		return []byte("sql文件执行成功")
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
		execSql, err := ExecSqlRead(input, dbBase)
		if err != nil {
			log.Error(err.Error())
		}
		fmt.Println(execSql)
	}
	return []byte("sql执行结束")
}

func ExecSql(params map[string]any, db *dbtool.BaseDb) (string, error) {
	sqlStr := params["sql"].(string)
	re := regexp.MustCompile(`^"+|"+$`)
	sqlStr = re.ReplaceAllString(sqlStr, "")
	switch {
	case strings.HasPrefix(strings.ToUpper(sqlStr), "SELECT"):
		rows, err := db.Raw(sqlStr).Rows()
		if err != nil {
			log.Error("sql执行失败: ", err)
			return "", err
		}
		defer rows.Close()
		//返回渲染格式化字符串
		return rendering(rows)
	default:
		// 对于增删改操作，需要判断是否强制执行
		_, ok := params["force"]
		if !ok {
			return "", errors.New("发现非查询语句,未输入强制执行命令,请确认！")
		}
		var rowsAffected int64
		var err error
		result := db.Exec(sqlStr)
		if err = result.Error; err != nil {
			log.Error("sql执行失败:" + err.Error())
			return "", err
		}
		rowsAffected = result.RowsAffected
		log.Info(fmt.Sprintf("Exec Success!,%d 行已操作", rowsAffected))
		return fmt.Sprintf("Exec Success!,%d 行已操作", rowsAffected), nil
	}
}
func ExecSqlRead(sqlStr string, db *dbtool.BaseDb) (string, error) {
	re := regexp.MustCompile(`^"+|"+$`)
	sqlStr = re.ReplaceAllString(sqlStr, "")
	switch {
	case strings.HasPrefix(strings.ToUpper(sqlStr), "SELECT"):
		rows, err := db.Raw(sqlStr).Rows()
		if err != nil {
			log.Error("sql执行失败: ", err)
			return "", err
		}
		defer rows.Close()
		//返回渲染格式化字符串
		return rendering(rows)
	default:
		// 对于增删改操作，需要判断是否强制执行
		var rowsAffected int64
		var err error
		result := db.Exec(sqlStr)
		if err = result.Error; err != nil {
			log.Error("sql执行失败:" + err.Error())
			return "", err
		}
		rowsAffected = result.RowsAffected
		log.Info(fmt.Sprintf("Exec Success!,%d 行已操作", rowsAffected))
		return fmt.Sprintf("Exec Success!,%d 行已操作", rowsAffected), nil
	}
}

func ExecSqlFile(sqlFilePath string, db *dbtool.BaseDb) error {
	_, err := os.Stat(sqlFilePath)
	if os.IsNotExist(err) {
		log.Error("数据库SQL文件不存在:", err)
		return err
	}
	sqls, _ := os.ReadFile(sqlFilePath)
	sqlArr := strings.Split(string(sqls), ";")
	for _, sql := range sqlArr {
		sql = strings.TrimSpace(sql)
		if sql == "" {
			continue
		}
		err := db.Exec(sql).Error
		if err != nil {
			log.Error("sql执行失败:" + err.Error())
			return err
		} else {
			log.Info(sql, "\t success!")
		}
	}
	return nil
}

func rendering(rows *sql.Rows) (string, error) {
	columns, err := rows.Columns()
	if err != nil {
		log.Error(fmt.Sprintf("错误的得到了字段集: %v", err))
		return "", err
	}
	// 创建一个 bytes.Buffer 来捕获输出
	buffer := new(bytes.Buffer)
	// 创建tablewriter.Table
	table := tablewriter.NewWriter(buffer)
	table.SetHeader(columns)
	// 为每列创建一个变量，并创建一个切片来保存这些变量的地址
	columnPtrs := make([]interface{}, len(columns))
	columnValues := make([]interface{}, len(columns))
	// 遍历每一行并写入数据
	for rows.Next() {
		for i := range columnPtrs {
			columnPtrs[i] = &columnValues[i]
		}
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
		table.Append(valueStrings)
	}
	// 获取捕获的输出
	table.Render()
	output := buffer.String()
	return output, nil
}
