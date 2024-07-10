package domain

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type RequestAdd struct {
	Passport string `json:"passportNumber"`
}

type ResponseAdd struct {
	PassportSeries int `json:"passportSeries"`
	PassportNumber int `json:"passportNumber"`
}

type User struct {
	ID             int64  `json:"id"`
	PassportSeries int    `json:"passportSeries"`
	PassportNumber int    `json:"passportNumber"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}

type UserSortParameters struct {
	PassportSeries int    `json:"passportSeries"`
	PassportNumber int    `json:"passportNumber"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}

type Metrics struct {
	User_ID  uuid.UUID     `json:"user_id"`
	FuncName string        `json:"func_name"`
	Time     time.Duration `json:"time"`
}
