package server

import (
	"crypto/rand"
	"encoding/hex"
	"image"
	_ "image/jpeg" // Register the JPEG decoder
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go_pakito/models"
)

const (
	avatarDir     = "/app/uploads/avatars"
	maxAvatarSize = 2 << 20 // 2 MB
	maxAvatarDim  = 2048    // pixels per side
)

func UploadAvatar(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := UserID(c)

		// Cut off the reading if the body exceeds the limit
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxAvatarSize+(1<<10))
		fh, err := c.FormFile("avatar")
		if err != nil || fh.Size > maxAvatarSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "archivo no recibido o mayor de 2 MB"})
			return
		}
		f, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo leer el archivo"})
			return
		}
		defer f.Close()

		// Check the format and dimensions WITHOUT decoding the entire image
		cfg, _, err := image.DecodeConfig(f)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato no soportado (solo PNG o JPEG)"})
			return
		}
		if cfg.Width > maxAvatarDim || cfg.Height > maxAvatarDim {
			c.JSON(http.StatusBadRequest, gin.H{"error": "la imagen supera 2048x2048 píxeles"})
			return
		}
		f.Seek(0, 0)
		img, _, err := image.Decode(f)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "imagen corrupta"})
			return
		}

		// Random name: the name provided by the customer is never used
		b := make([]byte, 16)
		rand.Read(b)
		name := hex.EncodeToString(b) + ".png"

		if err := os.MkdirAll(avatarDir, 0o755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"})
			return
		}
		out, err := os.Create(filepath.Join(avatarDir, name))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"})
			return
		}
		// Saving as a PNG removes metadata (EXIF, GPS) and any hidden content
		encErr := png.Encode(out, img)
		out.Close()
		if encErr != nil {
			os.Remove(out.Name())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"})
			return
		}

		var u models.User
		db.Select("id", "avatar").First(&u, uid)
		if err := db.Model(&models.User{}).Where("id = ?", uid).Update("avatar", name).Error; err != nil {
			os.Remove(out.Name())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"})
			return
		}
		if u.Avatar != nil { // Delete the previous photo
			os.Remove(filepath.Join(avatarDir, filepath.Base(*u.Avatar)))
		}
		c.JSON(http.StatusOK, gin.H{"avatar": name})
	}
}

func DeleteAvatar(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := UserID(c)
		var u models.User
		if err := db.Select("id", "avatar").First(&u, uid).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
			return
		}
		if u.Avatar != nil {
			os.Remove(filepath.Join(avatarDir, filepath.Base(*u.Avatar)))
			db.Model(&models.User{}).Where("id = ?", uid).Update("avatar", nil)
		}
		c.Status(http.StatusNoContent)
	}
}

func GetAvatar(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		var u models.User
		if err := db.Select("id", "avatar").First(&u, id).Error; err != nil || u.Avatar == nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Header("Cache-Control", "private, max-age=3600")
		c.File(filepath.Join(avatarDir, filepath.Base(*u.Avatar)))
	}
}