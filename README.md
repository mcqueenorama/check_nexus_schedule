Nagios check for the Jenkins Queue in Go.

Usage of ./jenkinsQueue:
  -c=20: critical level for job queue depth
  -h="http://ci.walmartlabs.com/jenkins": base url for jenkins  like http://ci.walmartlabs.com/jenkins
  -v=false: verbose output
  -w=10: warning level for job queue depth

Build it:

  go build

or:

  go build check_jenkins_queue.go

Build it in docker for another platform:

docker run -it -v /Users:/Users -w `pwd` google/golang go build check_jenkins_queue.go

Mess with the API:

curl -v  http://ci.walmartlabs.com/jenkins/queue/api/json|json_pp



