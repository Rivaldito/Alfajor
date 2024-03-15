package alfajor

import (
	"os"
)

func checIfExist(dir string) error {

	var err error

	if _, err = os.Stat(dir); err == nil {
		//The dir/file exist
		return nil
	}
	//The dir/file doesn't exist
	return err
}
