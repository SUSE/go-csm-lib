package provisioner_test

import (
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hpcloud/sidecar-extensions/go/extension/mysql/config"
	"github.com/hpcloud/sidecar-extensions/go/extension/mysql/provisioner"
	"github.com/pivotal-golang/lager/lagertest"
	"gopkg.in/caarlos0/env.v2"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("mysql-provisioner-test")

func getProvisioner() provisioner.MySQLProvisioner {

	mysqlConfig := config.MySQLConfig{}

	logger = lagertest.NewTestLogger("process-controller")
	err := env.Parse(&mysqlConfig)
	if err != nil {
		logger.Error("failed to loat environment variables", err)
	}

	if mysqlConfig.User == "" || mysqlConfig.Pass == "" || mysqlConfig.Host == "" || mysqlConfig.Port == "" {
		return nil
	}
	mysqlProvisioner := provisioner.NewGoSQL(logger, mysqlConfig)

	return mysqlProvisioner
}

func TestCreateDb(t *testing.T) {
	mysqlProvisioner := getProvisioner()
	if mysqlProvisioner == nil {
		t.Skip("Skipping test as not all env variables are set:'MYSQL_USER','MYSQL_PASS','MYSQL_HOST', 'MYSQL_PORT'")
	}

	dbName := "test_createdb"

	log.Println("Creating test database")
	err := mysqlProvisioner.CreateDatabase(dbName)

	if err != nil {
		log.Fatalln("Error creating database ", err)
	}
}

func TestCreateDbExists(t *testing.T) {
	mysqlProvisioner := getProvisioner()
	if mysqlProvisioner == nil {
		t.Skip("Skipping test as not all env variables are set:'MYSQL_USER','MYSQL_PASS','MYSQL_HOST', 'MYSQL_PORT'")
	}

	dbName := "test_createdb"

	log.Println("Testing if database exists")
	created, err := mysqlProvisioner.IsDatabaseCreated(dbName)
	if err != nil {
		log.Fatal(err)
	}
	if created {
		t.Log("Created true")
	} else {
		t.Log("Created false")
	}
}

func TestCreateUser(t *testing.T) {
	mysqlProvisioner := getProvisioner()
	if mysqlProvisioner == nil {
		t.Skip("Skipping test as not all env variables are set:'MYSQL_USER','MYSQL_PASS','MYSQL_HOST', 'MYSQL_PORT'")
	}

	dbName := "test_createdb"

	log.Println("Creating test user")
	err := mysqlProvisioner.CreateUser(dbName, "mytestUser", "mytestPass")
	if err != nil {
		t.Errorf("Error creating user %v", err)
	}
}

func TestCreateUserExists(t *testing.T) {
	mysqlProvisioner := getProvisioner()
	if mysqlProvisioner == nil {
		t.Skip("Skipping test as not all env variables are set:'MYSQL_USER','MYSQL_PASS','MYSQL_HOST', 'MYSQL_PORT'")
	}

	log.Println("Testing if user exists")
	created, err := mysqlProvisioner.IsUserCreated("mytestUser")
	if err != nil {
		t.Errorf("Error verifying user %v", err)
	}
	if created {
		t.Log("test user is created")
	} else {
		t.Log("test user was not created")
	}
}

func TestDeleteUser(t *testing.T) {
	mysqlProvisioner := getProvisioner()
	if mysqlProvisioner == nil {
		t.Skip("Skipping test as not all env variables are set:'MYSQL_USER','MYSQL_PASS','MYSQL_HOST', 'MYSQL_PORT'")
	}

	log.Println("Removing test user")
	err := mysqlProvisioner.DeleteUser("mytestUser")
	if err != nil {
		t.Errorf("Error deleting user %v", err)
	}
}

func TestDeleteTheDatabase(t *testing.T) {
	mysqlProvisioner := getProvisioner()
	if mysqlProvisioner == nil {
		t.Skip("Skipping test as not all env variables are set:'MYSQL_USER','MYSQL_PASS','MYSQL_HOST', 'MYSQL_PORT'")
	}

	dbName := "test_createdb"
	log.Println("Removing test database")

	err := mysqlProvisioner.DeleteDatabase(dbName)
	if err != nil {
		t.Errorf("Error deleting database %v", err)
	}
}
