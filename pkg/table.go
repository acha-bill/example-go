package pkg

import (
	"fmt"
	"sort"
)

var (
	ErrNotFound     = fmt.Errorf("not found")
	ErrAlreadyHasID = fmt.Errorf("already has an ID")
)

// Table is the interface that all tables must implement
type Table interface {
	// Name returns the name of the table
	Name() string
	// Insert inserts a new model into the database
	Insert(Model) error
	// Update updates an existing model in the database
	Update(Model) error
	// Delete deletes an existing model from the database
	Delete(PrimaryKey) error
	// Get returns a model from the database by its primary key
	Get(PrimaryKey) (Model, error)
	// Find returns a slice of models from the database that match the given function.
	Find(func(Model) bool) ([]Model, error)
}

type table struct {
	name   string
	lastID PrimaryKey
	data   map[PrimaryKey]Model
}

func (t *table) Name() string {
	return t.name
}

func (t *table) Insert(model Model) error {
	if model.GetID() != 0 {
		return ErrAlreadyHasID
	}
	t.lastID++
	model.SetID(t.lastID)
	t.data[t.lastID] = model
	return nil
}

func (t *table) Update(model Model) error {
	if _, ok := t.data[model.GetID()]; !ok {
		return ErrNotFound
	}
	t.data[model.GetID()] = model
	return nil
}

func (t *table) Delete(key PrimaryKey) error {
	if _, ok := t.data[key]; !ok {
		return ErrNotFound
	}
	delete(t.data, key)
	return nil
}

func (t *table) Get(key PrimaryKey) (Model, error) {
	if model, ok := t.data[key]; ok {
		return model, nil
	}
	return nil, ErrNotFound
}

func (t *table) Find(f func(Model) bool) ([]Model, error) {
	var models []Model
	for _, model := range t.data {
		if f(model) {
			models = append(models, model)
		}
	}
	sort.Slice(models, func(i, j int) bool {
		return models[i].GetID() < models[j].GetID()
	})
	return models, nil
}

var _ Table = &table{}
