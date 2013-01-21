package cadmis

import (
	"appengine"
	"fmt"
	//    "io"
	"io/ioutil"
	// "appengine/user"
	"appengine/datastore"
	"appengine/urlfetch"
	// "code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	// "io/ioutil"
	"net/http"
	"strconv"
	// "strings"
	// "time"
)

const (
	ProblemEntity string = "Problem"
)

type Rank struct {
	ExamId   int64
    ExamName string
	UserId   int64
    UserName string
	Score    float64
}

func handleComputeRank(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	uid, err1 := strconv.Atoi(r.FormValue("userId"))
	pid, err2 := strconv.Atoi(r.FormValue("examId"))
	if err1 != nil {
		uid = -1
	}
	if err2 != nil {
		pid = -1
	}
	score, err3 := computeRank(c, int64(uid), int64(pid))
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
		return
	}
	out, _ := json.Marshal(score)
	fmt.Fprint(w, string(out))
}

func handleInitExams(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	/*
	   client := urlfetch.Client(c)
	   resp, err := client.Get("http://api.iknow.jp/categories")
	   if err != nil {
	       http.Error(w, err.Error(), http.StatusInternalServerError)
	       //return err
	   }
	   defer resp.Body.Close()
	   body, err2 := ioutil.ReadAll(resp.Body)
	   if err2 != nil {
	       http.Error(w, err.Error(), http.StatusInternalServerError)
	       //return err
	   }
	   //dec := json.NewDecoder(strings.NewReader(stream));
	   //dec := json.NewDecoder(resp.Body)
	   var categories []map[string]interface{}
	   err3 := json.Unmarshal(body, &categories)
	   if err3 != nil {
	       http.Error(w, err.Error(), http.StatusInternalServerError)
	       //return err
	   }
	   for i := range categories {
	       fmt.Fprintf(w, "%s: %s\n", categories[i]["id"], categories[i]["name"])
	   }
	*/
	client := urlfetch.Client(c)
	resp, err := client.Get("http://api.iknow.jp/goals/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
	}
	var goals1 map[string]interface{}
	err3 := json.Unmarshal(body, &goals1)
	if err3 != nil {
		http.Error(w, err3.Error(), http.StatusInternalServerError)
	}
	var goals2 []interface{}
	goals2 = goals1["goals"].([]interface{})
	for i := range goals2 {
		//fmt.Fprintf(w, "%d: %s\n", int(items2[i]["id"].(float64)), items2[i]["uri"])

		var exam Exam
		exam.Id = int64(goals2[i].(map[string]interface{})["id"].(float64))
        exam.Name = goals2[i].(map[string]interface{})["title"].(string)
		//exam.Category = "English";
		// /*
		var er ExamResult
		if i < 5 {
			er.UserId = 1 // + int64(i)
			er.ExamId = exam.Id
		}
		// */
		client := urlfetch.Client(c)
		resp, err := client.Get("http://api.iknow.jp/goals/" + strconv.Itoa(int(exam.Id)) + "/items")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer resp.Body.Close()
		body, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}
		var items map[string]interface{}
		err3 := json.Unmarshal(body, &items)
		if err3 != nil {
			http.Error(w, err3.Error(), http.StatusInternalServerError)
			return
		}
		if items["items"] == nil {
			continue
		}
		items2 := items["items"].([]interface{})
		for j := range items2 {
			var problem Problem
			problem.Id = int64(items2[j].(map[string]interface{})["id"].(float64))
			problem.Score = 10
			//fmt.Fprintf(w, "%d: %s\n", int(items2[j]["id"].(float64)), items2[j]["uri"])
			key, err2 := datastore.Put(c, datastore.NewIncompleteKey(c, ProblemEntity, nil), &problem)
			if err2 != nil {
				http.Error(w, err2.Error(), http.StatusInternalServerError)
				return
			}
			exam.ProblemList = append(exam.ProblemList, key)
			// /*
			if i < 5 {
				var pr ProblemResult
				pr.ProblemId = problem.Id
				if i < 7 {
					pr.Correct = true
				} else {
					pr.Correct = false
				}
				key, err2 := datastore.Put(c, datastore.NewIncompleteKey(c, "ProblemResult", nil), &pr)
				if err2 != nil {
					http.Error(w, err2.Error(), http.StatusInternalServerError)
					return
				}
				er.ProblemList = append(er.ProblemList, key)
			}
			// */
		}
		_, err4 := datastore.Put(c, datastore.NewIncompleteKey(c, "Exam", nil), &exam)
		if err4 != nil {
			http.Error(w, err4.Error(), http.StatusInternalServerError)
			return
		}
		// /*
		if i < 5 {
			_, err5 := datastore.Put(c, datastore.NewIncompleteKey(c, "ExamResult", nil), &er)
			if err5 != nil {
				http.Error(w, err5.Error(), http.StatusInternalServerError)
				return
			}
		}
		// */
	}
}

func getExamResultScore(c appengine.Context, e ExamResult) (error, float64) {
	var score float64
	for i := range e.ProblemList {
		pr := new(ProblemResult)
		if err := datastore.Get(c, e.ProblemList[i], pr); err != nil {
			return err, score
		}
		if pr.Correct {
			q := datastore.NewQuery(ProblemEntity).Filter("Id =", pr.ProblemId).Limit(1)
			problems := make([]Problem, 0, 1)
			if _, err := q.GetAll(c, &problems); err != nil {
				return err, score
			}
            if (len(problems) > 0) {
			    score += problems[0].Score
            }
		}
	}
	return nil, score
}

func getExamScore(c appengine.Context, e Exam) (error, float64) {
	var score float64
	for i := range e.ProblemList {
		p := new(Problem)
		if err := datastore.Get(c, e.ProblemList[i], p); err != nil {
			return err, score
		}
		score += p.Score
	}
	return nil, score
}

func computeRank(c appengine.Context, userId int64, examId int64) ([]Rank, error) {
	score := make([]Rank, 0, 0)
	q := datastore.NewQuery("ExamResult")
	for t := q.Run(c); ; {
		var e ExamResult
		_, err := t.Next(&e)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return score, err
		}
		if (userId == -1 || e.UserId == userId) && (examId == -1 || e.ExamId == examId) {
			err, s := getExamResultScore(c, e)
			if err != nil {
				return score, err
			}
            var userName,examName string;
            var userProf []UserProfile
            var exam []Exam
	        q := datastore.NewQuery("UserProfile").Filter("UserId =", e.UserId)
	        if _, err := q.GetAll(c, &userProf); err == nil && len(userProf) > 0 {
		        userName = userProf[0].FirstName + " " + userProf[0].LastName
	        }
	        q2 := datastore.NewQuery("Exam").Filter("Id =", e.ExamId)
	        if _, err := q2.GetAll(c, &exam); err == nil && len(exam) > 0 {
		        examName = exam[0].Name
	        }
			score = append(score, Rank{e.ExamId, examName, e.UserId, userName, s})
		}
	}
	return score, nil
}

func handleGetExam(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	uid, err1 := strconv.Atoi(r.FormValue("userId"))
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}

	var entries []ExamResult
	var score float64
	score = 0
	q := datastore.NewQuery("ExamResult").Filter("UserId =", uid)
	if _, err := q.GetAll(c, &entries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for i := range entries {
		err, s := getExamResultScore(c, entries[i])
		if err != nil {
			return
		}
		score += s
	}
	if len(entries) > 0 {
		score = score / float64(len(entries))
	}
	score += 10 // boost level up!

	var exams []Exam
	q2 := datastore.NewQuery("Exam")
	if _, err := q2.GetAll(c, &exams); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var newExam Exam
	for i := range exams {
		var results []ExamResult
		q := datastore.NewQuery("ExamResult").Filter("UserId =", uid).Filter("ExamId =", exams[i].Id).Limit(1)
		if _, err := q.GetAll(c, &results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err1, diff1 := getExamScore(c, exams[i])
		err2, diff2 := getExamScore(c, newExam)
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusInternalServerError)
			return
		}
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}
		diff1 -= score
		diff2 -= score
		if diff1 < 0 {
			diff1 = -diff1
		}
		if diff2 < 0 {
			diff2 = -diff2
		}
		if len(results) == 0 && (i == 0 || diff1 < diff2) {
			newExam = exams[i]
			break
		}
	}
	out, _ := json.Marshal(newExam.Id)
	fmt.Fprint(w, string(out))
}
