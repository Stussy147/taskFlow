package main

import (
	"fmt"
	"net/http"
)

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusFound)
	})
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/createTask", createTask)
	http.HandleFunc("/deleteTask", deleteTask)
	http.HandleFunc("/updateState", updateState)
	http.HandleFunc("/editTask", editTask)
	http.HandleFunc("/admin", adminBoard)
	http.HandleFunc("/adminDeleteTask", adminDeleteTask)

	fmt.Println("Сервер запущен: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
