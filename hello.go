package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/satori/go.uuid"
)

// Option ...
type Option struct {
	Id   string
	Name string
}

// Answer ...
type Answer struct {
	Id_question string
	Id_option   string
}

// Voter ...
type Voter struct {
	Name    string
	Answers []Answer
}

// Question ...
type Question struct {
	Id      string
	Name    string
	Options []Option
}

// Voting ...
type Voting struct {
	Id        string
	Name      string
	Deadline  time.Time
	Questions []Question
	Voters    []Voter
}

// State ...
type State struct {
	v []*Voting
}

func getUUID() string {
	u, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return ""
	} else {
		return u.String()
	}
}

// HelloServer ...
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, string(getJSON()))
}

// AddVoteHandler ...
func AddVoteHandler(w http.ResponseWriter, req *http.Request) {
	AddVote()
	io.WriteString(w, "add vote...")
	fmt.Println(req.Method)
}

// AddVote ...
func AddVote() {
	var v Voting
	var qs [3]Question
	// ---
	v = Voting{}
	v.Id = getUUID()
	v.Name = "_V1"
	v.Deadline = time.Now()
	v.Voters = []Voter{}

	qs[0] = Question{getUUID(), "Q1", []Option{Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"}}}
	qs[1] = Question{getUUID(), "Q2", []Option{Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"}}}
	qs[2] = Question{getUUID(), "Q3", []Option{Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"}}}

	v.Questions = qs[:3]
	s.v = append(s.v, &v)
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
		var v *Voting

		// ---
		v = new(Voting)
		v.Id = getUUID()
		v.Name = "V1"
		v.Deadline = time.Now()
		v.Voters = []Voter{}

		v.Questions = []Question{
			Question{getUUID(), "123Q1", []Option{Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"}}},
			Question{getUUID(), "123Q2", []Option{Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"}}},
			Question{getUUID(), "123Q3", []Option{Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"}}},
		}
		s.v = append(s.v, v)
		// +++

		// ---
		v = new(Voting)
		v.Id = getUUID()
		v.Name = "V2"
		v.Deadline = time.Now()
		v.Voters = []Voter{}

		v.Questions = []Question{
			Question{getUUID(), "456Q1", []Option{Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"}}},
			Question{getUUID(), "456Q2", []Option{Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"}}},
			Question{getUUID(), "456Q3", []Option{Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"}}},
		}
		s.v = append(s.v, v)
		// +++

	}
}

// 1
func getJSON() []byte {
	b, err := json.Marshal(s.v)
	if err == nil {
		return b
	} else {
		return []byte("")
	}
}

var s = State{}

func main() {
	Init()

	mux := http.NewServeMux()
	mux.HandleFunc("/vote", HelloServer)

	mux.HandleFunc("/toVote", func(w http.ResponseWriter, req *http.Request) {
		//io.WriteString(w, "Hello...")

		if req.Method == "POST" {
			var f interface{}
			var b bytes.Buffer

			buf := make([]byte, 4)
			for {
				n, err := req.Body.Read(buf)
				if n > 0 {
					b.Write(buf[:n])
				}

				if err == io.EOF {
					break
				}
			}

			err := json.Unmarshal(b.Bytes(), &f)
			if err != nil {
				fmt.Println(err)
			} else {
				m := f.(map[string]interface{})

				for _, v := range s.v {
					if (v.Id == m["votingId"]) {
						fmt.Printf("VotingId: %v, voters: %v\n", v.Id, v.Voters);
						v.Voters = append(v.Voters, Voter{"Roman", []Answer{}})
					}
				}
			}
		} else {
			http.NotFound(w, req)
		}

		//fmt.Fprintf(w, "Hello... %v", req.Method)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
		} else {
			IndexHandler(w, req)
		}
	})
	mux.HandleFunc("/add", AddVoteHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":7777", mux)
}
