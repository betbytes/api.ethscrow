package utils

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ParseRequestBody(r *http.Request, reqType interface{}) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	if err = json.Unmarshal(reqBody, reqType); err != nil {
		return err
	}

	return nil
}

func DecodeHexString(data string) []byte {
	arr, err := hex.DecodeString(data)
	if err != nil {
		return nil
	}
	return arr
}
