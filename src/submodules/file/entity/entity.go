package file_entity

import "github.com/google/uuid"

type File struct {
	Name string
	Data []byte
}

func (f *File) GenerateName() error {
	f.Name = uuid.New().String()
	return nil
}
