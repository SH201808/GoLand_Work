package main

import (
	"LRU/client"
	"fmt"
	"reflect"
	"testing"
)

func TestClient(t *testing.T) {
	c := client.NewClient("127.0.0.1:8080")
	//set 如果没有就新建,有就更新
	c.Set("key1", "value1")
	c.Set("key2", "value2")
	c.Set("key3", "value3")
	//Get成功 value2
	want := "value2"
	key2, _ := c.Get("key2")
	if !reflect.DeepEqual(want, key2) {
		t.Errorf("expected:%v, got:%v", want, key2)
	}
	fmt.Printf("key2: %v\n", key2)
	//过期并返回错误
	key1, _ := c.Get("key1")
	fmt.Printf("key1: %v\n", key1)
	//删除 key2
	c.Delete("key2")
}