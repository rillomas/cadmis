package cadmis

import (
	"appengine"
	//"fmt"
	// "appengine/user"
	// "appengine/datastore"
	// "code.google.com/p/go.crypto/bcrypt"
	// "encoding/json"
	// "io/ioutil"
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

type UserInformationQuery struct {
	AccessTokenId int64
	UserId        int64
}

func handleUserProfileRequest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Infof("Method: %s Url:%s ContentLength: %d\n", r.Method, r.URL, r.ContentLength)

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

	// ユーザーの情報をもってきて返す

}
