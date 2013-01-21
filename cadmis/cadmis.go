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
	http.HandleFunc("/api/1/get_exam", handleGetExam)
	http.HandleFunc("/api/1/compute_rank", handleComputeRank)
	http.HandleFunc("/api/1/init_exams", handleInitExams)
}
