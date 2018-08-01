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

func (v *Voting) SetID() {
	v.ID = GetUUID()
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

// State ...
type State struct {
	v []*Voting
}

func (s *State) CreateTestData() {
	if len(s.v) == 0 {
		var v Voting

		v = Voting{Name: "Опрос - Выбор ноутбука", Deadline: time.Now(), Voters: []Voter{}}
		v.SetID()

		v.Questions = []Question{
			Question{
				GetUUID(), 
				"Какой объем оперативной памяти предпочтительнее?", 
				[]Option{
					Option{"1", "4 Гб"}, 
					Option{"2", "8 Гб"}, 
					Option{"3", "Более 8 Гб"},
				},
			},
			Question{
				GetUUID(), 
				"Фирма производитель?", 
				[]Option{
					Option{"1", "Asus"}, 
					Option{"2", "Acer"}, 
					Option{"3", "Apple"}, 
					Option{"4", "Другие"},
				},
			},
			Question{
				GetUUID(), 
				"Диагональ экрана?", 
				[]Option{
					Option{"1", "до 13\""}, 
					Option{"2", "от 13\" до 14\""}, 
					Option{"3", "15\" и более"},
				},
			},
		}
		s.v = append(s.v, &v)
	}
}