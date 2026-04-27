// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "postgres"
	dbname   = "test"
)

type zap_user struct {
	Name string `json:"name"`
	Phone int `json:"phone_number"`
}

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

// func getUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
// 	insertUser := `insert into "users"("name", "phone") values($1, $2)`
// 	_, e := db.Exec(insertUser, "Bob", 100)
// 	if e != nil {
// 		panic(e)
// 	}
//
// 	row, err := db.Query(`select "name" from "users" where "id"=$1`, 13)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	var name string
// 	row.Next()
// 	err = row.Scan(&name)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	fmt.Println(name)
//
// 	user_json := &zap_user{
// 		Name: name,
// 		Phone: 1,
// 	}
// 	user_json_marsh, _ := json.Marshal(user_json)
//
// 	w.Header().Set("Content-Type", "application/json")
// 	fmt.Fprint(w, string(user_json_marsh))
// }

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/", serveHome)
	// http.HandleFunc("/getUser", func(w http.ResponseWriter, r *http.Request) {
	// 	getUser(db, w, r)
	// })
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
