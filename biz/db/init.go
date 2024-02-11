package db

import (
	"web_chat/biz/db/mysql"
	"web_chat/biz/db/redis"
)

func Init() {
	mysql.Init()
	redis.Init()
}