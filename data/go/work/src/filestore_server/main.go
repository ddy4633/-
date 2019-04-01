package main

import (
	"./handler"
	"fmt"
	"net/http"
)

func main() {
	//建立简单的路由规则
	//静态资源路由
	http.Handle("/static/", http.StripPrefix("/static/",http.FileServer(http.Dir("./static"))))
	//上传
	http.HandleFunc("/file/upload",handler.UploadHandler)
	//上传成功
	http.HandleFunc("/file/upload/seccess",handler.UploadSeccessInfo)
	//获取元数据信息
	http.HandleFunc("/file/meta",handler.GetFileMetaHandler)
	//下载
	http.HandleFunc("/file/download",handler.DownloadHandler)
	//删除
	http.HandleFunc("/file/delete",handler.FileDeleteHandler)
	//更新重命名
	http.HandleFunc("/file/update",handler.FileMetaUpdateHandler)
	//用户注册
	http.HandleFunc("/user/signup",handler.SigupHandler)
	//用户登入
	http.HandleFunc("/user/signin",handler.SigninHandler)
	//用户信息
	http.HandleFunc("/user/info",handler.UserInfoHandler)
	//监听端口
	err := http.ListenAndServe(":8080",nil)
	if err != nil {
		fmt.Println("http.ListenAndServe Error ->",err)
	}
}
