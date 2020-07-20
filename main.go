package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type TblProfile struct {
	Id         int    `json:"id, omitempty"`
	Localroute string `json:"local_route, omitempty"`
}

func main() {
	router := mux.NewRouter()
	//endpoints
	router.HandleFunc("/pictures", GetPictures).Methods("GET")
	router.HandleFunc("/pictures/{id}", GetPictureById).Methods("GET")
	router.HandleFunc("/pictures", PostPicture).Methods("POST")
	router.HandleFunc("/pictures/{id}", PutPicture).Methods("PUT")
	router.HandleFunc("/pictures/{id}", DeletePicture).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))
}

func GetPictures(w http.ResponseWriter, request *http.Request) {
	db := DbConn()
	selDB, err := db.Query("SELECT * FROM tblprofile ORDER BY id DESC")
	CheckForError(err)

	tbl := TblProfile{}
	res := []TblProfile{}
	for selDB.Next() {
		var id int
		var local_route string
		err = selDB.Scan(&id, &local_route)
		if err != nil {
			panic(err.Error())
		}
		tbl.Id = id
		tbl.Localroute = local_route
		res = append(res, tbl)
	}
	defer db.Close()

	fmt.Println(res)
	json.NewEncoder(w).Encode(res)
}

func GetPictureById(w http.ResponseWriter, request *http.Request) {
	db := DbConn()
	params := mux.Vars(request)

	selDB, err := db.Query("SELECT * FROM tblprofile WHERE id = ?", params["id"])
	CheckForError(err)

	tbl := TblProfile{}
	res := []TblProfile{}
	for selDB.Next() {
		var id int
		var local_route string
		err = selDB.Scan(&id, &local_route)
		if err != nil {
			panic(err.Error())
		}
		tbl.Id = id
		tbl.Localroute = local_route
		res = append(res, tbl)
	}
	defer db.Close()

	fmt.Println(res)
	json.NewEncoder(w).Encode(res)
}

func PostPicture(w http.ResponseWriter, request *http.Request) {
	db := DbConn()

	var newProfile TblProfile
	_ = json.NewDecoder(request.Body).Decode(&newProfile)

	stmt, err := db.Prepare("INSERT tblprofile SET local_route=?")
	CheckForError(err)

	//Activar linea 89 en caso de querer usar POST con la encriptacion
	//newProfile.Localroute = RenameWithMD5(newProfile.Localroute)
	res, err := stmt.Exec(newProfile.Localroute)
	CheckForError(err)

	id, err := res.LastInsertId()
	CheckForError(err)

	newProfile.Id = int(id)

	defer db.Close()
	json.NewEncoder(w).Encode(newProfile)
}

func PutPicture(w http.ResponseWriter, request *http.Request) {
	db := DbConn()

	var newProfile TblProfile
	_ = json.NewDecoder(request.Body).Decode(&newProfile)

	params := mux.Vars(request)

	selDB, err := db.Query("SELECT * FROM tblprofile WHERE id = ?", params["id"])
	CheckForError(err)

	tbl := TblProfile{}
	res := []TblProfile{}
	for selDB.Next() {
		var id int
		var local_route string
		err = selDB.Scan(&id, &local_route)
		if err != nil {
			panic(err.Error())
		}
		tbl.Id = id
		tbl.Localroute = local_route
		res = append(res, tbl)
	}

	if len(res) == 0 {
		stmt, err := db.Prepare("INSERT tblprofile SET local_route=?")
		CheckForError(err)
		newProfile.Localroute = RenameWithMD5(newProfile.Localroute)
		res, err := stmt.Exec(newProfile.Localroute)
		CheckForError(err)
		id, err := res.LastInsertId()
		CheckForError(err)
		newProfile.Id = int(id)
	} else {
		stmt, err := db.Prepare("UPDATE tblprofile SET local_route=? WHERE id=?")
		CheckForError(err)
		newProfile.Localroute = RenameWithMD5(newProfile.Localroute)
		res, err := stmt.Exec(newProfile.Localroute, params["id"])
		CheckForError(err)
		id, err := res.LastInsertId()
		CheckForError(err)
		newProfile.Id = int(id)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(newProfile)
}

func DeletePicture(w http.ResponseWriter, request *http.Request) {
	db := DbConn()
	params := mux.Vars(request)

	stmt, err := db.Prepare("DELETE FROM tblprofile WHERE id=?")
	CheckForError(err)

	res, err := stmt.Exec(params["id"])
	CheckForError(err)
	affect, err := res.RowsAffected()
	CheckForError(err)

	defer db.Close()
	json.NewEncoder(w).Encode(affect)
}
