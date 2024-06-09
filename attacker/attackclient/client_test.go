package attackclient

import (
	"context"
	"testing"
)

func TestClient_Echo(t *testing.T) {
	client, err := Dial("http://127.0.0.1:10000")
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Close()

	res, err := client.Echo(context.Background(), "hello")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}
