package redistool

import "testing"

func TestCreateRedisClient(t *testing.T) {
	addr := "192.168.56.128:17001"
	password := "123456"
	db := 0
	client, err := CreateRedisClient(addr, password, db)
	if err != nil {
		t.Error(err)
	}
	t.Logf("redis client created %v", client)
	vel := client.Get(client.Context, "key")
	t.Logf("redis find value %v", vel)
}
