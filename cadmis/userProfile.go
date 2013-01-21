package cadmis

import (
	"appengine"
	//"fmt"
	// "appengine/user"
	"appengine/datastore"
	// "code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
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
	UserProfileEntity string = "UserProfile"
)

// ユーザーの情報
type UserProfile struct {
	Id         int64     // プロフィール自体のID
	UserId     int64     // 対応するユーザーのID
	FirstName  string    // 名前
	LastName   string    // 苗字
	Group      UserGroup // 所属するグループ
	SchoolName string    // 所属する学校名（もしあれば）
	JoinDate   time.Time // 加入日
}

type ApplyProfileRequest struct {
	Profile UserProfile
}

// ユーザープロフィールを追加する
func addUserProfile(c appengine.Context, userId int64, profile *UserProfile) error {
	uk := datastore.NewKey(c, UserEntity, "", userId, nil)
	user := User{}
	err := datastore.Get(c, uk, &user)
	if err != nil {
		return err
	}

	ik := datastore.NewIncompleteKey(c, UserProfileEntity, uk)
	profile.UserId = userId
	key, err := datastore.Put(c, ik, profile)
	if err != nil {
		return err
	}

	profile.Id = key.IntID()                // 生成されたIDを格納する
	_, err = datastore.Put(c, key, profile) // 再度格納
	if err != nil {
		return err
	}

	// ユーザー情報の鍵も追加して格納
	user.Profile = key
	_, err = datastore.Put(c, uk, &user)
	if err != nil {
		return err
	}

	c.Infof("Added user profile")
	return nil
}

// ユーザーをもってくる
func getUser(c appengine.Context, userId int64) (*User, error) {
	key := datastore.NewKey(c, UserEntity, "", userId, nil)

	user := User{}
	err := datastore.Get(c, key, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ユーザープロフィールをもってくる
func getUserProfile(c appengine.Context, profileId int64) (*UserProfile, error) {
	key := datastore.NewKey(c, UserProfileEntity, "", profileId, nil)

	profile := new(UserProfile)
	err := datastore.Get(c, key, profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// ユーザープロフィールを設定する
func setUserProfile(c appengine.Context, profile *UserProfile) error {
	var userId = profile.UserId
	uk := datastore.NewKey(c, UserEntity, "", userId, nil)
	user := User{}
	err := datastore.Get(c, uk, &user)
	if err != nil {
		return err
	}

	key := user.Profile
	_, err = datastore.Put(c, key, profile)
	if err != nil {
		return err
	}
	return nil
}

// ユーザー取得のリクエストを取得する
func processUserGetRequest(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	// リクエストをパースする
	accessTokenIdStr := r.FormValue("at")
	userIdStr := r.FormValue("ui")

	base := 0
	bitSize := 64
	userId, err := strconv.ParseInt(userIdStr, base, bitSize)
	if err != nil {
		c.Errorf("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	accessTokenId, err := strconv.ParseInt(accessTokenIdStr, base, bitSize)
	if err != nil {
		c.Errorf("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Infof("UserId: %d AccessTokenId: %d\n", userId, accessTokenId)

	// ユーザーに発行されたトークンが正しいか確かめる
	match, _, err := accessTokenMatches(c, userId, accessTokenId)
	if err != nil {
		c.Errorf("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !match {
		http.Error(w, "Access token does not match user.", http.StatusUnauthorized)
		return
	}

	// ユーザーをもってきて返す
	user, err := getUser(c, userId)
	if err != nil {
		c.Errorf("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile := new(UserProfile)
	if user.Profile == nil {
		// プロフィールがなかったので新規に追加する
		c.Infof("Profile for user %d not found. Creating profile.", userId)
		err = addUserProfile(c, userId, profile)
		if err != nil {
			c.Errorf("%s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// プロフィールがあったのでそれをかえす
		c.Infof("Profile for user %d found. Sending profile.", userId)
		err = datastore.Get(c, user.Profile, profile)
		if err != nil {
			c.Errorf("%s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	c.Infof("Sending profile %d", profile.Id)
	b, err := json.Marshal(profile)
	if err != nil {
		c.Errorf("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(b)
}

func processUserPostRequest(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.Errorf("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Infof("%s", string(buf))
	req := new(ApplyProfileRequest)
	err = json.Unmarshal(buf, req)
	if err != nil {
		// 受け付けられないリクエスト
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	profile := req.Profile
	c.Infof("Id: %d FirstName: %s LastName: %s", profile.Id, profile.FirstName, profile.LastName)
	err = setUserProfile(c, &profile)
	if err != nil {
		c.Errorf("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Infof("Set new user profile")
}

func handleUserProfileRequest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Infof("Method: %s Url:%s ContentLength: %d\n", r.Method, r.URL, r.ContentLength)

	if r.Method == "GET" {
		processUserGetRequest(c, w, r)
	} else if r.Method == "POST" {
		processUserPostRequest(c, w, r)
	}
}
