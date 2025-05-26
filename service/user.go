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

// 添加client db与fx框架配合
type UserService struct {
	client                                *arkruntime.Client
	db                                    *gorm.DB
	userpb.UnimplementedUserServiceServer // 推荐嵌入，避免版本不兼容
}

func NewUserService(db *gorm.DB, client *arkruntime.Client) *UserService {
	return &UserService{db: db, client: client}
}

// 实现 Register 方法
func (s *UserService) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	tr := otel.Tracer("taxin")
	ctx, span := tr.Start(ctx, "UserService.Register")
	defer span.End()

	//按道理来说应该在router层用结构体标签来参数校验，这里为了体现功能完整性在此进行参数校验
	if req.Password == "" {
		span.SetStatus(codes.Error, "password err")
		return nil, errors.New(PassWordNil)
	}

	var user model.UserModel
	err := s.db.Where("username = ?", req.Username).First(&user).Error
	if err == nil {
		// 用户已存在，注册幂等处理：直接返回 Token
		accessToken, refreshToken, err := utils.DeliverTokenByRPCService(user.UserId)
		if err != nil {
			return &userpb.RegisterResponse{Message: UserErr}, err
		}
		span.SetAttributes(attribute.String("user_id", user.UserId))
		user_id, _ := strconv.Atoi(user.UserId)
		return &userpb.RegisterResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			UserId:       int64(user_id),
			Message:      "用户已存在，返回Token",
		}, nil
	}

	UserId, _ := utils.GenerateUserID()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 没找到，说明用户还不存在，可以注册
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
			span.SetStatus(codes.Error, DbCreateErr)
			return &userpb.RegisterResponse{Message: DbCreateErr}, err
		}
	} else {
		span.SetStatus(codes.Error, DbFindErr)
		return &userpb.RegisterResponse{Message: DbFindErr}, err
	}

	accessToken, refreshToken, err := utils.DeliverTokenByRPCService(user.UserId)
	if err != nil {
		span.SetStatus(codes.Error, UserErr)
		return &userpb.RegisterResponse{Message: UserErr}, err
	}
	// 假设注册成功：
	span.SetAttributes(
		attribute.String("user_id", strconv.FormatUint(UserId, 10)),
		attribute.String("username", req.Username),
		attribute.String("event", "user_registered"),
	)

	return &userpb.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserId:       int64(UserId),
		Message:      "注册成功",
	}, nil
}

// 实现 Login 方法
func (s *UserService) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	// ✅ 这里写你的登录业务逻辑：校验账号密码，生成 token
	tr := otel.Tracer("user-service")
	ctx, span := tr.Start(ctx, "Login")
	defer span.End()

	var user model.UserModel
	err := s.db.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		span.SetStatus(codes.Error, UserNotFound)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &userpb.LoginResponse{Message: UserNotFound}, err
		} else {
			span.SetStatus(codes.Error, DbFindErr)
			return &userpb.LoginResponse{Message: DbFindErr}, err
		}
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		span.SetStatus(codes.Error, PassWordErr)
		return &userpb.LoginResponse{Message: PassWordErr}, errors.New(PassWordErr)
	}
	accessToken, refreshToken, err := utils.DeliverTokenByRPCService(user.UserId)
	if err != nil {
		span.SetStatus(codes.Error, UserErr)
		return &userpb.LoginResponse{Message: UserErr}, err
	}
	span.SetAttributes(
		attribute.String("user_id", user.UserId),
		attribute.String("username", req.Username),
		attribute.String("event", "user_login"),
	)
	return &userpb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "登录成功",
	}, nil
}

// 实现 GetUserInfo 方法
func (s *UserService) GetUserInfo(ctx context.Context, req *userpb.GetUserInfoRequest) (*userpb.GetUserInfoResponse, error) {
	// ✅ 这里写你的获取信息逻辑：根据 userId 查询数据库返回信息
	tr := otel.Tracer("user-service")
	ctx, span := tr.Start(ctx, "GetUserInfo")
	defer span.End()

	var user model.UserModel
	err := s.db.Where("user_id = ?", req.UserId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Error, UserNotFound)
			return &userpb.GetUserInfoResponse{User: nil}, errors.New(UserNotFound + err.Error())
		} else {
			span.SetStatus(codes.Error, DbFindErr)
			return &userpb.GetUserInfoResponse{User: nil}, errors.New(DbFindErr + err.Error())
		}
	}
	span.SetAttributes(
		attribute.String("user_id", user.UserId),
		attribute.String("username", user.UserName),
		attribute.String("event", "user_registered"),
	)
	return &userpb.GetUserInfoResponse{
		User: &userpb.User{
			Id:            user.Id,
			UserId:        user.UserId,
			Like:          user.Like,
			LikeEmbedding: user.LikeEmbedding.String(),
			CreateAt:      user.CreateAt.String(),
			UpdateAt:      user.UpdateAt.String(),
		},
	}, nil
}
