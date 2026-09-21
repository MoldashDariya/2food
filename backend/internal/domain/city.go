package domain

type City struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Label string `json:"label"`
}

type CreateCityInput struct {
	Name  string
	Label string
}

func NewCity(input CreateCityInput) *City {
	return &City{
		Name:  input.Name,
		Label: input.Label,
	}
}