package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Storage interface {
	DeleteUser(string) error
	GetUserById(string) (*User, error)
	//SaveDataToCache(string) (*User, error)
	CreateUserFromNATS(msgCh chan []byte) error
	GetAllUsers() (*User, error)
}

type PostgresStore struct {
	db *pgx.Conn
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "postgres://wb_user:wb_password@localhost:5434/wb_db"

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	//defer conn.Close(context.Background())
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: conn,
	}, nil
}

func (s *PostgresStore) CreateUserFromNATS(msgCh chan []byte) error {
	u := new(User)
	for {

		msg := <-msgCh
		err := json.Unmarshal(msg, u)
		_, err = s.db.Exec(context.Background(),
			"INSERT INTO orders (order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) "+
				"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
			u.OrderUid, u.TrackNumber, u.Entry, u.Delivery, u.Payment, u.Items, u.Locale, u.InternalSignature, u.CustomerID, u.DeliveryService, u.Shardkey, u.SmID, u.DateCreated, u.OofShard)
		if err != nil {
			return err
		}

		fmt.Println("Received and processed message from NATS")
		return nil
	}
}

func (s *PostgresStore) DeleteUser(id string) error {
	query := "delete from orders where order_uid = $1"
	_, err := s.db.Query(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) GetUserById(id string) (*User, error) {
	var user User
	query := `SELECT * FROM orders WHERE order_uid = $1;`
	rows, err := s.db.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&user.OrderUid,
			&user.TrackNumber,
			&user.Entry,
			&user.Delivery,
			&user.Payment,
			&user.Items,
			&user.Locale,
			&user.InternalSignature,
			&user.CustomerID,
			&user.DeliveryService,
			&user.Shardkey,
			&user.SmID,
			&user.DateCreated,
			&user.OofShard,
		); err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *PostgresStore) GetAllUsers() (*User, error) {

	var u *User

	query := `SELECT * FROM orders`
	_, err := s.db.Query(context.Background(), query)
	if err != nil {
		return u, err
	}
	return u, nil
}
