package models

import "sync"

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type TaskStore struct {
	mu     sync.Mutex
	tasks  map[int]Task
	nextID int
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}

func (s *TaskStore) Create(title string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}
	s.tasks[s.nextID] = task
	s.nextID++
	return task
}

func (s *TaskStore) GetAll() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task)
	}
	return result
}

func (s *TaskStore) GetByID(id int) (*Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, found := s.tasks[id]
	if !found {
		return nil, false
	}
	return &task, true
}

func (s *TaskStore) Update(id int, done bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, found := s.tasks[id]
	if !found {
		return false
	}
	task.Done = done
	s.tasks[id] = task
	return true
}

func (s *TaskStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, found := s.tasks[id]
	if !found {
		return false
	}
	delete(s.tasks, id)
	return true
}
