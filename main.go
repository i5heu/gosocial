package gosocial

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/microcosm-cc/bluemonday"
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

type GetCommentsResults struct {
	Id        int
	Name      template.HTML
	Title     template.HTML
	Text      template.HTML
	Upvotes   int
	Downvotes int
	Color     string
}

type GetCommentsResultsArray struct { //GetCommentsResults Array
	Comments []GetCommentsResults
}

var GetCommentsTemplate *template.Template

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
	case "GetComments":
		GetComments(w, jsondata)
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

func GetComments(w http.ResponseWriter, jsondata API2STRUCT) {

	ids, err := db.Query("SELECT  ID, Name, Title, Text, upvotes, downvotes FROM gosocial_comments WHERE slug = ? AND ModRelease = '1' ORDER BY submitTime DESC LIMIT 1000")
	defer ids.Close()
	checkErr(err)
	var TMP []GetCommentsResults
	var ClassSwitch bool = true
	var ClassTMP = "ProjectTableDark"

	for ids.Next() {
		var id int
		var name string
		var title string
		var text string
		var upvotes int
		var downvotes int
		_ = ids.Scan(&id, &name, &title, &text, &upvotes, &downvotes)
		checkErr(err)

		NameTMP := template.HTML(bluemonday.UGCPolicy().SanitizeBytes([]byte(name)))
		TitleTMP := template.HTML(bluemonday.UGCPolicy().SanitizeBytes([]byte(title)))
		TextTMP := template.HTML(bluemonday.UGCPolicy().SanitizeBytes([]byte(text)))

		if ClassSwitch == true {
			ClassTMP = "ProjectTableDark"
			ClassSwitch = false
		} else {
			ClassTMP = "ProjectTableBright"
			ClassSwitch = true
		}

		TMP = append(TMP, GetCommentsResults{id, NameTMP, TitleTMP, TextTMP, upvotes, downvotes, ClassTMP})
	}

	lists := GetCommentsResultsArray{TMP}

	GetCommentsTemplate.Execute(w, lists)

	fmt.Println("Api2Handler-GetComment")
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("\033[0;31m", err, "\033[0m")
		err = nil
	}
}
