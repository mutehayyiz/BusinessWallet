package controller

import (
	"BusinessWallet/storage"
	"encoding/json"
	"fmt"
	"net/http"
)

func ReturnResponse(w http.ResponseWriter, statusCode int, resp interface{}) {
	bytes, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprintf(w, string(bytes))
}

func ReturnError(w http.ResponseWriter, statusCode int, errMsg string) {
	resp := map[string]interface{}{
		"error": errMsg,
	}
	bytes, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprintf(w, string(bytes))
}

func DeleteAll(w http.ResponseWriter, r *http.Request) {
	err := storage.Delete()
	if err != nil {
		ReturnError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ReturnResponse(w, http.StatusOK, "ok")
}

func Seed(w http.ResponseWriter, r *http.Request) {
	storage.Seed()

	ReturnResponse(w, http.StatusOK, "ok")
}
