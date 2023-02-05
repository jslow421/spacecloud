package models

type PersonInSpace struct {
	UpdatedTime string   `json:"updatedTime"`
	People      []Person `json:"people"`
}
