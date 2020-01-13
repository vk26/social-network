package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	DB   *sql.DB
	Tmpl *template.Template
}

var Users []User

type User struct {
	Id           int
	Name         string
	Surname      string
	Birthday     string
	City         string
	About        string
	Email        string
	PasswordHash string
}

func UserAuthenticate(email string, password string) User {
	for _, u := range Users {
		if u.Email == email {
			err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
			if err != nil {
				return User{}
			}
			return u
		}
	}
	return User{}
}

func (h *Handler) LoginForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("frontend/templates/login.html"))
	tmpl.Execute(w, nil)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user := UserAuthenticate(email, password)
	if (User{}) == user {
		fmt.Fprintln(w, "Invalid email or password!")
	} else {
		tmpl := template.Must(template.ParseFiles("frontend/templates/user_page.html"))
		tmpl.Execute(w, user)
	}
}

func (h *Handler) SignupForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("frontend/templates/signup.html", "frontend/templates/layouts/base.html"))
	tmpl.Execute(w, nil)
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 14)
	result, err := h.DB.Exec(
		"INSERT INTO users (`name`, `surname`, `birthday`, `city`, `about`, `email`, `password_hash`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		r.FormValue("name"),
		r.FormValue("surname"),
		r.FormValue("birthday"),
		r.FormValue("city"),
		r.FormValue("about"),
		r.FormValue("email"),
		string(passwordHash),
		time.Now().Format("2006-01-02 15:04:05"),
		time.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		panic(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	userIDStr := strconv.FormatInt(int64(id), 10)
	http.Redirect(w, r, "/users/"+userIDStr, http.StatusFound)

}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Home page")
}

func (h *Handler) UserPage(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	user := &User{}
	row := h.DB.QueryRow("SELECT id, name, surname, birthday, city, about, email FROM users WHERE id = ?", id)

	err := row.Scan(&user.Id, &user.Name, &user.Surname, &user.Birthday, &user.City, &user.About, &user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl := template.Must(template.ParseFiles("frontend/templates/user_page.html", "frontend/templates/layouts/base.html"))
	tmpl.Execute(w, user)
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := os.Getenv("DB_MYSQL_USER")
	dbPass := os.Getenv("DB_MYSQL_PASSWORD")
	dbName := "social_dev"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {
	db := dbConn()

	handlers := &Handler{
		DB:   db,
		Tmpl: template.Must(template.ParseGlob("frontend/templates/*/*")),
	}

	mux := mux.NewRouter()
	mux.HandleFunc("/login", handlers.LoginForm).Methods("GET")
	mux.HandleFunc("/login", handlers.Login).Methods("POST")
	mux.HandleFunc("/signup", handlers.SignupForm).Methods("GET")
	mux.HandleFunc("/signup", handlers.Signup).Methods("POST")
	mux.HandleFunc("/users/{id}", handlers.UserPage)
	mux.HandleFunc("/", handlers.Home)

	assetsHandler := http.StripPrefix("/data/", http.FileServer(http.Dir("frontend/assets")))
	mux.PathPrefix("/data/").Handler(assetsHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server is listening ...")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
