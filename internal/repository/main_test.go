package repository

import (
	"fmt"
	"log"
	"os"
	"simple-order-go/pkg/config"
	"strconv"
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

	setUpDatabase(cfg)

	// setUpDocketTestEnv(cfg)

	test := m.Run()

	// tearDownDockerTestEnv()

	os.Exit(test)
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
			"POSTGRES_HOST_AUTH_METHOD=trust",
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

	rcPort := resource.GetPort("5432/tcp")
	port, err := strconv.Atoi(rcPort)
	if err != nil {
		log.Fatal("Couldn't set port")

	}
	cfg.Database.Port = port

	if err := pool.Retry(func() error {
		return setUpDatabase(cfg)
	}); err != nil {
		log.Fatal("Couldn't connect to DB: ", err)
	}
}

func setUpDatabase(cfg config.Config) error {

	testDB, err := database.InitDB(cfg.Database)
	if err != nil {
		return err
	}

	// testDB.AutoMigrate(&entity.Order{})
	// testDB.AutoMigrate(&entity.Item{})

	testOrderRepo = NewOrderRepository(testDB)

	return nil
}

func tearDownDockerTestEnv() {
	if err := pool.Purge(resource); err != nil {
		log.Fatal("Could not purge Docker: ", err)
	}
}
