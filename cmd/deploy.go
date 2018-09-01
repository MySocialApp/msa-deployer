package cmd

import (
	"github.com/spf13/cobra"
	"strings"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"github.com/spf13/viper"
	"strconv"
)

var s string

type Pipelines struct {
	Jobs string
}

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy <client id>",
	Short: "Deploy client ID applications and infrastructure",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//var pipelineId int
		log.Info("Deploying " + args[0] + " requested")

		// Check client exit and establish connection
		checkClientExist(args[0])
		git := gitlabConnection()
		// Make pipeline + get jobs + run desired job
		pipelineId := gitlabBuildPipeline(git, args[0])
		jobs := gitlabGetJob(git, pipelineId)
		gitlabRunJob(git, pipelineId, jobs, "deploy")
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}

func checkClientExist(clientId string) {
	clientFilename := "clients.csv"

	log.Infof("Check %s client exists", clientId)

	// read the whole file at once
	b, err := ioutil.ReadFile(clientFilename)
	if err != nil {
		panic(err)
	}
	s := string(b)

	//check whether s contains substring text
	if ! strings.Contains(s, clientId) {
		log.Fatalf("This client has not been found in %s", clientFilename)
	}
}

// gitlabConnection establish a gitlab connection
func gitlabConnection() *gitlab.Client {
	return gitlab.NewClient(nil, viper.GetString("gitlab_private_token"))
}

// gitlabBuildPipeline generate a pipeline from what has been configured in .gitlab-ci.yaml.
// All jobs for this pipeline will be generated
// Example: pipeline_id=$(curl -X POST -F "ref=deployer" -F "variables[client_id]=${client_id}" "https://gitlab.com/api/v4/projects/${gitlab_project_id}/trigger/pipeline?token=${gitlab_token}" | jq --raw-input '.id')
func gitlabBuildPipeline(git *gitlab.Client, clientId string) int {
	// Add forms to pipeline trigger
	customForms := make(map[string]string)
	customForms["client_id"] = clientId

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