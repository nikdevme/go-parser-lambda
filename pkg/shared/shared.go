package shared

type ProductParse struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Type  string `json:"type"`
	Size  string `json:"size"`
}

type ProductArray struct {
	Array [4]ProductParse `json:"array"`
}

type ProductDescription struct {
	Product
	Description1 string `json:"description-1"`
	Description2 string `json:"description-2"`
}

type Product struct {
	Page       string `json:"page"`
	Date       string `json:"date"`
	Name       string `json:"name"`
	Parameter1 string `json:"parameter-1"`
	Parameter2 string `json:"parameter-2"`
	Parameter3 string `json:"parameter-3"`
}
