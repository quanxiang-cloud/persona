package elasticsearch

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"git.internal.yunify.com/qxp/persona/pkg/config"
	"git.internal.yunify.com/qxp/persona/pkg/db"
	"testing"
	"time"
)

var (
	TestEsAPI   db.BackendStorage
	TestKey     = "Test_key_persona"
	TestValue   = "Test_value_persona"
	TestUserID  = "Test_userID_persona"
	TestVersion = "v1.0.1_persona"
	// 完成测试后需要清理的key
	CleanupKeys = make([]string, 0)
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var (
		configPath = flag.String("config", "../../../configs/config.yml", "-config 配置文件地址")
	)
	flag.Parse()
	err := config.Init(*configPath)
	if err != nil {
		panic(err)
	}
	client, err := NewEsClient(&config.Config.ES)
	if err != nil {
		panic(err)
	}
	TestEsAPI = &Elasticsearch{client: client, esConfig: &config.Config.ES}
	m.Run()
	// 清理测试数据
	fmt.Printf("Delete keys: %s", CleanupKeys)
	for _, i := range CleanupKeys {
		_ = TestEsAPI.DeleteData(&ctx, &i)
	}
	fmt.Println("Data cleanup success")
}

func TestPutValue(t *testing.T) {
	ctx := context.Background()
	fmt.Printf("Put key: [%s] value: [%s]\n", TestKey, TestUserID)
	err := TestEsAPI.Put(ctx, TestKey, TestUserID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetValue(t *testing.T) {
	ctx := context.Background()
	res, err := TestEsAPI.Get(ctx, TestKey)
	fmt.Printf("Got value: %v\n", res)
	if err != nil {
		t.Fatal(err)
	}
	if TestKey != res["key"] {
		t.Fatal(fmt.Sprintf("%s != %s", TestKey, res["key"]))
	}
}

func TestPutWithVersion(t *testing.T) {
	ctx := context.Background()
	key := fmt.Sprintf("%s_%s", TestKey, TestVersion)
	fmt.Printf("Put Version: [%s] Key: [%s] Value: [%s]", TestVersion, key, TestValue)
	err := TestEsAPI.PutWithVersion(ctx, TestVersion, key, TestValue)
	if err != nil {
		t.Fatal(err)
	}
	CleanupKeys = append(CleanupKeys, key)
}

func TestGetWithVersion(t *testing.T) {
	ctx := context.Background()
	key := fmt.Sprintf("%s_%s", TestKey, TestVersion)
	res, err := TestEsAPI.GetWithVersion(ctx, TestVersion, key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Got value: %v\n", res)
}

// 设置用户带版本的值
func TestUserPutWithVersion(t *testing.T) {
	ctx := context.Background()
	key := fmt.Sprintf("%s_%s_%s", TestKey, TestVersion, TestUserID)
	err := TestEsAPI.UserPutWithVersion(ctx, TestVersion, key, TestValue)
	if err != nil {
		t.Fatal(err)
	}
	CleanupKeys = append(CleanupKeys, key)
}

func TestUserGetWithVersion(t *testing.T) {
	ctx := context.Background()
	key := fmt.Sprintf("%s_%s_%s", TestKey, TestVersion, TestUserID)
	res, err := TestEsAPI.UserGetWithVersion(ctx, TestVersion, key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Got value: %v\n", res)
}

// 根据key更新数据
func TestUpdateData(t *testing.T) {
	ctx := context.Background()

	updateValue := map[string]interface{}{
		"value":   fmt.Sprintf("test_update_by_%s", TestKey),
		"version": "v1.0.1_test",
	}
	err := TestEsAPI.UpdateData(&ctx, &TestKey, updateValue)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 5)
	// 更新是否成功？
	data, err := TestEsAPI.GetData(&ctx, &TestKey)
	fmt.Printf("Got TestUpdateData data: %s", string(*data))
	if err != nil {
		t.Fatal(err)
	}
	if data == nil {
		t.Fatal("not found data")
	}
	var M = make(map[string]interface{})
	err = json.Unmarshal(*data, &M)
	if err != nil {
		t.Fatal(err)
	}
	if updateValue["value"] != M["value"] {
		t.Fatal("update fail")
	}
}

// 根据条件获取数据
func TestGetDataByKVs(t *testing.T) {
	ctx := context.Background()
	KVCondition := map[string]interface{}{
		"key": TestKey,
	}
	data, err := TestEsAPI.GetDataByKVs(&ctx, &KVCondition)
	fmt.Printf("Got TestGetDataByKVs data: %v", data)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 0 {
		t.Fatal("not fond data")
	}
}

// 根据key删除数据
func TestDeleteData(t *testing.T) {
	ctx := context.Background()
	_ = TestEsAPI.DeleteData(&ctx, &TestKey)
	time.Sleep(time.Second * 5)
	// 数据是否已删除
	data, _ := TestEsAPI.GetData(&ctx, &TestKey)
	if data != nil {
		t.Fatal("delete data fail")
	}
	fmt.Printf("Got TestDeleteData data: %v", data)
}

func TestMainOrder(t *testing.T) {
	t.Run("TestPutValue", TestPutValue)
	t.Run("TestGetValue", TestGetValue)
	t.Run("TestPutWithVersion", TestPutWithVersion)
	t.Run("TestGetWithVersion", TestGetWithVersion)
	t.Run("TestUserPutWithVersion", TestUserPutWithVersion)
	t.Run("TestUserGetWithVersion", TestUserGetWithVersion)
	t.Run("TestUpdateData", TestUpdateData)
	t.Run("TestGetDataByKVs", TestGetDataByKVs)
	t.Run("TestDeleteData", TestDeleteData)
}
