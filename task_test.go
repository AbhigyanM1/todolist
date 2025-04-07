package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-backend/cassandra"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	cassandra.Connect()
	r := gin.Default()
	r.POST("/tasks", AddTask)
	r.GET("/tasks", GetTasks)
	return r
}

func TestAddAndGetTasks(t *testing.T) {
	router := setupRouter()

	// Add a task
	task := map[string]interface{}{
		"title":     "Test from Go test",
		"completed": false,
	}
	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.Code)
	}

	// Get all tasks
	req2, _ := http.NewRequest("GET", "/tasks", nil)
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)

	if resp2.Code != http.StatusOK {
		t.Errorf("Expected 200 OK on GET, got %d", resp2.Code)
	}
}