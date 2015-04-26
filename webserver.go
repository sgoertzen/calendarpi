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
	case "saveuser":
		if !needskey(w) {
			saveAddForm(w, r)
		}
	case "savecalendar":
		if !needskey(w) {
			saveCalendarForm(w, r)
		}
	case "delete":
		if !needskey(w) {
			showDeleteForm(w, r, "", "")
		}
	case "confirmdelete":
		if !needskey(w) {
			performDelete(w, r)
		}
	case "selectcalendar":
		if !needskey(w) {
			showCalendarSelectPage(w, r)
		}
	case "logic": // TESTING ONLY.  REMOVE!
		user := GetUser("goertzs")
		log.Println("Starting on user ", user.Username)

		soapResults := getExchangeCalendarData(user)
		appointments := ParseAppointments(soapResults)

		events, err := getGCalAppointments(user)
		if err != nil {
			log.Fatal(err)
		}
		mergeEvents(appointments, events)
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
	http.Redirect(w, r, "selectcalendar?username="+user.Username, 302)
}

func showCalendarSelectPage(w http.ResponseWriter, r *http.Request) {
	m := r.URL.Query()
	username := m["username"][0]
	user := GetUser(username)

	data := map[string]interface{}{
		"Calendars": GetCalendarList(user),
		"Username": username,
	}
	showTemplatedFile(w, "html/calendarform.html", data)
}

func showAddForm(w http.ResponseWriter, r *http.Request) {
	showTemplatedFile(w, "html/entryform.html", nil)
}

func showDeleteForm(w http.ResponseWriter, r *http.Request, message string, username string) {
	if len(username) == 0 {
		m := r.URL.Query()
		usernames := m["username"]
		if len(usernames) >= 1 {
			username = usernames[0]
		}
	}
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

func saveCalendarForm(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	calendar := r.FormValue("calendar")
	user := GetUser(username)
	user.GCalid = calendar
	user.Save()
	redirectHome(w, r)
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
	if storedUser.Password == password { //|| string(Key()) == CreateKey(username) {
		DeleteUser(username)
	} else {
		// Allow for pw that matches encryption key
		key, _ := CreateKey(password)
		if string(key) == string(Key()) {
			DeleteUser(username)
		} else {
			showDeleteForm(w, r, "Incorrect Password", username)
			return
		}
	}
	redirectHome(w, r)
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
