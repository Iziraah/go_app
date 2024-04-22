package usergrpc

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/robotomize/gb-golang/homework/03-03-umanager/internal/database"
	"gitlab.com/robotomize/gb-golang/homework/03-03-umanager/pkg/pb"
)

var _ pb.UserServiceServer = (*Handler)(nil)

func New(usersRepository usersRepository, timeout time.Duration) *Handler {
	return &Handler{usersRepository: usersRepository, timeout: timeout}
}

type Handler struct {
	pb.UnimplementedUserServiceServer
	usersRepository usersRepository
	timeout         time.Duration
}

func (h Handler) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	var (
		id  uuid.UUID
		err error
	)
	if in.Id == "" {
		id = uuid.New()
	} else {
		id, err = uuid.Parse(in.Id)
	}

	if err != nil {
		return &pb.Empty{}, err
	}

	req := database.CreateUserReq{
		ID:       id,
		Username: in.Username,
		Password: in.Password,
	}
	_, err = h.usersRepository.Create(ctx, req)
	return &pb.Empty{}, err
}

func (h Handler) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, err
	}

	user, err := h.usersRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb.User{
		Id:        user.ID.String(),
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}

func (h Handler) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return &pb.Empty{}, err
	}
	req := database.CreateUserReq{
		ID:       id,
		Username: in.Username,
		Password: in.Password,
	}
	_, err = h.usersRepository.Create(ctx, req)
	return &pb.Empty{}, err
}

func (h Handler) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return &pb.Empty{}, err
	}

	err = h.usersRepository.DeleteByUserID(ctx, id)
	return &pb.Empty{}, err
}

func (h Handler) ListUsers(ctx context.Context, in *pb.Empty) (*pb.ListUsersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	users, err := h.usersRepository.FindAll(ctx)
	if err != nil {
		return &pb.ListUsersResponse{}, err
	}

	res := make([]*pb.User, len(users))
	for i, u := range users {
		res[i] = &pb.User{
			Id:        u.ID.String(),
			Username:  u.Username,
			Password:  u.Password,
			CreatedAt: u.CreatedAt.String(),
			UpdatedAt: u.UpdatedAt.String(),
		}
	}
	return &pb.ListUsersResponse{Users: res}, err
}