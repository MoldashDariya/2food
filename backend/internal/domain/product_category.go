package domain

type ProductCategory struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Label string `json:"label"`
}

type CreateProductCategoryInput struct {
	Name  string
	Label string
}

func NewProductCategory(input CreateProductCategoryInput) *ProductCategory {
	return &ProductCategory{
		Name:  input.Name,
		Label: input.Label,
	}
}
