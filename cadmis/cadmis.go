package cadmis

import (
	"appengine"
	"appengine/user"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/api/1/user", handleUserRequest)
	http.HandleFunc("/api/1/access_token", handleAccessTokenRequest)
}

// ユーザーに関するリクエストを処理する
func handleUserRequest(w http.ResponseWriter, r *http.Request) {
}

// ログイン用トークンに関するリクエストを処理する
func handleAccessTokenRequest(w http.ResponseWriter, r *http.Request) {
}

/*
func hostStaticFile() {
	path := req.URL.Path[1:]
	fmt.Printf("path: %s\n",path)
	http.ServeFile(w, req, path)	
}

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
	   fmt.Fprintf(w, "Please, <a href='/_ah/login_required'>login</a>.")
	} else {
		url, _ := user.LogoutURL(c, "/")
		fmt.Fprintf(w, `Welcome, %s! (<a href="%s">sign out</a>)`, u, url)
	}
}


func loginHandler(w http.ResponseWriter, r * http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	providers := map[string]string {
		"Google"   : "gmail.com",
	}

	c := appengine.NewContext(r)
	fmt.Fprintf(w, "Sign in at: ")
	for name, url := range providers {
		login_url, err := user.LoginURLFederated(c, "/", url)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "[<a href='%s'>%s</a>]", login_url, name)
	}
}
*/