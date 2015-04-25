package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type WebConfig interface {
	Port() string
	Certificate() string
	PrivateKey() string
}

var header string

func RunServer(conf WebConfig) {
	log.Println("Webserver starting on port " + conf.Port())
	headerbytes, err := ioutil.ReadFile("html/header.html")
	if err != nil {
		log.Fatal(err)
	}
	header = string(headerbytes)
	http.HandleFunc("/", handler)
	err = http.ListenAndServeTLS(":"+conf.Port(), conf.Certificate(), conf.PrivateKey(), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path[1:]

	switch path {
	case "savekey":
		saveKey(w, r)
	case "oauth2callback":
		handleOauthCallback(w, r)
	case "add":
		if !needskey(w) {
			showAddForm(w, r)
		}
	case "save":
		if !needskey(w) {
			saveAddForm(w, r)
		}
	case "delete":
		if !needskey(w) {
			showDeleteForm(w, r, "")
		}
	case "confirmdelete":
		if !needskey(w) {
			performDelete(w, r)
		}
	case "logic": // TESTING ONLY.  REMOVE!
		user := GetUser("goertzs")
		log.Println("Starting on user ", user.Username)
		//soapResults := getExchangeCalendarData(user)
		getGCalAppointments(user, "todo")
		//log.Println("soapresults", soapResults)
		//appointments := ParseAppointments(soapResults)
		//log.Println("len:", len(appointments))
		//processAppointments(user, appointments)
		redirectHome(w, r)
	case "":
		if !needskey(w) {
			showUserList(w, r)
		}
	default:
		showFile(w, path)
	}
}

func needskey(w http.ResponseWriter) bool {
	key := Key()
	if len(key) < 10 {
		data := map[string]interface{}{
			"Header": template.HTML(header),
		}
		showTemplatedFile(w, "html/keyform.html", data)
		return true
	}
	return false
}

func saveKey(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("encryptionkey")
	SetKey(key)
	err := DeserializeUsers()
	if err != nil {
		log.Println("Error while loading file.  Ignoring:", err)
	}
	redirectHome(w, r)
}

func handleOauthCallback(w http.ResponseWriter, r *http.Request) {
	user, err := handleOAuthResponse(w, r)
	if err != nil {
		return
	}
	showCalendarSelectPage(w, r, user)
}

func showCalendarSelectPage(w http.ResponseWriter, r *http.Request, user User) {
	data := map[string]interface{}{
		"Calendars": GetCalendarList(user),
	}
	showTemplatedFile(w, "html/calendarform.html", data)
}

func showAddForm(w http.ResponseWriter, r *http.Request) {
	showTemplatedFile(w, "html/entryform.html", nil)
}

func showDeleteForm(w http.ResponseWriter, r *http.Request, message string) {
	m := r.URL.Query()
	username := m["username"][0]
	data := map[string]interface{}{
		"Username": username,
		"Message":  message,
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
	data := map[string]interface{}{
		"Users": GetUsers(),
	}
	showTemplatedFile(w, "html/users.html", data)
}

func performDelete(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	storedUser := GetUser(username)
	if storedUser.Password == password || string(Key()) == password {
		DeleteUser(username)
		redirectHome(w, r)
		return
	}
	showDeleteForm(w, r, "Incorrect Password")

}

func redirectHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", 302)
}

func addCommonData(data map[string]interface{}) {
	data["Header"] = template.HTML(header)
}

func showTemplatedFile(w http.ResponseWriter, filename string, data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	addCommonData(data)

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
	if strings.HasSuffix(filename, ".js") {
		w.Header().Set("content-Type", "application/x-javascript")
	}
	if strings.HasSuffix(filename, ".css") {
		w.Header().Set("content-Type", "text/css")
	}
	fmt.Fprintf(w, string(html))
}
