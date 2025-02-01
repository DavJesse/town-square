package database // Global setup before tests
import (
	"fmt"
	"os"
	"testing"
)

const dbFile = "test.db"

// Setup function to delete the database file before each test
func deleteTestDb() {
	// Delete the database file if it exists
	if _, err := os.Stat(dbFile); err == nil {
		err := os.Remove(dbFile)
		if err != nil {
			fmt.Printf("Failed to delete the database file: %v\n", err)
		}
	}
}

// a special function used to setup testing
func TestMain(m *testing.M) {
	// Setup before tests
	fmt.Println("Global setup: init db handle")

	// delete the database file used in tests if available
	deleteTestDb()

	// inititlize all tables
	Init(dbFile)

	// with db opened, ensure it will be closed once tests are done
	defer db.Close()

	// create a single user that will be used in tests
	_, err = CreateUser("milton", "milton@mail.com", "mPass")
	if err != nil {
		error_s := fmt.Errorf("%w\n", err)
		fmt.Println(error_s)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Teardown after tests
	fmt.Println("Global teardown: after tests")
	deleteTestDb()
	os.Exit(code) // exit with the code from m.Run()
}
