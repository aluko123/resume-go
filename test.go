package main

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
	newResumeData.Education = []models.Education{}
	addMoreEducation := true
	for addMoreEducation {
		var educationEntry models.Education
		err := survey.Ask([]*survey.Question{
			{
				Name:      "institution",
				Prompt:    &survey.Input{Message: "Institution?"},
				Validate:  survey.Required,
				Transform: survey.Title,
			},

			{
				Name:      "degree",
				Prompt:    &survey.Input{Message: "Degree?"},
				Validate:  survey.Required,
				Transform: survey.Title,
			},
			{
				Name:      "startdate",
				Prompt:    &survey.Input{Message: "Start Date (e.g. Aug 2024)?"},
				Validate:  survey.Required,
				Transform: survey.Title,
			},
			{
				Name:      "enddate",
				Prompt:    &survey.Input{Message: "End Date (e.g. Aug 2024)?"},
				Validate:  survey.Required,
				Transform: survey.Title,
			},
			{
				Name:      "gpa",
				Prompt:    &survey.Input{Message: "GPA?"},
				Transform: survey.Title,
			},
		}, &educationEntry)

		if err != nil {
			fmt.Println("error getting education info %w", err)
		}

		educationEntry.Coursework = []string{}

		addCoursework := true
		for addCoursework {
			course := ""
			err := &survey.Input{
				Message: "Enter coursework (or press Enter to finish):",
			}
			survey.AskOne(err, &course)

			if course != "" {
				educationEntry.Coursework = append(educationEntry.Coursework, course)
			} else {
				addCoursework = false //user pressed enter
			}
		}

		//fmt.Println(newResumeData.BasicInfo)
		newResumeData.Education = append(newResumeData.Education, educationEntry)
		val := ""
		prompt := &survey.Input{
			Message: "Add another education entry? (True or false)",
		}
		survey.AskOne(prompt, &val)
		if val == "false" {
			addMoreEducation = false //user doesn't want more education entry
		}
	}

	bytes, err := json.MarshalIndent(newResumeData.Education, "", "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bytes))

}
