package domain

type RequestAdd struct {
	Passport string `json:"passportNumber"`
}

type ResponseAdd struct {
	PassportSeries int `json:"passportSeries"`
	PassportNumber int `json:"passportNumber"`
}

type User struct {
	ID             int    `json:"id"`
	PassportSeries int    `json:"passportSeries"`
	PassportNumber int    `json:"passportNumber"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}
