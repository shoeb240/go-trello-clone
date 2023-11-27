package service

import (
	"context"

	"github.com/shoeb240/go-trello-clone/model"
	"github.com/shoeb240/go-trello-clone/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	client   *mongo.Client
	userRepo *repository.UserRepository
}

func NewUserService(repository *repository.Repository) *UserService {
	return &UserService{
		client:   &repository.Client,
		userRepo: repository.User,
	}
}

func (s *UserService) Signup(ctx context.Context, userModel model.User) (primitive.ObjectID, error) {
	return s.userRepo.Signup(ctx, userModel)
}
