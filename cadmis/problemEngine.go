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
	"time"
)

type Entry struct {
    ProblemId int64
    UserId    int64
    Result    bool
    Date      time.Time
}

type Problem struct {
    Id       int64
    Category string
    Score    int64
}


func handleComputeRank(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    uid, err1 := strconv.Atoi(r.FormValue("userId"));
    pid, err2 := strconv.Atoi(r.FormValue("problemId"));
    if err1 != nil {
        uid = -1;
    }
    if err2 != nil {
        pid = -1;
    }
    score, err3 := computeRank(c, int64(uid), int64(pid))
    if err3 != nil {
        http.Error(w, err3.Error(), http.StatusInternalServerError)
        return;
    }
    out, _ := json.Marshal(score)
    fmt.Fprint(w, string(out))
}

func handleProblemList(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    q := datastore.NewQuery("Problem")
    var problems []Problem
    _, err := q.GetAll(c, &problems)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return;
    }
    ids := make([]int64, len(problems), len(problems))
    for i := range problems {
        ids[i] = problems[i].Id;
    }
    out, _ := json.Marshal(ids)
    fmt.Fprint(w, string(out))
}

func handleUserList(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    q := datastore.NewQuery("User")
    var users []User
    _, err := q.GetAll(c, &users)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return;
    }
    ids := make([]int64, len(users), len(users))
    for i := range users {
        ids[i] = users[i].Id;
    }
    out, _ := json.Marshal(ids)
    fmt.Fprint(w, string(out))
}


func handleWriteEntry(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    uid, err1 := strconv.Atoi(r.FormValue("userId"));
    pid, err2 := strconv.Atoi(r.FormValue("problemId"));
    res, err3 := strconv.ParseBool(r.FormValue("result"));
    if err1 != nil {
        http.Error(w, err1.Error(), http.StatusInternalServerError)
        return;
    }
    if err2 != nil {
        http.Error(w, err2.Error(), http.StatusInternalServerError)
        return;
    }
    if err3 != nil {
        http.Error(w, err3.Error(), http.StatusInternalServerError)
        return;
    }

    e := Entry{
        ProblemId: int64(pid),
        UserId:    int64(uid),
        Result:    res,
        Date:      time.Now(),
    }
    _, err4 := datastore.Put(c, datastore.NewIncompleteKey(c, "Entry", nil), &e)
    if err4 != nil {
        http.Error(w, err4.Error(), http.StatusInternalServerError)
        return;
    }
}

func handleInitProblems(w http.ResponseWriter, r *http.Request) {
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
    resp, err := client.Get("http://api.iknow.jp/goals/469230/items")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    defer resp.Body.Close()
    body, err2 := ioutil.ReadAll(resp.Body)
    if err2 != nil {
        http.Error(w, err2.Error(), http.StatusInternalServerError)
    }
    var items map[string][]map[string]interface{}
    err3 := json.Unmarshal(body, &items)
    if err3 != nil {
        http.Error(w, err3.Error(), http.StatusInternalServerError)
    }
    items2 := items["items"];
    for i := range items2 {
        //fmt.Fprintf(w, "%d: %s\n", int(items2[i]["id"].(float64)), items2[i]["uri"])
        updateProblem(c, int64(items2[i]["id"].(float64)), "English", 10);
    }
}

func updateProblem(c appengine.Context, problemId int64, category string, score int64) error {
    q := datastore.NewQuery("Problem").Filter("Id =", problemId).Limit(1)
    problems := make([]Problem, 0, 1)
    keys, err := q.GetAll(c, &problems)    
    if err == nil {
        problems := make([]Problem, 1, 1)
        problems[0].Id = problemId
        problems[0].Category = category;
        problems[0].Score = score;
        _, err2 := datastore.Put(c, datastore.NewIncompleteKey(c, "Problem", nil), &problems[0])
        if err2 != nil {
            return err2
        }
    } else {
        problems[0].Score = score;    
        _, err2 := datastore.Put(c, keys[0], &problems[0])
        if err2 != nil {
            return err2
        }
    }
    return nil;
}

func computeRank(c appengine.Context, userId int64, problemId int64) (int64, error) {
    var score int64
    q := datastore.NewQuery("Entry")
    for t := q.Run(c); ; {
        var e Entry;
        _, err := t.Next(&e)
        if err == datastore.Done {
            break;
        }
        if err != nil {
            return score, err
        }
        if e.Result && (userId == -1 || e.UserId == userId) && (problemId == -1 || e.ProblemId == problemId) {
            q := datastore.NewQuery("Problem").Filter("Id =", e.ProblemId).Limit(1)
            problems := make([]Problem, 0, 1)
            if _, err2 := q.GetAll(c, &problems); err != nil {
                return score, err2
            }
            score += problems[0].Score
        }
    }
    return score, nil
}

/*
func getSuitableProblem(c appengine.Context, userId int64) {
    q := datastore.NewQuery("Problem").Order("-Date").Limit(10)
    problems := make([]Problem, 0, 10)
    if _, err := q.GetAll(c, &problems); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    for (p : problems) {
        if (p.Result) {
            q := datastore.NewQuery("Problem");
            if _, err := q.GetAll(c, &problems); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            score += p.Score
        }
    }

}
*/
