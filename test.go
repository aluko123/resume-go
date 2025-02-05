package test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aluko123/resume-go/internal/models"
)

// var newResumeData map[string]interface{} = make(map[string]interface{})
var newResumeData models.Resume

func main() {
	//Basic Info
	err := survey.Ask([]*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Name?"},
			Validate: survey.Required,
			//Transform: survey.Title,
		},

		{
			Name:     "email",
			Prompt:   &survey.Input{Message: "Email?"},
			Validate: survey.Required,
			//Transform: survey.Title,
		},
		{
			Name:   "phone",
			Prompt: &survey.Input{Message: "Phone Number?"},
			//Transform: survey.Title,
		},
		{
			Name:   "linkedin",
			Prompt: &survey.Input{Message: "Linkedin?"},
			//Transform: survey.Title,
		},
		{
			Name:   "github",
			Prompt: &survey.Input{Message: "GitHub?"},
			//Transform: survey.Title,
		},
	}, &newResumeData.BasicInfo)

	if err != nil {
		fmt.Println("error getting basic info %w", err)
	}

	//fmt.Println(newResumeData.BasicInfo)

	bytes, err := json.MarshalIndent(newResumeData.BasicInfo, "", "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bytes))

}
