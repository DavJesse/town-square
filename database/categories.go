package database

import (
	"fmt"
	"strings"

	"forum/models"
)

// initializes the categories table with the predefined categories
func InitCategories() error {
	categories := []string{
		"agriculture",
		"arts",
		"education",
		"lifestyle",
		"technology",
		"culture",
		"science",
		"miscellaneous",
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// defer a rollback in case of failure
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for _, category := range categories {
		query := `INSERT OR IGNORE INTO categories (name) VALUES (?)`
		_, err := tx.Exec(query, category)
		if err != nil {
			return fmt.Errorf("failed to insert category %s: %w", category, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// fetchCategories retrieves all categories from the database
func FetchCategories() ([]models.Category, error) {
	query := `SELECT id, name FROM categories`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// validateCategories checks if the provided category IDs exist in the database
func ValidateCategories(categoryIDs []int) error {
	query := `SELECT COUNT(*) FROM categories WHERE id IN (?` + strings.Repeat(",?", len(categoryIDs)-1) + `)`
	args := make([]interface{}, len(categoryIDs))
	for i, id := range categoryIDs {
		args[i] = id
	}

	var count int
	err := db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return err
	}

	if count != len(categoryIDs) {
		return fmt.Errorf("one or more categories do not exist")
	}

	return nil
}
