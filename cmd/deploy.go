package cmd

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
	"os"
	"strconv"
	"strings"
)

var s string

type Pipelines struct {
	Jobs string
}

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy <client id> [app name]",
	Short: "Deploy client ID applications and application (optional)",
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		clientFileName := "clients.csv"

		// var pipelineId int
		log.Infof("Deploying %s requested", args[0])

		// Check client/app exist and establish connection
		checkClientAndAppExist(clientFileName, args)
		git := gitlabConnection()

		// Make pipeline + get jobs + run desired job
		pipelineId := gitlabBuildPipeline(git, args)
		jobs := gitlabGetJob(git, pipelineId)
		gitlabRunJob(git, pipelineId, jobs, "deploy")
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}

func checkClientAndAppExist(clientFileName string, args []string) {
	clientFound := 0
	appFound := 0
	clientId := args[0]
	app := ""

	if len(args) == 2 {
		app = args[1]
	}

	inFile, _ := os.Open(clientFileName)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		// select clientId line
		if strings.Contains(scanner.Text(), clientId) {
			clientFound = 1
			// select app line
			if strings.Contains(scanner.Text(), app) {
				appFound = 1
				log.Debugf("App %s found for client id %s: %s", app, clientId, scanner.Text())
			}
		}
	}

	if clientFound == 0 {
		log.Fatalf("Client %s has not been found in %s", clientId, clientFileName)
	}
	if app != "" && appFound == 0 {
		log.Fatalf("Application %s is not set for the client %s in %s", app, clientId, clientFileName)
	}
}

// gitlabConnection establish a gitlab connection
func gitlabConnection() *gitlab.Client {
	return gitlab.NewClient(nil, viper.GetString("gitlab_private_token"))
}

// gitlabBuildPipeline generate a pipeline from what has been configured in .gitlab-ci.yaml.
// All jobs for this pipeline will be generated
// Example: pipeline_id=$(curl -X POST -F "ref=deployer" -F "variables[client_id]=${client_id}" "https://gitlab.com/api/v4/projects/${gitlab_project_id}/trigger/pipeline?token=${gitlab_token}" | jq --raw-input '.id')
func gitlabBuildPipeline(git *gitlab.Client, args []string) int {
	// Add forms to pipeline trigger
	customForms := make(map[string]string)
	customForms["client_id"] = args[0]
	if len(args) >= 2 {
		customForms["app_name"] = args[1]
	}

	// Generate pipeline trigger
	opt := &gitlab.RunPipelineTriggerOptions{
		Token:		gitlab.String(viper.GetString("gitlab_pipeline_token")),
		Variables:	customForms,
		Ref: 		gitlab.String("master"), // todo: remove when going to prod
	}

	// Build pipeline
	project, _, err := git.PipelineTriggers.RunPipelineTrigger(
		viper.GetInt("gitlab_project_id"),
		opt)
	if err != nil {
		log.Fatalf("Wasn't able to create the gitlab pipeline: %s", err)
	}

	return project.ID
}

// gitlabGetJobId get jobs from a pipeline ID
// Example: job_id=$(curl --header "PRIVATE-TOKEN: ${gitlab_token}" "https://gitlab.com/api/v4/projects/${gitlab_project_id}/pipelines/${pipeline_id}/jobs" | jq --raw-input ".[] | select(.name == 'add-client') | .id")
func gitlabGetJob(git *gitlab.Client, pipelineId int) []*gitlab.Job {
	jobs, _, err := git.Jobs.ListPipelineJobs(
		viper.GetInt("gitlab_project_id"),
		pipelineId, &gitlab.ListJobsOptions{})
	if err != nil {
		log.Fatalf("Wasn't able to list jobs from gitlab pipeline: %s", err)
	}
	return jobs
}

// gitlabRunJob plays a job from a job name
// Example: curl -X POST --header "PRIVATE-TOKEN: ${gitlab_token}" -F ref=deployer "https://gitlab.com/api/v4/projects/${gitlab_project_id}/jobs/${job_id}/play"
func gitlabRunJob(git *gitlab.Client, pipelineId int, jobs []*gitlab.Job, jobName string) {
	var jobId int

	// Get pipeline ID and job ID
	for job := range jobs {
		if jobs[job].Name == jobName {
			jobId = jobs[job].ID
		}
	}

	// Play job
	_, _, err := git.Jobs.PlayJob(
		viper.GetInt("gitlab_project_id"),
		jobId,
		nil,
	)
	if err != nil {
		log.Fatalf("Wasn't able to play job %s id %s on pipeline %s: %s", jobName, strconv.Itoa(jobId), strconv.Itoa(pipelineId), err)
	}
	log.Info("Job successfully been launched")
	log.Infof("Job progression: https://gitlab.com/%s/-/jobs/%s", viper.GetString("gitlab_project_name"), strconv.Itoa(jobId))
}