package controls

import (
	"context"
	"em/db"
	"fmt"
	"em/enrich"
)

type Person struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Gender      string `json:"gender"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
}

func InsertPerson(person *Person) error {
	age, err := enrich.GetAge(person.Name)
    if err != nil {
        fmt.Println("Error with age:", err)
    }

    gender, err := enrich.GetGender(person.Name)
    if err != nil {
        fmt.Println("Error with gender", err)
    }

    nationality, err := enrich.GetNationality(person.Name)
    if err != nil {
        fmt.Println("Error with nationality", err)
    }

    person.Age = age
    person.Gender = gender
    person.Nationality = nationality

    _, err = db.Conn.Exec(context.Background(),
        "INSERT INTO people (name, surname, patronymic, gender, age, nationality) VALUES ($1, $2, $3, $4, $5, $6)",
        person.Name, person.Surname, person.Patronymic, person.Gender, person.Age, person.Nationality,
    )
    if err != nil {
        return fmt.Errorf("error with  inserting: %w", err)
    }

    return nil
}

func DeletePerson(id int) error {
	_, err := db.Conn.Exec(context.Background(), "DELETE FROM people where id=$1", id)

	if err != nil {
		fmt.Println("Error with deleting", err)
	}
	return err
}

func GetFilteredPersons(name, surname, gender, nationality, age, limit, offset string) ([]Person, error) {
	query := "SELECT id, name, surname, patronymic, gender, age, nationality FROM people WHERE 1=1"
	args := []interface{}{}
	argID := 1

	if name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argID)
		args = append(args, "%"+name+"%")
		argID++
	}
	if surname != "" {
		query += fmt.Sprintf(" AND surname ILIKE $%d", argID)
		args = append(args, "%"+surname+"%")
		argID++
	}
	if gender != "" {
		query += fmt.Sprintf(" AND gender = $%d", argID)
		args = append(args, gender)
		argID++
	}
	if nationality != "" {
		query += fmt.Sprintf(" AND nationality = $%d", argID)
		args = append(args, nationality)
		argID++
	}
	if age != "" {
		query += fmt.Sprintf(" AND age = $%d", argID)
		args = append(args, age)
		argID++
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	rows, err := db.Conn.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []Person
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Gender, &p.Age, &p.Nationality); err != nil {
			return nil, err
		}
		people = append(people, p)
	}

	return people, nil
}


func UpdatePerson(id,newAge int, newName, newSurname, newPatronymic,newGender,newNationality string) error {
	_, err := db.Conn.Exec(context.Background(), "UPDATE people SET name=$1, surname=$2, patronymic=$4, gender=$5,age=$6,nationality=$7 where id=$3", newName, newSurname, id,newPatronymic,newGender,newAge,newNationality)
	if err != nil {
		fmt.Println("Error with update", err)
	}
	return err
}
