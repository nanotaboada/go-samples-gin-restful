// Package tests provides integration and utility code to support automated
// testing of the application.
package tests

import (
	"errors"

	"github.com/nanotaboada/go-samples-gin-restful/model"
)

// MockPlayerService is a test double that implements service.PlayerService.
//
// # Go interfaces and implicit satisfaction
//
// In Go, a type satisfies an interface simply by having the right method
// signatures — no "implements" keyword or explicit declaration is required.
// MockPlayerService satisfies service.PlayerService because it defines all six
// methods with matching signatures.  The compiler verifies this at the call
// site (e.g. controller.NewPlayerController(mockService)).
//
// # Opt-in override pattern
//
// Each method has a corresponding Func field (e.g. CreateFunc).  When a test
// sets one of these fields, the method delegates to it; when the field is nil
// the method returns a safe zero-value default.  This lets each test override
// only the methods relevant to the scenario being tested, leaving the rest as
// no-ops, without creating a new type per scenario.
type MockPlayerService struct {
	CreateFunc                func(player *model.Player) error
	RetrieveAllFunc           func() ([]model.Player, error)
	RetrieveByIDFunc          func(id string) (model.Player, error)
	RetrieveBySquadNumberFunc func(squadNumber int) (model.Player, error)
	UpdateFunc                func(player *model.Player) error
	DeleteFunc                func(player *model.Player) error
}

// Create delegates to CreateFunc if set, otherwise returns nil (no-op success).
func (m *MockPlayerService) Create(player *model.Player) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(player)
	}
	return nil
}

// RetrieveAll delegates to RetrieveAllFunc if set, otherwise returns an empty slice.
func (m *MockPlayerService) RetrieveAll() ([]model.Player, error) {
	if m.RetrieveAllFunc != nil {
		return m.RetrieveAllFunc()
	}
	return []model.Player{}, nil
}

// RetrieveByID delegates to RetrieveByIDFunc if set, otherwise returns a zero-value Player.
func (m *MockPlayerService) RetrieveByID(id string) (model.Player, error) {
	if m.RetrieveByIDFunc != nil {
		return m.RetrieveByIDFunc(id)
	}
	return model.Player{}, nil
}

// RetrieveBySquadNumber delegates to RetrieveBySquadNumberFunc if set, otherwise returns a zero-value Player.
func (m *MockPlayerService) RetrieveBySquadNumber(squadNumber int) (model.Player, error) {
	if m.RetrieveBySquadNumberFunc != nil {
		return m.RetrieveBySquadNumberFunc(squadNumber)
	}
	return model.Player{}, nil
}

// Update delegates to UpdateFunc if set, otherwise returns nil (no-op success).
func (m *MockPlayerService) Update(player *model.Player) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(player)
	}
	return nil
}

// Delete delegates to DeleteFunc if set, otherwise returns nil (no-op success).
func (m *MockPlayerService) Delete(player *model.Player) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(player)
	}
	return nil
}

// Sentinel errors used by mock-assisted tests to simulate failure conditions
// that cannot be triggered naturally with a healthy in-memory SQLite database.
var (
	ErrDatabaseFailure = errors.New("database connection failed")
	ErrGenericError    = errors.New("generic internal error")
)
