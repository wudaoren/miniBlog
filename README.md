# 说明：
	miniBlog是一个非常简易、没有经过精心设计的博客系统，三方包只使用了gin和xorm，适合初学go web编程的爱好者理解go web开发的过程，
	欢迎大家多提意见和建议。
# 目录结构
	|---aseet				常用工具
	|---controller			控制器目录
	|---global				存放全局变量
	|---model				存放数据模型
	|---static				静态文件目录
		|---editor			富文本编辑器
		|---upload			文件上传目录
	|---views				html模板
	|---vendor				依赖包
	main.go					主程序
	config.json				配置文件
	
# 功能
	1、首页
	2、登录
	3、设置网站信息
	4、添加博客
	5、修改博客
	6、删除博客
	7、博客信息查看
	8、退出登录
	
# 使用方法
## 1、配置
	打开confi.json
```
	{
		"Host":":8000",//服务器监听的地址和端口
		"Mysql":{`
			"Host":"127.0.0.1:3306",//数据库服务器地址
			"Username":"root",		//数据库账号
			"Password":"",		//数据库密码
			"Database":"miniBlog"	//数据库
		}
	}
```
## 2、安装
	终端执行：miniBlog install
	默认账号和密码都是admin
## 3、启动
	终端执行：miniBlog
	
	