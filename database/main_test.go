package database // Global setup before tests
import (
	"fmt"
	"os"
	"testing"
	"time"
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
	_, err = CreateUser("milton", "milton@mail.com", "$2a$10$no3SHHuvF8C6gR.zUSlsJu46Hq8kIMRSqp5hpj/5b1XXMlROPpXIa")
	if err != nil {
		error_s := fmt.Errorf("%w\n", err)
		fmt.Println(error_s)
		os.Exit(1)
	}

	// create a single session that will be used for testing
	err = CreateSession("3cb056f2-3458-4137-863b-32751e68928e", 1, time.Now().Add(time.Hour))
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
