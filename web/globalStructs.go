package web

type StandartResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func contains[T comparable, A []T](s A, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
