package model

type User struct {
	Id int64
	Username string
	Password string
	TotalMoney float64
	UseMoney  float64
	Money   float64
	TotalScore float64
	UseScore  float64
	Score   float64
	CreatedTime int64
	UpdatedTime int64
}
func(*User) TableName() string {
	return "user"
}
