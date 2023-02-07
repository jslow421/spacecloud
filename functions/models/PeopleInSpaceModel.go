package models

type PeopleInSpace struct {
	UpdatedTime string   `json:"updatedTime"`
	People      []Person `json:"people"`
}
