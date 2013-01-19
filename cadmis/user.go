package cadmis

import (
	"appengine"
	"appengine/datastore"
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	UserEntity        string = "User"
	AccessTokenEntity string = "AccessToken"
)

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
	Profile      *datastore.Key // ユーザー情報への鍵
}

// アクセストークンのモデル
type AccessToken struct {
	Id     int64 // アクセストークン自体のID
	UserId int64 // アクセストークンを発行されたユーザーのID
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

// アクセストークンが発行済みかどうかを調べる
func accessTokenPublished(c appengine.Context, userId int64) (bool, *datastore.Query, error) {
	q := datastore.NewQuery(AccessTokenEntity).Limit(1).Filter("UserId =", userId)
	count, err := q.Count(c)
	return count > 0, q, err
}

func accessTokenMatches(c appengine.Context, userId, accessToken int64) (bool, *datastore.Query, error) {
	q := datastore.NewQuery(AccessTokenEntity).Limit(1).Filter("UserId =", userId).Filter("Id =", accessToken)
	count, err := q.Count(c)
	return count > 0, q, err
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

			// 既にトークンを発行してるかどうかを調べる
			published, query, err := accessTokenPublished(c, u.Id)
			if err != nil {
				c.Errorf("%s", err.Error())
				http.Error(w, "Error while authenticating", http.StatusInternalServerError)
				return
			}

			at := new(AccessToken)
			if published {
				c.Infof("Found published access token")
				// アクセストークンがもうあるのでそれをまた返す
				_, err = query.Run(c).Next(at)
				if err != nil {
					c.Errorf("%s", err.Error())
					http.Error(w, "Error while authenticating", http.StatusInternalServerError)
					return
				}
			} else {
				c.Infof("Creating new access token")
				// アクセストークンが存在しないので発行する
				at, err = publishAccessToken(c, u.Id)
				if err != nil {
					c.Errorf("%s", err.Error())
					http.Error(w, "Error while authenticating", http.StatusInternalServerError)
					return
				}
			}

			c.Infof("Sending access token %d to user %d", at.Id, at.UserId)
			b, err := json.Marshal(at)
			if err != nil {
				c.Errorf("%s", err.Error())
				http.Error(w, "Error while authenticating", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.Write(b)
			return
		}
	}

	// 認証失敗
	http.Error(w, "Wrong email address or password", http.StatusUnauthorized)
}

//新しいアクセストークンを発行する 
func publishAccessToken(c appengine.Context, userId int64) (*AccessToken, error) {
	at := AccessToken{
		UserId: userId,
	}

	ik := datastore.NewIncompleteKey(c, AccessTokenEntity, nil)
	key, err := datastore.Put(c, ik, &at)
	if err != nil {
		return nil, err
	}

	at.Id = key.IntID()                 // 生成されたIDを格納する
	_, err = datastore.Put(c, key, &at) // 再度格納
	if err != nil {
		return nil, err
	}

	return &at, nil
}
