package global

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

var Database *sqlx.DB
var Redis *redis.Client
var Router *gin.Engine
