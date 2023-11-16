package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	Client mongo.Client
	Board  *BoardRepository
	List   *ListRepository
	Card   *CardRepository
}

func NewRepositories(DB *mongo.Database) *Repository {
	return &Repository{
		Client: *DB.Client(),
		Board:  newBoardRepository(DB),
		List:   newListRepository(DB),
		Card:   newCardRepository(DB),
	}
}
