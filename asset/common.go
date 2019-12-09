// Copyright 2014 wudaoren.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package asset

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"math"
	"math/rand"
	"miniBlog/global"
	"miniBlog/model"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

//致命错误检查，title=如果错误输出的标题，e=要检查的错误输出
func CheckFatalError(title string, e error) {
	if e != nil {
		log.Fatal(title, e)
	}
}

//输出消息页面
const (
	MSG_SUCC = "success"
	MSG_ERR  = "error"
)

//输出消息跳转页面
//status=MSG_SUCC 或者 MSG_ERR，msg=提示消息
func HtmlMsg(c *gin.Context, status, msg string, gotoURL interface{}) {
	HtmlData(c, "msg.html", gin.H{
		"Status":  status,
		"Message": msg,
		"GotoURL": gotoURL,
	})
}

//渲染html页面，并增加全局使用的参数
func HtmlData(c *gin.Context, templFile string, data interface{}) {
	sess := UseSession(c)
	glob := &model.System{Id: 1}
	global.DB.Get(glob)
	token := fmt.Sprintf("%X", rand.Int63())
	sess.Set(token, true)
	c.HTML(200, templFile, gin.H{
		"Glob":  glob,
		"Data":  data,
		"Login": sess.Get("login"),
		"Token": token, //创建表单验证token,防止表单重复提交
	})
}

//分页查询
//sess=xorm的查询会话，page=查询页码，limt=每页数量,find=查询列表指针
func PageQuery(sess *xorm.Session, page int, limit int, find interface{}) int {
	newSess := sess.Clone()
	num, _ := newSess.Count()
	maxPage := int(math.Ceil(float64(num) / float64(limit)))
	//当页码超过最大值
	if page > maxPage {
		page = maxPage
	}
	if page > 0 {
		page = page - 1
	}
	e := sess.Limit(limit, page*limit).Find(find)
	if e != nil {
		log.Println("分页查询错误：", e)
	}
	return maxPage
}

//md5加密
func Md5(str string) string {
	data := md5.New()
	data.Write([]byte(str))
	cipherStr := data.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//调试输出
func Debug(args ...interface{}) {
	fmt.Println(args...)
}

//创建分页按钮
//url=跳转的url,pageCurrent=当前页，pageTotal=页面总数
func CreatePageMenu(url string, pageCurrent, pageTotal int) template.HTML {
	var html = `<div class="pageBox">`
	if pageCurrent > 1 {
		html += fmt.Sprintf(`<li class="previous"><a href="%s%d">首页</a>&nbsp;<a href="%s%d">上一页</a></li>`, url, 1, url, pageCurrent-1)
	}
	html += fmt.Sprintf(`<li class="info">%d/%d</li>`, pageCurrent, pageTotal)
	if pageCurrent < pageTotal {
		html += fmt.Sprintf(`<li class="next"><a href="%s%d">下一页</a>&nbsp;<a href="%s%d">尾页</a></li>`, url, pageCurrent+1, url, pageTotal)
	}
	html += "</div>"
	return template.HTML(html)
}
