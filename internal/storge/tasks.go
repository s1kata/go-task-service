package storge

import (
	"context"
	"database/sql"
	"errors"
	"http-practive/internal/models"
)
	type TaskStore interface{
		Create(ctx context.Context, input models.CreateTaskInput)(*models.Task, error)
		List(ctx context.Context,completed *bool)([]models.Task, error)
		GetByID(ctx context.Context, id int)(*models.Task, error)
		Update(ctx context.Context, id int, input models.UpdateTaskInput)(*models.Task,error)
		//Delete(ctx context.Context, id int)error
	}
	type PostgresTaskStore struct{
		db *sql.DB

	}
	var ErrTaskNotFound = errors.New("task not found")
	func NewPostgresTaskStore(db *sql.DB)*PostgresTaskStore{
		return &PostgresTaskStore{db:db}
	}

func (s *PostgresTaskStore) Create(ctx context.Context, input models.CreateTaskInput)(*models.Task, error){
	var task models.Task
	query := "INSERT INTO tasks (title,description) VALUES ($1,$2) RETURNING id, title,description,completed,created_at"
	err :=s.db.QueryRowContext(ctx, query, input.Title, input.Description).Scan(&task.ID, &task.Title, &task.Description, &task.Completed,&task.CreatedAt)
	if err != nil{
		return nil,err
	}
	return &task, nil
}
func (s *PostgresTaskStore)List(ctx context.Context, completed *bool)([]models.Task, error){
	query := "SELECT id,title,description,completed, created_at FROM tasks"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil{
		return nil,err
	}
	defer rows.Close()
	tasks := make([]models.Task, 0)
	for rows.Next(){
		
		var t models.Task
	if err := rows.Scan(&t.ID, &t.Title,&t.Description,&t.Completed,&t.CreatedAt); err != nil{
		return nil,err
	}
		tasks = append(tasks, t)
		
	}
	if err := rows.Err(); err != nil {
    return nil, err
}
	return tasks, nil
	
}

func (s *PostgresTaskStore)GetByID(ctx context.Context, id int)(*models.Task, error){
	var task models.Task
	query := "SELECT id, title,description,completed,created_at FROM tasks WHERE id = $1"
	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&task.ID,&task.Title,&task.Description,&task.Completed,&task.CreatedAt)
	if errors.Is(err, sql.ErrNoRows){
		return nil, ErrTaskNotFound
	}
	if err != nil{
		return nil,err
	}
	return &task,nil

}
func (s *PostgresTaskStore)Update(ctx context.Context, id int, input models.UpdateTaskInput)(*models.Task, error){
	var tasks models.Task
	query := "UPDATE tasks SET title = $1, description = $2, completed = $3 WHERE id = $4 RETURNING id,title,description,completed,created_at"
	err := s.db.QueryRowContext(ctx,query,input.Title,input.Description,input.Completed, id).Scan(&tasks.ID, tasks.Title,&tasks.Description,&tasks.Completed,&tasks.CreatedAt)
	if errors.Is(err, sql.ErrNoRows){
		return nil, err
	}
	if err != nil{
		return nil,err
	}
	return &tasks, nil
}