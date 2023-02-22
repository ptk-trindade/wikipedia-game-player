package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb"
)

var db *sql.DB

func getEnvVars() {
	err := godotenv.Load("credentials.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func connectDB() *sql.DB {
	server := strings.TrimSpace(os.Getenv("DB_SERVER"))
	port, _ := strconv.Atoi(strings.TrimSpace(os.Getenv("DB_PORT")))
	user := strings.TrimSpace(os.Getenv("DB_USER"))
	password := strings.TrimSpace(os.Getenv("DB_PASSWORD"))
	database := strings.TrimSpace(os.Getenv("DB_DATABASE"))

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	// fmt.Println(connString)
	var err error

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")

	return db
}

/*
-- input:
pageSrc: page source
pageDstList: list of pages destination [][2]string{pageName, linkWord}
*/
func insertPageLinks(pageSrc string, pageDstList [][2]string) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("No connection to database", err.Error())
	}

	// Create insert values list
	// valuesList := make([]string, 0, len(pageDstList))
	for _, pageDst := range pageDstList {
		// values := fmt.Sprintf("(@LC, @pageSrc, @pageName, @linkWord)", LC, pageSrc, pageDst[0], pageDst[1])

		// Create insert query
		tsql := "INSERT INTO wikigame.dbo.Link (LanguageCode, PageSrc, PageDst, LinkWord) VALUES (@LC, @pageSrc, @pageName, @linkWord);"

		// Execute query
		_, err = db.ExecContext(ctx, tsql, sql.Named("LC", LC), sql.Named("pageSrc", pageSrc), sql.Named("pageName", pageDst[0]), sql.Named("linkWord", pageDst[1]))
		if err != nil {
			fmt.Println(tsql)
			log.Fatal("Error execution query: ", err.Error())
		}
	}
}

/*
	NOT BEING USED / UNSAFE

-- input:
pageSrc: page source
pageDstList: list of pages destination [][2]string{pageName, linkWord}
This Method is faster but it is not safe against SQL injection.
*/
func insertPageLinksFast(pageSrc string, pageDstList [][2]string) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("No connection to database", err.Error())
	}

	// Create insert values list
	valuesList := make([]string, 0, len(pageDstList))
	for _, pageDst := range pageDstList {
		valuesList = append(valuesList, fmt.Sprintf("('%s', '%s', '%s', '%s')", LC, pageSrc, pageDst[0], pageDst[1]))
	}

	// Create insert query
	tsql := fmt.Sprintf("INSERT INTO wikigame.dbo.Link (LanguageCode, PageSrc, PageDst, LinkWord) VALUES %s;",
		strings.Join(valuesList, ","))

	// Execute query
	_, err = db.ExecContext(ctx, tsql)
	if err != nil {
		fmt.Println(tsql)
		log.Fatal("Error execution query: ", err.Error())
	}
}

func selectPageLinks(pageName string) *sql.Rows {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal("No connection to database", err.Error())
	}

	// Create select query
	tsql := "SELECT PageDst, LinkWord FROM wikigame.dbo.Link WHERE PageSrc = @pageName;"

	// Execute query
	rows, err := db.QueryContext(ctx, tsql, sql.Named("pageName", pageName))
	if err != nil {
		fmt.Println(tsql)
		log.Fatal("Error execution query", err.Error())
	}

	return rows
}
