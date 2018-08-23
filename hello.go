package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	_ "log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"regexp"
	"time"
)

// APIHandler ...
type APIHandler string

func (a APIHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Hello from ApiHandler...")
}

// HelloServer ...
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, string(getJSON()))
}

// AddVoteHandler ...
func AddVoteHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	v := AddVote()
	if b, err := json.Marshal(v); err == nil {
		io.WriteString(w, string(b))
	}
	fmt.Println(req.Method)
}

// AddVote ...
func AddVote() Voting {
	var v Voting
	var qs [3]Question
	// ---
	v = Voting{}
	v.ID = GetUUID()
	v.Name = "_V1"
	v.Deadline = time.Now()
	v.Voters = []Voter{}

	qs[0] = Question{GetUUID(), "Q1", []Option{Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"}}}
	qs[1] = Question{GetUUID(), "Q2", []Option{Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"}}}
	qs[2] = Question{GetUUID(), "Q3", []Option{Option{"1", "O1"}, Option{"2", "O2"}, Option{"3", "O3"}}}

	v.Questions = qs[:3]
	app.v = append(app.v, &v)

	return v
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
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

// 1
func getJSON() []byte {
	b, err := json.Marshal(app.v)
	if err == nil {
		return b
	}

	return []byte("")
}

var app = State{}

func main() {
	app.CreateTestData()

	mux := http.NewServeMux()
	var a APIHandler
	mux.Handle("/test/", a)
	mux.HandleFunc("/vote/", func(w http.ResponseWriter, req *http.Request) {
		var re = regexp.MustCompile(`(?is)vote/(.+)`)
		fmt.Println("Hello vote...")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		var url = req.URL.String()
		fmt.Printf("%v\n", req.RequestURI)

		if len(re.FindStringIndex(url)) > 0 {
			fmt.Println(re.FindString(url), "found at index", re.FindStringIndex(url)[0])
			http.NotFound(w, req)

			return
		}

		HelloServer(w, req)
	})

	mux.HandleFunc("/toVote", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
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

				for _, v := range app.v {
					if v.ID == m["votingId"] {
						fmt.Printf("VotingId: %v, voters: %v\n", v.ID, v.Voters)
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
			indexHandler(w, req)
		}
	})
	mux.HandleFunc("/add", AddVoteHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))
	http.ListenAndServe(":7777", mux)
}
