package main

import (
	"context"
	"fmt"
	"github.com/Camelia-hu/taxin/pb/systempb"
	"github.com/Camelia-hu/taxin/pb/userpb"
	"github.com/Camelia-hu/taxin/service"
	"github.com/Camelia-hu/taxin/utils"
	"github.com/joho/godotenv"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
	"os"
)

func newDatabase() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
}

func newModel() *arkruntime.Client {
	return arkruntime.NewClientWithApiKey(os.Getenv("ARK_API_KEY"))
}

func newGRPCServer() *grpc.Server {
	return grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
}

// 创建 Jaeger 追踪器
func newTracerProvider(lc fx.Lifecycle) (func(context.Context) error, error) {
	ctx := context.Background()
	tp, err := utils.InitTracer(ctx, "taxin")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	// 设置全局的追踪传播器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// 注册生命周期钩子
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			fmt.Println("Shutting down tracer provider")
			return tp.Shutdown(ctx)
		},
	})

	return tp.Shutdown, nil
}

func registerServices(
	s *grpc.Server,
	userSvc *service.UserService,
	systemSvc *service.SystemService,
) {
	userpb.RegisterUserServiceServer(s, userSvc)
	systempb.RegisterSystemServiceServer(s, systemSvc)
}

func startServer(lc fx.Lifecycle, s *grpc.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, _ := net.Listen("tcp", ":50051")
			go s.Serve(lis)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.GracefulStop()
			return nil
		},
	})
}

func main() {
	godotenv.Load()
	fx.New(
		fx.Provide(
			newDatabase,
			newModel,
			newGRPCServer,
			newTracerProvider,
			service.NewUserService,
			service.NewService,
		),
		fx.Invoke(
			registerServices,
			startServer,
		),
	).Run()
}

//func main() {
//	godotenv.Load()
//	someInit.Initdb()
//	someInit.ModelInit()
//	lis, err := net.Listen("tcp", ":50051")
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//
//	grpcServer := grpc.NewServer()
//	userpb.RegisterUserServiceServer(grpcServer, &service.UserService{})
//
//	log.Println("🚀 gRPC 服务已启动，监听端口 :50051")
//	if err := grpcServer.Serve(lis); err != nil {
//		log.Fatalf("failed to serve: %v", err)
//	}
//}

//package main
//
//import (
//"log"
//"net"
//
//"google.golang.org/grpc"
//
//"github.com/Camelia-hu/taxin/pb/userpb"
//"github.com/Camelia-hu/taxin/service"
//)
//
//func main() {
//	lis, err := net.Listen("tcp", ":50051")
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//
//	grpcServer := grpc.NewServer()
//	userpb.RegisterUserServiceServer(grpcServer, &service.UserService{})
//
//	log.Println("🚀 gRPC 服务已启动，监听端口 :50051")
//	if err := grpcServer.Serve(lis); err != nil {
//		log.Fatalf("failed to serve: %v", err)
//	}
//}
