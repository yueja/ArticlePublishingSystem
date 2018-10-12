package main

import (
	_ "test/routers"
	_ "test/models"
	"github.com/astaxie/beego"
)

func main() {
	beego.AddFuncMap("Firstpage",ShowFirstpage)
	beego.AddFuncMap("Nextpage",ShowNextpage)
	beego.Run()
}

//视图函数：获取上一页页码
/*
1.在视图中定义视图函数名，加管道  |   funcname
2.一般在main.go里面实现视图函数
3.在main.go函数里面把实现的函数与视图函数关联起来

 */
func  ShowFirstpage (index int)( int){
	return   index-1
}

func  ShowNextpage (index int)(int){
	return   index+1
}