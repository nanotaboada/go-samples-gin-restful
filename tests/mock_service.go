// Package tests provides integration and utility code to support automated
// testing of the application.
package tests

import (
	"errors"

	"github.com/nanotaboada/go-samples-gin-restful/model"
)

// MockPlayerService is a mock implementation of PlayerService for testing error scenarios
type MockPlayerService struct {
	CreateFunc                func(player *model.Player) error
	RetrieveAllFunc           func() ([]model.Player, error)
	RetrieveByIDFunc          func(id int) (model.Player, error)
	RetrieveBySquadNumberFunc func(squadNumber int) (model.Player, error)
	UpdateFunc                func(player *model.Player) error
	DeleteFunc                func(id int) error
}

// Create delegates to CreateFunc or returns nil
func (m *MockPlayerService) Create(player *model.Player) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(player)
	}
	return nil
}

// RetrieveAll delegates to RetrieveAllFunc or returns empty slice
func (m *MockPlayerService) RetrieveAll() ([]model.Player, error) {
	if m.RetrieveAllFunc != nil {
		return m.RetrieveAllFunc()
	}
	return []model.Player{}, nil
}

// RetrieveByID delegates to RetrieveByIDFunc or returns empty player
func (m *MockPlayerService) RetrieveByID(id int) (model.Player, error) {
	if m.RetrieveByIDFunc != nil {
		return m.RetrieveByIDFunc(id)
	}
	return model.Player{}, nil
}

// RetrieveBySquadNumber delegates to RetrieveBySquadNumberFunc or returns empty player
func (m *MockPlayerService) RetrieveBySquadNumber(squadNumber int) (model.Player, error) {
	if m.RetrieveBySquadNumberFunc != nil {
		return m.RetrieveBySquadNumberFunc(squadNumber)
	}
	return model.Player{}, nil
}

// Update delegates to UpdateFunc or returns nil
func (m *MockPlayerService) Update(player *model.Player) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(player)
	}
	return nil
}

// Delete delegates to DeleteFunc or returns nil
func (m *MockPlayerService) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// Common errors for testing
var (
	ErrDatabaseFailure = errors.New("database connection failed")
	ErrGenericError    = errors.New("generic internal error")
)
