package service

import (
	"context"

	"github.com/shoeb240/go-trello-clone/model"
	"github.com/shoeb240/go-trello-clone/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardService struct {
	repository *repository.BoardRepository
}

func NewBoardService(repository *repository.Repository) *BoardService {
	return &BoardService{
		repository: repository.Board,
	}
}

func (s *BoardService) FullBoard(ctx context.Context, ID string) (model.Board, error) {
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return model.Board{}, err
	}

	return s.repository.FindByID(ctx, objID)
}

func (s *BoardService) Create(ctx context.Context, boardModel model.Board) (primitive.ObjectID, error) {
	return s.repository.Create(ctx, boardModel)
}

func (s *BoardService) Update(ctx context.Context, boardModel model.Board) (string, error) {
	updateData := bson.M{
		"title": boardModel.Title,
	}
	return s.repository.Update(ctx, updateData)
}
