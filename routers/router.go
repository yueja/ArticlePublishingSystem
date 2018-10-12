package routers

import (
	"test/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.MainController{})
    //注意：当实现了自定义GET请求方法，请求将不会访问默认方法
	beego.Router("/login", &controllers.MainController{},"get:ShowLogin;post:HandleLogin")//自定义get请求方法
	//后台页面
	beego.Router("/index", &controllers.MainController{},"get:ShowIndex")
    //增加新闻
	beego.Router("/addArticle", &controllers.MainController{},"get:ShowAdd;post:HandleAdd")
    //新闻详情
	beego.Router("/content", &controllers.MainController{},"get:ShowContent")
//编辑
	beego.Router("/update", &controllers.MainController{},"get:ShowUpdate;post:HandleUpdate")
    //删除
	beego.Router("/delete", &controllers.MainController{},"get:HandleDelete")
   //增加类型
	beego.Router("/addType", &controllers.MainController{},"get:ShowAddType;post:HandleAddType")
    //删除类型
	beego.Router("/deleteType", &controllers.MainController{},"get:HandleDeleteType")
//退出登录
	beego.Router("/loginout", &controllers.MainController{},"get:LoginOut")
}
