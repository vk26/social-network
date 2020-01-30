package models

import (
	"database/sql"
	"time"
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

func (u *User) GetUserByID(db *sql.DB) error {
	row := db.QueryRow("SELECT id, name, surname, birthday, city, about, email FROM users WHERE id = ?", u.Id)

	err := row.Scan(&u.Id, &u.Name, &u.Surname, &u.Birthday, &u.City, &u.About, &u.Email)
	return err
}

func (u *User) GetUserByEmail(db *sql.DB) error {
	row := db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = ?", u.Email)

	err := row.Scan(&u.Id, &u.Email, &u.PasswordHash)
	return err
}

func (u *User) CreateUser(db *sql.DB) error {
	result, err := db.Exec(
		"INSERT INTO users (`name`, `surname`, `birthday`, `city`, `about`, `email`, `password_hash`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		u.Name,
		u.Surname,
		u.Birthday,
		u.City,
		u.About,
		u.Email,
		u.PasswordHash,
		time.Now().Format("2006-01-02 15:04:05"),
		time.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	u.Id = int(id)
	return err
}

func GetUsers(db *sql.DB, count, start int) ([]User, error) {
	rows, err := db.Query(
		"SELECT id, name, surname, birthday, city, about, email FROM users LIMIT ? OFFSET ?",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Name, &u.Surname, &u.Birthday, &u.City, &u.About, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func SearchUsers(db *sql.DB, nameSubstr string, count, start int) ([]User, error) {
	wildcardSubstr := nameSubstr + "%"
	rows, err := db.Query(
		"SELECT id, name, surname, birthday, city, about, email FROM users WHERE name LIKE ? OR surname LIKE ? LIMIT ? OFFSET ?",
		wildcardSubstr, wildcardSubstr, count, start)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Name, &u.Surname, &u.Birthday, &u.City, &u.About, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
