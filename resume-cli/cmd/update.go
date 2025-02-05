package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a resume",
	Long:  `Update an existing resume.`,
	Args:  cobra.ExactArgs(1), //require just the ID
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		//1. get current resume data
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/resumes/%s", id))
		if err != nil {
			return fmt.Errorf("error getting resume: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("error getting resume: status code %d", resp.StatusCode)
		}

		//begin update process
		var resumeData map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&resumeData)
		if err != nil {
			return fmt.Errorf("error decoding resume data: %w", err)
		}

		//2. prompt the user for updates
		var updates map[string]interface{} = make(map[string]interface{})

		fmt.Print("Enter Basic Info updates (JSON format, or press Enter to skip):\n")
		var basicInfoUpdates map[string]interface{} = make(map[string]interface{})
		if _, err := fmt.Scanln(); err == nil {
			decoder := json.NewDecoder(os.Stdin)
			if err := decoder.Decode(&basicInfoUpdates); err != nil {
				return fmt.Errorf("error decoding basic info updates: %w", err)
			}
		}
		if len(basicInfoUpdates) > 0 {
			updates["basic_info"] = basicInfoUpdates
		}

		//education updates
		fmt.Print("Enter Education updates (JSON format, or press Enter to skip):\n")
		var educationUpdates []interface{} = make([]interface{}, 0)
		if _, err := fmt.Scanln(); err == nil {
			decoder := json.NewDecoder(os.Stdin)
			if err := decoder.Decode(&educationUpdates); err != nil {
				return fmt.Errorf("error decoding basic education updates: %w", err)
			}
		}

		if len(educationUpdates) > 0 {
			updates["education"] = educationUpdates
		}

		//skills updates
		fmt.Print("Enter Skills updates (JSON format, or press Enter to skip):\n")
		var skillsUpdates []interface{} = make([]interface{}, 0)
		if _, err := fmt.Scanln(); err == nil {
			decoder := json.NewDecoder(os.Stdin)
			if err := decoder.Decode(&skillsUpdates); err != nil {
				return fmt.Errorf("error decoding skills updates: %w", err)
			}
		}

		if len(skillsUpdates) > 0 {
			updates["skills"] = skillsUpdates
		}

		//experience updates
		fmt.Print("Enter Experience updates (JSON format, or press Enter to skip):\n")
		var experienceUpdates []interface{} = make([]interface{}, 0)
		if _, err := fmt.Scanln(); err == nil {
			decoder := json.NewDecoder(os.Stdin)
			if err := decoder.Decode(&experienceUpdates); err != nil {
				return fmt.Errorf("error decoding experience updates: %w", err)
			}
		}

		if len(experienceUpdates) > 0 {
			updates["experience"] = experienceUpdates
		}

		//project updates
		fmt.Print("Enter Project updates (JSON format, or press Enter to skip):\n")
		var projectUpdates []interface{} = make([]interface{}, 0)
		if _, err := fmt.Scanln(); err == nil {
			decoder := json.NewDecoder(os.Stdin)
			if err := decoder.Decode(&projectUpdates); err != nil {
				return fmt.Errorf("error decoding project updates: %w", err)
			}
		}

		if len(projectUpdates) > 0 {
			updates["project"] = projectUpdates
		}

		//affliation updates
		fmt.Print("Enter Affiliation updates (JSON format, or press Enter to skip):\n")
		var affiliationUpdates []interface{} = make([]interface{}, 0)
		if _, err := fmt.Scanln(); err == nil {
			decoder := json.NewDecoder(os.Stdin)
			if err := decoder.Decode(&affiliationUpdates); err != nil {
				return fmt.Errorf("error decoding affiliation updates: %w", err)
			}
		}

		if len(affiliationUpdates) > 0 {
			updates["affiliation"] = affiliationUpdates
		}

		//3. Make PUT request
		updatesJSON, err := json.Marshal(updates)
		if err != nil {
			return fmt.Errorf("error marshaling updates: %w", err)
		}

		req, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:8080/resumes/%s", id), strings.NewReader(string(updatesJSON)))
		if err != nil {
			return fmt.Errorf("error creating requests %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("error updating resume: status code %d", resp.StatusCode)
		}

		fmt.Println("Resume update successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
