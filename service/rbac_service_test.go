package service_test

import (
	"golizilla/domain/model"
	"golizilla/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) Create(role *model.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleRepository) GetByID(id uint) (*model.Role, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByName(name string) (*model.Role, error) {
	args := m.Called(name)
	return args.Get(0).(*model.Role), args.Error(1)
}

func (m *MockRoleRepository) List() ([]model.Role, error) {
	args := m.Called()
	return args.Get(0).([]model.Role), args.Error(1)
}

func (m *MockRoleRepository) AssignPermission(roleID, permissionID uint) error {
	args := m.Called(roleID, permissionID)
	return args.Error(0)
}


// Other methods of RoleRepository can be mocked similarly

func TestRBACService_CreateRole(t *testing.T) {
	mockRepo := new(MockRoleRepository)
	rbacService := service.NewRBACService(mockRepo, nil, nil)

	mockRole := &model.Role{Name: "Admin", Description: "Admin Role"}
	mockRepo.On("Create", mockRole).Return(nil)

	role, err := rbacService.CreateRole("Admin", "Admin Role")

	assert.NoError(t, err)
	assert.Equal(t, "Admin", role.Name)
	assert.Equal(t, "Admin Role", role.Description)
	mockRepo.AssertExpectations(t)
}
