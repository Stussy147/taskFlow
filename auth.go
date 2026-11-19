package main

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const baseStyles = `<style>
	* { box-sizing: border-box; margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
	body { background: #f1f5f9; color: #0f172a; line-height: 1.5; }
	.center-page { display: flex; justify-content: center; align-items: center; min-height: 100vh; padding: 20px; }
	.auth-card { background: #fff; padding: 32px; border-radius: 12px; box-shadow: 0 4px 20px rgba(0,0,0,0.08); width: 100%; max-width: 400px; }
	.auth-card h2 { margin-bottom: 20px; font-size: 22px; text-align: center; }
	.stack-form { display: flex; flex-direction: column; gap: 12px; }
	input { padding: 10px 14px; border: 1px solid #cbd5e1; border-radius: 8px; font-size: 14px; outline: none; }
	input:focus { border-color: #3b82f6; }
	button, .btn-primary, .btn-secondary, .btn-danger { padding: 10px 16px; border-radius: 8px; border: none; font-size: 14px; font-weight: 500; cursor: pointer; text-decoration: none; text-align: center; display: inline-block; }
	.btn-primary, button[type="submit"] { background: #3b82f6; color: #fff; }
	.btn-primary:hover, button[type="submit"]:hover { background: #2563eb; }
	.btn-secondary { background: #e2e8f0; color: #334155; }
	.btn-danger { background: #ef4444; color: #fff; padding: 6px 12px; font-size: 12px; }
	.btn-sm { padding: 6px 12px; font-size: 12px; }
	.btn-row { display: flex; gap: 8px; }
	.link-hint { margin-top: 16px; text-align: center; font-size: 13px; color: #64748b; }
	.link-hint a { color: #3b82f6; text-decoration: none; }
	.container { max-width: 720px; margin: 40px auto; padding: 0 20px; }
	.header { display: flex; justify-content: space-between; align-items: center; background: #fff; padding: 20px 24px; border-radius: 12px; box-shadow: 0 2px 10px rgba(0,0,0,0.04); margin-bottom: 24px; }
	.create-box { background: #fff; padding: 20px 24px; border-radius: 12px; box-shadow: 0 2px 10px rgba(0,0,0,0.04); margin-bottom: 24px; display: flex; gap: 10px; }
	.create-box input { flex: 1; }
	.task-card { background: #fff; border-radius: 12px; padding: 20px 24px; margin-bottom: 12px; box-shadow: 0 2px 8px rgba(0,0,0,0.04); border: 1px solid #e2e8f0; }
	.task-top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
	.task-desc { color: #475569; font-size: 14px; margin-bottom: 16px; }
	.task-footer { display: flex; justify-content: space-between; align-items: center; border-top: 1px solid #f1f5f9; padding-top: 12px; }
	select { background: #eff6ff; color: #1d4ed8; border: 1px solid #bfdbfe; padding: 4px 10px; border-radius: 20px; font-size: 12px; font-weight: 500; cursor: pointer; outline: none; }
	.badge { background: #dbeafe; color: #1e40af; padding: 4px 10px; border-radius: 20px; font-size: 12px; font-weight: 500; }
</style>`

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="ru"><head><meta charset="UTF-8"><title>Регистрация</title>%s</head>
<body class="center-page">
	<div class="auth-card">
		<h2>Регистрация</h2>
		<form action="/register" method="POST" class="stack-form">
			<input name="name" placeholder="Имя пользователя" required>
			<input type="email" name="email" placeholder="Email" required>
			<input type="password" name="password" placeholder="Пароль" required>
			<button type="submit">Зарегистрироваться</button>
		</form>
		<div class="link-hint">Уже есть аккаунт? <a href="/login">Войти</a></div>
	</div>
</body></html>`, baseStyles)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		name := strings.TrimSpace(r.FormValue("name"))
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Ошибка шифрования пароля", http.StatusInternalServerError)
			return
		}

		var newID string
		query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id"
		err = db.QueryRow(query, name, email, string(hashedPassword)).Scan(&newID)
		if err != nil {
			http.Error(w, "Ошибка при сохранении в базу: "+err.Error(), http.StatusBadRequest)
			return
		}

		setAuthCookies(w, name, newID, "user")
		http.Redirect(w, r, "/profile", http.StatusFound)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="ru"><head><meta charset="UTF-8"><title>Вход</title>%s</head>
<body class="center-page">
	<div class="auth-card">
		<h2>Вход в TaskFlow</h2>
		<form action="/login" method="POST" class="stack-form">
			<input name="name" placeholder="Имя пользователя" required>
			<input type="password" name="password" placeholder="Пароль" required>
			<button type="submit">Войти</button>
		</form>
		<div class="link-hint">Нет аккаунта? <a href="/register">Регистрация</a></div>
	</div>
</body></html>`, baseStyles)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		name := strings.TrimSpace(r.FormValue("name"))
		password := r.FormValue("password")

		var dbID, dbPassword, dbRole string
		err := db.QueryRow("SELECT id, password, role FROM users WHERE name = $1", name).Scan(&dbID, &dbPassword, &dbRole)
		if err != nil {
			http.Error(w, "Пользователь с таким именем не найден", http.StatusUnauthorized)
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)); err != nil {
			http.Error(w, "Неверный пароль", http.StatusUnauthorized)
			return
		}

		setAuthCookies(w, name, dbID, dbRole)
		if dbRole == "admin" {
			http.Redirect(w, r, "/admin", http.StatusFound)
		} else {
			http.Redirect(w, r, "/profile", http.StatusFound)
		}
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	for _, cookieName := range []string{"user_name", "user_id", "role"} {
		http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "", Path: "/", MaxAge: -1})
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

func setAuthCookies(w http.ResponseWriter, name, id, role string) {
	http.SetCookie(w, &http.Cookie{Name: "user_name", Value: name, Path: "/"})
	http.SetCookie(w, &http.Cookie{Name: "user_id", Value: id, Path: "/"})
	http.SetCookie(w, &http.Cookie{Name: "role", Value: role, Path: "/"})
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	nameCookie, err := r.Cookie("user_name")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	idCookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	adminLink := ""
	if roleCookie, err := r.Cookie("role"); err == nil && roleCookie.Value == "admin" {
		adminLink = `<a href="/admin" class="btn-secondary btn-sm">Админ-панель</a>`
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="ru"><head><meta charset="UTF-8"><title>Профиль</title>%s</head>
<body>
<div class="container">
	<div class="header">
		<h2>Привет, %s!</h2>
		<div class="btn-row">%s <a href="/logout" class="btn-danger">Выйти</a></div>
	</div>
	<form action="/createTask" method="POST" class="create-box">
		<input name="nameTask" placeholder="Новая задача..." required>
		<input name="description" placeholder="Описание (необязательно)">
		<button type="submit">Создать</button>
	</form>`, baseStyles, html.EscapeString(nameCookie.Value), adminLink)

	query := "SELECT id, title, description, state FROM tasks WHERE user_id = $1 AND status = '0' ORDER BY id DESC"
	rows, err := db.Query(query, idCookie.Value)
	if err != nil {
		fmt.Fprint(w, "<p>Ошибка загрузки задач</p></div></body></html>")
		return
	}
	defer rows.Close()

	var id int
	var title, description, state string
	for rows.Next() {
		rows.Scan(&id, &title, &description, &state)
		safeTitle := html.EscapeString(title)
		safeDesc := html.EscapeString(description)
		safeState := html.EscapeString(state)

		fmt.Fprintf(w, `<div class="task-card">
			<div class="task-top">
				<h3>%s</h3>
				<form action="/updateState" method="POST">
					<input type="hidden" name="task_id" value="%d">
					<select name="state" onchange="this.form.submit()">
						<option value="%s" selected hidden>%s</option>
						<option value="новая">новая</option>
						<option value="в процессе">в процессе</option>
						<option value="выполнено">выполнено</option>
					</select>
				</form>
			</div>
			<p class="task-desc">%s</p>
			<div class="task-footer">
				<a href="/editTask?id=%d" class="btn-primary btn-sm">Редактировать</a>
				<form action="/deleteTask" method="POST">
					<input type="hidden" name="task_id" value="%d">
					<button type="submit" class="btn-danger">Удалить</button>
				</form>
			</div>
		</div>`, safeTitle, id, safeState, safeState, safeDesc, id, id)
	}
	fmt.Fprint(w, "</div></body></html>")
}
