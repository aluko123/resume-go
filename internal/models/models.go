package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Resume struct to start
type Resume struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ResumeID     int                `bson:"resume_id" json:"resume_id"`
	BasicInfo    BasicInfo          `bson:"basic_info" json:"basic_info"`
	Education    []Education        `bson:"education" json:"education"`
	Skills       Skills             `bson:"skills" json:"skills"`
	Experience   []Experience       `bson:"experience" json:"experience"`
	Projects     []Project          `bson:"projects" json:"projects"`
	Affiliations []Affiliation      `bson:"affiliations" json:"affiliations"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// define basic info struct containing personal information
type BasicInfo struct {
	Name     string `bson:"name" json:"name"`
	Email    string `bson:"email" json:"email"`
	Phone    string `bson:"phone" json:"phone"`
	LinkedIn string `bson:"linkedin" json:"linkedin"`
	GitHub   string `bson:"github" json:"github"`
}

// education structs
type Education struct {
	Institution string   `bson:"institution" json:"institution"`
	Degree      string   `bson:"degree" json:"degree"`
	StartDate   string   `bson:"start_date" json:"start_date"`
	EndDate     string   `bson:"end_date" json:"end_date"`
	GPA         float64  `bson:"gpa,omitempty" json:"gpa,omitempty"`
	Coursework  []string `bson:"coursework" json:"coursework"`
}

// skills
type Skills struct {
	ProgrammingLanguages []string `bson:"programming_languages" json:"programming_languages"`
	Technologies         []string `bson:"technologies" json:"technologies"`
}

// experience list
type Experience struct {
	Company     string   `bson:"company" json:"company"`
	Title       string   `bson:"title" json:"title"`
	Location    string   `bson:"location" json:"location"`
	StartDate   string   `bson:"start_date" json:"start_date"`
	EndDate     string   `bson:"end_date" json:"end_date"`
	Description []string `bson:"description" json:"description"`
}

// projects struct
type Project struct {
	Name         string   `bson:"name" json:"name"`
	Description  string   `bson:"description" json:"description"`
	Technologies []string `bson:"technologies" json:"technologies"`
	Link         string   `bson:"link,omitempty" json:"link,omitempty"`
}

// professional affiliations
type Affiliation struct {
	Organization string `bson:"organization" json:"organization"`
	Role         string `bson:"role" json:"role"`
	StartDate    string `bson:"start_date" json:"start_date"`
	EndDate      string `bson:"end_date,omitempty" json:"end_date,omitempty"`
}
