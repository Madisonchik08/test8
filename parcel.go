package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int64, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec(
		"Insert into parcel (client, status, address, created_at) values (?, ?, ?, ?)",
		p.Client, p.Status, p.Address, p.CreatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("error inserting parcel: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert id: %w", err)
	}
	// верните идентификатор последней добавленной записи
	return id, nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}
	row := s.db.QueryRow("Select number, client, status, address, created_at from parcel where number = ?", number)
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Parcel{}, fmt.Errorf("parcel with number %d not found", number)
		}
		return Parcel{}, fmt.Errorf("error getting parcel %d: %w", number, err)
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("Select number, client, status, address, created_at from parcel where client = ?", client)
	if err != nil {
		return nil, fmt.Errorf("error getting parcels for client %d: %w", client, err)
	}
	defer rows.Close()
	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		var p Parcel
		err = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan parcel row: %w", err)
		}
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan parcel rows: %w", err)
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	res, err := s.db.Exec("Update parcel set status = ? where number = ?", status, number)
	if err != nil {
		return fmt.Errorf("error updating parcel %d: %w", number, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected parcel %d: %w", number, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no parcel with number %d found or status already set to %s", number, status)
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	res, err := s.db.Exec("Update parcel set address = ? where number = ? and status = ?", address, number, ParcelStatusRegistered)
	if err != nil {
		return fmt.Errorf("error getting rows affected for parcel %d: %w", number, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected parcel %d: %w", number, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("cannot change address for parcel %d: not found or status is not '%s'", number, ParcelStatusRegistered)
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	res, err := s.db.Exec("Delete from parcel where number = ? and status = ?", number, ParcelStatusRegistered)
	if err != nil {
		return fmt.Errorf("error deleting parcel %d: %w", number, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected parcel %d: %w", number, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("cannot delete parcel %d: not found or status is not '%s'", number, ParcelStatusRegistered)
	}
	return nil
}
