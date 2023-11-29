package forwarder

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mock_database "github.com/lebleuciel/maani/pkg/database/mocks"
	"github.com/lebleuciel/maani/pkg/repository/user"
	"github.com/lebleuciel/maani/pkg/services/auth"
	"github.com/stretchr/testify/assert"
)

// initForwarderModuleWithMockDB function tests creating a new ForwarderModule and mockDatabase and returns instance of both
func initForwarderModuleWithMockDB(t *testing.T, authEnabled bool) (*Forwarder, *mock_database.MockDatabase) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := mock_database.NewMockDatabase(ctrl)
	userRepo, err := user.NewUserRepository(db)
	assert.Nil(t, err)
	authMod, err := auth.NewAuth(userRepo, "secret", "email", "panel", 50*time.Hour, 50*time.Hour)
	assert.Nil(t, err)
	userMod, err := NewForwarderModule(authMod, "http://store", 9000, 9001, "X-User", authEnabled)
	assert.Nil(t, err)
	assert.NotNil(t, userMod)
	return userMod, db
}

func TestNewUsersModule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	t.Run("empty_store_host", func(t *testing.T) {
		_, err := NewForwarderModule(nil, "", 0, 0, "", false)
		assert.NotNil(t, err)
		assert.Equal(t, ErrEmptyStoreHost, err)
	})
	t.Run("empty_admin_port", func(t *testing.T) {
		_, err := NewForwarderModule(nil, "http://store", 0, 0, "", false)
		assert.NotNil(t, err)
		assert.Equal(t, ErrEmptyAdminPort, err)
	})
	t.Run("empty_backend_port", func(t *testing.T) {
		_, err := NewForwarderModule(nil, "http://store", 9000, 0, "", false)
		assert.NotNil(t, err)
		assert.Equal(t, ErrEmptyBackendPort, err)
	})
	t.Run("empty_backend_port", func(t *testing.T) {
		_, err := NewForwarderModule(nil, "http://store", 9000, 9000, "", false)
		assert.NotNil(t, err)
		assert.Equal(t, ErrEmptyUserHeaderKey, err)
	})
	t.Run("valid", func(t *testing.T) {
		mod, err := NewForwarderModule(nil, "http://store", 9000, 9000, "X-UserKey", false)
		assert.Nil(t, err)
		assert.NotNil(t, mod)
	})
}

// TestForwarder_RegisterRoutes tests all routes functionalities
func TestForwarder_RegisterRoutes(t *testing.T) {
	forwarderMod, _ := initForwarderModuleWithMockDB(t, true)
	baseRecorder := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(baseRecorder)

	v1 := engine.Group("/api")
	forwarderMod.RegisterRoutes(v1)

	t.Run("list_all_files", func(t *testing.T) {
		req := httptest.NewRequest("GET", "https://store.foo/api/file/list", nil)
		recorder := httptest.NewRecorder()

		engine.ServeHTTP(recorder, req)
		assert.NotEqual(t, http.StatusNotFound, recorder.Code)
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("list_all_users", func(t *testing.T) {
		req := httptest.NewRequest("GET", "https://store.foo/api/user/list", nil)
		recorder := httptest.NewRecorder()

		engine.ServeHTTP(recorder, req)
		assert.NotEqual(t, http.StatusNotFound, recorder.Code)
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})
}
