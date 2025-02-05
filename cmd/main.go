package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aluko123/resume-go/internal/handlers"
	"github.com/aluko123/resume-go/internal/repository"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	//connect MongoDB
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// // After MongoDB connection
	// if err := client.Ping(context.Background(), nil); err != nil {
	// 	log.Fatal("Failed to connect to MongoDB:", err)
	// } else {
	// 	log.Println("Successfully connected to MongoDB")
	// }

	db := client.Database("resume_db")

	//intialize repository
	resumeRepo := repository.NewResumeRepository(db)

	//initialize handlers
	resumeHandler := handlers.NewResumeHandler(resumeRepo)

	//create router
	router := mux.NewRouter()

	//setup routes
	router.HandleFunc("/resumes", resumeHandler.CreateResume).Methods("POST")
	router.HandleFunc("/resumes", resumeHandler.ListResumes).Methods("GET")
	router.HandleFunc("/resumes/{id}", resumeHandler.GetResume).Methods("GET")
	router.HandleFunc("/resumes/{id}", resumeHandler.UpdateResume).Methods("PUT")
	router.HandleFunc("/resumes/{id}", resumeHandler.DeleteResume).Methods("DELETE")
	router.HandleFunc("/resumes/{id}/education", resumeHandler.UpdateEducation).Methods("PUT")
	router.HandleFunc("/resumes/{id}/experience", resumeHandler.UpdateExperience).Methods("PUT")

	//start server
	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
