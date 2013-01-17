package cadmis

import (
	"appengine"
	// "appengine/user"
	"appengine/datastore"
	"time"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"code.google.com/p/go.crypto/bcrypt"
)

// ユーザーが所属するグループ
type UserGroup int

const (
	FirstYearHighSchool UserGroup = iota// 高1
	SecondYearHighSchool // 高2
	ThirdYearHighSchool // 高3
	College // 大学生
	CramSchool // 予備校 
	Certified // 高認
	Adult // 社会人
)

const (
	UserEntity string = "User"
	UserInformationEntity string = "UserInformation"
)

// ハンドラを設定する
func init() {
	http.HandleFunc("/api/1/user", handleUserRequest)
	http.HandleFunc("/api/1/access_token", handleAccessTokenRequest)
}

// ユーザー追加リクエスト
type AddUserRequest struct {
	UserId string
	Password string
}

// ユーザーのモデル
type User struct {
	Id string
	PasswordHash []byte
	Information *datastore.Key // ユーザー情報への鍵
}

// ユーザーの名前 
type UserName struct {
	First string // 名前
	Last string // 苗字
}

// ユーザーの情報
type UserInformation struct {
	Id string // ユーザーのID
	UserName
	Group UserGroup // 所属するグループ
	SchoolName string // 所属する学校名（もしあれば）
	JoinDate time.Time // 加入日
}

//  ユーザーを追加する
func addUser(c appengine.Context, userId, password string) error {
	
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := User {
		Id : userId,
		PasswordHash : hash,
	}

	key := datastore.NewKey(c, UserEntity, userId, 0, nil) // userIDをキーにする
	_, err = datastore.Put(c, key, &user)
	if err != nil {
		return err
	}

	return nil
}

// 指定されたIDをもつユーザーを持ってくる
func getUser(c appengine.Context, userId string) (*User, error) {
	key := datastore.NewKey(c, UserEntity, userId, 0, nil)
	user := new(User)
	err := datastore.Get(c, key, user)
	if err != nil {
		return user,err
	}
	return user,nil
}

// ユーザーに関するリクエストを処理する
func handleUserRequest(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	// リクエストをパースする
	c.Infof("Method: %s Url:%s ContentLength: %d\n", r.Method, r.URL, r.ContentLength)
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.Errorf("%s",err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req := AddUserRequest{}
	json.Unmarshal(buf, &req)

	// 既におなじユーザーがいるかどうかを調べて、いなかったら追加する
	_, err = getUser(c, req.UserId)
	if err == nil {
		// 同じIDのユーザーが既に存在するので失敗
		http.Error(w, "User ID already taken", http.StatusNotAcceptable)
		return
	}

	switch(err) {
	case datastore.ErrNoSuchEntity:
		// ユーザーが重複しないので追加
		err = addUser(c, req.UserId, req.Password)
		if err != nil {
			c.Errorf("%s",err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// ユーザーの追加に成功
		c.Infof("added User: %s\n", req.UserId)
	default:
		c.Errorf("Error at getUser: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ログイン用トークンに関するリクエストを処理する
func handleAccessTokenRequest(w http.ResponseWriter, r *http.Request) {
}
