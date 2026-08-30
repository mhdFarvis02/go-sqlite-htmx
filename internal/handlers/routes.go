package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/dporkka/go-sqlite-htmx/internal/db"
)

type Application struct {
	DB        *db.DB
	Templates *template.Template
}

// Home displays the main page with all users
func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	users, err := app.DB.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to load users", http.StatusInternalServerError)
		log.Printf("Error loading users: %v", err)
		return
	}

	data := struct {
		Users []*db.User
	}{
		Users: users,
	}

	if err := app.Templates.ExecuteTemplate(w, "home.page.tmpl", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

// UserView renders a single user in read-only mode
func (app *Application) UserView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := app.DB.GetUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to load user", http.StatusInternalServerError)
			log.Printf("Error loading user: %v", err)
		}
		return
	}

	// Render just the user view partial
	tmpl := `<div id="user-{{.ID}}" class="user-card">
	<h3>{{.Name}}</h3>
	<p><strong>Email:</strong> {{.Email}}</p>
	<button hx-get="/user/{{.ID}}/edit" hx-target="#user-{{.ID}}" hx-swap="outerHTML">
		Edit
	</button>
</div>`

	t, err := template.New("user-view").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, user); err != nil {
		http.Error(w, "Failed to render", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

// UserEdit renders a single user in edit mode (form)
func (app *Application) UserEdit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := app.DB.GetUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to load user", http.StatusInternalServerError)
			log.Printf("Error loading user: %v", err)
		}
		return
	}

	// Render the edit form
	tmpl := `<form id="user-{{.ID}}" class="user-card edit-form"
	hx-post="/user/{{.ID}}/update"
	hx-target="#user-{{.ID}}"
	hx-swap="outerHTML">
	<div class="form-group">
		<label for="name-{{.ID}}">Name</label>
		<input type="text" id="name-{{.ID}}" name="name" value="{{.Name}}" required>
	</div>
	<div class="form-group">
		<label for="email-{{.ID}}">Email</label>
		<input type="email" id="email-{{.ID}}" name="email" value="{{.Email}}" required>
	</div>
	<div class="form-actions">
		<button type="submit">Save</button>
		<button type="button" hx-get="/user/{{.ID}}/view" hx-target="#user-{{.ID}}" hx-swap="outerHTML">
			Cancel
		</button>
	</div>
</form>`

	t, err := template.New("user-edit").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, user); err != nil {
		http.Error(w, "Failed to render", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

// UserUpdate handles the form submission and updates the user
func (app *Application) UserUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")

	if name == "" || email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	if err := app.DB.UpdateUser(id, name, email); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to update user: %v", err), http.StatusInternalServerError)
			log.Printf("Error updating user: %v", err)
		}
		return
	}

	// Fetch the updated user and render the view
	user, err := app.DB.GetUser(id)
	if err != nil {
		http.Error(w, "Failed to load updated user", http.StatusInternalServerError)
		return
	}

	// Render the view mode
	tmpl := `<div id="user-{{.ID}}" class="user-card">
	<h3>{{.Name}}</h3>
	<p><strong>Email:</strong> {{.Email}}</p>
	<button hx-get="/user/{{.ID}}/edit" hx-target="#user-{{.ID}}" hx-swap="outerHTML">
		Edit
	</button>
</div>`

	t, err := template.New("user-view").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, user); err != nil {
		http.Error(w, "Failed to render", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}
