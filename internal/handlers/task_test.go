package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task-api/internal/models"
	"testing"
)

func TestCreateTask(t *testing.T) {
	store := models.NewTaskStore()
	handler := NewTaskHandler(store)

	body := map[string]string{"title": "Test task"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateTask(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var task models.Task
	json.NewDecoder(w.Body).Decode(&task)

	if task.Title != "Test task" {
		t.Errorf("Expected title 'Test task', got '%s'", task.Title)
	}

	if task.ID != 1 {
		t.Errorf("Expected ID 1, got %d", task.ID)
	}
}

func TestCreateTaskEmptyTitle(t *testing.T) {
	store := models.NewTaskStore()
	handler := NewTaskHandler(store)

	body := map[string]string{"title": ""}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetTasks(t *testing.T) {
	store := models.NewTaskStore()
	handler := NewTaskHandler(store)

	store.Create("Task 1")
	store.Create("Task 2")

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var tasks []models.Task
	json.NewDecoder(w.Body).Decode(&tasks)

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
}

func TestGetTaskByID(t *testing.T) {
	store := models.NewTaskStore()
	handler := NewTaskHandler(store)

	store.Create("Test task")

	req := httptest.NewRequest(http.MethodGet, "/tasks?id=1", nil)
	w := httptest.NewRecorder()

	handler.GetTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var task models.Task
	json.NewDecoder(w.Body).Decode(&task)

	if task.Title != "Test task" {
		t.Errorf("Expected title 'Test task', got '%s'", task.Title)
	}
}

func TestGetTaskByInvalidID(t *testing.T) {
	store := models.NewTaskStore()
	handler := NewTaskHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/tasks?id=999", nil)
	w := httptest.NewRecorder()

	handler.GetTasks(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestUpdateTask(t *testing.T) {
	store := models.NewTaskStore()
	handler := NewTaskHandler(store)

	store.Create("Test task")

	body := map[string]bool{"done": true}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPatch, "/tasks?id=1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateTask(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestDeleteTask(t *testing.T) {
	store := models.NewTaskStore()
	handler := NewTaskHandler(store)

	store.Create("Test task")

	req := httptest.NewRequest(http.MethodDelete, "/tasks?id=1", nil)
	w := httptest.NewRecorder()

	handler.DeleteTask(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	tasks := store.GetAll()
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks after delete, got %d", len(tasks))
	}
}
