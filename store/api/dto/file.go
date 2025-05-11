package dto

type FilePartsDto struct {
	Parts []PartDiapasonDto `json:"parts"`
}

type PartDiapasonDto struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type WritePartsResultDto struct {
	Results []WriteResultDto `json:"results"`
}

type WriteResultDto struct {
	ID    string `json:"id"`
	Error string `json:"error"`
}
