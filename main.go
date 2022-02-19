package main

import (
	"github.com/WeiWeiCheng123/URLshortener/config"
	"github.com/WeiWeiCheng123/URLshortener/handler"
)

func main() {
	db_host := config.GetStr("DB_HOST")
	db_port := config.GetStr("DB_PORT")
	db_user := config.GetStr("DB_USERNAME")
	db_name := config.GetStr("DB_NAME")
	db_password := config.GetStr("DB_PASSWORD")
	sslmode := config.GetStr("DB_SSL_MODE")
	redis_host := config.GetStr("REDIS_HOST")
	redis_password := config.GetStr("REDIS_PASSWORD")
	redis_pool := config.GetInt("REDIS_POOL")

	handler.Build(
		db_host, db_port, db_name,
		db_user, db_password, sslmode,
		redis_host, redis_pool, redis_password)
}
