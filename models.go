package main

import "time"

// Voting ...
type Voting struct {
	ID        string
	Name      string
	Deadline  time.Time
	Questions []Question
	Voters    []Voter
}

// Question ...
type Question struct {
	ID      string
	Name    string
	Options []Option
}

// Option ...
type Option struct {
	ID   string
	Name string
}

// Answer ...
type Answer struct {
	IDQuestion string
	IDOption   string
}

// Voter ...
type Voter struct {
	Name    string
	Answers []Answer
}
