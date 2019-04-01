package handler

import (
	"../meta"
	"../util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)
//处理文件上传函数Uploadhandler
func UploadHandler(w http.ResponseWriter,r *http.Request) {
	//判断用户的请求方式是什么，返回相应的页面
	if r.Method == "GET" {
		//返回HTML页面
		data,err :=ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w,"internel server error")
			return
		}
		io.WriteString(w,string(data))
	}else if r.Method == "POST"{
		//接受文件流及其存储到本地的文件中去
		//FormFile -> 返回文件句柄，文件头，错误信息
		file,filehead,err := r.FormFile("file")
		if err != nil {
			fmt.Println("Failed to get data err ->",err)
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName:filehead.Filename,
			Location:"/tmp/" + filehead.Filename,
			UploadAt:time.Now().Format("2006-01-02 15:04:05"),
		}

		//创建一个本地的文件来接受文件流
		Newfile,err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Create file failed Error ->",err)
			return
		}
		defer Newfile.Close()

		//将临时接收的文件Copy到新的地方保存
		fileMeta.FileSize,err = io.Copy(Newfile,file)
		if err != nil {
			fmt.Printf("IO.Copy is Fail ->",err)
		}
		//偏移量
		Newfile.Seek(0,0)
		fileMeta.FileShall = util.FileSha1(Newfile)
		log.Println(fileMeta)
		//meta.UploadFileMeta(fileMeta)
		a := meta.UploadFileMetaDB(fileMeta)
		log.Println("当前获取的对象",a)
		if  a == false{
			http.Error(w,"获取元数据对象是失败",403)
			os.Exit(1)
		}
		url := "/file/upload/seccess" + "?filehash=" + fileMeta.FileShall
		//调用地址重写函数当执行完成后路由到指定的请求地址
		http.Redirect(w,r,url,http.StatusFound)
	}
}

//UploadSeccessInfo上传完成反馈
func UploadSeccessInfo(w http.ResponseWriter,r *http.Request) {
	//解析操作
	r.ParseForm()
	//取出hash的值
	filehash := r.Form["filehash"][0]
	io.WriteString(w,"Upload finished!!! hash_Value ->"+filehash)
	time.Sleep(3*time.Second)
	http.Redirect(w,r,"/static/view/home.html",302)
}

//实现元数据查询的信息API
func  GetFileMetaHandler(w http.ResponseWriter,r *http.Request) {
	//解析操作
	r.ParseForm()
	//取出hash的值
	filehash := r.Form["filehash"][0]
	//从存储的元数据中取出对应的Value
	//fMeta := meta.GetFileMeta(filehash)	//原来的方法
	fMeta,err := meta.GetFileMetaDB(filehash)
	if err != nil {
		http.Error(w,"Come Mysql date is Failed",500)
		return
	}
	//将取回的值转换成Json格式返回给客户端，返回数组类型、错误
	data,err := json.Marshal(fMeta)
	if err !=nil {
		//以下两个的含义其实是等价的
		http.Error(w,"json.Marshal is failed!",500)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

//提供下载函数
func DownloadHandler(w http.ResponseWriter,r *http.Request) {
	//解析操作
	r.ParseForm()
	//取出对应的hash值
	hash := r.Form.Get("filehash")
	//调用meta函数中的方法获取Metafile
	fhash := meta.GetFileMeta(hash)
	//通过获取Metafile元数据中的文件路径打开文件
	file,err := os.Open(fhash.Location)
	if err != nil {
		http.Error(w,"open file is failed"+fhash.FileName,500)
		return
	}
	defer file.Close()
	//小量的数据可以加载到内存中去,如果数据量大则使用流读取的方式
	//少量的部分读取,然后再刷新内存中的数据
	data,err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w,"ioutil.ReadAall is failed",500)
		return
	}
	//Http的响应头设置,可以让浏览器识别并且下载
	w.Header().Set("Contebt-Type","application/octect-stream")
	w.Header().Set("content-disposition","attachment;filename=\""+fhash.FileName+"\"")
	w.Write(data)
}

//更新元信息接口(重命名)
func FileMetaUpdateHandler(w http.ResponseWriter,r *http.Request) {
	//解析操作
	r.ParseForm()
	//获取表单的op的值
	optype := r.Form.Get("op")
	//获取到hash值
	filesha1 := r.Form.Get("filehash")
	//获取新命名文件的名称
	NewFileName := r.Form.Get("filename")

	//0操作代表修改文件名，后续1,2,3可以扩展
	if optype != "0" {
		http.Error(w,"OpType is Failed",403)
		return
	}else if r.Method != "POST" {
		http.Error(w,"You Should Use POST Request",405)
		return
	}
	//获取元数据
	curlMeta := meta.GetFileMeta(filesha1)
	//赋值重命名
	curlMeta.FileName = NewFileName
	//操作完成后重新调用上传使文件重新保存
	meta.UploadFileMeta(curlMeta)
	//装换成为Json反馈给客户端
	data,err := json.Marshal(curlMeta)
	if err != nil {
		http.Error(w,"Json is Failed",403)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//删除操作接口
func FileDeleteHandler(w http.ResponseWriter,r *http.Request) {
	//解析操作
	r.ParseForm()
	//获取表单的op的值
	fileDelete := r.Form.Get("filehash")
	//获取元信息
	Fmeta := meta.GetFileMeta(fileDelete)
	//调用系统传递文件路径删除文件
	err := os.Remove(Fmeta.Location)
	if err != nil {
		http.Error(w,"System Delete The File is Failed",500)
		return
	}
	//删除信息但是没有删除源文件
	meta.RemoteFileMeta(fileDelete)
	//返回客户端信息,ok
	w.WriteHeader(200)
}