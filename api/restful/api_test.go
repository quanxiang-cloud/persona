package restful

//import (
//	"bytes"
//	"encoding/json"
//	"flag"
//	"fmt"
//	"git.internal.yunify.com/qxp/persona/internal/persona"
//	"git.internal.yunify.com/qxp/persona/pkg/config"
//	"git.internal.yunify.com/qxp/persona/pkg/utils"
//	"io/ioutil"
//	"net/http"
//	"testing"
//)
//
//var (
//	BaseURL    = "http://127.0.0.1"
//	httpClient http.Client
//	KeyID      string
//)
//
//func TestMain(m *testing.M) {
//	var (
//		configPath = flag.String("config", "../../configs/config.yml", "-config 配置文件地址")
//	)
//	flag.Parse()
//	err := config.Init(*configPath)
//	if err != nil {
//		panic(err)
//	}
//	m.Run()
//}
//
//// TestCreateDataSet 创建数据集测试
//func TestCreateDataSet(t *testing.T) {
//	path := "/api/v1/structor/dataset/m/create"
//	url := fmt.Sprintf("%s%s%s", BaseURL, config.Config.Port, path)
//	reqData := persona.CreateDataSetReq{
//		Name:    "persona_test",
//		Tag:     "tag_1111",
//		Type:    1,
//		Content: "_test_content_test_content",
//	}
//	buf, err := utils.Struct2Bytes(&reqData)
//	if err != nil {
//		t.Fatal(err)
//	}
//	resp, err := httpClient.Post(url, "application/json", buf)
//	if err != nil {
//		t.Fatal(err)
//	}
//	if resp.StatusCode >= 400 {
//		t.Fatal(fmt.Sprintf("error response code: %d", resp.StatusCode))
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	fmt.Printf("Response: %s", string(body))
//}
//
//
//// TestGetDataSetByID 根据ID获取数据
//func TestGetDataSetByID(t *testing.T) {
//	path := "/api/v1/structor/dataset/m/get"
//	url := fmt.Sprintf("%s%s%s", BaseURL, config.Config.Port, path)
//	reqData := persona.GetDataSetReq{
//		ID: "0d6dccff-980f-49e8-ba7b-7f95b1520d04",
//	}
//	buf, err := utils.Struct2Bytes(&reqData)
//	if err != nil {
//		t.Fatal(err)
//	}
//	resp, err := httpClient.Post(url, "application/json", buf)
//	if err != nil {
//		t.Fatal(err)
//	}
//	if resp.StatusCode >= 400 {
//		t.Fatal(fmt.Sprintf("error response code: %d", resp.StatusCode))
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	fmt.Printf("Response: %s", string(body))
//}
//
//// TestUpdateDataSet 更新数据集
//func TestUpdateDataSet(t *testing.T) {
//	path := "/api/v1/structor/dataset/m/get"
//	url := fmt.Sprintf("%s%s%s", BaseURL, config.Config.Port, path)
//	reqData := persona.GetDataSetReq{
//		ID: "0d6dccff-980f-49e8-ba7b-7f95b1520d04",
//	}
//	buf, err := utils.Struct2Bytes(&reqData)
//	if err != nil {
//		t.Fatal(err)
//	}
//	resp, err := httpClient.Post(url, "application/json", buf)
//	if err != nil {
//		t.Fatal(err)
//	}
//	if resp.StatusCode >= 400 {
//		t.Fatal(fmt.Sprintf("error response code: %d", resp.StatusCode))
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	fmt.Printf("Response: %s", string(body))
//}
//
//// TestGetDataSetByCondition 根据条件获取数据
//func TestGetDataSetByCondition(t *testing.T) {
//	path := "/api/v1/structor/dataset/m/get"
//	url := fmt.Sprintf("%s%s%s", BaseURL, config.Config.Port, path)
//	reqData := persona.GetDataSetReq{
//		ID: "0d6dccff-980f-49e8-ba7b-7f95b1520d04",
//	}
//	buf, err := utils.Struct2Bytes(&reqData)
//	if err != nil {
//		t.Fatal(err)
//	}
//	resp, err := httpClient.Post(url, "application/json", buf)
//	if err != nil {
//		t.Fatal(err)
//	}
//	if resp.StatusCode >= 400 {
//		t.Fatal(fmt.Sprintf("error response code: %d", resp.StatusCode))
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	fmt.Printf("Response: %s", string(body))
//}
//
//// TestGetDataSetByIDHome 根据key获取数据集
//func TestGetDataSetByIDHome(t *testing.T) {
//	path := "/api/v1/structor/dataset/m/create"
//	url := fmt.Sprintf("%s:%s%s", BaseURL, config.Config.Port, path)
//	reqData := persona.CreateDataSetReq{
//		Name:    "persona_test",
//		Tag: "tag_1111",
//		Type: 1,
//		Content: "_test_content_test_content",
//	}
//	var buf bytes.Buffer
//	err := json.NewEncoder(&buf).Encode(reqData)
//	if err != nil {
//		t.Fatal(err)
//	}
//	resp, err := httpClient.Post(url, "application/json", &buf)
//	if err != nil {
//		t.Fatal(err)
//	}
//	if resp.StatusCode >= 400 {
//		t.Fatal(fmt.Sprintf("error response code: %d", resp.StatusCode))
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Printf("Response: %s", string(body))
//}
//
//func TestMainOrder(t *testing.T) {
//	t.Run("TestCreateDataSet", TestCreateDataSet)
//}
