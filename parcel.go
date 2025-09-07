package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :createdat)",
		sql.Named("client", p.Client), sql.Named("status", p.Status), sql.Named("address", p.Address),
		sql.Named("createdat", p.CreatedAt))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number", sql.Named("number", number))

	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return Parcel{}, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", sql.Named("client", client))
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return []Parcel{}, err
	}

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return []Parcel{}, err
		}
		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number", sql.Named("status", status), sql.Named("number", number))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))
	p := Parcel{}
	err := row.Scan(&p.Status)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if p.Status == "registered" {
		_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number", sql.Named("address", address), sql.Named("number", number))
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	}

	return fmt.Errorf("Нужен статус: registered. Ваш статус: %s", p.Status)
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))
	p := Parcel{}
	err := row.Scan(&p.Status)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if p.Status == "registered" {
		_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	}

	return fmt.Errorf("Нужен статус: registered. Ваш статус: %s", p.Status)

}
