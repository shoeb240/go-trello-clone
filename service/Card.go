package service

import (
	"context"

	"github.com/shoeb240/go-trello-clone/model"
	"github.com/shoeb240/go-trello-clone/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	updateData := bson.M{
		"title":       cardModel.Title,
		"description": cardModel.Description,
		"position":    cardModel.Position,
	}
	return s.cardRepo.Update(ctx, updateData)
}

func (s *CardService) MoveCard(ctx context.Context, cardMoveReq repository.CardMoveReq) (string, error) {
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

		cardDeleteFromList := repository.CardDeleteFromList{
			CardID:  cardMoveReq.CardID,
			BoardID: cardMoveReq.BoardID,
			ListID:  cardMoveReq.FromListID,
		}

		updateData := bson.M{
			"list_id":  cardMoveReq.ToListID,
			"position": cardMoveReq.ToPosition,
		}
		_, err = s.cardRepo.Update(sc, updateData)
		if err != nil {
			return err
		}

		cardModel, err := s.cardRepo.FindByID(sc, cardMoveReq.CardID)
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}

		err = s.boardRepo.AddCardToList(sc, cardModel)
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}

		err = s.boardRepo.RemoveCardFromList(sc, cardDeleteFromList)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})

	if err != nil {
		return "", err
	}

	return "moved successfully", nil
}

func (s *CardService) DeleteCard(ctx context.Context, cardObjID primitive.ObjectID) (string, error) {
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

		cardModel, err := s.cardRepo.FindByID(sc, cardObjID)
		if err != nil {
			session.AbortTransaction(sc)
			return err
		}

		cardDeleteFromList := repository.CardDeleteFromList{
			CardID:  cardModel.ID,
			BoardID: cardModel.BoardID,
			ListID:  cardModel.ListID,
		}

		_, err = s.cardRepo.Delete(sc, cardObjID)
		if err != nil {
			return err
		}

		err = s.boardRepo.RemoveCardFromList(sc, cardDeleteFromList)
		if err != nil {
			return err
		}

		return session.CommitTransaction(sc)
	})

	if err != nil {
		return "", err
	}

	return "deleted successfully", nil
}
