package utils

import (
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

var PORT string
var SESSION_KEY string
var DATABASE_URL string

func SetParams() error {
	err := godotenv.Load()
	PORT = os.Getenv("PORT")
	SESSION_KEY = os.Getenv("SESSION_KEY")
	DATABASE_URL = os.Getenv("DATABASE_URL")
	return err

}

type ErrorMessage struct {
	Message string `json:"error"`
}

type BasicData struct {
	Data any `json:"data"`
}

func Forbidden(w http.ResponseWriter) {
	JSON(w, http.StatusForbidden, "Access denied.")
}

func JSON(w http.ResponseWriter, status int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if _, err := w.Write(buf.Bytes()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func String(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, &BasicData{Data: msg})
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, ErrorMessage{
		Message: msg,
	})
}
