package main

import (
	"fmt"
	"html"
	"net/http"
)

func adminBoard(w http.ResponseWriter, r *http.Request) {
	idCookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	var role, adminName string
	err = db.QueryRow("SELECT role, name FROM users WHERE id = $1", idCookie.Value).Scan(&role, &adminName)
	if err != nil || role != "admin" {
		http.Error(w, "Доступ запрещен: требуются права администратора", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="ru"><head><meta charset="UTF-8"><title>Админ-панель</title>%s</head>
<body>
<div class="container">
	<div class="header">
		<h2>Админ-панель (%s)</h2>
		<div class="btn-row">
			<a href="/profile" class="btn-secondary btn-sm">Мои задачи</a>
			<a href="/logout" class="btn-danger">Выйти</a>
		</div>
	</div>`, baseStyles, html.EscapeString(adminName))

	query := `
		SELECT tasks.id, tasks.title, tasks.description, tasks.state, users.name 
		FROM tasks 
		JOIN users ON tasks.user_id = users.id 
		WHERE tasks.status = '0'
		ORDER BY tasks.id DESC`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Ошибка загрузки задач: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var id int
	var title, description, state, author string
	for rows.Next() {
		rows.Scan(&id, &title, &description, &state, &author)
		fmt.Fprintf(w, `<div class="task-card">
			<div class="task-top">
				<h3>%s</h3>
				<span class="badge">%s</span>
			</div>
			<p class="task-desc">%s</p>
			<div class="task-footer">
				<span style="font-size: 13px; color: #64748b;">Автор: <b>%s</b></span>
				<form action="/adminDeleteTask" method="POST">
					<input type="hidden" name="task_id" value="%d">
					<button type="submit" class="btn-danger">Удалить задачу</button>
				</form>
			</div>
		</div>`, html.EscapeString(title), html.EscapeString(state), html.EscapeString(description), html.EscapeString(author), id)
	}
	fmt.Fprint(w, "</div></body></html>")
}

func adminDeleteTask(w http.ResponseWriter, r *http.Request) {
	idCookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	var role string
	err = db.QueryRow("SELECT role FROM users WHERE id = $1", idCookie.Value).Scan(&role)
	if err != nil || role != "admin" {
		http.Error(w, "Доступ запрещен", http.StatusForbidden)
		return
	}

	taskID := r.FormValue("task_id")
	_, err = db.Exec("UPDATE tasks SET status = '1' WHERE id = $1", taskID)
	if err != nil {
		http.Error(w, "Ошибка удаления: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}
