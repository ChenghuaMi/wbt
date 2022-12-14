package model

type Order struct {
	Id    		int64
	Type        byte
	OrderNo     string
	Price       float64
	Num         int64
	UserId      int64
	Status      byte
	CreatedTime int64
	UpdatedTime int64
}

func(*Order) TableName() string {
	return "order"
}
