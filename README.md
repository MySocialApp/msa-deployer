# MSA Deployer [![Build Status](https://travis-ci.org/MySocialApp/msa-deployer.svg?branch=master)](https://travis-ci.org/MySocialApp/msa-deployer)

This tool is the one we're using at [MySocialApp](https://mysocialapp.io) (iOS and Android social app builder - SaaS) to simplify deployment.
It launch gitlab pipelines jobs and give a wizzard for some specific actions.

```
  +---------------+   +------------+  +---------------+
  | MSA Deployer  |-->| Gitlab API |  | Launch action |
  +------+--------+   +-----+------+  +---------------+
         |                  |                 ^
         v                  v                 |
    +---------+       +-----------+    +------+-------+
    | Wizzard |       | Gitlab CI |    | Docker Image |
    +----+----+       +----+------+    +--------------+
         |                  |                 ^
         v                  v                 |
  +--------------+     +---------+   +--------+-------+
  | Local Action |     | CD Jobs |-->| Clone MSA repo |
  +--------------+     +---------+   +----------------+
```

Some explaination:
* MSA Deployer: When this tool is run, you can perform local actions or remote (through GitLab CI/CD)
* Local Action: Some actions do not request GitLab API to be performed
* GitLab API: Config file need to be fulfill to authorize access
* Gitlab CI: Use .gitlab-ci.yml in your repository. Jobs will be used to make a pipeline
* CD Jobs: Jobs will be created to be played once dependencies will be satisfied
* Clone MSA repo: Git repository with submodules will be pulled
* Docker image: a dedicated image with dependencies has to be made and will be pulled
* Launch action: actions inside the docker image will be performed

## Configuration

To be able to use the deployer, you need to make this config file and fill with your own config:

```yaml
gitlab_project_id: ""
gitlab_pipeline_token: ""
gitlab_private_token: ""
```

## Usage

Simply run this to get all available options:

```
./msa-deployer --help
```

To deploy, look at:
```
./msa-deployer deploy --help
```

If you want to deploy an app for all clients:
```
./msa-deployer deploy all <your_app_name>
```
