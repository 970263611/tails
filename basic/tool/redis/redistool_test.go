package redistool

import "testing"

func TestCreateRedisClient(t *testing.T) {
	addr := "127.0.0.1:6379"
	password := "123456"
	db := 0
	client := CreateRedisClient(addr, password, db)
	t.Logf("redis client created %v", client)
	vel := client.Get(client.Context, "key")
	t.Logf("redis find value %v", vel)
}
