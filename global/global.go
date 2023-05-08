package global

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v6"
	"net/smtp"
)

var Database *sqlx.DB
var Redis *redis.Client
var Router *gin.Engine
var MinioClient *minio.Client
var SMTPAuth smtp.Auth
