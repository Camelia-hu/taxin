package service

import (
	"context"
	"errors"
	"github.com/Camelia-hu/taxin/model"
	"github.com/Camelia-hu/taxin/pb/userpb"
	"github.com/Camelia-hu/taxin/utils"
	"github.com/pgvector/pgvector-go"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
	"strconv"
	"time"
)

const (
	PassWordNil  = "密码不能为空"
	UserIdExist  = "用户名已存在，请重新登录"
	DbFindErr    = "数据库查询失败"
	DbCreateErr  = "数据库Create操作失败"
	UserErr      = "用户操作失败"
	UserNotFound = "用户不存在"
	PassWordErr  = "密码错误"
)

// 你自己的结构体，实现 gRPC 接口
type UserService struct {
	client                                *arkruntime.Client
	db                                    *gorm.DB
	userpb.UnimplementedUserServiceServer // 推荐嵌入，避免版本不兼容
}

func NewUserService(db *gorm.DB, client *arkruntime.Client) *UserService {
	return &UserService{db: db, client: client}
}

// 1. 实现 Register 方法
func (s *UserService) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	// ✅ 这里写你的注册业务逻辑：校验输入、查询数据库、插入新用户
	//按道理来说应该在router层用结构体标签来参数校验，这里为了体现功能完整性在此进行参数校验
	tr := otel.Tracer("user-service")
	_, span := tr.Start(ctx, "Register")
	defer span.End()

	if req.Password == "" {
		span.SetStatus(codes.Error, "password err")
		return nil, errors.New(PassWordNil)
	}
	var user model.UserModel
	err := s.db.Where("username = ?", req.Username).First(&user).Error
	if err == nil {
		return &userpb.RegisterResponse{Message: UserIdExist}, err
	}
	span.SetAttributes(attribute.String("user_id", user.UserId))
	var UserId uint64
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 没找到，说明用户还不存在，可以注册
		UserId, err = utils.GenerateUserID()
		like, err := utils.Embedding(ctx, req.Like, s.client)
		embeddingVector := like.Data[0].Embedding
		embedding := pgvector.NewVector(embeddingVector[:128])
		passWord, err := utils.HashPassword(req.Password)
		if err != nil {
			return &userpb.RegisterResponse{Message: UserErr}, err
		}
		err = s.db.Create(&model.UserModel{
			UserId:        strconv.FormatUint(UserId, 10),
			Password:      passWord,
			Like:          req.Like,
			LikeEmbedding: embedding,
			CreateAt:      time.Now(),
			UpdateAt:      time.Now(),
			UserName:      req.Username,
		}).Error
		if err != nil {
			return &userpb.RegisterResponse{Message: DbCreateErr}, err
		}

	} else {
		return &userpb.RegisterResponse{Message: DbFindErr}, err
	}
	accessToken, refreshToken, err := utils.DeliverTokenByRPCService(user.UserId)
	if err != nil {
		return &userpb.RegisterResponse{Message: UserErr}, err
	}
	// 假设注册成功：
	return &userpb.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserId:       int64(UserId),
		Message:      "注册成功",
	}, nil
}

// 2. 实现 Login 方法
func (s *UserService) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	// ✅ 这里写你的登录业务逻辑：校验账号密码，生成 token
	tr := otel.Tracer("user-service")
	_, span := tr.Start(ctx, "Login")
	defer span.End()

	var user model.UserModel
	err := s.db.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &userpb.LoginResponse{Message: UserNotFound}, err
		} else {
			return &userpb.LoginResponse{Message: DbFindErr}, err
		}
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		return &userpb.LoginResponse{Message: PassWordErr}, errors.New(PassWordErr)
	}
	accessToken, refreshToken, err := utils.DeliverTokenByRPCService(user.UserId)
	if err != nil {
		return &userpb.LoginResponse{Message: UserErr}, err
	}
	return &userpb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "登录成功",
	}, nil
}

// 3. 实现 GetUserInfo 方法
func (s *UserService) GetUserInfo(ctx context.Context, req *userpb.GetUserInfoRequest) (*userpb.GetUserInfoResponse, error) {
	// ✅ 这里写你的获取信息逻辑：根据 userId 查询数据库返回信息
	tr := otel.Tracer("user-service")
	_, span := tr.Start(ctx, "GetUserInfo")
	defer span.End()

	var user model.UserModel
	err := s.db.Where("user_id = ?", req.UserId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &userpb.GetUserInfoResponse{User: nil}, errors.New(UserNotFound + err.Error())
		} else {
			return &userpb.GetUserInfoResponse{User: nil}, errors.New(DbFindErr + err.Error())
		}
	}
	span.SetAttributes(attribute.String("user_id", user.UserId))
	return &userpb.GetUserInfoResponse{
		User: &userpb.User{
			Id:       user.Id,
			UserId:   user.UserId,
			Like:     user.Like,
			CreateAt: user.CreateAt.String(),
			UpdateAt: user.UpdateAt.String(),
		},
	}, nil
}
