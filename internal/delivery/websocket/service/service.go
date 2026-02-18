package service

type ParamsWsMessage struct {
	Line string `json:"line"`
}

type ParamsWsError struct {
	Error string `json:"error"`
}

type ParamsOutput struct {
	Message string `json:"message"`
}
