package gosocial

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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
	PWD   string
	APP   string
	ID    int
	Slug  string
	Name  string
	Title string
	Text  string
}

func ApiHandler(w http.ResponseWriter, r *http.Request, AdminHASH string) { //THIS ONE IS WORKING WITH jsondataRequests
	startAPI2 := time.Now()

	decoder := json.NewDecoder(r.Body)
	var jsondata API2STRUCT
	errSearch := decoder.Decode(&jsondata)

	switch {
	case errSearch == io.EOF: //if API reqest is empty, the server will return the js code
		fmt.Fprintf(w, "<h1>THIS IS UNDER DEV</h1>")
		fmt.Println("SERVE JS")
		return
	case errSearch != nil:
		bar := `{"Status":"` + errSearch.Error() + `"}`
		fmt.Fprintf(w, bar)
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

	fmt.Println("Api2Handler:", time.Since(startAPI2))
}

func WriteComment(w http.ResponseWriter, jsondata API2STRUCT) {
	db.Exec("INSERT INTO gosocial_comments(slug,Name,Title,Text) VALUES(?,?,?,?)", jsondata.Slug, jsondata.Name, jsondata.Title, jsondata.Text)

	fmt.Fprintf(w, `{"Status": "OK"}`)
	fmt.Println("Api2Handler-WriteComment")
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("\033[0;31m", err, "\033[0m")
		err = nil
	}
}
