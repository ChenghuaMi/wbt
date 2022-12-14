package svc

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"wbt/server/internal/config"
)

type ServiceContext struct {
	Config config.Config
	Redis *LoadRedis
	GoDb *gorm.DB

}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Redis: NewLoadRedis(c),
		GoDb: NewGoDb(c),
	}
}

type LoadRedis struct {
	C *redis.Client
	Ctx context.Context
}
func NewLoadRedis(c config.Config) *LoadRedis {
	return &LoadRedis{
		C:   newRedis(c),
		Ctx: context.Background(),
	}
}

func newRedis(c config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:               c.CacheRedis[0].Host,

		Password:           c.CacheRedis[0].Pass,
		DB:                 0,
	})
	return client
}

func NewGoDb(c config.Config) *gorm.DB {
	db,err := gorm.Open(mysql.Open(c.Mysql.DataUrl),&gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}