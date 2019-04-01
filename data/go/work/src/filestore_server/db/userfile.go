package db

import (
	dbsql "../db/mysql"
	"fmt"
	"log"
)
//通过账号密码完成User的登入
func UserLogin(username string,password string) bool {
	log.Println("开始数据库处理！")
	//调用自己定义的Mysql连接信息返回的句柄进行Prepare操作
	stmt,err := dbsql.DBConn().Prepare(
		"insert ignore into tbl_user(`user_name`,`user_pwd`)values(?,?)")
	if err != nil{
		fmt.Println("Failed to insert,Error ->",err)
		return false
	}
	defer stmt.Close()

 	//sql语句的执行
	//将这里的变量传入dbsql.DBConn().Prepare的占位符中,返回错误和结果
 	ret,err := stmt.Exec(username,password)
 	if err != nil {
 		fmt.Println("Failed to Exec insert,Error ->",err)
 		return false
	}
 	//判断是否重复注册，重复也是失败，成功返回Ture否则False
 	if rowsaffected,err := ret.RowsAffected();nil == err && rowsaffected >0 {
		return true
	}
 	return false
}

//校验用户名密码是否正确
func UserCheck(username string,encwd string) bool {
	//从数据库tble_user表中进行查询操作取出匹配的username
	stmt,err := dbsql.DBConn().Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		log.Println(err.Error())
		return false
	}
	//进行匹配查询操作：无法查询/info返回空则false
	info ,err := stmt.Query(username)
	log.Println("stmt.Query info->",info,"\nERROR ->",err)
	if err != nil {
		log.Println("stmt.Query ->",err.Error())
		return false
	}else if info == nil {
		log.Println("Username Not Found"+username)
		return false
	}
	//对比两个Password是否是一致
	Prows := dbsql.ParseRows(info)
	log.Println("UserCheck Prows ->",Prows)
	//Prows：大于0、且等于传入的加密过的passwd
	if len(Prows) >0 && string(Prows[0]["user_pwd"].([]byte)) == encwd{
		return true
	}else {
		log.Println(string(Prows[0]["user_pwd"].([]byte)))
		return false
	}
}

//Token数据库校验
func TokenIsValue(username string,token string) bool{
	stmt , err :=dbsql.DBConn().Prepare("select * from tbl_user_token where user_name=? limit 1")
	if err !=nil {
		log.Println("获取*sql错误 ->",err.Error())
		return false
	}
	//进行匹配操作是否存在，真Ture/假false
	info ,err := stmt.Query(username)
	if err != nil {
		log.Println("TokenisValue stmt.Query ->",err)
		return false
	}
	//进行匹配对比
	result := dbsql.ParseRows(info)
	if len(result) >0 && string(result[0]["user_token"].([]byte)) == token{
		return true
	}else {
		return false
	}

}

//刷新用户登入的Token
func UpdateToken(username string,token string) bool{
	stmt ,err :=dbsql.DBConn().Prepare(
		"replace into tbl_user_token(`user_name`,user_token)values(?,?)")
	if err !=nil {
		log.Println("UpdateToken Error ->",err.Error())
		return false
	}
	defer stmt.Close()

	//执行插入表的操作
	_,err = stmt.Exec(username,token)
	if err != nil {
		log.Println("UpdateToken.Exec -> Error",err.Error())
		return false
	}
	return true
}

//用户信息结构体
type UserInfo struct {
	Username 	string
	Email		string
	Phone		string
	SignupAt	string
	status		string
}

//查询mysql用户数据库信息
func GetUserInfo(username string) (UserInfo,error) {
	user := UserInfo{}
	stmt,err := dbsql.DBConn().Prepare(
		"select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		log.Println("GetUserInfo mysql is Failed ->",err)
		return	user,err
	}
	defer stmt.Close()
	//执行操作
	//Scan()可以将数据库中查询到的数据给外部的一个变量来存储
	err = stmt.QueryRow(username).Scan(&user.Username,&user.SignupAt)
	if err != nil {
		return user,err
	}
	return user,nil
}