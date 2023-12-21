package file

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/antchfx/htmlquery"
	"github.com/gin-gonic/gin"
	"github.com/lebleuciel/maani/models"
	"github.com/lebleuciel/maani/pkg/database"
	"github.com/lebleuciel/maani/pkg/helpers"
	repository "github.com/lebleuciel/maani/pkg/repository/file"
	"github.com/lebleuciel/maani/pkg/settings"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

// logger is a global variable for logging using Zap.
var logger *zap.SugaredLogger

// init initializes the Zap logger.
func init() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	logger = zapLogger.Sugar()
}

type FileService struct {
	st         settings.Settings
	repository *repository.FileRepository
}

func NewFileService(repo *repository.FileRepository, st settings.Settings) (*FileService, error) {
	if repo == nil {
		return nil, ErrNilFileRepo
	}
	return &FileService{
		st:         st,
		repository: repo,
	}, nil
}

func (f *FileService) SearchGoogle(c *gin.Context, isAdmin bool) ([]image.Image, error) {
	searchQuery := c.Query("q")
	maxImagesStr := c.Query("maxnum")
	maxImages, err := strconv.Atoi(maxImagesStr)
	if err != nil {
		log.Println("Error converting maxImages to int:", err)
	}
	url := fmt.Sprintf("https://www.google.com/search?q=%s&tbm=isch", searchQuery)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var images []image.Image
	count := 0

	for _, imgNode := range htmlquery.Find(doc, "//img") {
		if count >= maxImages {
			break
		}

		imgURL := htmlquery.SelectAttr(imgNode, "src")
		img, err := downloadImage(imgURL)
		if err != nil {
			log.Println("Error downloading image:", err)
			continue
		}

		images = append(images, img)
		count++
	}

	return images, nil
}

func downloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (f *FileService) GetFile(c *gin.Context, isAdmin bool) {
	tags := helpers.SplitBySpaceComma(c.PostFormArray("tags"))
	name := helpers.SplitBySpaceComma(c.PostFormArray("name"))

	tx, file, err := f.repository.GetEncryptedFile(name, tags)
	defer func() {
		if err != nil && tx != nil {
			if e, ok := err.(*pq.Error); !ok || e.Code != database.ErrSerializationFailure {
				rollbackErr := tx.Rollback()
				if rollbackErr != nil {
					err = fmt.Errorf("error while rolling back transaction. original error: %w", err)
				}
			}
		}
	}()
	if err != nil {
		if err.Error() == "file not found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File not found"})
		} else {
			logger.Errorw("failed to get encrypted file", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get encrypted file"})
		}
		return
	}

	// Define the path of the file to be retrieved
	filePath := fmt.Sprintf("%s/%s", f.st.BackendServer.FilePath, file.UUID)

	err = helpers.DecryptFile(filePath, []byte(f.st.BackendServer.EncryptKey))
	if err != nil {
		logger.Errorw("failed to decryptFile file", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decryptFile file"})

		if os.IsNotExist(err) {
			err = nil
			tx.Commit()
		}
		return
	}

	// Set the headers for the file transfer and return the file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name))
	c.Header("Content-Type", file.TypeId)
	c.Header("Content-Length", fmt.Sprintf("%d", file.Size))
	c.FileAttachment(filePath, file.Name)

	// Delete the file after sending it to the client
	err = os.Remove(filePath)
	if err != nil {
		logger.Errorw("failed to delete file", "error", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Get %s.%s successfully", file.Name, file.TypeId),
	})
	tx.Commit()
}

func (f *FileService) SaveFiles(c *gin.Context, isAdmin bool) {
	form, err := c.MultipartForm()
	if err != nil {
		logger.Errorw("failed to parse form", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse form"})
		return
	}
	if len(form.File["files"]) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field files is empty"})
		return
	}

	for _, file := range form.File["files"] {
		err := f.repository.IsValidFile(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	userId, err := strconv.Atoi(c.GetHeader(f.st.GatewayServer.UserIdHeaderKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can not parse user id from header"})
		return
	}

	tags := helpers.SplitBySpaceComma(c.PostFormArray("tags"))

	var message []string
	var error_message []string
	for _, file := range form.File["files"] {
		err := f.repository.SaveEncryptedFile(models.File{
			Name:    file.Filename,
			Size:    int(file.Size),
			TypeId:  file.Header.Get("Content-Type"),
			UserId:  userId,
			Content: file,
			Tags:    tags,
		})
		if err != nil {
			error_message = append(error_message, fmt.Sprintf("Can not save file, error: %s", err.Error()))
		} else {
			message = append(message, fmt.Sprintf("File %s saved successfully", file.Filename))
		}
	}

	if len(error_message) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": error_message})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func (f *FileService) GetFileList(c *gin.Context, isAdmin bool) {
	files, err := f.repository.GetFileList()
	if err != nil {
		logger.Errorw("failed to get file list", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can not get files list"})
		return
	}
	c.JSON(http.StatusOK, files)
}
