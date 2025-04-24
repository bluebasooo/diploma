package base_server

import (
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	answer := `
	{
		"message": "Hello, world!"
	}
	`

	w.Write([]byte(answer))
}
