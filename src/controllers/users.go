package controllers

import "net/http"

func UserCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create user"))
}

func UserGetAll(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get all users"))
}

func UserGetOne(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get one user"))
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update user"))
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete user"))
}
