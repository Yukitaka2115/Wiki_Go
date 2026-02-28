package config

import "encoding/json"

func ConvertStructToJSONStr(data interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func UnmarshalJSONStr(jsonStr string, target interface{}) error {
	if jsonStr == "" {
		return nil
	}
	return json.Unmarshal([]byte(jsonStr), target)
}
