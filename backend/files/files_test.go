package files

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lebleuciel/maani/models"
	mock_database "github.com/lebleuciel/maani/pkg/database/mocks"
	"github.com/lebleuciel/maani/pkg/repository/file"
	fileservice "github.com/lebleuciel/maani/pkg/services/file"
	"github.com/lebleuciel/maani/pkg/settings"
	"github.com/stretchr/testify/assert"
)

// initFilesModuleWithMockDB function tests creating a new ForwarderModule and mockDatabase and returns instance of both
func initFilesModuleWithMockDB(t *testing.T, authEnabled bool) (*Files, *mock_database.MockDatabase) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var st settings.Settings
	st.BackendServer.FilePath = "\tmp"
	db := mock_database.NewMockDatabase(ctrl)
	fileRepo, err := file.NewFileRepository(st, db)
	assert.Nil(t, err)
	fileService, err := fileservice.NewFileService(fileRepo, st, db)
	assert.Nil(t, err)
	fileMod, err := NewFileModule(fileService, fileRepo, false)
	assert.Nil(t, err)
	assert.NotNil(t, fileMod)
	return fileMod, db
}

// initFileModuleWithMockTransactoin function tests creating a new UserModule and mockTransaction and returns instance of both
func initFileModuleWithMockTransactoin(t *testing.T) *mock_database.MockTransaction {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tx := mock_database.NewMockTransaction(ctrl)
	return tx
}

func TestNewFilesModule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	t.Run("nil_file_service", func(t *testing.T) {
		_, err := NewFileModule(nil, nil, false)
		assert.NotNil(t, err)
		assert.Equal(t, ErrNilFileService, err)
	})
	t.Run("nil_file_repo", func(t *testing.T) {
		var st settings.Settings
		db := mock_database.NewMockDatabase(ctrl)
		fileRepo, err := file.NewFileRepository(st, db)
		assert.Nil(t, err)
		fileService, err := fileservice.NewFileService(fileRepo, st, db)
		assert.Nil(t, err)
		_, err = NewFileModule(fileService, nil, false)
		assert.NotNil(t, err)
		assert.Equal(t, ErrNilFileRepo, err)
	})
	t.Run("valid", func(t *testing.T) {
		var st settings.Settings
		db := mock_database.NewMockDatabase(ctrl)
		fileRepo, err := file.NewFileRepository(st, db)
		assert.Nil(t, err)
		fileService, err := fileservice.NewFileService(fileRepo, st, db)
		assert.Nil(t, err)
		mod, err := NewFileModule(fileService, fileRepo, false)
		assert.Nil(t, err)
		assert.NotNil(t, mod)
	})
}

// TestForwarder_RegisterRoutes tests all routes functionalities
func TestFiles_RegisterRoutes(t *testing.T) {
	fileMod, db := initFilesModuleWithMockDB(t, true)
	baseRecorder := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(baseRecorder)

	tx := initFileModuleWithMockTransactoin(t)
	db.EXPECT().NewSerializableTransaction(gomock.Any()).Return(tx, nil).AnyTimes()
	tx.EXPECT().GetFile(gomock.Any(), gomock.Any()).Return(models.File{}, nil).AnyTimes()
	tx.EXPECT().Commit().Return(nil).AnyTimes()

	v1 := engine.Group("/api")
	fileMod.RegisterRoutes(v1)

	t.Run("get_file", func(t *testing.T) {
		req := httptest.NewRequest("GET", "https://store.foo/api/file", nil)
		recorder := httptest.NewRecorder()

		engine.ServeHTTP(recorder, req)
		assert.NotEqual(t, http.StatusNotFound, recorder.Code)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
	t.Run("save_file", func(t *testing.T) {
		req := httptest.NewRequest("POST", "https://store.foo/api/file", nil)
		recorder := httptest.NewRecorder()

		engine.ServeHTTP(recorder, req)
		assert.NotEqual(t, http.StatusNotFound, recorder.Code)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}
