package storge

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"http-practive/internal/models"
	"strings"
)
	type TaskStore interface{
		Create(ctx context.Context, input models.CreateTaskInput)(*models.Task, error)
		List(ctx context.Context,completed *bool)([]models.Task, error)
		GetByID(ctx context.Context, id int)(*models.Task, error)
		Update(ctx context.Context, id int, input models.UpdateTaskInput)(*models.Task,error)
		Delete(ctx context.Context, id int)error
		SetCompleted(ctx context.Context, id int, completed bool)(*models.Task, error)
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
	conditions := []string{}
	args := []any{}
	if completed != nil{
		conditions = append(conditions,fmt.Sprintf("completed = $%d", len(args)+1))
		args = append(args, *completed)
	}
	if len(conditions)  > 0{
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := s.db.QueryContext(ctx, query,args...)
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
	err := s.db.QueryRowContext(ctx,query,input.Title,input.Description,input.Completed, id).Scan(&tasks.ID, &tasks.Title,&tasks.Description,&tasks.Completed,&tasks.CreatedAt)
	if errors.Is(err, sql.ErrNoRows){
		return nil, ErrTaskNotFound
	}
	if err != nil{
		return nil,err
	}
	return &tasks, nil
}

func (s *PostgresTaskStore)Delete(ctx context.Context,id int)(error){
	query := "DELETE FROM tasks WHERE id = $1"
	result,err := s.db.ExecContext(ctx,query,id)
	if err != nil{
		return err
	}
	n,err2 := result.RowsAffected()
	if err2 != nil{
		return err2
	}
	if n == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (s *PostgresTaskStore)SetCompleted(ctx context.Context, id int,completed bool)(*models.Task, error){
	var tasks models.Task
	query := "UPDATE tasks SET completed = $1 WHERE id = $2 RETURNING id,title,description,completed,created_at"
	err := s.db.QueryRowContext(ctx,query,completed,id).Scan(&tasks.ID,&tasks.Title,&tasks.Description,&tasks.Completed,&tasks.CreatedAt)
	if errors.Is(err,sql.ErrNoRows){
		return nil,ErrTaskNotFound
	}
	if err !=nil {
		return nil,err
	}
	return &tasks,nil
}