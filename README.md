# Trello Board Clone

This project is a Golang and MongoDB-based clone of the popular "Trello Board" project management software. It offers a set of secure and authenticated endpoints for user registration, login, and seamless management of boards, lists, and cards. With a user-friendly interface and robust functionality, this clone provides teams with a familiar and efficient platform to collaborate, organize tasks, and monitor project progress in a Trello-style environment.

## Installation

### Prerequisites

- Go 
- MongoDB 

## User Endpoints

### 1. Signup

### Endpoint: POST /signup

Request:
```
{
  "email": "user@example.com",
  "password": "userpassword"
}
```

Response:
```
{
  "id": "user_id",
  "email": "user@example.com"
}
```

### 2. Login

### Endpoint: POST /login

Request:

```
{
  "email": "user@example.com",
  "password": "userpassword"
}
```

Response:
```
{
  "token": "jwt_token"
}
```


## Board Endpoints

Note: All board endpoints require authentication.
### 1. Get Board

### Endpoint: GET /board/:boardID

Response:
```
{
  "id": "board_id",
  "user_id": "user_id",
  "title": "Board Title",
  "lists": [
    {
      "id": "list_id",
      "title": "List Title",
      "position": 1,
      "cards": ["card_id1", "card_id2"],
      "card_details": [
        {
          "id": "card_id1",
          "title": "Card Title",
          "description": "Card Description",
          "position": 1
        },
        {
          "id": "card_id2",
          "title": "Card Title",
          "description": "Card Description",
          "position": 2
        }
      ]
    }
  ]
}
```

### 2. Create Board

### Endpoint: POST /board

Request:
```
{
  "title": "Board Title"
}
```

Response:
```
{
  "id": "board_id",
  "user_id": "user_id",
  "title": "Board Title",
  "lists": []
}
```

### 4. Delete Board

### Endpoint: DELETE /board/:boardID

Response:
```
{
  "message": "Board deleted successfully"
}
```


# List Endpoints
### 1. Create List

### Endpoint: POST /list/:boardID

Request:
```
{
  "title": "List Title",
  "position": 1
}
```

Response:
```
{
  "id": "list_id",
  "title": "List Title",
  "position": 1,
  "cards": [],
  "card_details": []
}
```

### 2. Update List

### Endpoint: PATCH /list/:boardID/:listID

Request:
```
{
  "title": "New List Title"
}
```

Response:
```
{
  "id": "list_id",
  "title": "New List Title",
  "position": 1,
  "cards": [],
  "card_details": []
}
```

### 3. Delete List

### Endpoint: DELETE /list/:boardID/:listID

Response:
```
{
  "message": "List deleted successfully"
}
```

# Card Endpoints
### 1. Create Card

### Endpoint: POST /card

Request:
```
{
  "list_id": "list_id",
  "board_id": "board_id",
  "user_id": "user_id",
  "title": "Card Title",
  "description": "Card Description",
  "position": 1
}
```

Response:
```
{
  "id": "card_id",
  "list_id": "list_id",
  "board_id": "board_id",
  "user_id": "user_id",
  "title": "Card Title",
  "description": "Card Description",
  "position": 1
}
```


### 2. Update Card

### Endpoint: PATCH /card/:cardID

Request:
```
{
  "title": "New Card Title",
  "description": "New Card Description",
  "position": 2
}
```

Response:
```
{
  "id": "card_id",
  "list_id": "list_id",
  "board_id": "board_id",
  "user_id": "user_id",
  "title": "New Card Title",
  "description": "New Card Description",
  "position": 2
}
```

### 3. Move Card

### Endpoint: PATCH /card/:cardID/move

Request:
```
{
  "list_id": "new_list_id",
  "position": 3
}
```

Response:
```
{
  "id": "card_id",
  "list_id": "new_list_id",
  "board_id": "board_id",
  "user_id": "user_id",
  "title": "Card Title",
  "description": "Card Description",
  "position": 3
}
```

### 4. Delete Card with Transaction

### Endpoint: DELETE /card/:cardID

Response:
```
{
  "message": "Card deleted successfully"
}
```


# Models
### User Model
```
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email    string             `bson:"email" json:"email" validate:"required"`
	Password []byte             `bson:"password" json:"password" validate:"required"`
}
```

### Board Model
```
type Board struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID string             `bson:"user_id" json:"user_id" validate:"required"`
	Title  string             `bson:"title" json:"title" validate:"required"`
	Lists  []List             `bson:"lists" json:"lists" validate:"required"`
}
```

### List Model
```
type List struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title       string               `bson:"
}
```

### Card Model
```
type Card struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ListID      primitive.ObjectID `bson:"list_id" json:"list_id" validate:"required"`
	BoardID     primitive.ObjectID `bson:"board_id" json:"board_id" validate:"required"`
	UserID      string             `bson:"user_id" json:"user_id" validate:"required"`
	Title       string             `bson:"title" json:"title" validate:"required"`
	Description string             `bson:"description" json:"description"`
	Position    int                `bson:"position" json:"position"`
}
```














