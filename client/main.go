package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb_user "github.com/Camelia-hu/taxin/pb/userpb"
)

func main() {
	ctx := context.Background()
	tr := otel.Tracer("all-services-client")
	ctx, span := tr.Start(ctx, "ClientCallAllServices")
	defer span.End()

	// 连接到 gRPC 服务器
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		span.SetStatus(codes.Error, "Failed to connect to server")
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// 测试 UserService
	testUserService(ctx, conn)

	//SystemService用单元测试
}

func testUserService(ctx context.Context, conn *grpc.ClientConn) {
	tr := otel.Tracer("user-service-client")
	ctx, span := tr.Start(ctx, "TestUserService")
	defer span.End()

	// 创建 UserService 客户端
	client := pb_user.NewUserServiceClient(conn)

	// 测试注册
	registerReq := &pb_user.RegisterRequest{
		Password: "testpassword",
		Like:     "swimming",
		Username: "testuser2",
	}
	registerResp, err := client.Register(ctx, registerReq)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to register")
		log.Fatalf("Failed to register: %v", err)
	}
	span.SetStatus(codes.Error, "register ok")
	fmt.Printf("Register User ID: %d\n", registerResp.UserId)

	// 测试登录
	loginReq := &pb_user.LoginRequest{
		Username: "testuser2",
		Password: "testpassword",
	}
	loginResp, err := client.Login(ctx, loginReq)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to login")
		log.Fatalf("Failed to login: %v", err)
	}
	fmt.Printf("Login Access Token: %s\n", loginResp.AccessToken)

	// 测试获取用户信息
	ctxWithToken := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+loginResp.AccessToken)
	userInfoReq := &pb_user.GetUserInfoRequest{UserId: "568176012081253442"}
	userInfoResp, err := client.GetUserInfo(ctxWithToken, userInfoReq)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to get user info")
		log.Fatalf("Failed to get user info: %v", err)
	}
	fmt.Printf("User ID: %s\n", userInfoResp.User.UserId)
	fmt.Printf("Likes: %v\n", userInfoResp.User.Like)
	fmt.Printf("Like Embedding: %v\n", userInfoResp.User.LikeEmbedding)
	fmt.Printf("Create At: %s\n", userInfoResp.User.CreateAt)
	fmt.Printf("Update At: %s\n", userInfoResp.User.UpdateAt)
	//fmt.Printf("Username %s\n", userInfoResp.User.)

	// 添加自定义标签和事件
	span.SetAttributes(attribute.String("user_id", strconv.FormatInt(registerResp.UserId, 10)))
	span.AddEvent("User service tests completed successfully")
}
