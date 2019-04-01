package handler

import (
	"../db"
	"../util"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

//常量提供给password使用hash加密，盐
const pwd_salt  = "#890"

//注册用户信息
func SigupHandler(w http.ResponseWriter,r *http.Request)  {
	if r.Method == http.MethodGet{
		http.Redirect(w,r,"/static/view/signup.html",302)
		return
		//data,err := ioutil.ReadFile("/static/view/signup.html")
		//if err != nil {
		//	http.Error(w,"ReadFile static signup.html is Failed",404)
		//	return
		//}
		//w.Write(data)
		//return

	}
	//解析表单
	r.ParseForm()
	//获取Post请求的Username和password信息
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	//对账号密码进行判断
	if len(username)<3 || len(password) < 5{
		w.Write([]byte("Invaild parameter"))
		return
	}
	//使用util包中的hash加密passwd
	enc_passwd := util.Sha1([]byte(password+pwd_salt))
	//传递到userfile中去进行数据库处理
	info := db.UserLogin(username,enc_passwd)
	if info {
		w.Write([]byte("SUCCESS"))
		log.Printf("用户%s注册成功",username)
		return
	}else {
		w.Write([]byte("ERROR!"))
	}
}

//登入函数
func SigninHandler(w http.ResponseWriter,r *http.Request) {
	log.Println("当前已经跳转到SigninHandler")
	//地址重写到前端登录页面
	if r.Method == http.MethodGet {
		http.Redirect(w,r,"/static/view/signin.html",302)
		return
		//data, err := ioutil.ReadFile("/static/view/signin.html")
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	return
		//}
		//w.Write(data)
	}
	//1.校验用户名和密码是否正确,Token认证/Cookis认证
	//解析表单取出需要的用户名密码+并且加密
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPasswd := util.Sha1([]byte(password + pwd_salt))
	log.Println(username,password,"->",encPasswd)
	//调用方法开始检验
	Check := db.UserCheck(username,encPasswd)
	log.Println("Check ->",Check)
	if !Check {
		w.Write([]byte("FAILED"))
		log.Println("当前返回的Check",Check)
		return
	}

	//2.生成访问的凭证(Token)
	token := GetToken(username)
	updtoken := db.UpdateToken(username,token)
	if !updtoken {
		w.Write([]byte("FAILED"))
		return
	}

	//3.登入成功后重定向到首页
	//w.Write([]byte("http://" + r.Host +"/static/view/home.html"))
	resp := util.RespMsg{
		Code:0,
		Msg:"ok",
		Data: struct {
			Location string
			Username string
			Token string
		}{
			Location:"http://"+ r.Host +"/static/view/home.html",
			Username:username,
			Token:token,
		},
	}
	w.Write(resp.JsonBytes())
}

//查询用户的信息
func UserInfoHandler(w http.ResponseWriter,r *http.Request) {
	//解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")
	//验证token信息是否有效
	info := TokenisValid(username,token)	//还没有完成
	if !info {
		//跳转返回登入页面刷新时间戳
		http.RedirectHandler("/user/signin",302)
		return
	}
	//查询用户信息
	userinfo,err := db.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	//组装用户信息并且响应给用户
	resp := util.RespMsg{
		Code:0,
		Msg:"ok",
		Data:userinfo,
	}
	w.Write(resp.JsonBytes())
}

//验证Token是否有效
func TokenisValid(username string,token string) bool {
	log.Println("开始判断Token时间")
	//判断Token的时效性是否过期
	RetokenTime := token[32:]
	Retime := time.Now().Unix()
	//将16位的转换成unix()10进制+2小时过期时间和当前的时间戳进行比较
	if sten,err := strconv.ParseInt(RetokenTime,16,10);err == nil {
		if sten+7200 > Retime {
			log.Println("Token超时了跳转回登入页面")
			return false
		}
	}

	//判断Token是否是数据库中的user_token相同的Token
	info := db.TokenIsValue(username,token)
	if ! info {
		log.Println("数据库不存在该Token")
		return false
	}
	//两者是否一致
	return true
}

//生产Token，40位
func GetToken(username string)string{
	//md5(usernmae + time + token_salt)+timestamp[:]8
	time := fmt.Sprintf("%x",time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username+time+"_tokensalt"))
	return tokenPrefix + time[:8]
}