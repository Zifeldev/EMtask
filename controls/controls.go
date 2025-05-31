package controls

import (
	"context"
	"em/db"
	"em/enrich"
	"fmt"
	"log"
)

// Person represents a person entity
// @Description Person information with all fields
type Person struct {
	ID          int    `json:"id" example:"1"`
	Name        string `json:"name" example:"Dmitriy"`
	Surname     string `json:"surname" example:"Ushakov"`
	Patronymic  string `json:"patronymic,omitempty" example:"Vasilevich"`
	Gender      string `json:"gender"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
}

func InsertPerson(person *Person) error {
	log.Printf("DEBUG: Starting enrichment for person: %s", person.Name)

	age, err := enrich.GetAge(person.Name)
	if err != nil {
		log.Printf("ERROR: Age enrichment failed for %s: %v", person.Name, err)
		return fmt.Errorf("age enrichment failed: %w", err)
	}

	gender, err := enrich.GetGender(person.Name)
	if err != nil {
		log.Printf("ERROR: Gender enrichment failed for %s: %v", person.Name, err)
		return fmt.Errorf("gender enrichment failed: %w", err)
	}

	nationality, err := enrich.GetNationality(person.Name)
	if err != nil {
		log.Printf("ERROR: Nationality enrichment failed for %s: %v", person.Name, err)
		return fmt.Errorf("nationality enrichment failed: %w", err)
	}

	person.Age = age
	person.Gender = gender
	person.Nationality = nationality

	log.Printf("INFO: Enriched person data: Age=%d, Gender=%s, Nationality=%s",
		person.Age, person.Gender, person.Nationality)

	_, err = db.Conn.Exec(context.Background(),
		"INSERT INTO people (name, surname, patronymic, gender, age, nationality) VALUES ($1, $2, $3, $4, $5, $6)",
		person.Name, person.Surname, person.Patronymic, person.Gender, person.Age, person.Nationality,
	)
	if err != nil {
		log.Printf("ERROR: Database insert failed for %s %s: %v", person.Name, person.Surname, err)
		return fmt.Errorf("database insert failed: %w", err)
	}

	log.Printf("INFO: Successfully inserted person: %s %s (ID: %d)",
		person.Name, person.Surname, person.ID)
	return nil
}

func DeletePerson(id int) error {
	log.Printf("DEBUG: Attempting to delete person with ID: %d", id)

	_, err := db.Conn.Exec(context.Background(), "DELETE FROM people where id=$1", id)
	if err != nil {
		log.Printf("ERROR: Delete failed for ID %d: %v", id, err)
		return fmt.Errorf("database delete failed: %w", err)
	}

	log.Printf("INFO: Successfully deleted person with ID: %d", id)
	return nil
}

func GetFilteredPersons(name, surname, gender, nationality, age, limit, offset string) ([]Person, error) {
	log.Printf("DEBUG: Building query with filters: name=%s, surname=%s, gender=%s, nationality=%s, age=%s",
		name, surname, gender, nationality, age)

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

	log.Printf("DEBUG: Executing query: %s\nParameters: %v", query, args)

	rows, err := db.Conn.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("ERROR: Query execution failed: %v", err)
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var people []Person
	count := 0
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Gender, &p.Age, &p.Nationality); err != nil {
			log.Printf("ERROR: Row scanning failed: %v", err)
			return nil, fmt.Errorf("data scanning failed: %w", err)
		}
		people = append(people, p)
		count++
	}

	log.Printf("INFO: Retrieved %d persons matching filters", count)
	return people, nil
}

func UpdatePerson(id, newAge int, newName, newSurname, newPatronymic, newGender, newNationality string) error {
	log.Printf("DEBUG: Updating person ID %d with: name=%s, surname=%s, age=%d, gender=%s, nationality=%s",
		id, newName, newSurname, newAge, newGender, newNationality)

	_, err := db.Conn.Exec(context.Background(),
		"UPDATE people SET name=$1, surname=$2, patronymic=$4, gender=$5, age=$6, nationality=$7 WHERE id=$3",
		newName, newSurname, id, newPatronymic, newGender, newAge, newNationality)

	if err != nil {
		log.Printf("ERROR: Update failed for ID %d: %v", id, err)
	}

	log.Printf("INFO: Successfully updated person with ID: %d", id)
	return nil
}
