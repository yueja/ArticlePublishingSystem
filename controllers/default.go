package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"path"
	"test/models"
	"time"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	/*
		//对数据库插入数据
		//1.要有orm对象
		o := orm.NewOrm()
		//2.有一个要插入数据的结构体对象
		user := models.User{}
		//3.对结构体赋值
		user.Name = "1111"
		user.Pwd = "2222"
		//4.插入数据
		_, err := o.Insert(&user)
		if err != nil {
			beego.Info("插入失败", err)
			return
		}
	*/

	/*
			//对数据库查看操作
			//1.要有orm对象
			o := orm.NewOrm()
			//2.查询对象
			user := models.User{}
			//3。指定查询对象值
			//照Id查询
			user.Id=1
			//4.查询
			err:=o.Read(&user)

			//3。指定查询对象值
			//按照Name查询
			user.Name="1111"
			//4.查询
			err:=o.Read(&user,"Name")
			if err != nil {
				beego.Info("查询失败", err)
				return
			}
		     beego.Info("查询成功",user)
	*/

	/*//数据库更新操作
	//1.要有orm对象
	o := orm.NewOrm()
	//2.需要更新的结构体对象
	user := models.User{}
	//3。查到需要更新的对象值
	//照Id查询
	user.Id=1
	err:=o.Read(&user)
	//4.对数据重新赋值
	if err==nil{
		user.Name="张三"
		user.Pwd="123456"
		_,err=o.Update(&user)
		if err !=nil{
			beego.Info("更新失败",err)
			return
		}
	}
	beego.Info("更新成功",user)
	*/

	/*  //删除数据库信息
	//1.要有orm对象
	o := orm.NewOrm()
	//2.需要删除的结构体对象
	user := models.User{}
	//3。指定删除的对象值
	//照Id查询
	user.Id=3
	//4.删除
	_,err:=o.Delete(&user)
	if err!=nil{
		beego.Info("删除失败",err)
		return
	}
	beego.Info("删除成功",user)
	*/
	c.TplName = "register.html"
}

func (c *MainController) Post() {
	//1.拿到数据
	userName := c.GetString("userName")
	pwd := c.GetString("pwd")
	beego.Info("aaaaaaaa")
	//2.对数据校验
	if userName == "" || pwd == "" {
		beego.Info("数据不能为空")
		c.Redirect("/register", 302) //重定向
		return
	}
	//3.插入数据库
	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	user.Pwd = pwd
	_, err := o.Insert(&user)
	if err != nil {
		beego.Info("插入失败", err)
		c.Redirect("/register", 302) //重定向
		return
	}
	//4.返回登录界面
	c.Redirect("/login", 302) //重定向
}

func (c *MainController) ShowLogin() {
	//取出Cookie
	userName := c.Ctx.GetCookie("userName")
	beego.Info("userName is:", userName)
	if userName != "" {
		c.Data["userName"] = userName
		c.Data["checked"] = "checked"
	}

	c.TplName = "login.html"
}

//登录业务逻辑处理
func (c *MainController) HandleLogin() {
	//1	.拿到数据
	userName := c.GetString("userName")
	pwd := c.GetString("pwd")
	remember := c.GetString("remember") //记住用户名

	//2.对数据校验
	if userName == " " || pwd == " " {
		beego.Info("输入数据不合法")
		c.Redirect("/login", 302)
		return
	}
	//3.查询账号密码是否正确
	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Info("登录失败")
		c.Redirect("/login", 302)
		return
	}

	//Cookie设置
	if remember == "on" {
		c.Ctx.SetCookie("userName", userName, 200) //y有效时间为20秒
	} else {
		c.Ctx.SetCookie("userName", userName, -1) //时间为负数代表马上失效
	}
	//Session的设置
	c.SetSession("userName", userName)
	//4.跳转页面
	c.Redirect("/index", 302) //重定向
}

//显示首页内容
func (c *MainController) ShowIndex() {
	userName:= c.GetSession("userName")
	if userName == nil {
		c.Redirect("/login", 302) //重定向
		return
	}
	//传递用户名
	c.Data["userName"]=userName


	o := orm.NewOrm()
	id, err := c.GetInt("select")
	beego.Info(id)
	if err != nil {
		beego.Info("获取失败")
	}

	c.Data["typeid"] =id //文章类型ID

	var articles []models.Article //结构体数组
	//第一种：将所有数据一页显示
	/*_, err := o.QueryTable("Article").All(&articles) //指定查询某张表,并将所有数据放到结构体数组
	if err != nil {
		beego.Info("查询所有信息失败", err)
		return
	}*/

	//第二种：分页显示数据
	//获取数据总数，总页数，当前页码
	//1.获取数据总数
	count, err := o.QueryTable("Article").Filter("ArticleType__Id", id).Count() //Filter()根据下拉框的内容选取同类型的文章
	if err != nil {
		beego.Info("查询失败", err)
		return
	}
	c.Data["count"] = count

	//2.计算总页数
	pagesize := int64(2)                                       //设置每页显示容量
	pageCount := math.Ceil(float64(count) / float64(pagesize)) //math.Ceil此函数用来获得一个比传过来的浮点数最近还比这个数大的整数
	c.Data["pageCount"] = pageCount
	//3.分页显示数据设置

	index, err := c.GetInt("pageIndex") //当前页码
	if err != nil {
		index = 1
	}
	if index <= 0 {
		index = 1
	}
	if index > int(pageCount) {
		index = int(pageCount)
	}
	c.Data["index"] = index //传递当前页码
	start := (int64(index) - 1) * pagesize
	o.QueryTable("Article").Limit(pagesize, start).RelatedSel("ArticleType").All(&articles)
	c.Data["articles"] = articles //将数据传给视图   .RelatedSel("ArticleType")用于显示类别

	//获取类型数据
	var articleType []models.ArticleType
	_, err = o.QueryTable("ArticleType").All(&articleType)
	if err != nil {
		beego.Info("获取类型失败")
		return
	}
	c.Data["articleType"] = articleType //将数据传给视图

	c.TplName = "index.html"
}

//显示增加文章界面
func (c *MainController) ShowAdd() {
	userName:= c.GetSession("userName")
	if userName == nil {
		c.Redirect("/login", 302) //重定向
		return
	}
	//传递用户名
	c.Data["userName"]=userName

	//获取类型数据
	o := orm.NewOrm()
	var articleType []models.ArticleType
	_, err := o.QueryTable("ArticleType").All(&articleType)
	if err != nil {
		beego.Info("获取类型失败")
		return
	}
	c.Data["articleType"] = articleType //将数据传给视图

	c.TplName = "add.html"
}

//处理添加文章界面数据
func (c *MainController) HandleAdd() {
	//1.拿到数据
	artiName := c.GetString("articleName")
	id, err := c.GetInt("select")
	if err != nil {
		beego.Info("查询类型序号错误")
		return
	}
	artiContent := c.GetString("content")
	beego.Info(artiName, artiContent)
	f, h, err := c.GetFile("uploadname")
	defer f.Close()

	//1.限定图片格式
	fileext := path.Ext(h.Filename) //获取文件格式,后缀名
	beego.Info(fileext)
	if fileext != ".jpg" && fileext != ".png" {
		beego.Info("上传文件格式错误")
		return
	}
	//2.限定大小
	if h.Size > 50000000 {
		beego.Info("上传文件过大")
	}
	//3.需要对文件重命名，防止文件名重复
	filename := time.Now().Format("2016-01-02 15:04:05") + fileext //Format("2016-03-12-13-34-26")将时间转换成字符串，后面是显示显示类型

	if err != nil {
		beego.Info("上传图片失败")
		return
	} else {
		c.SaveToFile("uploadname", "./static/img/"+filename) //h.filename也可以取出文件名字
	}
	//2.判断数据是否合法
	if artiName == "" || artiContent == "" {
		beego.Info("添加文章错误")
		return
	}
	//3.插入数据
	o := orm.NewOrm()
	arti := models.Article{}
	arti.Aname = artiName
	arti.Acontent = artiContent
	arti.Aimg = "./static/img/" + filename
	//查找type对象
	artiType := models.ArticleType{Id: id}
	o.Read(&artiType)
	arti.ArticleType = &artiType

	_, err = o.Insert(&arti)
	if err != nil {
		beego.Info("插入数据失败", err)
		return
	}
	//4.插入成功，返回文章界面
	c.Redirect("/addArticle", 302)
}

//显示内容详情页
func (c *MainController) ShowContent() {
	userName:= c.GetSession("userName")
	if userName == nil {
		c.Redirect("/login", 302) //重定向
		return
	}
	//传递用户名
	c.Data["userName"]=userName


	//1.获取文章Id
	id, err := c.GetInt("id")
	if err != nil {
		beego.Info("获取文章Id错误", err)
		return
	}

	//2.查询数据库获取数据
	o := orm.NewOrm()
	arti := models.Article{Id: id}
	err = o.Read(&arti)
	if err != nil {
		beego.Info("查询数据库获取数据错误", err)
		return
	}
	//3.传递数据给视图
	c.Data["article"] = arti //将数据传给视图
	c.TplName = "content.html"
}

//显示编辑界面
func (c *MainController) ShowUpdate() {
	userName:= c.GetSession("userName")
	if userName == nil {
		c.Redirect("/login", 302) //重定向
		return
	}
	//传递用户名
	c.Data["userName"]=userName

	//1.获取文章Id
	id, err := c.GetInt("id")
	if err != nil {
		beego.Info("获取文章Id错误", err)
		return
	}

	//2.查询数据库获取数据
	o := orm.NewOrm()
	arti := models.Article{Id: id}
	err = o.Read(&arti)
	if err != nil {
		beego.Info("查询数据库获取数据错误", err)
		return
	}
	//3.传递数据给视图
	c.Data["article"] = arti //将数据传给视图

	c.TplName = "update.html"
}

//编辑
func (c *MainController) HandleUpdate() {
	//1拿到数据
	var filename string
	id, _ := c.GetInt("id")
	artiName := c.GetString("articleName")
	artiCount := c.GetString("content")
	f, h, _ := c.GetFile("uploadname")
	/*if err != nil{
		beego.Info("上传图片失败")
		c.TplName = "update.html"
		return
	}*/
	defer f.Close()

	//1.限定图片格式
	fileext := path.Ext(h.Filename) //获取文件格式,后缀名
	beego.Info(fileext)

	if fileext != ".jpg" && fileext != ".png" {
		beego.Info("上传文件格式错误")
		return
	}
	//2.限定大小
	if h.Size > 50000000 {
		beego.Info("上传文件过大")
	}
	//3.需要对文件重命名，防止文件名重复
	filename = time.Now().Format("2016-01-02 15:04:05") + fileext //Format("2016-03-12-13-34-26")将时间转换成字符串，后面是显示显示类型
	c.SaveToFile("uploadname", "./static/img/"+filename)          //h.filename也可以取出文件名字

	//2.对数据进行判断或处理
	if artiName == "" || artiCount == "" {
		beego.Info("更新数据失败")
		return
	}

	//3.进行更新操作
	o := orm.NewOrm()
	arti := models.Article{Id: id}
	/*if err != nil {
		beego.Info("查询数据错误", err)
		return
	}*/
	arti.Aname = artiName
	arti.Acontent = artiCount
	arti.Aimg = "./static/img/" + filename
	err := o.Read(&arti)
	_, err = o.Update(&arti, "Aname", "Acontent", "Aimg")
	if err != nil {
		beego.Info("更新失败", err)
		return
	}
	//4.返回列表页面
	c.Redirect("/index", 302)
}

//删除
func (c *MainController) HandleDelete() {
	//1拿到数据
	id, err := c.GetInt("id")
	if err != nil {
		beego.Info("获取id数据错误")
		return
	}
	//2.执行删除操作
	o := orm.NewOrm()
	arti := models.Article{Id: id}
	err = o.Read(&arti)
	if err != nil {
		beego.Info("查询错误", err)
		return
	}
	_, err = o.Delete(&arti)
	if err != nil {
		beego.Info("删除失败", err)
		return
	}
	//3.返回列表页面
	c.Redirect("/index", 302)
}

//显示增加类型
func (c *MainController) ShowAddType() {
	userName:= c.GetSession("userName")
	if userName == nil {
		c.Redirect("/login", 302) //重定向
		return
	}
	//传递用户名
	c.Data["userName"]=userName

	o := orm.NewOrm()
	var artiType []models.ArticleType
	_, err := o.QueryTable("ArticleType").All(&artiType)
	if err != nil {
		beego.Info("没有获取到类型数据", err)
	}
	c.Data["artiType"] = artiType //将数据传给视图
	c.TplName = "addType.html"
}

//处理增加类型传输的信息
func (c *MainController) HandleAddType() {
	TypecName := c.GetString("typeName")
	if TypecName == "" {
		beego.Info("处理增加类型失败")
		c.Redirect("/addType", 302)
	}
	o := orm.NewOrm()
	artiType := models.ArticleType{}
	//3.对结构体赋值
	artiType.Tname = TypecName
	//4.插入数据
	_, err := o.Insert(&artiType)
	if err != nil {
		beego.Info("插入类型失败", err)
		return
	}
	c.Redirect("/addType", 302)
}

//删除类型
func (c *MainController) HandleDeleteType(){
	//1拿到数据
	id, err := c.GetInt("id")
	if err != nil {
		beego.Info("获取id数据错误")
		return
	}
	//2.执行删除操作
	o := orm.NewOrm()
	artiType:= models.ArticleType{Id: id}
	err = o.Read(&artiType)
	if err != nil {
		beego.Info("查询错误", err)
		return
	}
	_, err = o.Delete(&artiType)
	if err != nil {
		beego.Info("删除失败", err)
		return
	}
	//3.返回列表页面
	c.Redirect("/addType", 302)
}

func (c *MainController)LoginOut(){
	c.DelSession("userName")
	c.Redirect("/login",302)
}
