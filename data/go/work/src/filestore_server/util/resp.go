package util

import (
	"encoding/json"
	"log"
)

type RespMsg struct {
	Code int			`json:"code`	//定义状态码
	Msg string			`josn:msg`		//提示信息
	Data interface{}	`json:"data"`	//接受任意类型的数据
}

//生产response的对象
func NewResMsg(code int,msg string,data interface{}) *RespMsg{
	return &RespMsg{
		Code:code,
		Msg: msg,
		Data: data,
	}
}
//将调用的对象转换JSON格式的二进制数组
func (resp *RespMsg) JsonBytes() []byte {
	date,err :=json.Marshal(resp)
	if err != nil {
		log.Println("转换json.Marshal失败 ->",err)
	}
	return date
}