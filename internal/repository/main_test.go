package repository

import (
	"fmt"
	"log"
	"os"
	"simple-order-go/internal/entity"
	"simple-order-go/pkg/config"
	"strconv"
	"testing"
	"time"

	database "simple-order-go/pkg/db"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/gorm"
)

var (
	testDB        *gorm.DB
	testOrderRepo *OrderRepository
	pool          *dockertest.Pool
	resource      *dockertest.Resource
)

const useDocker = true

func TestMain(m *testing.M) {
	cfg := config.LoadConfig("../../app.yaml")

	var test int
	if useDocker {
		setUpDocketTestEnv(cfg)
		test = m.Run()
		tearDownDockerTestEnv()
	} else {
		setUpDatabase(cfg)
		test = m.Run()
	}

	os.Exit(test)
}

func setUpDocketTestEnv(cfg config.Config) {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatal("Couldn't construct pool: ", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatal("Couldn't connect to Docker: ", err)
	}

	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
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

	rcPort := resource.GetPort("5432/tcp")
	port, err := strconv.Atoi(rcPort)
	if err != nil {
		log.Fatal("Couldn't set port")
	}
	cfg.Database.Port = port

	// Takes a few seconds to start up
	// Increase the delay if it fails to connect
	time.Sleep(3 * time.Second)

	if err := pool.Retry(func() error {
		retries := 10
		err = setUpDatabase(cfg)

		for err != nil {
			if retries > 1 {
				retries--
				time.Sleep(5 * time.Second)
				err = setUpDatabase(cfg)
				if err != nil {
					continue
				} else {
					break
				}
			}
		}

		return err
	}); err != nil {
		log.Fatal("Couldn't connect to DB: ", err)
	}
}

func setUpDatabase(cfg config.Config) error {
	fmt.Println("config: ", cfg)
	var err error
	testDB, err = database.InitDB(cfg.Database)
	if err != nil {
		return err
	}

	if useDocker {
		err = testDB.AutoMigrate(&entity.Order{})
		if err != nil {
			log.Fatal("Couldn't create table orders")
		}

		err = testDB.AutoMigrate(&entity.Item{})
		if err != nil {
			log.Fatal("Couldn't create table items")
		}
	}

	testOrderRepo = NewOrderRepository(testDB)

	return nil
}

func tearDownDockerTestEnv() {
	if err := pool.Purge(resource); err != nil {
		log.Fatal("Could not purge Docker: ", err)
	}
}
