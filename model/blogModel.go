// Copyright 2014 wudaoren.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package model

type Blog struct {
	Id         int    `xorm:"int(11) pk autoincr"` //主键
	Title      string `xorm:"varchar(255)"`        //标题
	Keyword    string `xorm:"varchar(255)"`        //关键词
	Content    string `xorm:"text"`                //博客内容
	CreateDate string `xorm:"date"`                //发布日期（可以修改）
	CreateTime string `xorm:"datetime created"`    //创建时间
	UpdateTime string `xorm:"datetime updated"`    //修改时间

	HTMLContent interface{} `xorm:"-"` //页面输出html代码，该字段不保存到数据库
}
