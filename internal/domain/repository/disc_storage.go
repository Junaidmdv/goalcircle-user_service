package repository

import (
	"context"
	"io"
)

type FileMetadata struct {
	Filename    string
	ContentType string
	Size        int64
} 

type FilePath string 

const (
	Logo="/logo"  
	Avatar="/avatar" 
	Docs="/docs"
)

type FileStorage interface {
	UploadFile(context.Context, FilePath, io.Reader, *FileMetadata) (string, error)
	DeleteFile(context.Context, string) error
	GetURL(context.Context, string) (string, error)
}



