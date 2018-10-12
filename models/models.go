package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id int
	Name string
	Pwd string
	Article [] *Article  `orm:"rel(m2m)"`
}
//文章结构体
type Article struct{
	Id int  `orm:"pk;auto"`//设置主键，自增
	Aname string `orm:"size(20)"`  //文章名称,设置长度
	Atime time.Time  `orm:"auto_now"`//文章时间
	Acount int `orm:"default(0);null"`//阅读量
	Acontent string //文章内容
	Aimg string //图片


	ArticleType *ArticleType `orm:"rel(fk)"`  //外键
	User [] * User  `orm:"reverse(many)"`
}
//类型表
type ArticleType struct{
	Id int
	Tname string
	Article [] *Article `orm:"reverse(many)"`  //一对多
}

func init(){
	//设置数据库基本信息
	orm.RegisterDataBase("default","mysql","root:root123@tcp(127.0.0.1:3306)/test?charset=utf8")
	//映射model数据
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	//生成表
	orm.RunSyncdb("default",false,true)
}