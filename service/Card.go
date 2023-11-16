package service

import (
	"context"

	"github.com/shoeb240/go-trello-clone/model"
	"github.com/shoeb240/go-trello-clone/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type CardService struct {
	client    *mongo.Client
	cardRepo  *repository.CardRepository
	boardRepo *repository.BoardRepository
}

func NewCardService(repository *repository.Repository) *CardService {
	return &CardService{
		client:    &repository.Client,
		cardRepo:  repository.Card,
		boardRepo: repository.Board,
	}
}

func (s *CardService) Create(ctx context.Context, cardModel model.Card) (model.Card, error) {
	client := s.client
	session, err := client.StartSession()
	if err != nil {
		return model.Card{}, err
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		err := session.StartTransaction()
		if err != nil {
			return err
		}

		cardModelID, err := s.cardRepo.Create(sc, cardModel)
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}
		cardModel.ID = cardModelID

		err = s.boardRepo.AddCardToList(sc, cardModel)
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}

		return session.CommitTransaction(sc)
	})

	if err != nil {
		return cardModel, err
	}

	return cardModel, nil
}

func (s *CardService) Update(ctx context.Context, cardModel model.Card) (string, error) {
	return s.cardRepo.Update(ctx, cardModel)
}
