// Copyright 2018 The hhwe. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package webgo implements a simple web application package.
package main

import (
	"gopkg.in/mgo.v2"
	"net/http"
)

func main() {
	// connect to the database
	db, err := mgo.Dial("localhost")
	if err != nil {
		logger.Fatal("cannot dial mongo", err)
	}
	defer db.Close() // clean up when we’re done

	db.SetMode(mgo.Monotonic, true)
	ensureIndex(db)

	// testing
	http.HandleFunc("/test/error",
		Chain(ErrorTest, Method("GET", "POST"), DBSession(db), Logging()))
	http.HandleFunc("/test/comments",
		Chain(Comment, Method("GET", "POST"), DBSession(db), Logging()))

	// home page
	http.HandleFunc("/",
		Chain(Index, Method("GET"), DBSession(db), Logging()))

	// login and sign in
	http.HandleFunc("/register",
		Chain(Register, Method("GET", "POST"), DBSession(db), Logging()))
	http.HandleFunc("/login",
		Chain(Login, Method("GET", "POST"), DBSession(db), Logging()))
	http.HandleFunc("/logout",
		Chain(Logout, Method("GET", "POST"), DBSession(db), Logging()))

	// Restful API
	http.HandleFunc("/users", Chain(UserApi, Auth(),
		Method("GET", "POST", "PUT", "DELETE"), DBSession(db), Logging()))

	// shopping cart settlement
	http.HandleFunc("/cart",
		Chain(Cart, Method("GET", "POST"), Logging()))

	// handle static file with file server
	http.Handle("/static/", http.StripPrefix(
		"/static/", http.FileServer(http.Dir("static/"))))

	logger.Println(` 
* Environment: production
  WARNING: Do not use the development server in a production environment.
* Debug mode: off
* Running on http://127.0.0.1:8000/ (Press CTRL+C to quit)
`)
	logger.Fatal(http.ListenAndServe(":8000", nil))

}
