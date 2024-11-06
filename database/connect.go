package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func DatabaseConnection(connStr string) (*sql.DB, error) {
	var err error

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to the database: %v", err)
	}

	//create user table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT,
			password TEXT,
			email TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	//alter Table to add Role

	_, err = db.Exec(`
		ALTER TABLE users ADD COLUMN IF NOT EXISTS role TEXT DEFAULT 'user'
	`)
	if err != nil {
		log.Fatal(err)
	}

	//create adming user if not exists
	hashedPassword, err := hashPassword("admin")

	_, err = db.Exec(`
			INSERT INTO users (username, password, email, role)
			SELECT 'admin', $1, 'admin@admin.com', 'admin'
			WHERE NOT EXISTS (
				SELECT 1 FROM users WHERE username = 'admin'
			)
		`, hashedPassword)

	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name TEXT,
			description TEXT,
			price REAL,
			stock INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	//create order table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER,
			status TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	//create order item table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id INTEGER,
			product_id INTEGER,
			quantity INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	// Alter order item table to add foreign key for order_id if not already exist
	_, err = db.Exec(`
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1
                FROM information_schema.table_constraints
                WHERE constraint_name = 'fk_order_id'
            ) THEN
                ALTER TABLE order_items ADD CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders(id);
            END IF;
        END $$;
    `)
	if err != nil {
		log.Println("Error adding fk_order_id:", err)
	}

	// Alter order item table to add foreign key for product_id
	_, err = db.Exec(`
			DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_product_id'
    ) THEN
        ALTER TABLE order_items ADD CONSTRAINT fk_product_id FOREIGN KEY (product_id) REFERENCES products(id);
    END IF;
END $$
		`)
	if err != nil {
		log.Fatal(err)

	}

	//add user foriegn key to order table
	_, err = db.Exec(`
		DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_user_id'
    ) THEN
        ALTER TABLE orders ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id);
    END IF;
END $$;
	`)

	DB = db
	return db, nil

}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
