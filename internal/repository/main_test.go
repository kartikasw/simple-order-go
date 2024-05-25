package repository

import (
	"fmt"
	"log"
	"os"
	"simple-order-go/pkg/config"
	"testing"

	database "simple-order-go/pkg/db"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	testOrderRepo *OrderRepository
	pool          *dockertest.Pool
	resource      *dockertest.Resource
)

func TestMain(m *testing.M) {
	cfg := config.LoadConfig("../../app.yaml")

	setUpDocketTestEnv(cfg)
	connectToDB()

	test := m.Run()

	tearDownDockerTestEnv()

	os.Exit(test)
}

func connectToDB() {
	cfg := config.LoadConfig("../../app.yaml")
	testDB, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatal("Couldn't connect to DB: ", err)
	}

	testOrderRepo = NewOrderRepository(testDB)
}

func setUpDocketTestEnv(cfg config.Config) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal("Couldn't construct pool: ", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatal("Couldn't connect to Docker: ", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", cfg.Database.User),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", cfg.Database.Password),
			fmt.Sprintf("POSTGRES_DB=%s", cfg.Database.Name),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	if err != nil {
		log.Fatal("Couldn't start resource: ", err)
	}

	resource.Expire(120)

	resource.GetPort(fmt.Sprintf("tcp://%s:%d", cfg.Database.Host, cfg.Database.Port))

	if err := pool.Retry(func() error {
		var err error
		_, err = database.InitDB(cfg.Database)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		tearDownDockerTestEnv()
		log.Fatal("Couldn't connect to DB: ", err)
	}
}

func tearDownDockerTestEnv() {
	if err := pool.Purge(resource); err != nil {
		log.Fatal("Could not purge Docker: ", err)
	}
}
