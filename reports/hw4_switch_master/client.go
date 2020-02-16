package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"syreclabs.com/go/faker"
)

func main() {
	db, err := sql.Open("mysql", os.Getenv("SOCIAL_APP_MYSQL_DSN"))
	counterSuccess := 0
	counterFailure := 0
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		db.Close()
		log.Println("Count of success insertion in DB:", counterSuccess)
		log.Println("Count of failure insertion in DB:", counterFailure)
		os.Exit(1)
	}()

	log.Println("Start inserting data ...")
	for {
		record := fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')",
			escapeSingleQuotes(faker.Name().FirstName()),
			escapeSingleQuotes(faker.Name().LastName()),
			faker.Date().Birthday(18, 60).Format("2006-01-02"),
			escapeSingleQuotes(faker.Address().City()),
			strings.Join(faker.Lorem().Words(3)[:], ", "),
			faker.Avatar().String(),
			faker.Internet().SafeEmail(),
			"$2a$14$xSWEm06Ro.735D/cNTKHE.VcQ1pB5lgEPFmNJItz.yi5nGO8Z0bse",
			time.Now().Format("2006-01-02 15:04:05"),
			time.Now().Format("2006-01-02 15:04:05"),
		)

		sqlRaw := "INSERT INTO users (`name`, `surname`, `birthday`, `city`, `about`, `avatar`, `email`, `password_hash`, `created_at`, `updated_at`) VALUES "
		sqlRaw += record
		_, err := db.Exec(sqlRaw)

		if err != nil {
			counterFailure++
			log.Println("Insertion is failure. Failure counter:", counterFailure)
		} else {
			counterSuccess++
			log.Println("Insertion is succeed. Success counter:", counterSuccess)
		}
	}
}

func escapeSingleQuotes(str string) string {
	return strings.Replace(str, "'", "\\'", -1)
}
