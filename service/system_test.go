package service

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net"
	"os"
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
	//============ 测试准备 ============
	testContent := []byte("This is a test file content for gRPC streaming")
	tmpFile, err := os.CreateTemp("", "testfile-*.txt")
	require.NoError(t, err, "创建临时文件失败")
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(testContent)
	require.NoError(t, err, "写入测试内容失败")
	tmpFile.Close()

	//============ 服务端配置 ============
	lis := bufconn.Listen(bufferSize)
	grpcServer := grpc.NewServer()

	systemService := NewService()
	pb.RegisterSystemServiceServer(grpcServer, systemService)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Logf("服务端异常退出: %v", err)
		}
	}()
	defer grpcServer.GracefulStop()

	//============ 客户端配置 ============
	ctx := context.Background()

	// 使用新的客户端创建方式
	conn, err := grpc.NewClient(
		"passthrough:///bufnet", // 关键修改点
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "创建客户端失败")
	defer conn.Close()

	client := pb.NewSystemServiceClient(conn)

	//============ 测试执行 ============
	stream, err := client.SendFile(ctx, &pb.FileRequest{
		FilePath: tmpFile.Name(),
	})
	require.NoError(t, err, "创建流失败")

	//============ 结果验证 ============
	var receivedData []byte
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err, "接收数据块失败")
		receivedData = append(receivedData, chunk.Content...)
	}

	assert.Equal(t, testContent, receivedData, "接收内容与源文件不一致")
}
