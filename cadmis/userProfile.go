package cadmis

import (
	// "appengine"
	//"fmt"
	// "appengine/user"
	// "appengine/datastore"
	// "code.google.com/p/go.crypto/bcrypt"
	// "encoding/json"
	// "io/ioutil"
	"net/http"
	// "strconv"
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

func handleUserProfileRequest(w http.ResponseWriter, r *http.Request) {
}
