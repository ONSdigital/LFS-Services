package ws

import (
	"services/types"
	"sync"
)

type FileUploads struct {
	files map[string]*types.WSMessage
	mux   *sync.Mutex
}

var uploads FileUploads
var cm = &sync.Mutex{}

func NewFileUploads() FileUploads {
	cm.Lock()
	defer cm.Unlock()

	if uploads.files == nil {
		uploads.files = make(map[string]*types.WSMessage)
		uploads.mux = &sync.Mutex{}
	}
	return uploads
}

func (up *FileUploads) Status(fileName string) *types.WSMessage {
	m, ok := up.files[fileName]
	if !ok {
		return &types.WSMessage{
			Filename:     fileName,
			Percentage:   0,
			Status:       types.UploadError,
			ErrorMessage: "fileName not found",
		}
	}
	return m
}

func (up *FileUploads) Add(fileName string) *types.WSMessage {
	up.mux.Lock()
	defer up.mux.Unlock()

	m := types.NewWSMessage(fileName, persistStatus)
	up.files[fileName] = m
	return m

}

func persistStatus() {

}
