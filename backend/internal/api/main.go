package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server running")

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(map[string]string{
		"msg": "Server running.",
	})
	if err != nil {
		json.NewEncoder(w).Encode("Internal Server Error in /test")
	}

}
