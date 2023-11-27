package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	Client mongo.Client
	User   *UserRepository
	Board  *BoardRepository
	List   *ListRepository
	Card   *CardRepository
}

func NewRepositories(DB *mongo.Database) *Repository {
	return &Repository{
		Client: *DB.Client(),
		User:   newUserRepository(DB),
		Board:  newBoardRepository(DB),
		List:   newListRepository(DB),
		Card:   newCardRepository(DB),
	}
}
