package file

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/kabukky/journey/conversion"
	"github.com/kabukky/journey/dao/scheme"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var ErrContentPathReadError = errors.New("content path read error")

var DefaultFs afero.Fs

func InitFS() (err error) {
	DefaultFs, err = LocalFs()
	return
}

func GetFs() afero.Fs {
	return DefaultFs
}

func LocalFs() (fs afero.Fs, err error) {
	var stat os.FileInfo
	contentPath, _ := os.Getwd()
	contentPath = filepath.Join(contentPath, "content")
	if stat, err = os.Stat(contentPath); err != nil || !stat.IsDir() {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		contentPath = filepath.Join(dir, "content")
		if stat, err = os.Stat(contentPath); err != nil || !stat.IsDir() {
			err = ErrContentPathReadError
			return
		}
	}

	fs = afero.NewBasePathFs(afero.NewOsFs(), contentPath)
	return
}

func Put(db *gorm.DB, fileName string, content []byte, user *scheme.User) (file *scheme.File, err error) {
	fileHash := md5.Sum(content)
	fileName = conversion.XssFilter(fileName)
	fileObj := &scheme.File{}
	fileObj.Name = filepath.Base(fileName)
	fileObj.Path = filepath.Dir(fileName)
	fileObj.Hash = hex.EncodeToString(fileHash[0:])
	fileObj.Size = int64(len(content))
	fileObj.MIME = http.DetectContentType(content)
	fileObj.AbsolutePath = filepath.Join(fileObj.Path, fileObj.Name)
	if user != nil {
		fileObj.CreatedBy = &user.ID
	}

	f, err := GetFs().OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer f.Close()

	if err != nil {
		return
	}
	_, err = f.Write(content)
	if err != nil {
		return
	}

	if err = db.Create(&fileObj).Error; err != nil {
		return
	}
	file = fileObj
	return
}

func Delete(db *gorm.DB, fileName string) (err error) {
	var fileData *scheme.File
	if fileData, err = Info(db, fileName); err != nil {
		return
	}
	err = GetFs().Remove(fileData.AbsolutePath)
	return
}
func Info(db *gorm.DB, fileName string) (fileData *scheme.File, err error) {
	err = db.Model(&scheme.File{}).Where("absolute_path = ?", fileName).First(&fileData).Error
	return

}
func Read(db *gorm.DB, fileName string) (file *scheme.File, content []byte, err error) {
	var fileData *scheme.File
	if fileData, err = Info(db, fileName); err != nil {
		return
	}
	content = make([]byte, 0)
	f, err := GetFs().OpenFile(fileData.AbsolutePath, os.O_RDONLY, os.ModePerm)
	defer f.Close()

	if err != nil {
		return
	}

	content, err = ioutil.ReadAll(f)
	return
}
