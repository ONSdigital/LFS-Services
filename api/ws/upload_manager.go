package ws

import (
	"fmt"
	"services/types"
	"sync"
)

const (
	Idle = iota
	UploadStarted
	UploadFinished
	UploadError
	UploadCancelled
	NotFound
)

type FileUploads struct {
	files map[string]types.WSMessage
	mux   *sync.Mutex
}

var uploads FileUploads
var cm = &sync.Mutex{}

func NewFileUploads() FileUploads {
	cm.Lock()
	defer cm.Unlock()

	if uploads.files == nil {
		uploads.files = make(map[string]types.WSMessage)
		uploads.mux = &sync.Mutex{}
	}
	return uploads
}

func (up *FileUploads) Add(fileName string) *types.WSMessage {
	up.mux.Lock()
	defer up.mux.Unlock()

	m := types.WSMessage{
		Filename:   fileName,
		Percentage: 0,
		Status:     Idle,
	}
	up.files[fileName] = m
	return &m
}

func (up *FileUploads) SetPercentage(fileName string, percentage float64) error {
	up.mux.Lock()
	defer up.mux.Unlock()

	m, ok := up.files[fileName]
	if !ok {
		return fmt.Errorf(" %s not found in list of files", fileName)
	}
	m.Percentage = percentage
	up.files[fileName] = m
	return nil
}

func (up *FileUploads) SetUploadStarted(fileName string) error {
	up.mux.Lock()
	defer up.mux.Unlock()

	m, ok := up.files[fileName]
	if !ok {
		return fmt.Errorf(" %s not found in list of files", fileName)
	}
	m.Status = UploadStarted
	up.files[fileName] = m
	return nil
}

func (up *FileUploads) SetUploadFinished(fileName string) error {
	up.mux.Lock()
	defer up.mux.Unlock()

	m, ok := up.files[fileName]
	if !ok {
		return fmt.Errorf(" %s not found in list of files", fileName)
	}
	m.Status = UploadFinished
	m.Percentage = 100
	up.files[fileName] = m
	return nil
}

func (up *FileUploads) SetUploadError(fileName string) error {
	up.mux.Lock()
	defer up.mux.Unlock()

	m, ok := up.files[fileName]
	if !ok {
		return fmt.Errorf(" %s not found in list of files", fileName)
	}
	m.Status = UploadError
	up.files[fileName] = m
	return nil
}

func (up *FileUploads) SetUploadCancelled(fileName string) error {
	up.mux.Lock()
	defer up.mux.Unlock()

	m, ok := up.files[fileName]
	if !ok {
		return fmt.Errorf(" %s not found in list of files", fileName)
	}
	m.Status = UploadCancelled
	up.files[fileName] = m
	return nil
}

func (up FileUploads) Status(fileName string) *types.WSMessage {

	i, ok := up.files[fileName]
	if ok {
		return &i
	}
	return &types.WSMessage{
		Filename:   fileName + " Unknown",
		Percentage: 0,
		Status:     NotFound,
	}

}
