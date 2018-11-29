// Copyright 2014 wudaoren.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package controller

import (
	"fmt"
	. "miniBlog/asset"
	. "miniBlog/global"
	"miniBlog/model"
	"path"

	"github.com/gin-gonic/gin"
)

//系统首页
func SystemIndex(c *gin.Context) {
	var query struct {
		Page   int
		Search string
	}
	if e := c.ShouldBindQuery(&query); e != nil {
		HtmlMsg(c, MSG_ERR, "提交的数据有误", nil)
		return
	}
	//如果传入页码则默认为第一页
	if query.Page == 0 {
		query.Page = 1
	}
	//创建要查询的列表数据结构
	find := make([]struct {
		Id         int
		Title      string
		Keyword    string
		CreateDate string
	}, 0)
	//查询数据
	sess := DB.Table("Blog").Desc("CreateDate", "Id")
	if query.Search != "" {
		sess.Where("Title like ?", "%"+query.Search+"%")
	}
	PageTotal := PageQuery(sess, query.Page, 10, &find)
	//渲染页面
	HtmlData(c, "systemIndex.html", gin.H{
		"List":     find,
		"PageHtml": CreatePageMenu(fmt.Sprintf("/?Search=%s&Page=", query.Search), query.Page, PageTotal),
	})
}

//登录页面
func SystemLoginGet(c *gin.Context) {
	HtmlData(c, "systemLogin.html", nil)
}

//公共部分验证
func SystemPublicCheck(c *gin.Context) {
	sess := UseSession(c)
	//如果是post方式请求则检查form表单里面的token是否合法，如果合法则删除该token，保证token只使用一次，页就避免了表单重复提交
	if c.Request.Method == "POST" {
		if c.Query("img") != "" {
			return
		}
		if token := c.PostForm("Token"); token != "" && sess.Get(token) != nil {
			sess.Del(token)
		} else {
			HtmlMsg(c, MSG_ERR, "非法提交数据", "/")
			c.Abort()
		}
	}
}

//登录验证,如果没有登录则返回系统首页
func SystemCheckLogin(c *gin.Context) {
	sess := UseSession(c)
	if sess.Get("login") == nil {
		HtmlMsg(c, MSG_ERR, "你还没有登录", "/")
		c.Abort()
	}
}

//登录页面数据提交
func SystemLoginPost(c *gin.Context) {
	var req struct {
		Username string `binding:"required,min=5"` //
		Password string `binding:"required,min=5"` //
	}
	if e := c.ShouldBind(&req); e != nil {
		HtmlMsg(c, MSG_ERR, "提交的数据有误", nil)
		return
	}
	system := &model.System{Id: 1}
	DB.Get(system)
	req.Password = Md5(req.Password)
	if system.Username == req.Username && req.Password == system.Password {
		sess := UseSession(c)
		sess.Set("login", true)
		HtmlMsg(c, MSG_SUCC, "登录成功", "/")
	} else {
		HtmlMsg(c, MSG_ERR, "登录失败", nil)
	}
}

//退出登录
func SystemLogoutGet(c *gin.Context) {
	sess := UseSession(c)
	sess.Clear()
	HtmlMsg(c, MSG_SUCC, "已退出登录", "/")
}

//个人设置页面
func SystemSetGet(c *gin.Context) {
	HtmlData(c, "systemSet.html", nil)
}

//个人设置页面
func SystemSetPost(c *gin.Context) {
	var req struct {
		BlogName   string `binding:"required"`        //
		Keyword    string `binding:"required"`        //
		Username   string `binding:"required,min=5"`  //
		Password   string `binding:"omitempty,min=5"` //
		RePassword string `binding:"omitempty,min=5"` //
	}
	if e := c.ShouldBind(&req); e != nil {
		HtmlMsg(c, MSG_ERR, "提交的数据有误", nil)
		return
	}
	if req.Password != req.RePassword {
		HtmlMsg(c, MSG_ERR, "两次输入的密码不一致", nil)
		return
	}
	system := &model.System{
		BlogName: req.BlogName,
		Keyword:  req.Keyword,
		Username: req.Username,
		Password: req.Password,
	}
	if system.Password != "" {
		system.Password = Md5(system.Password)
	}
	if l, e := DB.Where("Id=?", 1).Update(system); l == 0 || e != nil {
		HtmlMsg(c, MSG_ERR, "设置失败", nil)
		return
	}
	HtmlMsg(c, MSG_SUCC, "设置成功", "")

}

//上传图片
func SystemUploadImage(c *gin.Context) {
	file, _ := c.FormFile("filename")
	ext := path.Ext(file.Filename)
	//只允许上传jpg格式
	if ext != ".jpg" {
		c.JSON(200, gin.H{"errno": 1})
		return
	}
	//只允许上传3M以内大小的文件
	if file.Size > 3*1024*1024 {
		c.JSON(200, gin.H{"errno": 2})
		return
	}
	//将上传的文件保存到上传文件目录
	dst := UPLOAD_DIR + file.Filename
	if e := c.SaveUploadedFile(file, "."+dst); e != nil {
		c.JSON(200, gin.H{"errno": 3})
		return
	}
	c.JSON(200, gin.H{
		"errno": 0,
		"data":  []string{dst},
	})
}
