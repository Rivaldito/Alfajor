package alfajor

import (
	"reflect"
)

type AlfajorParams struct {
	RootDir         string `default:"Alfajor"`
	PrettyAlfajor   bool
	SeparateAlfajor bool
}

type Alfajor struct {
	AlfajorParams AlfajorParams
}

func Newalfajor(AlfajorParams AlfajorParams) Alfajor {

	typ := reflect.TypeOf(AlfajorParams)

	if AlfajorParams.RootDir == "" {
		f, _ := typ.FieldByName(rootDir)
		AlfajorParams.RootDir = f.Tag.Get("default")
	}

	Alfajor := Alfajor{AlfajorParams: AlfajorParams}
	return Alfajor
}

func (Alfajor Alfajor) StartAlfajorManager() {

	err := checIfExist(Alfajor.AlfajorParams.RootDir)
	if err == nil {
		return
	}
	Alfajor.mkdir(Alfajor.AlfajorParams.RootDir)
}
