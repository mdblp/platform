package food

import (
	"github.com/tidepool-org/platform/data"
	"github.com/tidepool-org/platform/structure"
)

const (
	SmallMeal   = "small"
	MediumMeal  = "medium"
	LargeMeal   = "large"
	Snack       = "yes"
	NotASnack   = "no"
	FatMeal     = "yes"
	NotAFatMeal = "no"
)

func MealSize() []string {
	return []string{
		SmallMeal,
		MediumMeal,
		LargeMeal,
	}
}

func IsASnack() []string {
	return []string{
		Snack,
		NotASnack,
	}
}

func IsFat() []string {
	return []string{
		FatMeal,
		NotAFatMeal,
	}
}

type Meal struct {
	Meal  *string `json:"meal,omitempty" bson:"meal,omitempty"`
	Snack *string `json:"snack,omitempty" bson:"snack,omitempty"`
	Fat   *string `json:"fat,omitempty" bson:"fat,omitempty"`
}

func ParseMeal(parser structure.ObjectParser) *Meal {
	if !parser.Exists() {
		return nil
	}
	datum := NewMeal()
	parser.Parse(datum)
	return datum
}

func NewMeal() *Meal {
	return &Meal{}
}

func (m *Meal) Parse(parser structure.ObjectParser) {
	m.Meal = parser.String("meal")
	m.Snack = parser.String("snack")
	m.Fat = parser.String("fat")
}

func (m *Meal) Validate(validator structure.Validator) {
	if m.Meal != nil {
		validator.String("meal", m.Meal).Exists().OneOf(MealSize()...)
	}
	if m.Snack != nil {
		validator.String("snack", m.Snack).Exists().OneOf(IsASnack()...)
	}
	if m.Fat != nil {
		validator.String("fat", m.Fat).Exists().OneOf(IsFat()...)
	}
}

func (m *Meal) Normalize(normalizer data.Normalizer) {}
