package student

import (
	"context"

	"example.com/pz2-grpc/gen/studentpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	studentpb.UnimplementedStudentServiceServer
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Ping(_ context.Context, req *studentpb.PingRequest) (*studentpb.PingResponse, error) {
	msg := req.GetMessage()
	if msg == "" {
		msg = "ping"
	}
	return &studentpb.PingResponse{
		Message: "Server received: " + msg,
	}, nil
}

func (s *Service) GetStudentByID(_ context.Context, req *studentpb.GetStudentRequest) (*studentpb.GetStudentResponse, error) {
	id := req.GetId()
	if id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid student id")
	}
	st, err := s.repo.GetByID(id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "student not found")
	}
	return &studentpb.GetStudentResponse{Student: st}, nil
}

func (s *Service) ListStudents(_ context.Context, _ *emptypb.Empty) (*studentpb.ListStudentsResponse, error) {
	return &studentpb.ListStudentsResponse{
		Students: s.repo.List(),
	}, nil
}
