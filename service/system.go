package service

import (
	"io"
	"os"

	pb "github.com/Camelia-hu/taxin/pb/systempb" // 替换为实际路径
)

type SystemService struct {
	baseDir string
	pb.UnimplementedSystemServiceServer
}

func NewService() *SystemService {
	return &SystemService{baseDir: os.Getenv("BASE_DIR")}
}

func (s *SystemService) SendFile(req *pb.FileRequest, stream pb.SystemService_SendFileServer) error {
	filePath := s.baseDir + "/" + req.FilePath
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 1024*4) // 4KB chunks
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&pb.FileChunk{
			Content: buf[:n],
		}); err != nil {
			return err
		}
	}
	return nil
}
