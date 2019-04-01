package meta

import (
	sqldb "../db"
	"log"
)

//文件上传元信息
type FileMeta struct {
	FileShall string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}
//定义全局变量
var Metas map[string]FileMeta

//定义init函数首次运行会执行一次初始化
func init() {
	Metas = make(map[string]FileMeta)
}
//新增或者更FileMeta的元信息
func UploadFileMeta(meta FileMeta) {
	Metas[meta.FileShall] = meta
}
//新增信息加入到mysql中
func UploadFileMetaDB(meta FileMeta) bool{
	log.Println("当前传入的值",meta)
	return sqldb.OnFileUploadFinished(meta.FileShall,meta.FileName,meta.FileSize,meta.Location)
}
//通过mysql获取文件元数据信息
func GetFileMetaDB(meta string) (FileMeta,error){
	dfile,err := sqldb.GetFileMeta(meta)
	if err != nil {
		return FileMeta{},err
	}
	ffile := FileMeta{
		FileSize:dfile.FileSize.Int64,
		FileName:dfile.FileName.String,
		FileShall:dfile.FileHash,
		Location:dfile.FilePath.String,
	}
	return ffile,nil
}

//获取元数据信息
func GetFileMeta(fileshall string) FileMeta {
	return Metas[fileshall]
}
//获取批量的文件元信息列表
func GetListFileMetas(count int) []FileMeta {
	FMetaArray := make([]FileMeta,len(Metas))
	for _,v := range Metas {
		FMetaArray = append(FMetaArray,v)
	}
	return FMetaArray
}

//删除元数据信息
func RemoteFileMeta(fileshal string) {
	delete(Metas,fileshal)
}
