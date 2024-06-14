package entity

const (
	Like = "right"
	Pass = "left"
)

type Discovery struct {
	ID       int64    `json:"id"`
	FullName string   `json:"fullname"`
	Age      int64    `json:"age"`
	Gender   string   `json:"gender"`
	Bio      string   `json:"bio"`
	Location string   `json:"location"`
	Interest string   `json:"interest"`
	Photos   []string `json:"photos"`
}
