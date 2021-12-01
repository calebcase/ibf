package cmd

import (
	"encoding/json"
	"os"

	ibf "github.com/calebcase/ibf/lib"
	"github.com/zeebo/errs"
)

func create(path string, set *ibf.IBF) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		err = errs.Combine(err, file.Close())
	}()

	return json.NewEncoder(file).Encode(set)
}

func open(path string) (set *ibf.IBF, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errs.Combine(err, file.Close())
	}()

	set = &ibf.IBF{}

	return set, json.NewDecoder(file).Decode(set)
}
