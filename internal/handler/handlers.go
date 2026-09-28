package handler

import (
	"encoding/json"
	"http-practive/internal/models"
	"http-practive/internal/storge"
	"log"
	"net/http"
	"strconv"
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
	log.Printf("NOT FOUND HIT: method=%q path=%q", r.Method, r.URL.Path)
	ResponseWithErrorJSON(w, http.StatusNotFound, "not-found")
}

func (h *Handler)Get(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		ResponseWithErrorJSON(w, http.StatusMethodNotAllowed, "405")
		return
	}
	var completed *bool
	path := r.URL.Query().Get("completed")
	if path != ""{
		boolVal,err := strconv.ParseBool(path)
		if err != nil{
			ResponseWithErrorJSON(w, http.StatusBadRequest, "400")
			return
		}
		completed = &boolVal
	}
	tasks,err := h.store.List(r.Context(),completed)
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
		ResponseWithErrorJSON(w,http.StatusMethodNotAllowed,"method not allowed" )
		return
	}
	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil{
		ResponseWithErrorJSON(w, http.StatusBadRequest, "invalid id")
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
	idstr := r.PathValue("id")
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
func (h *Handler)TaskDelete(w http.ResponseWriter,r *http.Request){
	log.Printf("Delete hit: method=%q path=%q",r.Method, r.URL.Path)
	if r.Method != http.MethodDelete{
		ResponseWithErrorJSON(w,http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	idstr := r.PathValue("id")
	id,err := strconv.Atoi(idstr)
	if err != nil{
		ResponseWithErrorJSON(w, http.StatusBadRequest, "invalid id")
		return
	}
	err = h.store.Delete(r.Context(), id)
	if err != nil{
		WriteError(w,err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}