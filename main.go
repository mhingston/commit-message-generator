package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/mhingston/azoai"
)

func getFormattedBranchName() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	branchName := strings.TrimSpace(string(out))

	re := regexp.MustCompile(`^feature/`)
	if re.MatchString(branchName) {
		return re.ReplaceAllString(branchName, ""), nil
	}
	return "", nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: commit-message-generator <messagePath>")
		os.Exit(1)
	}

	messagePath := os.Args[1]

	originalMessage, err := os.ReadFile(messagePath)
	if err != nil {
		log.Fatalf("Error reading message file: %v", err)
	}
	fmt.Printf("Original Message: %s\n", string(originalMessage))

	apiKey := os.Getenv("AZURE_OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set the AZURE_OPENAI_API_KEY environment variable.")
		os.Exit(0)
	}

	apibaseUrl := os.Getenv("AZURE_OPENAI_API_ENDPOINT")
	if apibaseUrl == "" {
		fmt.Println("Please set the AZURE_OPENAI_API_ENDPOINT environment variable.")
		os.Exit(0)
	}

	apiVersion := os.Getenv("AZURE_OPENAI_API_VERSION")
	if apiVersion == "" {
		fmt.Println("Please set the AZURE_OPENAI_API_VERSION environment variable.")
		os.Exit(0)
	}

	deployment := os.Getenv("AZURE_OPENAI_API_DEPLOYMENT")
	if deployment == "" {
		fmt.Println("Please set the AZURE_OPENAI_API_DEPLOYMENT environment variable.")
		os.Exit(0)
	}

	systemPrompt := "Create a concise, meaningful git commit message (max 40 characters) based on the provided git diff."

	cmd := exec.Command("git", "diff", "--cached")
	diff, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error getting git diff: %v", err)
	}

	if len(diff) == 0 {
		fmt.Println("No changes to commit.")
		os.Exit(0)
	}

	var newMessage string
	if strings.TrimSpace(string(originalMessage)) == "" {
		resp, err := azoai.InvokeOpenAIRequest(azoai.OpenAIRequest{
			SystemPrompt: systemPrompt,
			Message:      string(diff),
			ApiBaseUrl:   apibaseUrl,
			APIKey:       apiKey,
			APIVersion:   apiVersion,
			Deployment:   deployment,
		})
		if err != nil {
			log.Fatalf("Error invoking OpenAI request: %v", err)
		}

		if resp != "" {
			newMessage = resp
		} else {
			newMessage = string(originalMessage)
		}
	} else {
		newMessage = string(originalMessage)
	}

	branchName, err := getFormattedBranchName()
	if err != nil {
		log.Fatalf("Error getting branch name: %v", err)
	}

	if branchName != "" {
		newMessage = fmt.Sprintf("%s: %s", branchName, newMessage)
	}

	err = os.WriteFile(messagePath, []byte(newMessage), 0644)
	if err != nil {
		log.Fatalf("Error writing new message to file: %v", err)
	}

	fmt.Printf("Updated commit message: %s\n", newMessage)
}
