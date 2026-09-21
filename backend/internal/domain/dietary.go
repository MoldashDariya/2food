package domain

type DietaryTag string

const (
	DietaryTagHalal      DietaryTag = "halal"
	DietaryTagVegetarian DietaryTag = "vegetarian"
	DietaryTagVegan      DietaryTag = "vegan"
	DietaryTagGlutenFree DietaryTag = "gluten_free"
	DietaryTagKosher     DietaryTag = "kosher"
)
