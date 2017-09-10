package gosocial

import (
	"database/sql"
	"fmt"
	"net/http"
)

var db *sql.DB

func Init(dbTMP *sql.DB) {
	db = dbTMP
	CreateTable()
	fmt.Println("HALLO - A")

}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HALLO")
	fmt.Println("HALLO - A")
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("\033[0;31m", err, "\033[0m")
		err = nil
	}
}
