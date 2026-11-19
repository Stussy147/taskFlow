package main

import (
	"fmt"
	"html"
	"net/http"
	"strings"
)

func createTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}

	idCookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	r.ParseForm()
	nameTask := strings.TrimSpace(r.FormValue("nameTask"))
	description := strings.TrimSpace(r.FormValue("description"))
	if nameTask == "" {
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}

	query := "INSERT INTO tasks (title, description, user_id, state, status) VALUES ($1, $2, $3, $4, $5)"
	_, err = db.Exec(query, nameTask, description, idCookie.Value, "новая", 0)
	if err != nil {
		http.Error(w, "Ошибка при сохранении в базу: "+err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/profile", http.StatusFound)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	idCookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	taskID := r.FormValue("task_id")
	query := "UPDATE tasks SET status = '1' WHERE id = $1 AND user_id = $2"
	_, err = db.Exec(query, taskID, idCookie.Value)
	if err != nil {
		http.Error(w, "Ошибка удаления: "+err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/profile", http.StatusFound)
}

func updateState(w http.ResponseWriter, r *http.Request) {
	idCookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	taskID := r.FormValue("task_id")
	stateValue := r.FormValue("state")
	query := "UPDATE tasks SET state = $1 WHERE id = $2 AND user_id = $3"
	_, err = db.Exec(query, stateValue, taskID, idCookie.Value)
	if err != nil {
		http.Error(w, "Ошибка обновления: "+err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/profile", http.StatusFound)
}

func editTask(w http.ResponseWriter, r *http.Request) {
	idCookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	userID := idCookie.Value

	if r.Method == http.MethodGet {
		taskID := r.URL.Query().Get("id")

		var oldTitle, oldDesc string
		err := db.QueryRow("SELECT title, description FROM tasks WHERE id = $1 AND user_id = $2", taskID, userID).Scan(&oldTitle, &oldDesc)
		if err != nil {
			http.Error(w, "Задача не найдена", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="ru">
<head><meta charset="UTF-8"><title>Редактирование</title>%s</head>
<body class="center-page">
	<div class="auth-card">
		<h2>Редактирование задачи</h2>
		<form action="/editTask" method="POST" class="stack-form">
			<input type="hidden" name="task_id" value="%s">
			<input type="text" name="title" value="%s" placeholder="Название" required>
			<input type="text" name="description" value="%s" placeholder="Описание">
			<div class="btn-row">
				<button type="submit" class="btn-primary">Сохранить</button>
				<a href="/profile" class="btn-secondary">Отмена</a>
			</div>
		</form>
	</div>
</body></html>`, baseStyles, html.EscapeString(taskID), html.EscapeString(oldTitle), html.EscapeString(oldDesc))

	} else if r.Method == http.MethodPost {
		r.ParseForm()
		taskID := r.FormValue("task_id")
		title := strings.TrimSpace(r.FormValue("title"))
		description := strings.TrimSpace(r.FormValue("description"))

		query := "UPDATE tasks SET title = $1, description = $2 WHERE id = $3 AND user_id = $4"
		_, err = db.Exec(query, title, description, taskID, userID)
		if err != nil {
			http.Error(w, "Ошибка обновления: "+err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/profile", http.StatusFound)
	}
}
