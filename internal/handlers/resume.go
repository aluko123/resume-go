package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/aluko123/resume-go/internal/models"
	"github.com/aluko123/resume-go/internal/repository"
)

// handle HTTP requests for resumes
type ResumeHandler struct {
	repo *repository.ResumeRepository
}

// create a new handler
func NewResumeHandler(repo *repository.ResumeRepository) *ResumeHandler {
	return &ResumeHandler{
		repo: repo,
	}
}

// handle POST requests to create a new resume
func (h *ResumeHandler) CreateResume(w http.ResponseWriter, r *http.Request) {
	//log.Println("CreateResume handler called")

	var resume models.Resume
	if err := json.NewDecoder(r.Body).Decode(&resume); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//log.Printf("Attempting to create resume: %+v", resume)

	if err := h.repo.Create(r.Context(), &resume); err != nil {
		log.Printf("Error creating resume: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully create resume with ID: %d", resume.ResumeID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resume)
}

// handles GET requests
func (h *ResumeHandler) GetResume(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//log.Printf("VARS: %v", vars)
	id := vars["id"]
	//log.Printf("Resume_ID: %s", id)

	resume, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resume)
}

// update resume PUT request
func (h *ResumeHandler) UpdateResume(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var resume models.Resume
	if err := json.NewDecoder(r.Body).Decode(&resume); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//ensure url ID matches resume ID
	objectID, _ := primitive.ObjectIDFromHex(id)
	resume.ID = objectID

	if err := h.repo.Update(r.Context(), &resume); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resume)
}

// update education using PUT request
func (h *ResumeHandler) UpdateEducation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var education []models.Education
	if err := json.NewDecoder(r.Body).Decode(&education); err != nil {
		log.Printf("Error decoding: %v", education)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateEducation(r.Context(), id, education); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

// update skills using PUT request
func (h *ResumeHandler) UpdateSkills(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var skills models.Skills
	if err := json.NewDecoder(r.Body).Decode(&skills); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateSkills(r.Context(), id, skills); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)

}

// update experience using PUT request
func (h *ResumeHandler) UpdateExperience(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var experience []models.Experience
	if err := json.NewDecoder(r.Body).Decode(&experience); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateExperience(r.Context(), id, experience); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

// delete resume using DELETE request
func (h *ResumeHandler) DeleteResume(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusNoContent)
}

// list resumes GET request for list of resume
func (h *ResumeHandler) ListResumes(w http.ResponseWriter, r *http.Request) {
	resumes, err := h.repo.ListAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resumes)
}

func (h *ResumeHandler) ListAllResumes(w http.ResponseWriter, r *http.Request) {
	log.Println("ListResumes handler called") // Add this

	resumes, err := h.repo.ListAll(r.Context())
	if err != nil {
		log.Printf("Error in ListAll: %v", err) // Add this
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Found %d resumes", len(resumes)) // Add this

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resumes)
}
