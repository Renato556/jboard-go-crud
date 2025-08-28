package routers

import (
	http2 "jboard-go-crud/src/controllers"
	"net/http"
)

func NewJobsController(jobHandler *http2.JobHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", jobHandler.CreateJob)
	mux.HandleFunc("PUT /jobs", jobHandler.UpdateJob)
	mux.HandleFunc("DELETE /jobs", jobHandler.DeleteJob)
	return mux
}
