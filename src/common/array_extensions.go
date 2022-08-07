package common

import (
	"strconv"
	"strings"
)

func ToEnumerableIds(values string) ([]int, error) {
	param := strings.Split(values, ",")
	ids := make([]int, 0)
	duplicate := make(map[int]bool)
	for _, id := range param {
		if !IsEmptyString(id) {
			value, err := strconv.Atoi(id)
			if err != nil {
				return nil, err
			}
			if !duplicate[value] {
				ids = append(ids, value)
				duplicate[value] = true
			}
		}
	}
	return ids, nil
}

func IsEmptyString(value string) bool {
	return value == ""
}
