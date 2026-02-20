package handlers

import (
	"encoding/json"
	"net/http"
)

type ExternalTask struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func (h *TaskHandler) GetExternalTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch external tasks"})
		return
	}
	defer resp.Body.Close()

	var externalTasks []ExternalTask
	if err := json.NewDecoder(resp.Body).Decode(&externalTasks); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to parse external tasks"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(externalTasks)
}
