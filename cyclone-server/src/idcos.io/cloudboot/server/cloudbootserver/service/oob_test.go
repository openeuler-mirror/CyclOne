package service

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestOObConfig(t *testing.T) {
	var mapper map[string][]*OOBUser

	err := json.Unmarshal([]byte(OOBVendorConfig), &mapper)
	if err != nil {
		fmt.Println(err.Error())
	}

	for k, v := range mapper {
		fmt.Println(fmt.Sprintf("%s: %s", k, v))
	}

	v := mapper["default"]
	fmt.Println(v)
}
