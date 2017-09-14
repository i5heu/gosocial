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
	"github.com/russross/blackfriday"
)

var db *sql.DB

func Init(dbTMP *sql.DB) {
	db = dbTMP
	CreateTable()
	fmt.Println("HALLO - A")

}

type API2STRUCT struct {
	PWD         string
	APP         string
	ID          int
	Slug        string
	Name        string
	Title       string
	Text        string
	ModerateNum int
}

type GetCommentsResults struct {
	Id         int
	Slug       template.HTML
	Name       template.HTML
	Title      template.HTML
	Text       template.HTML
	Upvotes    int
	Downvotes  int
	ModRelease int
}

type GetCommentsResultsArray struct { //GetCommentsResults Array
	Status   string
	Comments []GetCommentsResults
}

var GetCommentsTemplate *template.Template

func ApiHandler(w http.ResponseWriter, r *http.Request, AdminHASH string) (AppMethod string, Title string, Text string) { //THIS ONE IS WORKING WITH jsondataRequests
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
		GetComments(w, r, jsondata, AdminHASH)
	case "WriteComment":
		Title, Text = WriteComment(w, jsondata)
	case "ModerateComment":
		ModerateComment(w, r, jsondata, AdminHASH)
	default:
		fmt.Fprintf(w, `{"Status":"ERROR - NO METHOD"}`)
	}

	fmt.Println("Api2Handler:", time.Since(startAPI2))
	AppMethod = jsondata.APP
	return
}

func ModerateComment(w http.ResponseWriter, r *http.Request, jsondata API2STRUCT, AdminHASH string) {

	if jsondata.PWD != AdminHASH {
		fmt.Fprintf(w, `{"Status": "NOT LOGED IN"}`)
		fmt.Println("Api2Handler-ModerateComment-NOTLOGEDIN")
		return
	}
	fmt.Println(jsondata.ModerateNum, jsondata.ID)
	db.Exec("UPDATE `gosocial_comments` SET `ModRelease` = ? WHERE `gosocial_comments`.`ID` = ?", jsondata.ModerateNum, jsondata.ID)

	fmt.Fprintf(w, `{"Status": "OK"}`)
	fmt.Println("Api2Handler-ModerateComment")
}

func WriteComment(w http.ResponseWriter, jsondata API2STRUCT) (Title string, Text string) {
	db.Exec("INSERT INTO gosocial_comments(slug,Name,Title,Text) VALUES(?,?,?,?)", jsondata.Slug, jsondata.Name, jsondata.Title, jsondata.Text)

	Title = jsondata.Title
	Text = jsondata.Text
	fmt.Fprintf(w, `{"Status": "OK"}`)
	fmt.Println("Api2Handler-WriteComment")
	return
}

func GetComments(w http.ResponseWriter, r *http.Request, jsondata API2STRUCT, AdminHASH string) {

	var ids *sql.Rows
	var err error
	if jsondata.PWD == AdminHASH {
		ids, err = db.Query("SELECT  ID, slug ,Name, Title, Text, upvotes, downvotes, ModRelease FROM gosocial_comments WHERE ModRelease = 0 ORDER BY submitTime DESC LIMIT 1000")
	} else {

		ids, err = db.Query("SELECT  ID, slug, Name, Title, Text, upvotes, downvotes, ModRelease FROM gosocial_comments WHERE slug = ? AND ModRelease = 1 ORDER BY submitTime DESC LIMIT 1000", jsondata.Slug)
	}
	defer ids.Close()
	checkErr(err)
	var TMP []GetCommentsResults

	for ids.Next() {
		var id int
		var slug string
		var name string
		var title string
		var text string
		var upvotes int
		var downvotes int
		var ModRelease int
		_ = ids.Scan(&id, &slug, &name, &title, &text, &upvotes, &downvotes, &ModRelease)
		checkErr(err)

		NameTMP := template.HTML(bluemonday.UGCPolicy().SanitizeBytes([]byte(name)))
		TitleTMP := template.HTML(bluemonday.UGCPolicy().SanitizeBytes([]byte(title)))
		TextTMP := template.HTML(blackfriday.MarkdownCommon(bluemonday.UGCPolicy().SanitizeBytes([]byte(text))))
		SlugTMP := template.HTML(bluemonday.UGCPolicy().SanitizeBytes([]byte(slug)))

		TMP = append(TMP, GetCommentsResults{id, SlugTMP, NameTMP, TitleTMP, TextTMP, upvotes, downvotes, ModRelease})
	}

	lists := GetCommentsResultsArray{"OK", TMP}

	jData, err := json.Marshal(lists)
	checkErr(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)

	fmt.Println("Api2Handler-GetComment")
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("\033[0;31m", err, "\033[0m")
		err = nil
	}
}
