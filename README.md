I was just tired of loading the webpage.

```sh
go get github.com/brettbuddin/jenkins-status

export JENKINS_BASE_URL=http://<JENKINS_HOSTNAME>
export JENKINS_USER=<USERNAME>
export JENKINS_TOKEN=<TOKEN>

jenkins-status <JOB_NAME>
...
```
