package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"kayak-backend/api"
	"kayak-backend/docs"
	"kayak-backend/global"
	"log"
	"os"
)

func InitSql(Addr string, Port int, User string, Password string, Database string) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Addr,
		Port,
		User,
		Password,
		Database)

	fmt.Println(psqlInfo)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	global.Database = db
}

func InitRedis(Addr string, Port int, Password string) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Addr, Port),
		Password: Password,
		DB:       0,
	})
	global.Redis = rdb
}

func InitLog(Path string) {
	f, err := os.OpenFile(Path+"log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	ginLog, err := os.OpenFile(Path+"gin_log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	gin.DefaultWriter = io.MultiWriter(ginLog, os.Stdout)
}

func InitMinio(Addr string, Port int, AccessKey string, SecretKey string, UseSSL bool) {
	client, err := minio.New(fmt.Sprintf("%s:%d", Addr, Port), AccessKey, SecretKey, UseSSL)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err != nil {
		log.Fatalln("创建 MinIO 客户端失败", err)
		return
	}
	global.MinioClient = client
	log.Printf("创建 MinIO 客户端成功")

}

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	InitSql(viper.GetString("PostgresHost"), viper.GetInt("PostgresPort"),
		viper.GetString("PostgresUsername"), viper.GetString("PostgresPassword"),
		viper.GetString("PostgresDatabase"))
	InitRedis(viper.GetString("RedisHost"), viper.GetInt("RedisPort"),
		viper.GetString("RedisPassword"))
	InitLog(viper.GetString("LogPath"))
	InitMinio(viper.GetString("MinioHost"), viper.GetInt("MinioPort"),
		viper.GetString("MinioAccessKey"), viper.GetString("MinioSecretKey"),
		false)
	docs.SwaggerInfo.BasePath = viper.GetString("DocsPath")
}

// @title Kayak Backend API
// @version 0.0.2
// @license null
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Token
func main() {
	LoadConfig()
	global.Router = gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, global.TokenHeader)
	corsConfig.ExposeHeaders = append(corsConfig.ExposeHeaders, "Date")

	global.Router.Use(cors.New(corsConfig))
	global.Router.Use(global.Authenticate)
	global.Router.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("doc.json")))
	api.InitRoute()
	err := global.Router.Run("0.0.0.0:9000")
	if err != nil {
		return
	}
}
