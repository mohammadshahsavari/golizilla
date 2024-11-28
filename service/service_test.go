package service

import (
	"testing"

	"golizilla/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRoleRepo struct {
	mock.Mock
}

// Mock implementation of CreateRole
func (m *MockRoleRepo) CreateRole(role *model.Role) error {
	return m.Called(role).Error(0)
}

// Mock implementation of GetRoleByID
func (m *MockRoleRepo) GetRoleByID(id uint) (*model.Role, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Role), args.Error(1)
}

// Mock implementation of GetAllRoles
func (m *MockRoleRepo) GetAllRoles() ([]*model.Role, error) {
	args := m.Called()
	return args.Get(0).([]*model.Role), args.Error(1)
}

// Mock implementation of DeleteRole
func (m *MockRoleRepo) DeleteRole(id uint) error {
	return m.Called(id).Error(0)
}

func TestCreateRole(t *testing.T) {
	mockRepo := new(MockRoleRepo)
	mockRepo.On("CreateRole", mock.Anything).Return(nil)

	service := NewRoleService(mockRepo)
	role, err := service.CreateRole("Admin", "Administrator role")

	assert.NoError(t, err)
	assert.Equal(t, "Admin", role.Name)
}

func TestGetAllRoles(t *testing.T) {
	// Arrange
	mockRepo := new(MockRoleRepo)
	service := NewRoleService(mockRepo)

	roles := []*model.Role{
		{ID: 1, Name: "Admin", Description: "Administrator role"},
		{ID: 2, Name: "User", Description: "User role"},
	}

	mockRepo.On("GetAllRoles").Return(roles, nil)

	// Act
	result, err := service.GetAllRoles()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Admin", result[0].Name)
	mockRepo.AssertExpectations(t)
}
