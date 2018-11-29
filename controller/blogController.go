// Copyright 2014 wudaoren.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package controller

import (
	"html/template"
	. "miniBlog/asset"
	. "miniBlog/global"
	"miniBlog/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//查看博客信息
func BlogInfoGet(c *gin.Context) {
	//从param里面读取博客Id
	id, _ := strconv.Atoi(c.Param("id"))
	blog := &model.Blog{Id: id}
	if find, e := DB.Get(blog); !find || e != nil {
		HtmlMsg(c, MSG_ERR, "博客不存在", "/")
		return
	}
	blog.HTMLContent = template.HTML(blog.Content)
	HtmlData(c, "blogInfo.html", blog)
}

//创建博客页面
func BlogCreateGet(c *gin.Context) {
	HtmlData(c, "blogCreate.html", nil)
}

//创建博客
func BlogCreatePost(c *gin.Context) {
	var req struct {
		Title      string `binding:"required"` //标题
		Keyword    string `binding:"required"` //关键词
		Content    string `binding:"required"` //博客内容
		CreateDate string //发布日期（可以修改）
	}
	if e := c.ShouldBind(&req); e != nil {
		HtmlMsg(c, MSG_ERR, "提交的数据有误", nil)
		return
	}
	if req.CreateDate == "" {
		req.CreateDate = time.Now().Format("2006-01-02")
	}
	blog := &model.Blog{
		Title:      req.Title,
		Keyword:    req.Keyword,
		Content:    req.Content,
		CreateDate: req.CreateDate,
	}
	if id, e := DB.InsertOne(blog); id < 0 || e != nil {
		HtmlMsg(c, MSG_ERR, "博客添加失败", nil)
		return
	}
	HtmlMsg(c, MSG_SUCC, "添加成功", "/blog/create")
}

//编辑博客页面
func BlogUpdateGet(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	blog := &model.Blog{Id: id}
	if find, _ := DB.Get(blog); !find {
		HtmlMsg(c, MSG_ERR, "博客不存在", nil)
		return
	}
	HtmlData(c, "blogUpdate.html", blog)
}

//编辑博客
func BlogUpdatePost(c *gin.Context) {
	var req struct {
		Id         int    `binding:"required"` //id
		Title      string `binding:"required"` //标题
		Keyword    string `binding:"required"` //关键词
		Content    string `binding:"required"` //博客内容
		CreateDate string `binding:"required"` //发布日期（可以修改）
	}
	if e := c.ShouldBind(&req); e != nil {
		Debug("修改博客错误：", e)
		HtmlMsg(c, MSG_ERR, "提交的数据有误", nil)
		return
	}
	blog := &model.Blog{
		Id:         req.Id,
		Title:      req.Title,
		Keyword:    req.Keyword,
		Content:    req.Content,
		CreateDate: req.CreateDate,
	}
	if l, e := DB.Where("Id=?", req.Id).Update(blog); l == 0 || e != nil {
		HtmlMsg(c, MSG_ERR, "修改失败", nil)
		return
	}
	HtmlMsg(c, MSG_SUCC, "修改成功", "")
}

//删除博客
func BlogDeleteGet(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	blog := &model.Blog{Id: id}
	if l, e := DB.Where("Id=?", id).Delete(blog); l == 0 || e != nil {
		HtmlMsg(c, MSG_ERR, "删除失败", "")
		return
	}
	HtmlMsg(c, MSG_SUCC, "删除成功", "/")
}
