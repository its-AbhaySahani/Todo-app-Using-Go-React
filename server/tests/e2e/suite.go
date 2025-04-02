package e2e

import (
    "context"
    "database/sql"
    "fmt"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/docker/go-connections/nat"
    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/handler"
    "github.com/its-AbhaySahani/Todo-app-Using-Go-React/tests/e2e/helpers"
    "github.com/stretchr/testify/suite"
    tc "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

type E2ETestSuite struct {
    suite.Suite
    container tc.Container
    db        *sql.DB
    server    *httptest.Server
    client    *helpers.TestClient
}

func initializeTestDatabase(db *sql.DB) error {
    // Create tables
    queries := []string{
        `CREATE TABLE IF NOT EXISTS users (
            id VARCHAR(255) PRIMARY KEY,
            username VARCHAR(255) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
        `CREATE TABLE IF NOT EXISTS todos (
            id VARCHAR(255) PRIMARY KEY,
            task VARCHAR(255) NOT NULL,
            description TEXT,
            done BOOLEAN DEFAULT false,
            important BOOLEAN DEFAULT false,
            user_id VARCHAR(255),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id)
        )`,
    }

    for _, query := range queries {
        if _, err := db.Exec(query); err != nil {
            return fmt.Errorf("failed to execute query: %v", err)
        }
    }
    return nil
}

func (s *E2ETestSuite) SetupSuite() {
    var err error

    // Start MySQL container
    ctx := context.Background()
    req := tc.ContainerRequest{
        Image:        "mysql:8",
        ExposedPorts: []string{"3306/tcp"},
        Env: map[string]string{
            "MYSQL_ROOT_PASSWORD": "testpass",
            "MYSQL_DATABASE":     "testdb",
        },
        WaitingFor: wait.ForAll(
            wait.ForLog("port: 3306  MySQL Community Server"),
            wait.ForSQL("3306/tcp", "mysql", func(host string, port nat.Port) string {
                return fmt.Sprintf("root:testpass@tcp(%s:%s)/testdb", host, port.Port())
            }),
        ).WithDeadline(time.Minute * 2),
    }

    s.container, err = tc.GenericContainer(ctx, tc.GenericContainerRequest{
        ContainerRequest: req,
        Started:         true,
    })
    if err != nil {
        s.T().Fatalf("Failed to start container: %v", err)
    }

    // Get container connection details
    mappedPort, err := s.container.MappedPort(ctx, "3306")
    if err != nil {
        s.T().Fatalf("Failed to get container port: %v", err)
    }

    hostIP, err := s.container.Host(ctx)
    if err != nil {
        s.T().Fatalf("Failed to get container host: %v", err)
    }

    // Setup database connection
    dsn := fmt.Sprintf("root:testpass@tcp(%s:%s)/testdb?parseTime=true", hostIP, mappedPort.Port())
    s.db, err = sql.Open("mysql", dsn)
    if err != nil {
        s.T().Fatalf("Failed to connect to database: %v", err)
    }

    // Initialize test database
    if err := initializeTestDatabase(s.db); err != nil {
        s.T().Fatalf("Failed to initialize test database: %v", err)
    }

    // Setup HTTP test server
    router := mux.NewRouter()
    handler.SetupRoutes(router, s.db)
    s.server = httptest.NewServer(router)

    // Initialize test client
    s.client = helpers.NewTestClient(s.server.URL)
}

func (s *E2ETestSuite) TearDownSuite() {
    if s.server != nil {
        s.server.Close()
    }

    if s.db != nil {
        s.db.Close()
    }

    if s.container != nil {
        ctx := context.Background()
        if err := s.container.Terminate(ctx); err != nil {
            s.T().Errorf("Failed to terminate container: %v", err)
        }
    }
}

// Helper method to run the test suite
func RunE2ETests(t *testing.T) {
    suite.Run(t, new(E2ETestSuite))
}