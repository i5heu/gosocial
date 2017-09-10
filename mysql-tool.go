package gosocial

import "database/sql"

func MysqlRow_BACK_STRING5(db *sql.DB, comand string) (back string, back2 string, back3 string, back4 string, back5 string) {
	row, err := db.Query(comand)
	checkErr(err)
	defer row.Close()

	row.Next()
	_ = row.Scan(&back, &back2, &back3, &back4, &back5)

	return
}
