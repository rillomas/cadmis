package cadmis

import (
	"net/http"
)

// ハンドラを設定する
func init() {
	http.HandleFunc("/api/1/user", handleUserRequest)
	http.HandleFunc("/api/1/user_profile", handleUserProfileRequest)
	http.HandleFunc("/api/1/access_token", handleAccessTokenRequest)
	http.HandleFunc("/api/1/write_entry", handleWriteEntry)
	http.HandleFunc("/api/1/compute_rank", handleComputeRank)
	http.HandleFunc("/api/1/problem_list", handleProblemList)
	http.HandleFunc("/api/1/user_list", handleUserList)
	http.HandleFunc("/api/1/init_problems", handleInitProblems)
}
