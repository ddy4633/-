package db

import (
	mydb "../db/mysql"
	"database/sql"
	"fmt"
	"log"
)
//文件上传完成，成功存储到Mysql中
func OnFileUploadFinished(FileHash string,FileName string,FileSize int64,FilePath string) bool {
	//返回stmt操作sql,返回要操作的数据库列表
	log.Println("当前已经执行到了dbfile文件")
	//Prepare优先使用，防止SQL注入攻击，可以实现自定义参数的查询
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file(`file_hash`,`file_name`,`file_size`," +
			"`file_path`,`status`) values(?,?,?,?,1)")
	log.Println("当前取出来的stmt值",stmt)
	if err != nil {
		fmt.Println("Failed to Prepare statement ERROR ->",err.Error())
		return false
	}
	log.Println("从数据中返回的对象",stmt)
	defer stmt.Close()
	//返回数据库影响或插入相关的结果
	ret,err :=stmt.Exec(FileHash,FileName,FileSize,FilePath)
	log.Println("拿到的数据",ret)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}else if rf,err := ret.RowsAffected();nil == err{	//判断是否hash被重复的插入
		if rf <= 0 {
			fmt.Println("File with hash：%s has been upload before",FileHash)
		}
		return true
	}
	return false
}

//定义一个结构体表示一个表的字段
type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FilePath sql.NullString
}

//查询：从mysql获取元数据信息
func GetFileMeta(filehash string) (*TableFile,error){
	//获取查询相关的对象
	stmt,err := mydb.DBConn().Prepare(
		"select file_hash,file_path,file_name,file_size from tbl_file "+
			"where file_hash=? and status=1 limit 1")
	//如果有错误则返回异常
	if err != nil {
		log.Println(err.Error())
		return nil,err
	}
	defer stmt.Close()

	dfile := TableFile{}
	//QueryRow返回参数对象给dfile赋值scan中是用&取值
	err =stmt.QueryRow(filehash).Scan(
		&dfile.FileHash,&dfile.FilePath,&dfile.FileName,&dfile.FileSize)
	if err !=nil {
		log.Println(err.Error())
	}
	return &dfile,nil
}