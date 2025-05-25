package model

import (
	"database/sql/driver"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pgvector/pgvector-go"
	"strings"
	"time"
)

type UserModel struct {
	Id            int64           `gorm:"column:id"`
	UserId        string          `gorm:"column:user_id"`
	UserName      string          `gorm:"column:username"`
	Password      string          `gorm:"column:password"`
	Like          string          `gorm:"column:like"`
	LikeEmbedding pgvector.Vector `gorm:"column:like_embedding;type:vector(128)"`
	CreateAt      time.Time       `gorm:"column:create_at"`
	UpdateAt      time.Time       `gorm:"column:update_at"`
}

func (UserModel) TableName() string {
	return "user"
}

type Vector []float64

// 实现 driver.Valuer 接口：告诉 GORM 如何将它转换成数据库能接受的格式
func (v Vector) Value() (driver.Value, error) {
	strs := make([]string, len(v))
	for i, val := range v {
		strs[i] = fmt.Sprintf("%f", val)
	}
	return fmt.Sprintf("'%s'", "{"+strings.Join(strs, ",")+"}"), nil
}

type MyClaims struct {
	Id string `json:"user_id"`
	jwt.RegisteredClaims
}
