package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"os/exec"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	DBName     = "stumpfworks_nas"
	DBUser     = "stumpfworks"
	ConfigDir  = "/etc/stumpfworks-nas"
	PasswordFile = "/etc/stumpfworks-nas/.db-password"
)

var (
	dbHost     = flag.String("host", "localhost", "PostgreSQL host")
	dbPort     = flag.Int("port", 5432, "PostgreSQL port")
	adminUser  = flag.String("admin-user", "postgres", "PostgreSQL admin user")
	skipCreate = flag.Bool("skip-create", false, "Skip database/user creation (only generate password)")
	verbose    = flag.Bool("verbose", false, "Verbose output")
)

func main() {
	flag.Parse()

	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║     StumpfWorks NAS - PostgreSQL Database Setup          ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Check if PostgreSQL is installed
	if err := checkPostgreSQL(); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		fmt.Println("\nPlease install PostgreSQL first:")
		fmt.Println("  sudo apt install postgresql postgresql-client")
		os.Exit(1)
	}

	// Generate secure password
	password, err := generatePassword(32)
	if err != nil {
		fmt.Printf("❌ Failed to generate password: %v\n", err)
		os.Exit(1)
	}

	// Save password to file
	if err := savePassword(password); err != nil {
		fmt.Printf("❌ Failed to save password: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Password saved to %s\n", PasswordFile)

	if *skipCreate {
		fmt.Println("\n✅ Password generated successfully (database creation skipped)")
		os.Exit(0)
	}

	// Connect to PostgreSQL as admin
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=postgres sslmode=disable",
		*dbHost, *dbPort, *adminUser)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Printf("❌ Failed to connect to PostgreSQL: %v\n", err)
		fmt.Println("\nTry running as postgres user:")
		fmt.Printf("  sudo -u postgres %s\n", os.Args[0])
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("❌ Failed to ping PostgreSQL: %v\n", err)
		fmt.Println("\nMake sure PostgreSQL is running:")
		fmt.Println("  sudo systemctl start postgresql")
		os.Exit(1)
	}

	fmt.Println("\n📦 Setting up database...")

	// Create user if not exists
	if err := createUser(db, DBUser, password); err != nil {
		fmt.Printf("❌ Failed to create user: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ User '%s' created/updated\n", DBUser)

	// Create database if not exists
	if err := createDatabase(db, DBName, DBUser); err != nil {
		fmt.Printf("❌ Failed to create database: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Database '%s' created\n", DBName)

	// Test connection with new credentials
	if err := testConnection(*dbHost, *dbPort, DBUser, password, DBName); err != nil {
		fmt.Printf("⚠️  Warning: Could not verify connection: %v\n", err)
	} else {
		fmt.Println("✓ Database connection verified")
	}

	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║  ✅ PostgreSQL setup complete!                            ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Database credentials:")
	fmt.Printf("  Host:     %s\n", *dbHost)
	fmt.Printf("  Port:     %d\n", *dbPort)
	fmt.Printf("  Database: %s\n", DBName)
	fmt.Printf("  User:     %s\n", DBUser)
	fmt.Printf("  Password: (saved in %s)\n", PasswordFile)
	fmt.Println()
}

// checkPostgreSQL checks if PostgreSQL is installed
func checkPostgreSQL() error {
	cmd := exec.Command("psql", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("PostgreSQL not found")
	}
	return nil
}

// generatePassword generates a cryptographically secure random password
func generatePassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// savePassword saves the password to a file
func savePassword(password string) error {
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(ConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write password file
	if err := os.WriteFile(PasswordFile, []byte(password), 0600); err != nil {
		return fmt.Errorf("failed to write password file: %w", err)
	}

	return nil
}

// createUser creates a PostgreSQL user
func createUser(db *sql.DB, username, password string) error {
	// Check if user exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname = $1)", username).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if exists {
		// Update password
		query := fmt.Sprintf("ALTER USER %s WITH PASSWORD '%s'", username, password)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to update user password: %w", err)
		}
	} else {
		// Create user
		query := fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", username, password)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	return nil
}

// createDatabase creates a PostgreSQL database
func createDatabase(db *sql.DB, dbname, owner string) error {
	// Check if database exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbname).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	if exists {
		// Grant privileges
		query := fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", dbname, owner)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to grant privileges: %w", err)
		}
		return nil
	}

	// Create database
	query := fmt.Sprintf("CREATE DATABASE %s OWNER %s", dbname, owner)
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	return nil
}

// testConnection tests the database connection
func testConnection(host string, port int, user, password, dbname string) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Ping()
}
