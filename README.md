# Blog application using REST services and Kubernetes


## Docker images can be build by going to svc folder and doing docker build microservices/api0.1 for social graph service, post service or user service

# Social Graph Service

## Social graph service comprises of following endpoints
### /follower/{username}
fetches the followers for the provided username
### /following/{username}
fetches the following of provided username
### /follow
Post method to follow a user
### /unfollow
Post method to unfollow a user

# Feed Service

### Feed service uses NATS pub sub notification for notifying whenever a user has added or updated his post
## NATS
NATS using publisher subscriber method to notify about event changes

# Post Service

### Post service uses in memory database and a docker postgres image to save post data

## Middlewares
  Middlewares are used for tracking and metrics information. Metrics is used from go-kit and opentracing api is used for tracing the api calls
