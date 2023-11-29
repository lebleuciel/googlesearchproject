package users

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mock_database "github.com/lebleuciel/maani/pkg/database/mocks"
	"github.com/lebleuciel/maani/pkg/repository/user"
	userservice "github.com/lebleuciel/maani/pkg/services/user"
	"github.com/lebleuciel/maani/pkg/settings"
	"github.com/stretchr/testify/assert"
)

// initUsersModuleWithMockDB function tests creating a new ForwarderModule and mockDatabase and returns instance of both
func initUsersModuleWithMockDB(t *testing.T, authEnabled bool) (*Users, *mock_database.MockDatabase) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var st settings.Settings
	db := mock_database.NewMockDatabase(ctrl)
	userRepo, err := user.NewUserRepository(db)
	assert.Nil(t, err)
	userService, err := userservice.NewUserService(userRepo, st)
	assert.Nil(t, err)
	userMod, err := NewUserModule(userService, userRepo, false)
	assert.Nil(t, err)
	assert.NotNil(t, userMod)
	return userMod, db
}

func TestNewUsersModule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	t.Run("nil_user_service", func(t *testing.T) {
		_, err := NewUserModule(nil, nil, false)
		assert.NotNil(t, err)
		assert.Equal(t, ErrNilUserService, err)
	})
	t.Run("nil_user_repo", func(t *testing.T) {
		var st settings.Settings
		db := mock_database.NewMockDatabase(ctrl)
		userRepo, err := user.NewUserRepository(db)
		assert.Nil(t, err)
		userService, err := userservice.NewUserService(userRepo, st)
		assert.Nil(t, err)
		_, err = NewUserModule(userService, nil, false)
		assert.NotNil(t, err)
		assert.Equal(t, ErrNilUserRepo, err)
	})
	t.Run("valid", func(t *testing.T) {
		var st settings.Settings
		db := mock_database.NewMockDatabase(ctrl)
		userRepo, err := user.NewUserRepository(db)
		assert.Nil(t, err)
		userService, err := userservice.NewUserService(userRepo, st)
		assert.Nil(t, err)
		mod, err := NewUserModule(userService, userRepo, false)
		assert.Nil(t, err)
		assert.NotNil(t, mod)
	})
}

// TestForwarder_RegisterRoutes tests all routes functionalities
func TestUsers_RegisterRoutes(t *testing.T) {
	userMod, db := initUsersModuleWithMockDB(t, true)
	baseRecorder := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(baseRecorder)
	db.EXPECT().GetUserList().Return(nil, nil).AnyTimes()
	v1 := engine.Group("/api")
	userMod.RegisterRoutes(v1)

	t.Run("save_user", func(t *testing.T) {
		req := httptest.NewRequest("GET", "https://store.foo/api/user/list", nil)
		recorder := httptest.NewRecorder()

		engine.ServeHTTP(recorder, req)
		assert.NotEqual(t, http.StatusNotFound, recorder.Code)
		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}
