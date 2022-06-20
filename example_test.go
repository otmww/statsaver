package manager

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	fileDir := os.TempDir() + "stat_writer_test"
	log.Println(fileDir)
	serverName := "local"
	filePrefix := "test_file"
	test := ""
	rotationSec := 1
	maxBackups := 0

	m, err := NewManager(fileDir, serverName, filePrefix, test, int64(rotationSec), maxBackups)
	for i := 0; i < 10; i++ {
		exampleData := append([]byte("example row "), []byte(fmt.Sprintf("%d", i))...)
		err = m.Save(exampleData)
		assert.NoError(t, err)
	}
	workDone := make(chan struct{})
	time.AfterFunc(time.Second*2, func() {
		files, err := filepath.Glob(fileDir + "/*20??-*.log")
		assert.NoError(t, err)
		log.Println("Files in folder ", fileDir)
		for _, file := range files {
			log.Println(file)
			content, err := ioutil.ReadFile(file)
			assert.NoError(t, err)
			log.Println("file content:", string(content))

		}
		workDone <- struct{}{}
	})
	<-workDone
	err = os.RemoveAll(fileDir + "/")
	assert.NoError(t, err)
}
