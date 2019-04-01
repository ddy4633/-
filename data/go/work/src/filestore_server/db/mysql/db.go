package mysql

import (
	"database/sql"
	"log"

	//下划线的意思是：匿名导入,导入驱动后会初始化自己注册到database中去
	_"github.com/go-sql-driver/mysql"
	"fmt"
	"os"
)

var db *sql.DB
func init() {
	//打开SQL的连接，创建协程安全的对象为长连接
	db , _ = sql.Open("mysql","golang:qaz520133@tcp(docker.poph163.com:3306)/fileserver?charset=utf8")
	log.Println("当前的DB连接完成",db)
	//最大活跃连接数
	db.SetConnMaxLifetime(1000)
	//是否存活
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to Connect to Mysql,Error ->",err)
		os.Exit(1)
	}
}

//实现password的对比
func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns,_ := rows.Columns()
	scanArgs := make([]interface{},len(columns))
	values := make([]interface{},len(columns))
	for i := range values{
		scanArgs[i] = &values[i]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{},0)
	for rows.Next(){
		err := rows.Scan(scanArgs...)
		checkerr(err)

		for j,cool := range values {
			if cool != nil {
				record[columns[j]] = cool
			}
		}
		records = append(records,record)
	}
	return records
}

//错误检查
func checkerr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

//返回数据库的连接对象
func DBConn() *sql.DB {
	return db
}