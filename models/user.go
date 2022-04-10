package models

type Address struct {
	City    string `json:"city" bson:"city"`
	State   string `json:"state" bson:"state"`
	Pincode int    `json:"pincode" bson:"pincode"`
}

type User struct {
	Name    string  `json:"name" bson:"name"`
	Age     int     `json: "age" bson:"age"`
	Address Address `json:"address" bson:"address"`
}
