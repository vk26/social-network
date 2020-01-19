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
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	DB   *sql.DB
	Tmpl *template.Template
}

var Users []User

var (
	key          = []byte(os.Getenv("SESSIONS_KEY"))
	sessionStore = sessions.NewCookieStore(key)
)

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

func (h *Handler) LoginForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("frontend/templates/login.html"))
	tmpl.Execute(w, nil)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user := &User{}
	row := h.DB.QueryRow("SELECT id, email, password_hash FROM users WHERE email = ?", email)
	row.Scan(&user.Id, &user.Email, &user.PasswordHash)

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
	}

	session, _ := sessionStore.Get(r, "social_app")
	session.Values["userID"] = user.Id
	session.Save(r, w)

	userIDStr := strconv.FormatInt(int64(user.Id), 10)
	http.Redirect(w, r, "/users/"+userIDStr, http.StatusFound)
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

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Authentication check ...")
		session, _ := sessionStore.Get(r, "social_app")
		userID, ok := session.Values["userID"].(int)
		if !ok || userID == 0 {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
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

	siteMux := mux.NewRouter().PathPrefix("/").Subrouter()
	siteMux.HandleFunc("/login", handlers.LoginForm).Methods("GET")
	siteMux.HandleFunc("/login", handlers.Login).Methods("POST")
	siteMux.HandleFunc("/signup", handlers.SignupForm).Methods("GET")
	siteMux.HandleFunc("/signup", handlers.Signup).Methods("POST")

	siteMux.HandleFunc("/", handlers.Home)

	userMux := mux.NewRouter()
	userMux.HandleFunc("/users/{id}", handlers.UserPage)
	userHandler := authMiddleware(userMux)
	siteMux.Handle("/users/{id}", userHandler)

	assetsHandler := http.StripPrefix("/data/", http.FileServer(http.Dir("frontend/assets")))
	siteMux.PathPrefix("/data/").Handler(assetsHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: siteMux,
	}

	fmt.Println("Server is listening ...")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
