package utils

import (
	"github.com/jinzhu/copier"
)

func ConvertDtoToEntity[T, D any](dto D, opts ...copier.Option) (*T, error) {
	entity := new(T)
	err := copier.CopyWithOption(entity, dto, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	})
	if err != nil {
		return nil, err
	}
	return entity, nil
}
