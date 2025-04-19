package dto

type FilePartsDto struct {
	Parts []PartDiapasonDto `json:"parts"`
}

type PartDiapasonDto struct {
	Start int `json:"start"`
	End   int `json:"end"`
}
