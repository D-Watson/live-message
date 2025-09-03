package dbs

import (
	"context"
	"fmt"

	cf "github.com/D-Watson/live-safety/conf"
	"github.com/D-Watson/live-safety/dbs"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"gorm.io/gorm"
)

func InitMongo(ctx context.Context) {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb+srv://zhuanglinlin0321:pegufsOue8Fsdp9E@cluster0.smx7bpm.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	fmt.Println(err)
	fmt.Println(client)
}

var (
	MysqlEngine *gorm.DB
	RedisEngine *redis.Client
)

func InitMysql() error {
	var err error
	MysqlEngine, err = dbs.GetDBClient(cf.GlobalConfig.DB.Mysql)
	if err != nil {
		return err
	}
	return nil
}
func InitRedis(ctx context.Context) error {
	var err error
	RedisEngine, err = dbs.InitRedisCli(ctx)
	if err != nil {
		return err
	}
	return nil
}
