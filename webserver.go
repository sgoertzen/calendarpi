package main

import (
	"github.com/sgoertzen/xchango"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./scripts"))))
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
	case "changepassword":
		if !needskey(w) {
			showPasswordForm(w, r)
		}
	case "savenewpassword":
		if !needskey(w) {
			savePasswordForm(w, r)
		}
	case "sync":
		if !needskey(w) {
			syncUser(w, r)
		}
	case "":
		if !needskey(w) {
			showUserList(w, r)
		}
	default:
		http.NotFound(w, r)
	}
}

func syncUser(w http.ResponseWriter, r *http.Request) {
	m := r.URL.Query()
	usernames := m["username"]
	if len(usernames) == 0 {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	username := usernames[0]
	user := GetUser(username)
	Sync(user)
	redirectHome(w, r)
}

func needskey(w http.ResponseWriter) bool {
	key := storedKey()
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
	err := DeserializeUsers(key)
	if err != nil {
		log.Println("Unable to decrypt the file.  If you have forgotten the password just manually delete the saved users file.", err)
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
		"Calendars": GetGCalendarList(user),
		"Username":  username,
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
	user := User{
		ExUser: &xchango.ExchangeUser{
			Username: username,
			Password: password,
		},
		Username: username,
		Password: password,
		State:    exchangeLoginCaptured}
	user.Save()

	cal, err := xchango.GetExchangeCalendar(user.ExUser)
	if err != nil {
		user.State = registererror
		user.Save()
	} else {
		user.ExCal = cal
		user.State = exchangeLoginVerified
		user.Save()
		tryOAuth2(w, r, user)
	}
}

func saveCalendarForm(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	calendar := r.FormValue("calendar")
	user := GetUser(username)
	user.GCalid = calendar
	user.Save()
	redirectHome(w, r)
}

func showPasswordForm(w http.ResponseWriter, r *http.Request) {
	username := ""
	m := r.URL.Query()
	usernames := m["username"]
	if len(usernames) >= 1 {
		username = usernames[0]
	}
	log.Printf("Showing password change form for: '%s'", username)
	data := map[string]interface{}{
		"Username": username,
	}
	showTemplatedFile(w, "html/changepassword.html", data)
}

func savePasswordForm(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	oldpassword := r.FormValue("oldpassword")
	newpassword := r.FormValue("newpassword")
	confirmpassword := r.FormValue("confirmpassword")
	user := GetUser(username)
	
	if user == (User{}) {
		log.Printf("Unable to find a user named %s", username)
		return
	}
	if user.Password != oldpassword {
		log.Printf("Password does not match old password for user named %s", username)
		return
	}
	if newpassword != confirmpassword {
		log.Printf("New passwords do not match: '%s'", username)
		return
	}
	
	log.Printf("Changing password for: '%s'", username)
	user.LastSync = time.Time{}
	user.State = registered
	user.Password = newpassword
	user.Save()
	
	cal, err := xchango.GetExchangeCalendar(user.ExUser)
	if err != nil || cal == nil {
		log.Println("Unable to get the exchange client", err)
		user.State = registererror
		user.Save()
	} else {
		log.Printf("Got exchange calendar of Folder: '%s' and Key: '%s", cal.Folderid, cal.Changekey)
		user.ExCal = cal
		user.State = registered
		user.Save()
	}
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
	if storedUser.Password == password {
		DeleteUser(username)
	} else {
		// Allow for pw that matches encryption key
		key, _ := CreateKey(password)
		if string(key) == string(storedKey()) {
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
		log.Println("Error while parsing template file", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Println("Error while showing list ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
