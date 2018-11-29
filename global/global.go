// Copyright 2014 wudaoren.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package global

import (
	"github.com/go-xorm/xorm"
)

//配置参数
var Config struct {
	Host  string
	Mysql struct {
		Host     string
		Username string
		Password string
		Database string
	}
}

//数据库引擎
var DB *xorm.Engine

//上传文件保存目录
var UPLOAD_DIR = "/static/upload/"
