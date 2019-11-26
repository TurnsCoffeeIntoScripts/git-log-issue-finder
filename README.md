# git-log-ticket-finder

This small program is used to extract a list of Jira issues from the output of a git log.

## Content:
1. [Pipeline Configuration](#pipeline_configuration)
2. [Task Configuration](#task_configuration)
3. [Contact](#contact)

## <a name="pipeline_configuration" href="pipeline_configuration">Pipeline Configuration</a>

Here's an example of a Concourse job that uses git-log-ticket-finder

```yml
jobs:
  - name: find-jira /* You can use whatever name you like */
    serial: true
    public: false
    plan: 
      - in_parallel:
        - get: <GIT_REPOSITORY_RESOURCE>
          /* Add 'passed' and/or 'trigger' configuration if needed */
      - task: git-log-ticket-finder
        file: <PATH_TO_YML_TASK_CONFIGURATION>  
```

## <a name="task_configuration" href="task_configuration">Task Configuration</a>

To properly configure git-log-ticket-finder, it should be done as a task and not directly as a pipeline resource. 

Here's what the task's yaml file should look like

```yml
plateform: linux
image_resource:
  type: docker-image
  source:
    repository: turnscoffeeintoscripts/git-log-ticket-finder
    tag: latest

params:
  DESTINATION: 'git-repo'
  TICKETS_FILTER: 'ABC,DEFG,XYZ'
  FILENAME: gltf-output-filename.txt
    
  inputs:
    - name: git-repo
    
    
  outputs:
    <OUTPUT_IF_ANY_ARE_NEEDED>
    
  run:
    path: /bin/sh
    args:
        - <PATH_TO_SH_SCRIPT> 
```
The yaml configuration should contain three parameters ('params'). The destination (DESTINATION) parameter should contain
the name of the actual git repository folder. The tickets filter (TICKETS_FILTER) is a comma-separated list of Jira
project keys. Finally the filename (FILENAME) is the name of the file in which the result of the command will be written.
This last feature is useful if the result is needed as input of another job.

And now here's what the shell script should look like:

```bash
#!/bin/bash

# Force exit when a pipeline/single-command returns with non-zero status
set -e

destination=${DESTINATION}
tickets=${TICKETS_FILTER}
resultFile=${FILENAME}

if [[ -d ${destination} ]]; then
    cd ${destination}

    log=$(git log --pretty=oneline)

    # Launch the Go exec 'gitLogTicketFinder'
    gltfResult=$(gitLogTicketFinder --tickets ${tickets} --content "$log")

    if [[ -f "$resultFile" ]]; then
        rm -f ${resultFile}
    fi

    echo ${gltfResult} >> ${resultFile}
    echo ${gltfResult}

    cd ../
else
    echo "Git repo $destination does not exist or could not be found"
    exit 1
fi
```

## <a name="contact" href="contact">Contact</a>
If you have any questions/comments please send them at: turns.coffee.into.scripts@gmail.com.

You may also submit pull-requests on github at: https://github.com/TurnsCoffeeIntoScripts/git-log-ticket-finder 
to the branch 'master'.