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
