// Copyright 2014 wudaoren.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package asset

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

//这是将session数据保存在缓存里的一个简单session包
const (
	SESSION_NAME = "WUDAOREN"
)

var sessEngine *SesionEngine

//每隔30秒清理一下过期session
func init() {
	sessEngine = new(SesionEngine)
	sessEngine.data = new(sync.Map)
	sessEngine.maxLifeTime = 60 * 30
	go func() {
		for {
			time.Sleep(time.Second * 30)
			sessEngine.gc()
		}
	}()
}

//使用session时使用
func UseSession(c *gin.Context) *MemSession {
	sessionId, _ := c.Cookie(SESSION_NAME)
	if sessionId == "" {
		randId := rand.Int63()
		key := fmt.Sprintf("%d-%x", time.Now().Unix(), randId)
		md5Ctx := md5.New()
		md5Ctx.Write([]byte(key))
		cipherStr := md5Ctx.Sum(nil)
		sessionId = hex.EncodeToString(cipherStr)
		c.SetCookie(SESSION_NAME, url.QueryEscape(sessionId), 3600*24*365, "/", "", false, false)

	}
	return sessEngine.getSession(sessionId)
}

//默认session engine
type SesionEngine struct {
	maxLifeTime int64 //最长保存时间（秒）
	visited     int   //访问次数
	data        *sync.Map
}

type sessionData struct {
	sess *MemSession
	time int64
}

func (this *SesionEngine) getSession(sessionId string) *MemSession {
	data := new(sessionData)
	if res, ok := this.data.Load(sessionId); ok {
		data = res.(*sessionData)
	}
	overTime := this.isTimeout(data)
	if overTime > this.maxLifeTime {
		data.sess = new(MemSession)
		this.data.Store(sessionId, data)
	} else if overTime > this.maxLifeTime/2 && overTime < this.maxLifeTime {
		this.data.Store(sessionId, data)
	}
	return data.sess
}

//
func (this *SesionEngine) isTimeout(data *sessionData) int64 {
	nowTimestemp := time.Now().Unix()
	overTime := nowTimestemp - data.time //小于0表示正常，大于0表示超时
	data.time = nowTimestemp
	return overTime
}

//
func (this *SesionEngine) gc() {
	this.data.Range(func(k, v interface{}) bool {
		if data, ok := v.(*sessionData); ok && this.isTimeout(data) >= 0 {
			this.data.Delete(k)
		}
		return true
	})
}

type MemSession struct {
	sync.Map
}

func (this *MemSession) Set(key string, value interface{}) {
	this.Store(key, value)
}

func (this *MemSession) Get(key string) interface{} {
	res, _ := this.Load(key)
	return res
}

func (this *MemSession) Del(key string) {
	this.Delete(key)
}

func (this *MemSession) Clear() {
	this.Map = sync.Map{}
}
