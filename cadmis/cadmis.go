package cadmis

import (
//	"appengine"
//	"appengine/user"
	"fmt"
	"net/http"
//	"bytes"
	"encoding/json"
	"io/ioutil"
)

func init() {
	http.HandleFunc("/api/1/user", handleUserRequest)
	http.HandleFunc("/api/1/access_token", handleAccessTokenRequest)
}

type AddUserRequest struct {
	UserId string
	Password string
}

// ユーザーに関するリクエストを処理する
func handleUserRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Method: %s Url:%s ContentLength: %d\n", r.Method, r.URL, r.ContentLength)
	// fmt.Printf("Content: %s\n", s)
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(r.Body)
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	req := AddUserRequest{}
	json.Unmarshal(buf, &req)

	fmt.Printf("UserId: %s Password: %s\n", req.UserId, req.Password)

	// for k,v := range r.Form {
	// 	fmt.Printf("k:%s v:%s\n", k, v)
	// }
	// id := r.FormValue("userId")
	// pass := r.FormValue("password")
	// fmt.Printf("id: %s pass:%s\n", id, pass)

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