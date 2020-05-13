# Git Log Issue Finder (GLIF)

This small program is used to extract a list of Jira issues from the output of a git log. This require an already checked-out
repository. 

Version: 2.0.0

## Content:
1. [Usage](#usage)
2. [Pipeline Configuration](#pipeline_configuration)
3. [Task Configuration](#task_configuration)
4. [Contact](#contact)

## <a name="usage" href="usage">Usage</a>
While this program was designed with the idea of integrating it to a *<a href="https://concourse-ci.org/" target="_blank">Concourse</a>* pipeline it can also be used a stand-alone
command-line tool.

Let's start by looking at the '--help' command:
```bash
$> glif --help
  -force-fetch
        Force a 'git fetch' operation on the specified repository
  -repl
        Enter the Read-Eval-Print-Loop
  -script string
        The glif script file to execute
  -semver-latest
        script.DiffLatestSemver
  -semver-latest-builds
        script.DiffLatestSemverWithLatestBuilds
  -semver-latest-rcs
        script.DiffLatestSemverWithLatestRCs
  -tickets string
        The Jira tickets regex used to search the repo's log (default "*")

$> 
```
While the help command does not specify it, it's useful to note that every parameter can specified with one or two hyphen.

### The 'script' parameter
This parameter is simple in itself as it's only a path to the script to be interpreted by glif. You can read the documentation
of the glif scripts here [glif doc](glif_doc/README.md) and you can see examples here
[examples](examples/README.md). In the examples you will also find more details about `semver-latest`, `semver-latest-builds` and
`semver-latest-rcs`.

### The 'tickets' parameter
The tickets parameter can be specified in the command line or it can be specified within the glif
script. It represents the name of the Jira issues to match.

The following example assumes that you are working within the directory of your git repository.
The most basic usage requires two (2) parameters; here's an example with the output:
```bash
$> glif --tickets="ABC,XYZ" --semver-latest
[ABC-001, ABC-007, XYZ-9246, ABC-045, ABC-0245, XYZ-007]
$> 
```

## <a name="pipeline_configuration" href="pipeline_configuration">Pipeline Configuration</a>

Now, here's an example of a Concourse job that uses git-log-issue-finder

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

```yaml
#TODO
```
...

And now here's what the shell script should look like:

```bash
#!/bin/bash

set -e

# TODO
```

## <a name="contact" href="contact">Contact</a>
If you have any questions/comments please send them at: turns.coffee.into.scripts@gmail.com.

You may also submit pull-requests on github at: https://github.com/TurnsCoffeeIntoScripts/git-log-issue-finder 
to the branch 'master'.