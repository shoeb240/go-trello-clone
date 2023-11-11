package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	Board *BoardRepository
	List  *ListRepository
}

func NewRepositories(DB *mongo.Database) *Repository {
	return &Repository{
		Board: newBoardRepository(DB),
		List:  newListRepository(DB),
	}
}
