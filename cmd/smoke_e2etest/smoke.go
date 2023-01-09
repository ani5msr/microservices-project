package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	net_url "net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	om "github.com/ani5msr/microservices-project/pkg/object_model"
	. "github.com/ani5msr/microservices-project/pkg/test_util"
)

var (
	microUrl   string
	microToken = os.Getenv("DELINKCIOUS_TOKEN")
	httpClient = http.Client{}
)

func getPosts() {
	req, err := http.NewRequest("GET", string(microUrl)+"/posts", nil)
	Check(err)

	req.Header.Add("Access-Token", microToken)
	r, err := httpClient.Do(req)
	Check(err)

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		Check(errors.New(r.Status))
	}

	var glr om.GetPostResult
	body, err := ioutil.ReadAll(r.Body)

	err = json.Unmarshal(body, &glr)
	Check(err)

	log.Println("======= Posts =======")
	for _, post := range glr.Posts {
		log.Println(fmt.Sprintf("title: '%s', url: '%s', status: '%s', description: '%s'", post.Title,
			post.Url,
			post.Status,
			post.Description))
	}
}

func getFollowing() {
	req, err := http.NewRequest("GET", string(microUrl)+"/following", nil)
	Check(err)

	req.Header.Add("Access-Token", microToken)
	r, err := httpClient.Do(req)
	Check(err)

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		Check(errors.New(r.Status))
	}

	body, err := ioutil.ReadAll(r.Body)
	Check(err)

	log.Println("======= Following =======")
	log.Println(string(body))
}

func addPost(url string, title string) {
	params := net_url.Values{}
	params.Add("url", url)
	params.Add("title", title)
	qs := params.Encode()

	log.Println("===== Add Post ======")
	log.Println(fmt.Sprintf("Adding new post - title: '%s', url: '%s'", title, url))

	url = fmt.Sprintf("%s/posts?%s", microUrl, qs)
	req, err := http.NewRequest("POST", url, nil)
	Check(err)

	req.Header.Add("Access-Token", microToken)
	r, err := httpClient.Do(req)
	Check(err)
	if r.StatusCode != http.StatusOK {
		defer r.Body.Close()
		bodyBytes, err := ioutil.ReadAll(r.Body)
		Check(err)
		message := r.Status + " " + string(bodyBytes)
		Check(errors.New(message))
	}
}

func deletePost(url string) {
	params := net_url.Values{}
	params.Add("url", url)
	qs := params.Encode()

	url = fmt.Sprintf("%s/posts?%s", microUrl, qs)
	req, err := http.NewRequest("DELETE", url, nil)
	Check(err)

	req.Header.Add("Access-Token", microToken)
	r, err := httpClient.Do(req)
	Check(err)
	if r.StatusCode != http.StatusOK {
		defer r.Body.Close()
		bodyBytes, err := ioutil.ReadAll(r.Body)
		Check(err)
		message := r.Status + " " + string(bodyBytes)
		Check(errors.New(message))
	}
}

func main() {
	result, err := exec.Command("kubectl", "config", "current-context").CombinedOutput()
	Check(err)
	currContext := string(result[:len(result)-1])
	fmt.Println("Checking which platform the cluster is running on... kubectl context:", currContext)

	var tempUrl []byte
	if strings.HasPrefix(currContext, "minikube") {
		tempUrl, err = exec.Command("minikube", "service", "api-gateway", "--url").CombinedOutput()
		Check(err)
		fmt.Println("Running on minikube")
	} else if strings.HasPrefix(currContext, "gke") {
		filter := "jsonpath='{.status.loadBalancer.ingress[0].ip}'"
		tempUrl, err = exec.Command("kubectl", "get", "svc", "api-gateway",
			"-o", filter).CombinedOutput()
		Check(err)
		fmt.Println("Running on GKE")

	} else if strings.HasSuffix(currContext, "eksctl.io") {
		filter := "jsonpath='{.status.loadBalancer.ingress[0].hostname}'"
		tempUrl, err = exec.Command("kubectl", "get", "svc", "api-gateway", "-o", filter).CombinedOutput()
		Check(err)
		fmt.Println("Running on AWS")

	}

	//if err != nil {
	//	fmt.Println("Guessing running on KIND")
	//	go func() {
	//		exec.Command("kubectl", "port-forward", "svc/api-gateway", "5000:80")
	//	}()
	//	time.Sleep(time.Second * 3)
	//	tempUrl = []byte("http://localhost:5000/")
	//}

	microUrl = string(tempUrl[:len(tempUrl)-1]) + "/v1"
	if !strings.HasPrefix(microUrl, "http") {
		microUrl = "http://" + microUrl[1:]
	}

	fmt.Printf("url: '%s'\n", microUrl)

	// Get following
	getFollowing()

	// Delete post
	deletePost("https://github.com/the-ani5")

	// Get posts
	getPosts()

	// Add a new post
	addPost("https://github.com/the-ani5", "fb")

	// Get posts again
	getPosts()

	// Wait a little and get posts again
	time.Sleep(time.Second * 3)
	getPosts()
}
