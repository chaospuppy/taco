package taco

import "net/http"

func NewServer(config Config) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		config,
	)
	var handler http.Handler = mux
	return handler
}

func addRoutes(mux *http.ServeMux, config Config) error {
	mux.HandleFunc("/taco", handleTacos)
	return nil
}

func handleTacos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applicaiton/json")
	// if err := json.NewEncoder(w).Encode(message)
}
