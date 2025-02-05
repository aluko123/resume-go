package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aluko123/resume-go/internal/models"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new resume",
	Long:  `Create a new resume.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		//var newResumeData map[string]interface{} = make(map[string]interface{})
		var newResumeData models.Resume

		//Basic Info
		err := survey.Ask([]*survey.Question{
			{
				Name:      "name",
				Prompt:    &survey.Input{Message: "Name?"},
				Validate:  survey.Required,
				Transform: survey.Title,
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
			return fmt.Errorf("error getting basic info %w", err)
		}

		fmt.Print("Enter Basic Info (JSON format):\n")
		var basicInfo map[string]interface{} = make(map[string]interface{})
		decoder := json.NewDecoder(os.Stdin)
		if err := decoder.Decode(&basicInfo); err != nil {
			return fmt.Errorf("error decoding basic info: %w", err)
		}
		//newResumeData["basic_info"] = basicInfo

		fmt.Print("Enter Education (JSON format, or press Enter to skip):\n")
		var educationUpdates []interface{} = make([]interface{}, 0)
		if _, err := fmt.Scanln(); err == nil {
			decoder := json.NewDecoder(os.Stdin)
			if err := decoder.Decode(&educationUpdates); err != nil {
				return fmt.Errorf("error decoding education: %w", err)
			}
		}
		if len(educationUpdates) > 0 {
			newResumeData["education"] = educationUpdates
		}

		fmt.Print("Enter Experience (JSON format, or press Enter to skip):\n")
		var experienceUpdates []interface{} = make([]interface{}, 0)
		if _, err := fmt.Scanln(); err == nil {
			decoder := json.NewDecoder(os.Stdin)
			if err := decoder.Decode(&experienceUpdates); err != nil {
				return fmt.Errorf("error decoding experience: %w", err)
			}
		}
		if len(experienceUpdates) > 0 {
			newResumeData["experience"] = experienceUpdates
		}

		fmt.Print("Enter Skills (JSON format, or press Enter to skip):\n")
		var skillsUpdates map[string]interface{} = make(map[string]interface{})
		if _, err := fmt.Scanln(); err == nil {
			decoder := json.NewDecoder(os.Stdin)
			if err := decoder.Decode(&skillsUpdates); err != nil {
				return fmt.Errorf("error decoding skills: %w", err)
			}
		}
		if len(skillsUpdates) > 0 {
			newResumeData["skills"] = skillsUpdates
		}

		fmt.Print("Enter Projects (JSON format, or press Enter to skip):\n")
		var projectsUpdates []interface{} = make([]interface{}, 0)
		if _, err := fmt.Scanln(); err == nil {
			decoder := json.NewDecoder(os.Stdin)
			if err := decoder.Decode(&projectsUpdates); err != nil {
				return fmt.Errorf("error decoding projects: %w", err)
			}
		}
		if len(projectsUpdates) > 0 {
			newResumeData["projects"] = projectsUpdates
		}

		fmt.Print("Enter Affiliations (JSON format, or press Enter to skip):\n")
		var affiliationsUpdates []interface{} = make([]interface{}, 0)
		if _, err := fmt.Scanln(); err == nil {
			decoder := json.NewDecoder(os.Stdin)
			if err := decoder.Decode(&affiliationsUpdates); err != nil {
				return fmt.Errorf("error decoding affiliations: %w", err)
			}
		}
		if len(affiliationsUpdates) > 0 {
			newResumeData["affiliations"] = affiliationsUpdates
		}

		newResumeJSON, err := json.Marshal(newResumeData)
		if err != nil {
			return fmt.Errorf("error marshaling new resume data: %w", err)
		}

		resp, err := http.Post("http://localhost:8080/resumes", "application/json", bytes.NewReader(newResumeJSON))
		if err != nil {
			return fmt.Errorf("error creating resume: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated { // Expect 201 Created
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("error reading response body: %w", err)
			}
			bodyString := string(bodyBytes)
			return fmt.Errorf("error creating resume: status code %d, response body: %s", resp.StatusCode, bodyString)
		}

		fmt.Println("Resume created successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
