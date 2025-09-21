package utils

type ErrorResponse struct {
	Error string 
	Message string 
	Code string
}

type SuccessResponse struct {
	Data interface{}
	Message string
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := ErrorResponse{
		Error: http.StatusText(statusCode)
		Message: message,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to write error respponse: %v", err)
	}
}

func writeSuccess() {

}