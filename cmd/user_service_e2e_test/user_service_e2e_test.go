package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/ani5msr/microservices-project/pkg/db_utils"
	om "github.com/ani5msr/microservices-project/pkg/object_model"
	"github.com/ani5msr/microservices-project/pkg/user_client"
	_ "github.com/lib/pq"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func initDB() {
	db, err := db_utils.RunLocalDB("user_manager")
	if err != nil {
		return
	}

	tables := []string{"sessions", "users"}
	for _, table := range tables {
		err = db_utils.DeleteFromTableIfExist(db, table)
		check(err)
	}
}

func runServer(ctx context.Context) {
	// Build the server if needed
	_, err := os.Stat("./user_service")
	if os.IsNotExist(err) {
		out, err := exec.Command("go", "build", ".").CombinedOutput()
		log.Println(out)
		check(err)
	}

	cmd := exec.CommandContext(ctx, "./user_service")
	err = cmd.Start()
	check(err)
}

func killServer(ctx context.Context) {
	ctx.Done()
}

func main() {
	initDB()

	ctx := context.Background()
	defer killServer(ctx)
	runServer(ctx)

	// Run some tests with the client
	cli, err := user_client.NewClient("localhost:7070")
	check(err)

	err = cli.Register(om.User{"gg@gg.com", "gi"})
	check(err)
	log.Print("gi has registered successfully")

	session, err := cli.Login("gigi", "secret")
	check(err)
	log.Print("gi has logged in successfully. the session is: ", session)

	err = cli.Logout("gigi", session)
	check(err)
	log.Print("gigi has logged out successfully.")

}
