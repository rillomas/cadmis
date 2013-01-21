package cadmis

import (
	"net/http"
)

// ハンドラを設定する
func init() {
	http.HandleFunc("/api/1/user", handleUserRequest)
	http.HandleFunc("/api/1/user_profile", handleUserProfileRequest)
	http.HandleFunc("/api/1/access_token", handleAccessTokenRequest)
	http.HandleFunc("/api/1/exam", handleExamRequest)
}
