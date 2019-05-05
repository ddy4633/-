package db

import (
	dbsql "../db/mysql"
	"log"
)

func UserSignup(username string,passwd string)bool{
	stmt,err := dbsql.DBConn().Prepare(
		"insert ignore into tbl_user(`user_name`,`user_pwd`)value(?,?)")
	if err !=nil {
		log.Println("UserSignup Failed ERROR ->",err)
		return false
	}
	defer stmt.Close()
	result,err := stmt.Exec(username,passwd)
	if err != nil{
		log.Println("UserSignup.exec Failed Error ->",err)
		return false
	}
}
