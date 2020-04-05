package main

import (
	"fmt"
	. "go-networks/util"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

type ProxyHandler struct {}
func(* ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//fmt.Println(r.RequestURI)
	//	/a?b=123
	defer func() {
		if err := recover();err != nil{
			w.WriteHeader(500)
			log.Println(err)
		}
	}()
	fmt.Println(r.URL) //    /a?b=123
	fmt.Println(r.URL.Path)  //   /a

	for k,v := range ProxyConfigs{
		if matched,_:= regexp.MatchString(k,r.URL.Path);matched == true{
			//RequestUrl(w,r,v)
			target,_ := url.Parse(v)
			//fmt.Println(target)
			proxy:=httputil.NewSingleHostReverseProxy(target)
			proxy.ServeHTTP(w,r)
			return
		}
	}

	w.Write([]byte("default index"))
}


func main()  {
	http.ListenAndServe(":8080",&ProxyHandler{})

}
