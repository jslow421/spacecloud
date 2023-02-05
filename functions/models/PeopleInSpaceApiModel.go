package models

type PersonInSpaceApiResponse struct {
	Message string   `json:"message"`
	People  []Person `json:"people"`
}
