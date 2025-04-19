package utils

import (
	"bytes"
	"encoding/json"
)

func PrettyString(strBytes []byte) ([]byte, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, strBytes, " ", "    "); err != nil {
		return nil, err
	}
	return prettyJSON.Bytes(), nil
}
