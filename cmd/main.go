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
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
	"net/http"
	_ "net/http/pprof"
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

// UnaryServerInterceptor
// 创建 Jaeger 追踪器
func newTracerProvider() (*trace.TracerProvider, error) {
	ctx := context.Background()
	tp, err := utils.Inittracer(ctx, "taxin")
	if err != nil {
		return nil, fmt.Errorf("failed to init tracer: %w", err)
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

func registerTracerShutdown(lc fx.Lifecycle, tp *trace.TracerProvider) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			fmt.Println("shutting down tracer")
			return tp.Shutdown(ctx)
		},
	})
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
			lis, err := net.Listen("tcp", ":50051")
			if err != nil {
				return err
			}
			go s.Serve(lis)

			// 启动 pprof 服务
			go func() {
				fmt.Println("Starting pprof server on :6060")
				err := http.ListenAndServe(":6060", nil)
				if err != nil {
					fmt.Printf("Failed to start pprof server: %v\n", err)
				}
			}()

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
			registerTracerShutdown,
			registerServices,
			startServer,
		),
	).Run()
}
