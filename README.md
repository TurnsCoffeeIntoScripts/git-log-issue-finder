# Git Log Issue Finder (GLIF)

This small program is used to extract a list of Jira issues from the output of a git log. This require an already checked-out
repository. 

## Content:
1. [Pipeline Configuration](#pipeline_configuration)
2. [Task Configuration](#task_configuration)
3. [Contact](#contact)

## <a name="pipeline_configuration" href="pipeline_configuration">Pipeline Configuration</a>

Here's an example of a Concourse job that uses git-log-issue-finder

```yml
jobs:
  - name: find-jira /* You can use whatever name you like */
    serial: true
    public: false
    plan: 
      - in_parallel:
        - get: <GIT_REPOSITORY_RESOURCE>
          /* Add 'passed' and/or 'trigger' configuration if needed */
      - task: git-log-issue-finder
        file: <PATH_TO_YML_TASK_CONFIGURATION>  
```

## <a name="task_configuration" href="task_configuration">Task Configuration</a>

To properly configure git-log-issue-finder, it should be done as a task and not directly as a pipeline resource. 

Here's what the task's yaml file should look like

```yml
/* TODO */
```
The yaml configuration should contain three parameters ('params'). The destination (DESTINATION) parameter should contain
the name of the actual git repository folder. The tickets filter (TICKETS_FILTER) is a comma-separated list of Jira
project keys. Finally the filename (FILENAME) is the name of the file in which the result of the command will be written.
This last feature is useful if the result is needed as input of another job.

And now here's what the shell script should look like:

```bash
#!/bin/bash

# TODO
```

## <a name="contact" href="contact">Contact</a>
If you have any questions/comments please send them at: turns.coffee.into.scripts@gmail.com.

You may also submit pull-requests on github at: https://github.com/TurnsCoffeeIntoScripts/git-log-issue-finder 
to the branch 'master'.