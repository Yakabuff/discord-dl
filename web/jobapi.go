package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/yakabuff/discord-dl/models"
)

func (web Web) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	respondwithJSON(w, http.StatusCreated, web.JobQueue.Jobs)
}

func (web Web) GetJobByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "jobID"))
	j := web.JobQueue.Jobs[id]
	respondwithJSON(w, http.StatusCreated, j)
}

// func (web Web) GetJobBySnowflake(w http.ResponseWriter, r *http.Request) {

// }

func (web Web) SubmitJob(w http.ResponseWriter, r *http.Request) {
	var j models.JobArgs
	json.NewDecoder(r.Body).Decode(&j)
	// job := job.NewJob(j)
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
	// if job.Args.Mode != models.CHANNEL && job.Args.Mode != models.GUILD {
	// 	respondWithError(w, http.StatusInternalServerError, "Invalid mode")
	// 	return
	// }
	// err := web.JobQueue.Enqueue(job)

	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Job could not be submitted. Queue is full")
	// }
	// respondwithJSON(w, http.StatusCreated, map[string]string{"message": "job successfully created"})

}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondwithJSON(w, code, map[string]string{"message": msg})
}

// respondwithJSON write json response format
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (web Web) CancelJob(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "jobID"))
	j := web.JobQueue
	err := j.CancelJob(id)
	if err != nil {
		respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Succesfully deleted job"})
	} else {
		respondWithError(w, http.StatusInternalServerError, "Failed to cancel job. Invalid job ID")
	}
}

func (web Web) GetJobProgress(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "jobID"))
	j := web.JobQueue.Jobs[id]
	if j != nil {
		respondwithJSON(w, http.StatusCreated, map[string]string{"progress": fmt.Sprintf("%d", j.Progress)})
	} else {
		respondWithError(w, http.StatusInternalServerError, "Invalid job. Could not fetch progress")
	}
}

func (web Web) ShowJobPanel(w http.ResponseWriter, r *http.Request) {
	log.Println("hi")
	jobs := web.JobQueue.GetAllJobs()

	tmpl, err := template.ParseFS(templates, "static/job.html")
	if err != nil {
		log.Println(err)
	}
	tmpl.Execute(w, jobs)
}
