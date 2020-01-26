package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"social-network/models"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
	Tmpl   *template.Template
}

var (
	sessionKey   = []byte(os.Getenv("SESSIONS_KEY"))
	sessionStore = sessions.NewCookieStore(sessionKey)
)

type ctxKey string

const (
	currentUserKey ctxKey = "currentUserKey"
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

func init() {
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24 * 15, // 15 days
		HttpOnly: true,
	}
}

func (a *App) Initialize(dbDriver, dbUser, dbPassword, dbName string) {
	var err error
	a.DB, err = sql.Open(dbDriver, dbUser+":"+dbPassword+"@/"+dbName)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.Tmpl = template.Must(template.ParseGlob("frontend/templates/*"))
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	siteMux := mux.NewRouter().PathPrefix("/").Subrouter()
	siteMux.HandleFunc("/login", a.LoginForm).Methods("GET")
	siteMux.HandleFunc("/login", a.Login).Methods("POST")
	siteMux.HandleFunc("/logout", a.Logout).Methods("GET")
	siteMux.HandleFunc("/signup", a.SignupForm).Methods("GET")
	siteMux.HandleFunc("/signup", a.Signup).Methods("POST")

	userMux := siteMux.PathPrefix("/").Subrouter()
	userMux.HandleFunc("/users", a.UsersList).Methods("GET")
	userMux.HandleFunc("/users/{id:[0-9]+}", a.UserPage).Methods("GET")
	userMux.Use(a.authMiddleware)

	siteMux.Use(a.getCurrentUserMiddleware)

	assetsHandler := http.StripPrefix("/data/", http.FileServer(http.Dir("frontend/assets")))
	siteMux.PathPrefix("/data/").Handler(assetsHandler)

	a.Router = siteMux
}

func (a *App) LoginForm(w http.ResponseWriter, r *http.Request) {
	err := a.Tmpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user := models.User{Email: email}
	user.GetUserByEmail(a.DB)
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

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, "social_app")
	session.Values["userID"] = nil
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (a *App) SignupForm(w http.ResponseWriter, r *http.Request) {
	a.Tmpl.ExecuteTemplate(w, "signup.html", nil)
}

func (a *App) Signup(w http.ResponseWriter, r *http.Request) {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 14)
	user := models.User{
		Name:         r.FormValue("name"),
		Surname:      r.FormValue("surname"),
		Birthday:     r.FormValue("birthday"),
		City:         r.FormValue("city"),
		About:        r.FormValue("about"),
		Email:        r.FormValue("email"),
		PasswordHash: string(passwordHash),
	}
	err := user.CreateUser(a.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := sessionStore.Get(r, "social_app")
	session.Values["userID"] = int(user.Id)
	session.Save(r, w)

	userIDStr := strconv.FormatInt(int64(user.Id), 10)
	http.Redirect(w, r, "/users/"+userIDStr, http.StatusFound)

}

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Home page")
}

func (a *App) UserPage(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	user := models.User{Id: id}
	err := user.GetUserByID(a.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"user":        user,
		"currentUser": context.Get(r, currentUserKey),
	}
	a.Tmpl.ExecuteTemplate(w, "user_page.html", data)
}

func (a *App) UsersList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	page, _ := strconv.Atoi(params["page"])
	count, _ := strconv.Atoi(params["count"])
	if count == 0 {
		count = 15
	}
	start := page * count
	users, err := models.GetUsers(a.DB, start, count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"users":       users,
		"currentUser": context.Get(r, currentUserKey),
	}
	a.Tmpl.ExecuteTemplate(w, "users_list.html", data)
}

func (a *App) authMiddleware(next http.Handler) http.Handler {
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

func (a *App) getCurrentUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Getting current user ...")
		session, _ := sessionStore.Get(r, "social_app")
		userID, ok := session.Values["userID"].(int)
		if ok && userID != 0 {
			user := models.User{Id: userID}
			err := user.GetUserByID(a.DB)
			if err == nil {
				context.Set(r, currentUserKey, user)
			}
		}
		fmt.Println(context.Get(r, currentUserKey))
		next.ServeHTTP(w, r)
	})
}

func main() {
	a := App{}
	a.Initialize(
		"mysql",
		os.Getenv("DB_MYSQL_USER"),
		os.Getenv("DB_MYSQL_PASSWORD"),
		"social_dev",
	)

	a.Run(":8080")
}
