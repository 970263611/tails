package db_conn_num

import (
	"fmt"
	"testing"
)

func TestDoHandler(t *testing.T) {
	params := map[string]any{
		"host":     "127.0.0.1",
		"port":     15431,
		"username": "postgres",
		"password": "postgres",
		"dbname":   "postgres",
	}
	resp := doHandler(params)
	fmt.Println(string(resp))
}
