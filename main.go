package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type ParkingLot struct {
	db       *sql.DB
	capacity int
}

type Car struct {
	SlotNumber     int
	RegistrationNo string
}

func NewParkingLot(db *sql.DB) *ParkingLot {
	return &ParkingLot{
		db: db,
	}
}

func (pl *ParkingLot) InitializeDatabase() error {
	createParkingLotsTable := `
	CREATE TABLE IF NOT EXISTS parking_lots (
		id INT AUTO_INCREMENT PRIMARY KEY,
		capacity INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	createParkingSlotsTable := `
	CREATE TABLE IF NOT EXISTS parking_slots (
		slot_number INT PRIMARY KEY,
		registration_no VARCHAR(50),
		is_occupied BOOLEAN DEFAULT FALSE,
		parked_at TIMESTAMP NULL,
		INDEX idx_registration (registration_no),
		INDEX idx_occupied (is_occupied)
	)`

	createHistoryTable := `
	CREATE TABLE IF NOT EXISTS parking_history (
		id INT AUTO_INCREMENT PRIMARY KEY,
		registration_no VARCHAR(50) NOT NULL,
		slot_number INT NOT NULL,
		parked_at TIMESTAMP NOT NULL,
		left_at TIMESTAMP NULL,
		hours_parked INT DEFAULT 0,
		charge_amount DECIMAL(10,2) DEFAULT 0.00,
		INDEX idx_registration_history (registration_no)
	)`

	tables := []string{createParkingLotsTable, createParkingSlotsTable, createHistoryTable}

	for _, table := range tables {
		if _, err := pl.db.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}

	return nil
}

func (pl *ParkingLot) CreateParkingLot(capacity int) error {

	if _, err := pl.db.Exec("DELETE FROM parking_slots"); err != nil {
		return fmt.Errorf("failed to clear parking slots: %v", err)
	}

	if _, err := pl.db.Exec("DELETE FROM parking_lots"); err != nil {
		return fmt.Errorf("failed to clear parking lots: %v", err)
	}

	if _, err := pl.db.Exec("INSERT INTO parking_lots (capacity) VALUES (?)", capacity); err != nil {
		return fmt.Errorf("failed to create parking lot: %v", err)
	}

	for i := 1; i <= capacity; i++ {
		if _, err := pl.db.Exec("INSERT INTO parking_slots (slot_number, is_occupied) VALUES (?, FALSE)", i); err != nil {
			return fmt.Errorf("failed to initialize slot %d: %v", i, err)
		}
	}

	pl.capacity = capacity
	fmt.Printf("Created a parking lot with %d slots\n", capacity)
	return nil
}

func (pl *ParkingLot) Park(registrationNo string) error {
	var existingSlot int
	err := pl.db.QueryRow("SELECT slot_number FROM parking_slots WHERE registration_no = ? AND is_occupied = TRUE", registrationNo).Scan(&existingSlot)
	if err == nil {
		fmt.Printf("Car %s is already parked in slot %d\n", registrationNo, existingSlot)
		return nil
	}

	var slotNumber int
	err = pl.db.QueryRow("SELECT slot_number FROM parking_slots WHERE is_occupied = FALSE ORDER BY slot_number LIMIT 1").Scan(&slotNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Sorry, parking lot is full")
			return nil
		}
		return fmt.Errorf("failed to find available slot: %v", err)
	}

	if _, err := pl.db.Exec("UPDATE parking_slots SET registration_no = ?, is_occupied = TRUE, parked_at = NOW() WHERE slot_number = ?", registrationNo, slotNumber); err != nil {
		return fmt.Errorf("failed to allocate slot: %v", err)
	}

	if _, err := pl.db.Exec("INSERT INTO parking_history (registration_no, slot_number, parked_at) VALUES (?, ?, NOW())", registrationNo, slotNumber); err != nil {
		return fmt.Errorf("failed to record parking history: %v", err)
	}

	fmt.Printf("Allocated slot number: %d\n", slotNumber)
	return nil
}

func (pl *ParkingLot) Leave(registrationNo string, hours int) error {
	var slotNumber int
	err := pl.db.QueryRow("SELECT slot_number FROM parking_slots WHERE registration_no = ? AND is_occupied = TRUE", registrationNo).Scan(&slotNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Registration number %s not found\n", registrationNo)
			return nil
		}
		return fmt.Errorf("failed to find car: %v", err)
	}

	charge := pl.calculateCharge(hours)

	if _, err := pl.db.Exec("UPDATE parking_slots SET registration_no = NULL, is_occupied = FALSE, parked_at = NULL WHERE slot_number = ?", slotNumber); err != nil {
		return fmt.Errorf("failed to free slot: %v", err)
	}

	if _, err := pl.db.Exec("UPDATE parking_history SET left_at = NOW(), hours_parked = ?, charge_amount = ? WHERE registration_no = ? AND left_at IS NULL", hours, charge, registrationNo); err != nil {
		return fmt.Errorf("failed to update parking history: %v", err)
	}

	fmt.Printf("Registration number %s with Slot Number %d is free with Charge $%d\n", registrationNo, slotNumber, charge)
	return nil
}

func (pl *ParkingLot) calculateCharge(hours int) int {
	if hours <= 2 {
		return 10
	}
	return 10 + (hours-2)*10
}

func (pl *ParkingLot) Status() error {
	rows, err := pl.db.Query("SELECT slot_number, registration_no FROM parking_slots WHERE is_occupied = TRUE ORDER BY slot_number")
	if err != nil {
		return fmt.Errorf("failed to get status: %v", err)
	}
	defer rows.Close()

	fmt.Println("Slot No.\tRegistration No.")
	for rows.Next() {
		var slotNumber int
		var registrationNo string
		if err := rows.Scan(&slotNumber, &registrationNo); err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}
		fmt.Printf("%d\t\t%s\n", slotNumber, registrationNo)
	}

	return nil
}

func (pl *ParkingLot) ProcessCommand(command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil
	}

	switch parts[0] {
	case "create_parking_lot":
		if len(parts) != 2 {
			return fmt.Errorf("invalid create_parking_lot command")
		}
		capacity, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid capacity: %v", err)
		}
		return pl.CreateParkingLot(capacity)

	case "park":
		if len(parts) != 2 {
			return fmt.Errorf("invalid park command")
		}
		return pl.Park(parts[1])

	case "leave":
		if len(parts) != 3 {
			return fmt.Errorf("invalid leave command")
		}
		hours, err := strconv.Atoi(parts[2])
		if err != nil {
			return fmt.Errorf("invalid hours: %v", err)
		}
		return pl.Leave(parts[1], hours)

	case "status":
		return pl.Status()

	default:
		return fmt.Errorf("unknown command: %s", parts[0])
	}
}

func (pl *ParkingLot) ProcessFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		command := strings.TrimSpace(scanner.Text())
		if command == "" {
			continue
		}

		if err := pl.ProcessCommand(command); err != nil {
			fmt.Printf("Error processing command '%s': %v\n", command, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	return nil
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		os.Exit(1)
	}

	filename := os.Args[1]

	dsn := "root@tcp(localhost:3306)/parking_system?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	parkingLot := NewParkingLot(db)

	if err := parkingLot.InitializeDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	if err := parkingLot.ProcessFile(filename); err != nil {
		log.Fatal("Failed to process file:", err)
	}
}
