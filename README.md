# Bartr - Tinder for Item Swapping

A full-stack web application that allows users to swipe through items posted by others and match when both users are interested in each other's items. Built with a Go backend and vanilla JavaScript frontend.

## Quick Start

### Prerequisites

- A recent version of Go

- A modern web browser

### Setup & Running the Backend (Go)

1. Clone the repository:
  
  ```bash
  git clone https://github.com/notLeoHirano/bartr.git
  cd bartr
  ```

2. Install dependencies:

  ```bash
  go mod download
  ```

3. Run the backend server:
  
  ```bash
  go run main.go
  ```

The server will start on http://localhost:8080. If it doesnt, be sure to change the API_URL in the example frontend:

### Running the Frontend (JavaScript)

Simply open the index.html file in your browser.

### Running Tests

To run all tests:

  ``` bash
  go test -v
  ```

To run a specific test:
  
  ``` bash
  go test -v -run TestMatchCreation
  ```

## API Endpoints

### Authentication

All endpoints except /register and /login require a JWT token to be passed in the Authorization: Bearer YOUR_TOKEN header.

| Method | Endpoint   | Description                | Auth Required |
|--------|------------|----------------------------|---------------|
| POST   | /register  | Create a new user account  | No            |
| POST   | /login     | Login and get a JWT token  | No            |

### Items

| Method | Endpoint     | Description                                                        | Auth Required |
|--------|-------------|--------------------------------------------------------------------|---------------|
| GET    | /items      | List all items (excludes user's own items and items already swiped on) | Yes           |
| POST   | /items      | Create a new item                                                   | Yes           |
| DELETE | /items/:id  | Delete one of your items                                           | Yes           |

### Swipes & Matches

| Method | Endpoint    | Description                 | Auth Required |
|--------|------------|-----------------------------|---------------|
| POST   | /swipes    | Record a swipe (left or right) | Yes           |
| GET    | /matches   | Get all your matches         | Yes           |

### Comments

| Method | Endpoint                 | Description                 | Auth Required |
|--------|--------------------------|-----------------------------|---------------|
| GET    | /matches/:id/comments    | Get all comments for a specific match | Yes |
| POST   | /comments                | Add a comment to a match    | Yes           |

## API Examples

### Register

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice",
    "email": "alice@example.com",
    "password": "password123"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "password123"
  }'
```


### Response

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "Alice",
    "email": "alice@example.com"
  }
}
```

### Create an Item

```bash
curl -X POST http://localhost:8080/items \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Vintage Guitar",
    "description": "1960s Fender Stratocaster in excellent condition",
    "category": "Musical Instruments"
  }'
```

### Swipe on an Item

```bash
curl -X POST http://localhost:8080/swipes \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "item_id": 2,
    "direction": "right"
  }'
```

### Get Your Matches

```bash
curl -X GET http://localhost:8080/matches \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## How Matching Works

1. User 1 posts "Item A".

2. User 2 posts "Item B".

3. User 1 swipes right on "Item B".

4. The system checks if User 2 has already swiped right on any of User 1's items. In this case, no.

5. User 2 swipes right on "Item A".

6. The system checks if User 1 has already swiped right on any of User 2's items. In this case, yes (from step 3).

7. A match is created between "Item A" and "Item B". Both users can now see this match and add comments.
