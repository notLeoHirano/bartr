package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/notLeoHirano/bartr/database"
	"github.com/notLeoHirano/bartr/handlers"
	"github.com/notLeoHirano/bartr/models"
	"github.com/notLeoHirano/bartr/service"
	"github.com/notLeoHirano/bartr/store"
)

var testHandler *handlers.Handler
var testDB *database.DB

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Exit(code)
}

func setupTest(t *testing.T) {
	var err error
	testDB, err = database.New(":memory:")
	if err != nil {
		t.Fatal("Failed to open test database:", err)
	}

	if err := testDB.Init(); err != nil {
		t.Fatal("Failed to initialize test database:", err)
	}

	st := store.New(testDB.DB)
	svc := service.New(st)
	testHandler = handlers.New(svc)
}

func teardownTest() {
	if testDB != nil {
		testDB.Close()
	}
}

func TestCreateItem_Success(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	router := gin.New()
	router.POST("/items", testHandler.CreateItem)

	item := models.Item{
		UserID:      1,
		Title:       "Vintage Guitar",
		Description: "1960s Fender Stratocaster",
		Category:    "Musical Instruments",
	}

	body, _ := json.Marshal(item)
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response models.Item
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Title != item.Title {
		t.Errorf("Expected title '%s', got '%s'", item.Title, response.Title)
	}

	if response.ID == 0 {
		t.Error("Expected item ID to be set")
	}
}

func TestCreateItem_MissingTitle(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	router := gin.New()
	router.POST("/items", testHandler.CreateItem)

	item := models.Item{
		UserID:      1,
		Description: "Missing title",
	}

	body, _ := json.Marshal(item)
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "title is required" {
		t.Errorf("Expected error 'title is required', got '%s'", response["error"])
	}
}

func TestGetItems_Success(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	// Insert test item
	testDB.Exec("INSERT INTO items (user_id, title, description, category) VALUES (?, ?, ?, ?)",
		1, "Test Item", "Test Description", "Test Category")

	router := gin.New()
	router.GET("/items", testHandler.GetItems)

	req, _ := http.NewRequest("GET", "/items", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var items []models.ItemWithOwner
	json.Unmarshal(w.Body.Bytes(), &items)

	if len(items) == 0 {
		t.Error("Expected at least one item")
	}

	if items[len(items)-1].Title != "Test Item" {
		t.Errorf("Expected title 'Test Item', got '%s'", items[0].Title)
	}
}

func TestCreateSwipe_InvalidDirection(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	router := gin.New()
	router.POST("/swipes", testHandler.CreateSwipe)

	swipe := models.Swipe{
		UserID:    1,
		ItemID:    1,
		Direction: "invalid",
	}

	body, _ := json.Marshal(swipe)
	req, _ := http.NewRequest("POST", "/swipes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["error"] != "Invalid request body" {
		t.Errorf("Expected direction error, got '%s'", response["error"])
	}
}

// fake auth
func makeAuthRouter(handler gin.HandlerFunc, route string, method string, userID int) *gin.Engine {
	r := gin.New()
	// Inject the faked userID into the context before the handler runs
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})

	// Attach the actual route and handler
	switch method {
	case "POST":
		r.POST(route, handler)
	case "DELETE":
		r.DELETE(route, handler)
	case "GET":
		r.GET(route, handler)
	}
	return r
}

// New simplified request executor
func performRequest(router *gin.Engine, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}



func TestMatchCreation(t *testing.T) {
    setupTest(t)
    defer teardownTest()

    // --- Insert items and retrieve their ACTUAL IDs ---
    result1, _ := testDB.Exec("INSERT INTO items (user_id, title, description, category) VALUES (?, ?, ?, ?)", 1, "Alice's Guitar", "Vintage guitar", "Instruments")
    item1ID, _ := result1.LastInsertId()
    
    result2, _ := testDB.Exec("INSERT INTO items (user_id, title, description, category) VALUES (?, ?, ?, ?)", 2, "Bob's Keyboard", "Synthesizer", "Instruments")
    item2ID, _ := result2.LastInsertId()

    // --- User 2 (Bob) swipes right on Alice's guitar (item 1) ---
    // Use a MAP to explicitly set the snake_case JSON keys (item_id, direction)
    body1Map := map[string]interface{}{
        "item_id": int(item1ID),
        "direction": "right",
    }
    body1, _ := json.Marshal(body1Map)
    
    bobRouter := makeAuthRouter(testHandler.CreateSwipe, "/swipes", "POST", 2)
    w1 := performRequest(bobRouter, "POST", "/swipes", body1)

    if w1.Code != http.StatusCreated {
        t.Fatalf("Bob swipe failed: expected 201, got %d. Body: %s", w1.Code, w1.Body.String())
    }

    // --- User 1 (Alice) swipes right on Bob's keyboard (item 2) ---
    // Use a MAP to explicitly set the snake_case JSON keys (item_id, direction)
    body2Map := map[string]interface{}{
        "item_id": int(item2ID),
        "direction": "right",
    }
    body2, _ := json.Marshal(body2Map)

    aliceRouter := makeAuthRouter(testHandler.CreateSwipe, "/swipes", "POST", 1)
    w2 := performRequest(aliceRouter, "POST", "/swipes", body2)
    
    if w2.Code != http.StatusCreated {
        t.Fatalf("Alice swipe failed: expected 201, got %d. Body: %s", w2.Code, w2.Body.String())
    }

    // --- Verification ---
    var count int
    testDB.QueryRow("SELECT COUNT(*) FROM matches").Scan(&count)
    if count != 1 {
        t.Errorf("Expected 1 match, got %d", count)
    }

    // Verify match details
    var user1ID, user2ID, dbItem1ID, dbItem2ID int
    err := testDB.QueryRow("SELECT user1_id, user2_id, item1_id, item2_id FROM matches LIMIT 1").
        Scan(&user1ID, &user2ID, &dbItem1ID, &dbItem2ID)
    if err != nil {
        t.Fatalf("Failed to query match: %v", err)
    }

    // Check for both possible orderings (User 1/Item 1, User 2/Item 2) or vice versa
    if !((user1ID == 1 && user2ID == 2 && dbItem1ID == int(item1ID) && dbItem2ID == int(item2ID)) ||
        (user1ID == 2 && user2ID == 1 && dbItem1ID == int(item2ID) && dbItem2ID == int(item1ID))) {
        t.Errorf("Match details incorrect: (user1=%d, user2=%d, item1=%d, item2=%d)", user1ID, user2ID, dbItem1ID, dbItem2ID)
    }
}

func TestDeleteItem_Success(t *testing.T) {
    setupTest(t)
    defer teardownTest()

    // Insert item owned by user 1
    result, _ := testDB.Exec("INSERT INTO items (user_id, title, description, category) VALUES (?, ?, ?, ?)",
        1, "Test Item", "To be deleted", "Misc")
    id, _ := result.LastInsertId()
    itemID := strconv.FormatInt(id, 10)

    // Create router with faked UserID: 1
    deleteRouter := makeAuthRouter(testHandler.DeleteItem, "/items/:id", "DELETE", 1)
    
    w := performRequest(deleteRouter, "DELETE", "/items/"+itemID, nil)
    if w.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d. Body: %s", w.Code, w.Body.String()) // FIX: Check body for details
    }

    // Verify deletion
    var count int
    testDB.QueryRow("SELECT COUNT(*) FROM items WHERE id = ?", id).Scan(&count)
    if count != 0 {
        t.Error("Expected item to be deleted")
    }
}

func TestDeleteItem_NotFound(t *testing.T) {
    setupTest(t)
    defer teardownTest()

    deleteRouter := makeAuthRouter(testHandler.DeleteItem, "/items/:id", "DELETE", 1)

    w := performRequest(deleteRouter, "DELETE", "/items/99999", nil)
    // Expect 403 Forbidden based on the handler's error response logic
    if w.Code != http.StatusForbidden { 
        t.Errorf("Expected 403, got %d. Body: %s", w.Code, w.Body.String())
    }

    // Check the body content to ensure it's the expected permission error
    var response map[string]string
    json.Unmarshal(w.Body.Bytes(), &response)
    expectedError := "item not found or you don't have permission to delete it"
    if response["error"] != expectedError {
        t.Errorf("Expected error '%s', got '%s'", expectedError, response["error"])
    }
}