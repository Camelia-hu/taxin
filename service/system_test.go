package service

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/Camelia-hu/taxin/pb/systempb" // 替换为实际路径
)

const bufferSize = 1024 * 1024 // 1MB内存缓冲区

// 测试SendFile服务方法
func TestSendFile(t *testing.T) {
	os.Setenv("BASE_DIR", os.TempDir()) //  设置 baseDir

	testContent := []byte("This is a test file content for gRPC streaming")
	tmpFile, err := os.CreateTemp(os.Getenv("BASE_DIR"), "testfile-*.txt") //  创建在 baseDir 中
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(testContent)
	require.NoError(t, err)
	tmpFile.Close()

	lis := bufconn.Listen(bufferSize)
	grpcServer := grpc.NewServer()

	systemService := NewService()
	pb.RegisterSystemServiceServer(grpcServer, systemService)

	go grpcServer.Serve(lis)
	defer grpcServer.GracefulStop()

	ctx := context.Background()

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewSystemServiceClient(conn)

	stream, err := client.SendFile(ctx, &pb.FileRequest{
		FilePath: filepath.Base(tmpFile.Name()), //  只传相对路径
	})
	require.NoError(t, err)

	var receivedData []byte
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		receivedData = append(receivedData, chunk.Content...)
	}

	assert.Equal(t, testContent, receivedData)
}
