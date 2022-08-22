package utils

import (
	"strconv"
)

//生成分布式ID
func DistributedID() (string, error) {
	id, err := sonyFlake.NextID()
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(id)), nil
}
