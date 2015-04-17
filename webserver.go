package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type WebConfig interface {
	Port() string
	Certificate() string
	PrivateKey() string
}

func RunServer(conf WebConfig) {
	log.Println("Webserver starting on port " + conf.Port())
	http.HandleFunc("/", handler)
	err := http.ListenAndServeTLS(":"+conf.Port(), conf.Certificate(), conf.PrivateKey(), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path[1:]


	switch path {
	case "css":
		showFile(w, "css/styles.css")
	case "favicon.ico":
		showFile(w, "images/favicon.ico")
	case "savekey":
		saveKey(w,r)
	case "oauth2callback":
		handleOauthCallback(w, r)
	case "add":
		if needskey(w) { return }
		showAddForm(w, r)
	case "save":
		if needskey(w) { return }
		saveAddForm(w, r)
	case "delete":
		if needskey(w) { return }
		showDeleteForm(w, r, "")
	case "confirmdelete":
		if needskey(w) { return }
		performDelete(w, r)
	default:
		if needskey(w) { return }
		showUserList(w, r)
	}
}

func needskey(w http.ResponseWriter) bool {
	key := Key()
	if len(key) < 5 {
		showFile(w, "html/keyform.html")
		return true
	}
	return false
}

func saveKey(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("encryptionkey")
	SetKey(key)
	redirectHome(w, r)
}

func handleOauthCallback(w http.ResponseWriter, r *http.Request) {
	handleOAuthResponse(w, r)
	redirectHome(w, r)
}

func showAddForm(w http.ResponseWriter, r *http.Request) {
	showFile(w, "html/entryform.html")
}

func showDeleteForm(w http.ResponseWriter, r *http.Request, message string) {
	m := r.URL.Query()
	// TODO: Need to pull username from either query or from form value
	// r.FormValue("username")
	username := m["username"][0]
	data := struct {
		Username string
		Message  string
	}{
		username,
		message,
	}
	showTemplatedFile(w, "html/deleteform.html", data)
}

func saveAddForm(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user := User{Username: username, Password: password, State: exchangeLoginCaptured}
	user.Save()
	getFolderAndChangeKey(user)
	tryOAuth2(w, r, user)
}

func showUserList(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Users []User
	}{
		GetUsers(),
	}
	showTemplatedFile(w, "html/users.html", data)
}

func performDelete(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	storedUser := GetUser(username)
	if storedUser.Password == password {
		DeleteUser(username)
		redirectHome(w, r)
		return
	}
	showDeleteForm(w, r, "Incorrect Password")

}

func redirectHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", 302)
}

func showTemplatedFile(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		log.Printf("Error is", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Error while showing list ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func showFile(w http.ResponseWriter, filename string) {
	html, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, string(html))
}
