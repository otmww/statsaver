package manager

import (
	"encoding/json"
	"github.com/otmww/base/logging"
	"github.com/otmww/base/logmanagement"
	"io"
	"log"
	"os"
	"time"
)

var ManagerInstance *Manager

var sep = []byte("\n")

type Manager struct {
	writer io.WriteCloser
}

func NewManager(workDir, serverName, filePrefix, test string, rotationPeriodSec int64, maxBackups int) {
	if err := os.MkdirAll(workDir, os.ModePerm); err != nil {
		log.Fatal("Fail to create work dir ", err)
	}
	writer, err := logging.MakeLJWriter(
		serverName,
		workDir,
		filePrefix,
		rotationPeriodSec,
		test,
		maxBackups,
	)
	if err != nil {
		log.Fatal("Fail to create lj writer ", err)
	}
	ManagerInstance = &Manager{
		writer: writer,
	}
}

func (m *Manager) Save(data []byte) error {
	data = append(data, sep...)
	if _, err := m.writer.Write(data); err != nil {
		return err
	}
	return nil

}

func (m *Manager) StartUploadJob(upload bool, accessKey, secretKey, region, endpoint, bucket, project, fileMask string, maxAtOnce, delaySeconds int) {
	if upload {
		return
	}
	uj := &logmanagement.UploadJob{}
	uj.Gzip = true
	uj.AwsCreds.AccessKeyId = accessKey
	uj.AwsCreds.SecretKey = secretKey
	uj.AWSRegion = region
	uj.AWSEndpont = endpoint
	uj.ForcePathStyle = true //cuz of beeline certificate type
	uj.Bucket = bucket
	uj.Project = project
	uj.Mask = fileMask
	uj.MaxAtOnce = maxAtOnce
	uj.DelaySeconds = delaySeconds
	b, _ := json.MarshalIndent(&uj, " ", " ")
	log.Println(string(b))
	go func() {
		for {
			log.Println("Starting upload")
			err := uj.UploadOnce()
			if err != nil {
				log.Println("An error occured while uploading files", err)
			}
			if uj.DelaySeconds <= 0 {
				time.Sleep(time.Minute)
			} else {
				time.Sleep(time.Second * time.Duration(uj.DelaySeconds))
			}

		}
	}()

}
