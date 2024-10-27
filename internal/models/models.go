package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Question struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Options []Option `json:"options"`
}

type Option struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
}

type Answer struct {
	QuestionID int `json:"questionId"`
	OptionID   int `json:"optionId"`
}

type Submission struct {
	UserName string   `json:"userName"`
	Answers  []Answer `json:"answers"`
}

type Result struct {
	UserName string `json:"userName"`
	Score    int    `json:"score"`
}

type SubmissionResponse struct {
	Message            string  `json:"message"`
	Score              int     `json:"score"`
	Rank               float64 `json:"rank"`
	TotalQuestionCount int     `json:"totalQuestionCount"`
}

func (s Submission) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.UserName, validation.Required),
		validation.Field(&s.Answers, validation.Required, validation.Each(validation.By(func(value interface{}) error {
			answer, ok := value.(Answer)
			if !ok {
				return validation.ErrInInvalid
			}
			return validation.ValidateStruct(&answer,
				validation.Field(&answer.QuestionID, validation.Required),
				validation.Field(&answer.OptionID, validation.Required),
			)
		}))),
	)
}
