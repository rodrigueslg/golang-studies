package model

type Task struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

func (Task) TableName() string {
	return "task"
}
