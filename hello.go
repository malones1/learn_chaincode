package main

import (
	"fmt"
	"encoding/json"
	"time"
	"github.com/satori/go.uuid"
	"io"
	"net/http"
	"os"
)

type Option struct {
	Id string
	Name string
}

type Answer struct {
	Id_question string
	Id_option string
}

type Voter struct {
	Name string
	Answers []Answer
}

type Question struct {
	Id string
	Name string
	Options []Option
}

type Voting struct {
	Id string
	Name string
	Deadline time.Time
	Questions []Question
	Voters []Voter
}

type State struct {
	v []Voting
}

func getUuid() string {
	u, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return ""
	} else {
		return u.String()
	}
}

// ---
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, string(getJson()))
}

func AddVoteHandler(w http.ResponseWriter, req *http.Request) {
	AddVote()
	io.WriteString(w, "add vote...")
	fmt.Println(req.Method)
}

func AddVote() {
	var v Voting
	var qs [3]Question
	// ---
	v = Voting{}
	v.Id = getUuid()
	v.Name = "_V1"
	v.Deadline = time.Now()
	v.Voters = []Voter{}

	qs[0] = Question{getUuid(), "Q1", []Option{ Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"} }}
	qs[1] = Question{getUuid(), "Q2", []Option{ Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"} }}
	qs[2] = Question{getUuid(), "Q3", []Option{ Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"} }}

	v.Questions = qs[:3]
	s.v = append(s.v, v)
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	file, err := os.Open("index.html")
	if err != nil {
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return
	}

	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		return
	}

	w.Header().Set(
  	"Content-Type",
    "text/html",
  )
  io.WriteString(w, string(bs))
}

// +++

func Init() {
	if len(s.v) == 0 {
		var v Voting

		// ---
		v = Voting{}
		v.Id = getUuid()
		v.Name = "V1"
		v.Deadline = time.Now()
		v.Voters = []Voter{}

		v.Questions = []Question{
				Question{getUuid(), "123Q1", []Option{ Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"} }},
				Question{getUuid(), "123Q2", []Option{ Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"} }},
				Question{getUuid(), "123Q3", []Option{ Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"} }},
		}
		s.v = append(s.v, v)
		// +++

		// ---
		v = Voting{}
		v.Id = getUuid()
		v.Name = "V2"
		v.Deadline = time.Now()
		v.Voters = []Voter{}

		v.Questions = []Question{
				Question{getUuid(), "456Q1", []Option{ Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"} }},
				Question{getUuid(), "456Q2", []Option{ Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"} }},
				Question{getUuid(), "456Q3", []Option{ Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"} }},
		}
		s.v = append(s.v, v)
		// +++

	}
}

func getJson() []byte {
	b, err := json.Marshal(s.v)
	if err == nil {
		return b
	} else {
		return []byte("")
	}
}

var s State = State{}

func main() {
	Init()

	http.HandleFunc("/vote", HelloServer)
	http.HandleFunc("/add", AddVoteHandler)
	http.HandleFunc("/i", IndexHandler)
	http.ListenAndServe(":7777", nil)
}
