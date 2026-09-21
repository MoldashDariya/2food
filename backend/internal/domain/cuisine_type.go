package domain

type CuisineType struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Label string `json:"label"`
}

type CreateCuisineTypeInput struct {
	Name  string
	Label string
}

func NewCuisineType(input CreateCuisineTypeInput) *CuisineType {
	return &CuisineType{
		Name:  input.Name,
		Label: input.Label,
	}
}
