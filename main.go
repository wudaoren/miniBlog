// Copyright 2014 wudaoren.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package main

import (
	"os"
	"time"

	//标准库
	"encoding/json"
	"fmt"
	"io/ioutil"

	//myblog包
	"miniBlog/asset"
	"miniBlog/controller"
	"miniBlog/global"
	"miniBlog/model"

	//三方包
	"github.com/gin-gonic/gin"         //导入gin框架
	_ "github.com/go-sql-driver/mysql" //使用xorm必须要先导入mysql驱动
	"xorm.io/core"                     //导入xorm框架
	"xorm.io/xorm"                     //导入xorm框架
)

//这里是整个程序的入口
func main() {
	//第一步：读取配置文件
	ReadConfig()
	//第二步：先创建数据库
	CreateDatabase()
	//如果启动程序时执行  "miniBlog install" 则安装数据库
	if len(os.Args) == 2 && os.Args[1] == "install" {
		Install()
		return
	}
	//否则创建web服务
	CreateWebServer()
}

//安装数据库
func Install() {
	var e error
	var tables = []interface{}{
		new(model.Blog),
		new(model.System),
	}
	//删除原有数据库
	global.DB.DropTables(tables...)
	//创建新数据库
	e = global.DB.Sync2(tables...)
	asset.CheckFatalError("数据库创建失败：", e)
	//创建管理账号
	_, e = global.DB.InsertOne(&model.System{
		Id:       1,
		BlogName: "悟道人的博客",
		Keyword:  "常常写写博客",
		Username: "admin",
		Password: asset.Md5("admin"),
	})
	asset.CheckFatalError("初始数据创建失败：", e)
	asset.Debug("miniBlog已经安装成功！")
}

//读取配置文件，方便我们修改数据库账号信息等可变参数，配置文件为json格式，务必注意
func ReadConfig() {
	//从文件读取配置文件
	bt, e := ioutil.ReadFile("config.json")
	asset.CheckFatalError("配置文件读取错误：", e)
	//将配置文件解析后并赋值给Config结构体
	e = json.Unmarshal(bt, &global.Config)
	asset.CheckFatalError("配置文件格式错误：", e)
}

//增加注释
//创建数据库,如果数据库不存在则创建数据库
func CreateDatabase() {
	var conf = global.Config.Mysql
	var e error
	//创建数据库引擎
	global.DB, e = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", conf.Username, conf.Password, conf.Host, conf.Database))
	asset.CheckFatalError("数据库链接错误：", e)
	//设置xorm的时间为本地时间
	global.DB.DatabaseTZ = time.Local
	//设置orm映射方式为表和字段名一致
	global.DB.SetMapper(core.SameMapper{})
	global.DB.ShowSQL(true)
	e = global.DB.Ping()
	asset.CheckFatalError("数据库链接失败：", e)
}

//创建web服务
func CreateWebServer() {
	//创建一个gin实例
	server := gin.New()
	//监听静态文件目录，当浏览器访问 http://abc.com/static的时候，将直接访问该目录的文件
	server.Static("/static", "./static")
	//加载views目录下的所有html模板文件
	server.LoadHTMLGlob("./views/*")
	//使用gin自带的异常恢复中间件，避免出现异常时程序退出
	server.Use(gin.Recovery())
	//公共部分验证（注意：所有的请求都会经过该中间件）
	server.Use(controller.SystemPublicCheck)
	//博客首页
	server.GET("/", controller.SystemIndex)
	//博客详情页
	server.GET("/info/:id", controller.BlogInfoGet)
	//登录页面
	server.GET("/login", controller.SystemLoginGet)
	server.POST("/login", controller.SystemLoginPost)
	//登录验证（注意：下面的所有的请求都会经过该中间件进行登录判断）
	server.Use(controller.SystemCheckLogin)
	//以下路由必须要登录后才可以使用
	blog := server.Group("/blog")
	{
		//添加博客
		blog.GET("/create", controller.BlogCreateGet)
		blog.POST("/create", controller.BlogCreatePost)
		//修改博客
		blog.GET("/update/:id", controller.BlogUpdateGet)
		blog.POST("/update", controller.BlogUpdatePost)
		//删除博客
		blog.GET("/delete/:id", controller.BlogDeleteGet)
	}
	system := server.Group("/system")
	{
		//上传图片
		system.POST("/uploadimage", controller.SystemUploadImage)
		//用户设置
		system.GET("/set", controller.SystemSetGet)
		system.POST("/set", controller.SystemSetPost)
		//退出登录
		system.GET("/logout", controller.SystemLogoutGet)
	}
	//监听服务器端口//
	server.Run(global.Config.Host)

}
