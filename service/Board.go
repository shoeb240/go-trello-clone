package service

import (
	"context"

	"github.com/shoeb240/go-trello-clone/model"
	"github.com/shoeb240/go-trello-clone/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BoardService struct {
	client    *mongo.Client
	boardRepo *repository.BoardRepository
	cardRepo  *repository.CardRepository
}

func NewBoardService(repository *repository.Repository) *BoardService {
	return &BoardService{
		client:    &repository.Client,
		boardRepo: repository.Board,
		cardRepo:  repository.Card,
	}
}

func (s *BoardService) FullBoard(ctx context.Context, ID string) (model.Board, error) {
	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return model.Board{}, err
	}

	return s.boardRepo.FindByID(ctx, objID)
}

func (s *BoardService) Create(ctx context.Context, boardModel model.Board) (primitive.ObjectID, error) {
	return s.boardRepo.Create(ctx, boardModel)
}

func (s *BoardService) Update(ctx context.Context, boardModel model.Board) (string, error) {
	updateData := bson.M{
		"title": boardModel.Title,
	}
	return s.boardRepo.Update(ctx, updateData)
}

func (s *BoardService) Delete(ctx context.Context, boardObjID primitive.ObjectID) (string, error) {
	client := s.client
	session, err := client.StartSession()
	if err != nil {
		return "", err
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		err := session.StartTransaction()
		if err != nil {
			return err
		}

		_, err = s.boardRepo.Delete(sc, boardObjID)
		if err != nil {
			return err
		}

		_, err = s.cardRepo.DeleteCardsByBoardID(sc, boardObjID)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})

	if err != nil {
		return "", err
	}

	return "deleted", nil
}
