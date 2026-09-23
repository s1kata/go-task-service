package handler

import (
	"encoding/json"
	"http-practive/internal/models"
	"http-practive/internal/storge"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct{
	store storge.TaskStore
}


func HandlePing(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		ResponseWithErrorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	ResponseWithJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func HandleHello(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		ResponseWithErrorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	name := r.URL.Query().Get("name")

	if name == ""{
		ResponseWithErrorJSON(w, http.StatusBadRequest,"name required")
		return
	}
	ResponseWithJSON(w, http.StatusOK, map[string]string{"hello":name})
	
}

func HandleNotFound(w http.ResponseWriter, r *http.Request){
	ResponseWithErrorJSON(w, http.StatusNotFound, "not-found")
}

func (h *Handler)Get(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		ResponseWithErrorJSON(w, http.StatusMethodNotAllowed, "405")
		return
	}
	tasks,err := h.store.List(r.Context(), nil)
	if err != nil{
		WriteError(w,err)
		return
	}
	ResponseWithJSON(w, http.StatusOK, tasks)
	
}
func NewHandler(store storge.TaskStore) *Handler{
	return &Handler{store: store}
}
func (h *Handler) TaskByID(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		ResponseWithErrorJSON(w,http.StatusMethodNotAllowed,"405" )
		return
	}
	path := r.URL.Path
	idstr := strings.TrimPrefix(path, "/tasks/")
	id, err := strconv.Atoi(idstr)
	if err != nil{
		ResponseWithErrorJSON(w, http.StatusBadRequest, "400")
		return
	}
	task, err := h.store.GetByID(r.Context(),id)
	if err != nil{
		WriteError(w, err)
		return
	}
	ResponseWithJSON(w,http.StatusOK, task)
}
func (h *Handler)Update(w http.ResponseWriter, r *http.Request){
	var input models.UpdateTaskInput
	if r.Method != http.MethodPatch{
		ResponseWithErrorJSON(w,http.StatusMethodNotAllowed,"method not allowed")
		return
	}
	path := r.URL.Path
	idstr := strings.TrimPrefix(path, "/tasks/")
	id, err := strconv.Atoi(idstr)
	if err != nil || id <= 0{
		ResponseWithErrorJSON(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil{
		ResponseWithErrorJSON(w, http.StatusBadRequest, "invalid json")
		return
	}
	if input.Title == ""{
		ResponseWithErrorJSON(w, http.StatusBadRequest, "title required")
		return
	}
	task,err := h.store.Update(r.Context(),id,input)
	if err != nil{
	WriteError(w,err)
	return
	}
	ResponseWithJSON(w,http.StatusOK, task)
}
func (h *Handler)TaskCreate(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		ResponseWithErrorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var input models.CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil{
		ResponseWithErrorJSON(w, http.StatusBadRequest, "invalid json")
		return 
	}
	if input.Title == ""{
		ResponseWithErrorJSON(w, http.StatusBadRequest, "title required")
		return 
	}
	task,err := h.store.Create(r.Context(), input)
	if err != nil{
		WriteError(w,err)
		return
	}
	ResponseWithJSON(w,http.StatusCreated, task)
}