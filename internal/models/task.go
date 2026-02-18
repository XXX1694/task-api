package models

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type TaskStore struct {
	tasks  []Task
	nextID int
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  make([]Task, 0),
		nextID: 1,
	}
}

func (s *TaskStore) Create(title string) Task {
	task := Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}
	s.tasks = append(s.tasks, task)
	s.nextID++
	return task
}

func (s *TaskStore) GetAll() []Task {
	return s.tasks
}

func (s *TaskStore) GetByID(id int) (*Task, bool) {
	for i := range s.tasks {
		if s.tasks[i].ID == id {
			return &s.tasks[i], true
		}
	}
	return nil, false
}

func (s *TaskStore) Update(id int, done bool) bool {
	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks[i].Done = done
			return true
		}
	}
	return false
}
