package uam

import "fmt"

func wrapErr(err error) error {
	if err == nil {
		return nil
	}
	return wrapErrString(err.Error())
}

func wrapErrString(err string) error {
	if err == "" {
		return nil
	}
	return fmt.Errorf("[UAM]%s", err)
}
