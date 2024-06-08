package dbmodel

import (
	"github.com/prysmaticlabs/prysm/v4/testing/assert"
	"testing"
)

func init() {
	DbInit(DbConfig{
		Host:   "127.0.0.1",
		Port:   3306,
		User:   "root",
		Passwd: "12345678",
		DbName: "eth",
	})
}

func TestGetRewardListByValidatorIndex(t *testing.T) {
	list := GetRewardListByValidatorIndex(0)
	assert.Equal(t, len(list), 0)
}
