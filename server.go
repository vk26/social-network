package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var Users []User

type User struct {
	Name         string
	Surname      string
	Login        string
	Age          uint
	City         string
	Interests    string
	Email        string
	PasswordHash string
}

type signupForm struct {
	Name      string
	Surname   string
	Login     string
	Age       uint
	City      string
	Interests string
	Email     string
	Password  string
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

func UserRegistrate(data signupForm) (User, error) {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	user := User{
		Name:         data.Name,
		Surname:      data.Surname,
		Login:        data.Login,
		Age:          data.Age,
		City:         data.City,
		Interests:    data.Interests,
		Email:        data.Email,
		PasswordHash: string(passwordHash),
	}

	Users = append(Users, user)
	return user, nil
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("frontend/html/login.html"))
	tmpl.Execute(w, nil)
}

func handleLoginPost(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user := UserAuthenticate(email, password)
	if (User{}) == user {
		fmt.Fprintln(w, "Invalid email or password!")
	} else {
		tmpl := template.Must(template.ParseFiles("frontend/html/user_page.html"))
		tmpl.Execute(w, user)
	}
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("frontend/html/signup.html"))
	tmpl.Execute(w, nil)
}

func handleSignupPost(w http.ResponseWriter, r *http.Request) {
	age, _ := strconv.ParseUint(r.FormValue("age"), 10, 32)
	signupData := signupForm{
		Name:      r.FormValue("name"),
		Surname:   r.FormValue("surname"),
		Login:     r.FormValue("login"),
		Age:       uint(age),
		City:      r.FormValue("city"),
		Interests: r.FormValue("interests"),
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
	}
	user, err := UserRegistrate(signupData)
	if err != nil {
		fmt.Fprintln(w, "Can not create user: ", err)
	} else {
		tmpl := template.Must(template.ParseFiles("frontend/html/user_page.html"))
		tmpl.Execute(w, user)
	}

}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Home page")
}

func handleUserPage(w http.ResponseWriter, r *http.Request) {
	userLogin := mux.Vars(r)["user"]

	for _, u := range Users {
		if u.Login == userLogin {
			tmpl := template.Must(template.ParseFiles("frontend/html/user_page.html"))
			tmpl.Execute(w, u)
		}
	}
	fmt.Fprintln(w, "Can not find user with login: ", userLogin)
}

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/login", handleLogin).Methods("GET")
	mux.HandleFunc("/login", handleLoginPost).Methods("POST")
	mux.HandleFunc("/signup", handleSignup).Methods("GET")
	mux.HandleFunc("/signup", handleSignupPost).Methods("POST")
	mux.HandleFunc("/{user}", handleUserPage)
	mux.HandleFunc("/", handleHome)
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
