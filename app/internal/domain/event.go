package domain

import "time"

type Event struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	StartDate   time.Time         `json:"start_date"`
	EndDate     time.Time         `json:"end_date"`
	Categories  []EventCategories `json:"categories"`
	Cover       string            `json:"cover"`
}

type EventCategories struct {
	ID          string `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Slug        string `json:"slug"`
}
