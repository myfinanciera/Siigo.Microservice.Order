package mapper

import "github.com/mashingan/smapping"

func Adap[T any, T2 any](target *T, dataSource T2) error {
	dataMapped := smapping.MapFields(dataSource)
	err := smapping.FillStruct(target, dataMapped)
	return err
}
