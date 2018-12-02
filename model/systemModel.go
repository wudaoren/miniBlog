// Copyright 2014 wudaoren.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package model

type System struct {
	Id         int    `xorm:"int(11) pk autoincr"` //主键
	BlogName   string `xorm:"varchar(255)"`        //博客名称
	Keyword    string `xorm:"varchar(255)"`        //关键词
	Username   string `xorm:"varchar(32)"`         //管理员账号
	Password   string `xorm:"varchar(32)"`         //管理员登录密码
	UpdateTime string `xorm:"datetime updated"`    //修改时间
}
