package util

import (
"bytes"
"encoding/json"
"fmt"
)

func GetObjFormatStr(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return fmt.Sprintf("%+v", obj)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", obj)
	}
	return out.String()
}
