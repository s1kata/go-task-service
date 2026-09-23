package handler

import (
	"encoding/json"
	"errors"
	"http-practive/internal/storge"
	"log"
	"net/http"
)

func ResponseWithJSON(w http.ResponseWriter, status int, payload any){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil{
		log.Println("encode response",err)
	}
}
func ResponseWithErrorJSON(w http.ResponseWriter, status int, msg string){
	ResponseWithJSON(w, status, map[string]string{"error":msg} )
}
func WriteError(w http.ResponseWriter, err error){
	switch{
	case errors.Is(err, storge.ErrTaskNotFound):
		ResponseWithErrorJSON(w,http.StatusNotFound,"task not found")
	default:
		log.Println("internal error",err)
		ResponseWithErrorJSON(w, http.StatusInternalServerError, "intnernail server error")
	}
	
}
