package gosocial

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var db *sql.DB

func Init(dbTMP *sql.DB) {
	db = dbTMP
	CreateTable()
	fmt.Println("HALLO - A")

}

type API2STRUCT struct {
	PWD    string
	APP    string
	ID     int
	Title1 string
	Text1  string
}

func ApiHandler(w http.ResponseWriter, r *http.Request, AdminHASH string) { //THIS ONE IS WORKING WITH jsondataRequests
	startAPI2 := time.Now()

	APILogin := false

	decoder := json.NewDecoder(r.Body)
	var jsondata API2STRUCT
	errSearch := decoder.Decode(&jsondata)
	if errSearch != nil {
		fmt.Fprintf(w, `{"Status":"ERROR"}`)
		fmt.Println(errSearch)
		checkErr(errSearch)
		return
	}

	switch jsondata.APP {
	case "WriteComment":
		WriteComment(w, jsondata)
	default:
		fmt.Fprintf(w, `{"Status":"ERROR - NO METHOD"}`)
	}

	fmt.Println("Api2Handler:", time.Since(startAPI2), APILogin)
}

func WriteComment(w http.ResponseWriter, jsondata API2STRUCT) {
	fmt.Println("Api2Handler")
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("\033[0;31m", err, "\033[0m")
		err = nil
	}
}
