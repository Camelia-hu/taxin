package utils

import (
	"github.com/sony/sonyflake"
	"log"
)

var sf *sonyflake.Sonyflake

func init() {
	sf = sonyflake.NewSonyflake(sonyflake.Settings{})
	if sf == nil {
		log.Fatalf("初始化 sonyflake 失败")
	}
}

func GenerateUserID() (uint64, error) {
	return sf.NextID()
}
