package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/ani5msr/microservices-project/pkg/db_utils"
	om "github.com/ani5msr/microservices-project/pkg/object_model"
	"github.com/ani5msr/microservices-project/pkg/post_manager_client"
	_ "github.com/lib/pq"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func initDB() {
	db, err := db_utils.RunLocalDB("post_manager")
	check(err)

	tables := []string{"tags", "posts"}
	for _, table := range tables {
		err = db_utils.DeleteFromTableIfExist(db, table)
		check(err)
	}
}

// Build and run a service in a target directory
func runService(ctx context.Context, targetDir string, service string) {
	// Save and restore later current working dir
	wd, err := os.Getwd()
	check(err)
	defer os.Chdir(wd)

	// Build the server if needed
	_, err = os.Stat("./" + service)
	if os.IsNotExist(err) {
		out, err := exec.Command("go", "build", ".").CombinedOutput()
		log.Println(out)
		check(err)
	}

	cmd := exec.CommandContext(ctx, "./"+service)
	err = cmd.Start()
	check(err)
}

func runPostService(ctx context.Context) {
	// Set environment
	err := os.Setenv("MAX_LINKS_PER_USER", "10")
	check(err)

	runService(ctx, ".", "post_service")
}

func runSocialGraphService(ctx context.Context) {
	runService(ctx, "../social_graph_service", "post_service")
}

func killServer(ctx context.Context) {
	ctx.Done()
}

func main() {
	initDB()

	ctx := context.Background()
	defer killServer(ctx)
	runSocialGraphService(ctx)
	runPostService(ctx)

	// Run some tests with the client
	cli, err := post_manager_client.NewClient("localhost:8080")
	check(err)

	posts, err := cli.GetPost(om.GetPostRequest{Username: "ani5msr"})
	check(err)
	log.Print("posts:", posts)

	err = cli.AddPost(om.AddPostRequest{Username: "ani5msr",
		Url:   "https://github.com/ani5msr",
		Title: "Github",
		Tags:  map[string]bool{"programming": true}})
	check(err)
	posts, err = cli.GetPost(om.GetPostRequest{Username: "ani5msr"})
	check(err)
	log.Print("posts:", posts)

	err = cli.UpdatePost(om.UpdatePostRequest{Username: "ani5msr",
		Url:         "https://github.com/ani5msr",
		Description: "Most of my open source code is here"},
	)

	check(err)
	posts, err = cli.GetPost(om.GetPostRequest{Username: "ani5msr"})
	check(err)
	log.Print("posts:", posts)
}
