package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

// Random data generators
var (
	firstNames = []string{
		"John", "Jane", "Michael", "Emily", "David", "Sarah", "James", "Emma",
		"Robert", "Olivia", "William", "Ava", "Richard", "Isabella", "Thomas", "Sophia",
		"Charles", "Mia", "Daniel", "Charlotte", "Matthew", "Amelia", "Donald", "Harper",
		"Mark", "Evelyn", "Paul", "Abigail", "Steven", "Elizabeth", "Andrew", "Sofia",
		"Joshua", "Avery", "Kenneth", "Ella", "Kevin", "Scarlett", "Brian", "Grace",
		"George", "Chloe", "Timothy", "Victoria", "Ronald", "Madison", "Edward", "Luna",
		"Jason", "Hannah", "Jeffrey", "Lily", "Ryan", "Layla", "Jacob", "Zoey",
		"Gary", "Penelope", "Nicholas", "Riley", "Eric", "Nora", "Jonathan", "Lillian",
		"Stephen", "Aubrey", "Larry", "Stella", "Justin", "Hazel", "Scott", "Ellie",
		"Brandon", "Violet", "Benjamin", "Aurora", "Samuel", "Savannah", "Frank", "Audrey",
		"Gregory", "Brooklyn", "Raymond", "Bella", "Alexander", "Claire", "Patrick", "Skylar",
		"Jack", "Lucy", "Dennis", "Paisley", "Jerry", "Everly", "Tyler", "Anna",
		"Aaron", "Caroline", "Henry", "Nova", "Douglas", "Genesis", "Peter", "Emilia",
		"Adam", "Kennedy", "Nathan", "Samantha", "Zachary", "Maya", "Walter", "Willow",
	}

	lastNames = []string{
		"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
		"Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas",
		"Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson", "White",
		"Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson", "Walker", "Young",
		"Allen", "King", "Wright", "Scott", "Torres", "Nguyen", "Hill", "Flores",
		"Green", "Adams", "Nelson", "Baker", "Hall", "Rivera", "Campbell", "Mitchell",
		"Carter", "Roberts", "Gomez", "Phillips", "Evans", "Turner", "Diaz", "Parker",
		"Cruz", "Edwards", "Collins", "Reyes", "Stewart", "Morris", "Morales", "Murphy",
		"Cook", "Rogers", "Gutierrez", "Ortiz", "Morgan", "Cooper", "Peterson", "Bailey",
		"Reed", "Kelly", "Howard", "Ramos", "Kim", "Cox", "Ward", "Richardson",
		"Watson", "Brooks", "Chavez", "Wood", "James", "Bennett", "Gray", "Mendoza",
		"Ruiz", "Hughes", "Price", "Alvarez", "Castillo", "Sanders", "Patel", "Myers",
	}

	streets = []string{
		"Main St", "Oak Ave", "Pine Rd", "Maple Dr", "Elm St", "Cedar Ln",
		"Washington Blvd", "Park Ave", "Lake Dr", "Hill St", "Forest Rd", "River St",
		"Sunset Blvd", "Broadway", "Church St", "School Rd", "Mill St", "Bridge St",
		"Water St", "Spring St", "Garden Ave", "Valley Rd", "Mountain Dr", "Beach Ave",
		"Ocean Dr", "Harbor St", "Station Rd", "College Ave", "Market St", "Union St",
	}

	cities = []string{
		"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia",
		"San Antonio", "San Diego", "Dallas", "San Jose", "Austin", "Jacksonville",
		"Fort Worth", "Columbus", "Charlotte", "San Francisco", "Indianapolis", "Seattle",
		"Denver", "Boston", "Nashville", "Detroit", "Portland", "Las Vegas",
		"Memphis", "Louisville", "Baltimore", "Milwaukee", "Albuquerque", "Tucson",
	}

	infos = []string{
		"Regular customer", "Premium member", "VIP customer", "New customer",
		"Long-time client", "Gold member", "Silver member", "Bronze member",
		"Platinum customer", "Loyal customer", "Frequent buyer", "Preferred customer",
		"Elite member", "Standard customer", "Active member", "Valued customer",
	}
)

func randomName() string {
	return fmt.Sprintf("%s %s",
		firstNames[rand.Intn(len(firstNames))],
		lastNames[rand.Intn(len(lastNames))])
}

func randomAddress() string {
	streetNum := rand.Intn(9999) + 1
	return fmt.Sprintf("%d %s, %s",
		streetNum,
		streets[rand.Intn(len(streets))],
		cities[rand.Intn(len(cities))])
}

func randomMobile() int {
	// Generate valid 10-digit mobile number (1000000000 to 9999999999)
	return rand.Intn(9000000000) + 1000000000
}

func randomInfo() string {
	return infos[rand.Intn(len(infos))]
}

func randomAmount(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomDate(daysBack int) time.Time {
	daysAgo := rand.Intn(daysBack)
	return time.Now().AddDate(0, 0, -daysAgo)
}

// init loads .env file before running tests
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
}

// TestCreateUsers - Test case 1: Create 50 users
func TestCreateUsers(t *testing.T) {
	// Initialize database connection
	initDb()
	defer db.Close()

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	numUsers := 50 // Create 50 users

	successCount := 0
	failCount := 0

	for i := 1; i <= numUsers; i++ {
		user := User{
			Name:    randomName(),
			Mobile:  randomMobile(),
			Address: randomAddress(),
			Info:    randomInfo(),
			// referred_by is ignored - will be NULL
		}

		// Marshal user to JSON
		jsonData, err := json.Marshal(user)
		if err != nil {
			t.Logf("Failed to marshal user %d: %v", i, err)
			failCount++
			continue
		}

		// Create request
		req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Logf("Failed to create request for user %d: %v", i, err)
			failCount++
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		rr := httptest.NewRecorder()

		// Call handler
		handler := http.HandlerFunc(createUser)
		handler.ServeHTTP(rr, req)

		// Check status code
		if status := rr.Code; status == http.StatusCreated {
			successCount++
			if i%25 == 0 {
				t.Logf("Progress: Created %d/%d users", i, numUsers)
			}
		} else {
			failCount++
			t.Logf("Failed to create user %d: %s - Response: %s", i, user.Name, rr.Body.String())
		}
	}

	t.Logf("✅ Successfully created %d users", successCount)
	if failCount > 0 {
		t.Logf("❌ Failed to create %d users", failCount)
	}
}

// TestCreateHandouts - Test case 2: Get users from DB and create handouts
func TestCreateHandouts(t *testing.T) {
	// Initialize database connection
	initDb()
	defer db.Close()

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Get actual user IDs from database
	rows, err := db.Query("SELECT id FROM users ORDER BY id")
	if err != nil {
		t.Fatalf("Failed to query users: %v", err)
	}
	defer rows.Close()

	var userIds []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			t.Fatalf("Failed to scan user id: %v", err)
		}
		userIds = append(userIds, id)
	}

	if len(userIds) == 0 {
		t.Fatal("No users found in database. Run TestCreateUsers first.")
	}

	t.Logf("Found %d users in database", len(userIds))

	numHandouts := 75 // Create 75 handouts

	successCount := 0
	failCount := 0

	for i := 1; i <= numHandouts; i++ {
		// Pick a random user ID from the actual user IDs
		userId := userIds[rand.Intn(len(userIds))]

		handout := HandoutUpdate{
			Date:   randomDate(365), // Random date within last year
			Amount: randomAmount(1000.00, 50000.00),
			UserId: userId,
		}

		// Marshal handout to JSON
		jsonData, err := json.Marshal(handout)
		if err != nil {
			t.Logf("Failed to marshal handout %d: %v", i, err)
			failCount++
			continue
		}

		// Create request
		req, err := http.NewRequest("POST", "/handouts", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Logf("Failed to create request for handout %d: %v", i, err)
			failCount++
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		rr := httptest.NewRecorder()

		// Call handler
		handler := http.HandlerFunc(createHandout)
		handler.ServeHTTP(rr, req)

		// Check status code
		if status := rr.Code; status == http.StatusOK {
			successCount++
			if i%25 == 0 {
				t.Logf("Progress: Created %d/%d handouts", i, numHandouts)
			}
		} else {
			failCount++
			t.Logf("Failed to create handout %d for UserId %d: Response: %s", i, handout.UserId, rr.Body.String())
		}
	}

	t.Logf("✅ Successfully created %d handouts", successCount)
	if failCount > 0 {
		t.Logf("❌ Failed to create %d handouts", failCount)
	}
}

// TestCreateCollections - Test case 3: Get handouts from DB and create collections
func TestCreateCollections(t *testing.T) {
	// Initialize database connection
	initDb()
	defer db.Close()

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Get actual handout IDs from database
	rows, err := db.Query("SELECT id FROM handouts ORDER BY id")
	if err != nil {
		t.Fatalf("Failed to query handouts: %v", err)
	}
	defer rows.Close()

	var handoutIds []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			t.Fatalf("Failed to scan handout id: %v", err)
		}
		handoutIds = append(handoutIds, id)
	}

	if len(handoutIds) == 0 {
		t.Fatal("No handouts found in database. Run TestCreateHandouts first.")
	}

	t.Logf("Found %d handouts in database", len(handoutIds))

	numCollections := 125 // Create 125 collections

	successCount := 0
	failCount := 0

	for i := 1; i <= numCollections; i++ {
		// Pick a random handout ID from the actual handout IDs
		handoutId := handoutIds[rand.Intn(len(handoutIds))]

		// Create a map instead of Collection struct to avoid sending ID and total_paid
		collectionData := map[string]interface{}{
			"date":      randomDate(365), // Random date within last year
			"amount":    randomAmount(500.00, 10000.00),
			"handoutId": handoutId,
			// total_paid is ignored - it's computed
			// id, created_at, updated_at are auto-generated
		}

		// Marshal collection to JSON
		jsonData, err := json.Marshal(collectionData)
		if err != nil {
			t.Logf("Failed to marshal collection %d: %v", i, err)
			failCount++
			continue
		}

		// Create request
		req, err := http.NewRequest("POST", "/collections", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Logf("Failed to create request for collection %d: %v", i, err)
			failCount++
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		rr := httptest.NewRecorder()

		// Call handler
		handler := http.HandlerFunc(createCollection)
		handler.ServeHTTP(rr, req)

		// Check status code
		if status := rr.Code; status == http.StatusOK {
			successCount++
			if i%25 == 0 {
				t.Logf("Progress: Created %d/%d collections", i, numCollections)
			}
		} else {
			failCount++
			t.Logf("Failed to create collection %d for HandoutId %d: Response: %s", i, handoutId, rr.Body.String())
		}
	}

	t.Logf("✅ Successfully created %d collections", successCount)
	if failCount > 0 {
		t.Logf("❌ Failed to create %d collections", failCount)
	}
}
