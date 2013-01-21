package cadmis

import (
	"appengine"
	//"fmt"
	// "appengine/user"
	"appengine/datastore"
	"appengine/urlfetch"
	// "code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	// "time"
)

// ユーザーの一つの問題への回答
type ProblemResult struct {
	ProblemId int64 // 解いた問題のID
	Correct   bool  // 正解、不正解
}

// ユーザーの一つの試験への回答
type ExamResultMessage struct {
	UserId       int64           // 解いたユーザーのID
	ExamId       int64           // 解いた試験のID
	StartTime    int64           // 開始時刻
	FinishedTime int64           // 終了時刻
	ProblemList  []ProblemResult // 解いた問題のリスト
}

// ユーザーの一つの試験への回答
type ExamResult struct {
	UserId       int64            // 解いたユーザーのID
	ExamId       int64            // 解いた試験のID
	StartTime    int64            // 開始時刻
	FinishedTime int64            // 終了時刻
	ProblemList  []*datastore.Key // 解いた問題のリスト
}

type StoreExamResultRequest struct {
	Result ExamResultMessage
}

type ProblemStatistics struct {
	Id      int64   `json:"id"`
	Total   int64   `json:"total"`
	Correct int64   `json:"correct"`
	Answer  bool    `json:"answer"`
	Score   float64 `json:"score"`
}

type CalculatedProblemScore struct {
	Id    int64   `json:"id"`
	Score float64 `json:"score"`
}

// 問題
type Problem struct {
	Id    int64
	Score float64
}

// 試験
type Exam struct {
	Id          int64
	ProblemList []*datastore.Key
}

// 問題の統計情報
// type ProblemStatistics struct {
// 	ProblemId    int64 // 問題のID
// 	TimesSolved  int64 // 解かれた回数
// 	TimesCorrect int64 // 正解だった回数
// }

// 試験の統計情報
// type ExamStatistics struct {
// 	ExamId          int64 // 試験のID
// 	TimesChallenged int64 // 挑戦された回数
// }

const (
	ProblemResultEntity string = "ProblemResult"
	ExamResultEntity    string = "ExamResult"
	// ProblemStatisticsEntity string = "ProblemStatistics"
	// ExamStatisticsEntity    string = "ExamStatistics"
)

// 試験の提出リクエストを処理する
func processExamPostRequest(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.Errorf("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Infof(string(buf))
	req := new(StoreExamResultRequest)
	err = json.Unmarshal(buf, req)
	if err != nil {
		// 受け付けられないリクエスト
		c.Errorf(err.Error())
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	exam := req.Result
	storeExam := ExamResult{
		UserId:       exam.UserId,
		ExamId:       exam.ExamId,
		StartTime:    exam.StartTime,
		FinishedTime: exam.FinishedTime,
	}

	// 試験結果を格納する
	c.Infof("Storing exam result")
	examKey := datastore.NewIncompleteKey(c, ExamResultEntity, nil)
	examKey, err = datastore.Put(c, examKey, &storeExam)
	if err != nil {
		c.Errorf("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Infof("Storing problem result")
	problemKeyList := []*datastore.Key{}
	for _, p := range exam.ProblemList {
		// 問題の結果を格納する
		pKey := datastore.NewIncompleteKey(c, ProblemResultEntity, examKey)
		pKey, err = datastore.Put(c, pKey, &p)
		if err != nil {
			c.Errorf("%s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		problemKeyList = append(problemKeyList, pKey)
	}

	storeExam.ProblemList = problemKeyList
	_, err = datastore.Put(c, examKey, &storeExam)
	if err != nil {
		c.Errorf("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 結果をサーバーに送りつけて評価してもらう
	err = reEvaluateProblemScore(c, exam)
	if err != nil {
		c.Errorf("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func gatherProblemStatistics(c appengine.Context, result ExamResultMessage) ([]ProblemStatistics, error) {
	problemNum := len(result.ProblemList)
	statList := make([]ProblemStatistics, problemNum)
	for i, p := range result.ProblemList {
		var q = datastore.NewQuery(ProblemResultEntity).Filter("ProblemId=", p.ProblemId)

		count, err := q.Count(c)
		if err != nil {
			return nil, err
		}

		var correctCount int64
		t := q.Run(c)
		for {
			var p ProblemResult
			_, err := t.Next(&p)
			if err != nil {
				if err == datastore.Done {
					// 最後まで到達したので抜ける
					break
				} else {
					return nil, err
				}
			}

			if p.Correct {
				correctCount++
			}
		}

		stat := &statList[i]
		stat.Id = p.ProblemId
		stat.Total = int64(count)
		stat.Correct = correctCount
		stat.Answer = p.Correct
		stat.Score = 10

		// c.Infof("total %d answer %d", stat.total, stat.answer)
	}
	return statList, nil
}

func storeReevaluatedProblemScore(c appengine.Context, calculated []CalculatedProblemScore) error {
	for _, cp := range calculated {
		var q = datastore.NewQuery(ProblemEntity).Limit(1).Filter("Id=", cp.Id)

		//keyList := *[]datastore.Key{}

		t := q.Run(c)
		for {
			var p Problem
			key, err := t.Next(&p)
			if err != nil {
				if err == datastore.Done {
					// 最後まで到達したので抜ける
					break
				} else {
					return err
				}
			}

			p.Score = cp.Score
			_, err = datastore.Put(c, key, &p)
			if err != nil {
				return err
			}

			//keyList = append(keyList, key)
		}

		// for _, k := range keyList {
		// 	datastore.
		// }
	}
	return nil
}

func reEvaluateProblemScore(c appengine.Context, result ExamResultMessage) error {
	statList, err := gatherProblemStatistics(c, result)
	if err != nil {
		return err
	}

	b, err := json.Marshal(statList)
	if err != nil {
		return err
	}

	targetUrl := "http://180.37.181.90:8080/score"
	data := string(b)
	userId := result.UserId

	c.Infof("%s", data)

	values := url.Values{}
	values.Set("user", strconv.FormatInt(userId, 10))
	values.Set("data", data)

	client := urlfetch.Client(c)
	res, err := client.PostForm(targetUrl, values)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	c.Infof(string(buf))
	calculated := []CalculatedProblemScore{}
	err = json.Unmarshal(buf, &calculated)
	if err != nil {
		return err
	}

	// 再計算されたスコアをストアする
	err = storeReevaluatedProblemScore(c, calculated)
	if err != nil {
		return nil
	}

	return nil
}

func handleExamRequest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Infof("Method: %s Url:%s ContentLength: %d\n", r.Method, r.URL, r.ContentLength)

	if r.Method == "GET" {
	} else if r.Method == "POST" {
		processExamPostRequest(c, w, r)
	}
}
