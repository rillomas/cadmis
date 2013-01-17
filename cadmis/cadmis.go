package cadmis

import (
	"appengine"
	//"fmt"
	// "appengine/user"
	"appengine/datastore"
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// ユーザーが所属するグループ
type UserGroup int

const (
	FirstYearHighSchool  UserGroup = iota // 高1
	SecondYearHighSchool                  // 高2
	ThirdYearHighSchool                   // 高3
	College                               // 大学生
	CramSchool                            // 予備校 
	Certified                             // 高認
	Adult                                 // 社会人
)

const (
	UserEntity            string = "User"
	UserInformationEntity string = "UserInformation"
)

// ハンドラを設定する
func init() {
	http.HandleFunc("/api/1/user", handleUserRequest)
	http.HandleFunc("/api/1/access_token", handleAccessTokenRequest)
}

// 認証情報
type AuthenticationInformation struct {
	Email    string
	Password string
}

// ユーザーのモデル
type User struct {
	Id           int64 // 自動生成される一意なID
	Email        string
	PasswordHash []byte
	Information  *datastore.Key // ユーザー情報への鍵
}

// ユーザーの名前 
type UserName struct {
	FirstName string // 名前
	LastName  string // 苗字
}

// ユーザーの情報
type UserInformation struct {
	Id int64 // 対応するユーザーのID
	UserName
	Group      UserGroup // 所属するグループ
	SchoolName string    // 所属する学校名（もしあれば）
	JoinDate   time.Time // 加入日
}

//  ユーザーを追加する
func addUser(c appengine.Context, email, password string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := User{
		Email:        email,
		PasswordHash: hash,
	}

	ik := datastore.NewIncompleteKey(c, UserEntity, nil)
	key, err := datastore.Put(c, ik, &user)
	if err != nil {
		return err
	}

	user.Id = key.IntID()                 // 生成されたIDを格納する
	_, err = datastore.Put(c, key, &user) // 再度格納
	if err != nil {
		return err
	}

	return nil
}

// 指定されたメールアドレスをもつユーザーがいるかどうかを調べる
func userExists(c appengine.Context, email string) (bool, error) {
	q := datastore.NewQuery(UserEntity).Limit(1).Filter("Email =", email)
	count, err := q.Count(c)
	return count > 0, err
}

// ユーザーに関するリクエストを処理する
func handleUserRequest(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	// リクエストをパースする
	c.Infof("Method: %s Url:%s ContentLength: %d\n", r.Method, r.URL, r.ContentLength)
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.Errorf("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	info := AuthenticationInformation{}
	err = json.Unmarshal(buf, &info)
	if err != nil {
		// 受け付けられないリクエスト
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 既におなじユーザーがいるかどうかを調べて、いなかったら追加する
	exists, err := userExists(c, info.Email)
	if err != nil {
		c.Errorf("Error at userExists: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		// 同じIDのユーザーが既に存在するので失敗
		http.Error(w, "Email address already in use", http.StatusNotAcceptable)
		return
	} else {
		// ユーザーが重複しないので追加
		err = addUser(c, info.Email, info.Password)
		if err != nil {
			c.Errorf("%s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// ユーザーの追加に成功
		c.Infof("added User: %s\n", info.Email)
	}
}

// ログイン用トークンに関するリクエストを処理する
func handleAccessTokenRequest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// リクエストをパースする
	c.Infof("Method: %s Url:%s ContentLength: %d\n", r.Method, r.URL, r.ContentLength)
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.Errorf("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	info := AuthenticationInformation{}
	err = json.Unmarshal(buf, &info)
	if err != nil {
		// 受け付けられないリクエスト
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 同じメールアドレスのユーザーを探す
	q := datastore.NewQuery(UserEntity).Limit(1).Filter("Email =", info.Email)
	count, err := q.Count(c)
	if err != nil {
		c.Errorf("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if count <= 0 {
		// 指定されたメールアドレスを持つユーザーがいない
		http.Error(w, "Wrong Email address or password", http.StatusUnauthorized)
		return
	}

	// 見つけたユーザーを認証可能かチェックする
	t := q.Run(c)
	for {
		var u User
		_, err := t.Next(&u)
		if err != nil {
			if err == datastore.Done {
				// 最後まで到達したので抜ける
				break
			} else {
				c.Errorf("%s", err.Error())
				http.Error(w, "Error while authenticating", http.StatusInternalServerError)
				return
			}
		}

		// ハッシュの一致比較
		err = bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(info.Password))
		if err == nil {
			// 認証成功
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("accesstoken"))
			return
		}
	}

	// 認証失敗
	http.Error(w, "Wrong email address or password", http.StatusUnauthorized)
}
