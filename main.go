package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"reflect"

	"github.com/mattn/go-sqlite3"
)

//go:embed _webroot
var webrootFS embed.FS

func main() {
	sql.Register("sqlite3_with_spatialite",
		&sqlite3.SQLiteDriver{
			Extensions: []string{"mod_spatialite"},
		})
	db, err := sql.Open("sqlite3_with_spatialite", "./spatia.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	// _, err = db.Exec(`SELECT InitSpatialMetaData()`)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	rows, err := db.Query(`select a.objectid,a.wadmkk from spatia a, spatia b
		where b.wadmkk=? and Touches(a.geometry, b.geometry)`,
		"Kab. Boyolali",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			fmt.Println(errRow.Error())
		}
	}()

	types, err := rows.ColumnTypes()
	if err != nil {
		log.Fatal(err)
	}
	cols := make([]string, len(types))
	resPtr := make([]any, len(types))
	for i := 0; i < len(types); i++ {
		resPtr[i] = reflect.New(types[i].ScanType()).Interface()
		cols[i] = types[i].Name()
	}

	results := []map[string]any{}

	for rows.Next() {
		if err = rows.Scan(
			resPtr...,
		); err != nil {
			log.Fatal(err)
		}
		result := make(map[string]any)
		for i := 0; i < len(types); i++ {
			// switch resPtr[i].(type) {
			// case *sql.NullString:
			// 	fmt.Println("OK")
			// }
			// result[cols[i]] = reflect.Indirect(reflect.ValueOf(resPtr[i])).Interface()
			result[cols[i]] = reflect.ValueOf(resPtr[i]).Elem().Interface()
		}
		results = append(results, result)
	}
	fmt.Println(results)

	mux := http.NewServeMux()

	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("OK"))
	// })

	uiFs, err := fs.Sub(webrootFS, "_webroot")
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(uiFs))))

	server := http.Server{Addr: ":8080", Handler: mux}
	log.Fatal(server.ListenAndServe())
}
