package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type web1handler struct {

}

func (web1handler) GetIp(request *http.Request) string{
	ips := request.Header.Get("x-forwarded-for")
	if ips != ""{
		ips_list := strings.Split(ips,",")
		if len(ips_list) > 0 && ips_list[0] != ""{
			return ips_list[0]
		}
	}
	return request.RemoteAddr
}

func(this web1handler) ServeHTTP(writer http.ResponseWriter,request *http.Request){
	auth := request.Header.Get("Authorization")
	fmt.Println(auth)
	if auth == ""{
		writer.Header().Set("WWW-Authenticate",`Basic realm="您必须输入用户名和密码"`)
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	//fmt.Println(auth)
	// Basic c3VubG9uZzoxMjM=
	auth_list :=  strings.Split(auth," ")
	if len(auth_list) == 2  && auth_list[0] == "Basic"{
		res,err := base64.StdEncoding.DecodeString(auth_list[1])
		if err == nil && string(res) == "sunlong:123456" {
			writer.Write([]byte(fmt.Sprintf("<h1>web1,来自:%s</h1>",this.GetIp(request))))
			return
		}
	}
	writer.Write([]byte("用户名密码错误"))
}

type web2handler struct {}

func(web2handler) ServeHTTP(writer http.ResponseWriter, request *http.Request)  {
	time.Sleep(5*time.Second)
	writer.Write([]byte("<h1>web2</h1>"))
}


func main(){
	c :=  make(chan os.Signal)
	go(func() {
		http.ListenAndServe(":9091",web1handler{})
	})()

	go(func() {
		http.ListenAndServe(":9092",web2handler{})
	})()

	signal.Notify(c,os.Interrupt)
	s := <- c
	log.Println(s)
}