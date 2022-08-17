//go:build integration

package aliyun

import (
	"net/http"
	"path/filepath"
	"reflect"
	"testing"

	config2 "github.com/eleven26/goss/config"
	"github.com/eleven26/goss/core"
	"github.com/eleven26/goss/utils"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	fs "github.com/eleven26/go-filesystem"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	storage core.Storage

	bucket *oss.Bucket

	key          = "test/foo.txt"
	testdata     string
	fooPath      string
	localFooPath string
)

func init() {
	config2.ReadInUserHomeConfig()

	d := NewDriver()
	storage = d.Storage()
	bucket = storage.(*Storage).store.Bucket

	testdata = filepath.Join(utils.RootDir(), "testdata")
	fooPath = filepath.Join(testdata, "foo.txt")
	localFooPath = filepath.Join(testdata, "foo1.txt")

	viper.Set("show_progress_bar", false)
}

func setUp() {
	err := bucket.PutObjectFromFile(key, fooPath)
	if err != nil {
		panic(err)
	}
}

func tearDown() {
	deleteLocal()
	deleteRemote()
}

func deleteRemote() {
	err := bucket.DeleteObject(key)
	if err != nil {
		panic(err)
	}
}

func deleteLocal() {
	exists, _ := fs.Exists(localFooPath)
	if exists {
		_ = fs.Delete(localFooPath)
	}
}

func TestPut(t *testing.T) {
	defer tearDown()

	err := storage.Put(key, fooPath)
	assert.Nil(t, err)

	exists, err := bucket.IsObjectExist(key)
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestGet(t *testing.T) {
	setUp()
	defer tearDown()

	content, err := storage.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, content, "foo")

	content, err = storage.Get(key + "not_exists")
	assert.Empty(t, content)
	assert.Equal(t, http.StatusNotFound, err.(oss.ServiceError).StatusCode)
}

func TestSave(t *testing.T) {
	setUp()
	defer tearDown()

	err := storage.Save(key, localFooPath)
	assert.Nil(t, err)
	assert.Equal(t, "foo", fs.MustGetString(localFooPath))
}

func TestSize(t *testing.T) {
	setUp()
	defer tearDown()

	size, err := storage.Size(key)

	var expectedSize int64 = 3
	assert.Nil(t, err)
	assert.Equal(t, expectedSize, size)

	var s int64 = 0
	size, err = storage.Size(key + "not_exists")
	assert.Equal(t, s, size)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.(oss.ServiceError).StatusCode)
}

func TestDelete(t *testing.T) {
	setUp()

	err := storage.Delete(key)
	assert.Nil(t, err)

	exists, err := bucket.IsObjectExist(key)
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestExists(t *testing.T) {
	setUp()
	defer tearDown()

	exists, err := storage.Exists(key)
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestFiles(t *testing.T) {
	setUp()
	defer tearDown()

	files, err := storage.Files("test/")
	assert.Nil(t, err)
	assert.Len(t, files, 1)

	var expectedSize int64 = 3
	assert.Equal(t, key, files[0].Key())
	assert.Equal(t, expectedSize, files[0].Size())
}

func TestStorage(t *testing.T) {
	s := storage.Storage()

	assert.Equal(t, "aliyun.Storage", reflect.TypeOf(s).Elem().String())
}