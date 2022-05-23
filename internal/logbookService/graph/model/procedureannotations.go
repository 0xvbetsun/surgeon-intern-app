package qlmodel

type ProcedureAnnotations struct {
	Methods []Method `json:"methods"`
}

type Method struct {
	DisplayName string `json:"displayName"`
	Aids        []Aid  `json:"aids"`
}

type Aid struct {
	DisplayName string  `json:"displayName"`
	Routes      []Route `json:"routes"`
}

type Route struct {
	DisplayName string    `json:"displayName"`
	Purposes    []Purpose `json:"purposes"`
}

type Purpose struct {
	DisplayName string    `json:"displayName"`
	Findings    []Finding `json:"findings"`
}

type Finding struct {
	DisplayName string `json:"displayName"`
}
