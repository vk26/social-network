package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"syreclabs.com/go/faker"
)

var (
	batchSize        = 10000
	insertIterations = 100
)

func main() {
	db, err := sql.Open("mysql", os.Getenv("SOCIAL_APP_MYSQL_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Start inserting data ...")
	for j := 0; j < insertIterations; j++ {
		records := []string{}
		recordsRaw := ""
		for i := 0; i < batchSize; i++ {
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
			records = append(records, record)
		}
		recordsRaw = strings.Join(records[:], ",")
		sqlRaw := "INSERT INTO users (`name`, `surname`, `birthday`, `city`, `about`, `avatar`, `email`, `password_hash`, `created_at`, `updated_at`) VALUES "
		sqlRaw += recordsRaw
		_, err := db.Exec(sqlRaw)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Data insertion is completed!")

}

func escapeSingleQuotes(str string) string {
	return strings.Replace(str, "'", "\\'", -1)
}
