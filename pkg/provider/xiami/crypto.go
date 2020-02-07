package xiami

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/winterssy/ghttp"
)

const (
	APPKey = "23649156"
)

var (
	reqHeader = map[string]interface{}{
		"appId":      200,
		"platformId": "h5",
	}
)

func signPayload(token string, model interface{}) ghttp.Params {
	payload := map[string]interface{}{
		"header": reqHeader,
		"model":  model,
	}
	requestBytes, _ := json.Marshal(payload)
	data := map[string]string{
		"requestStr": string(requestBytes),
	}
	dataBytes, _ := json.Marshal(data)
	dataStr := string(dataBytes)

	t := time.Now().UnixNano() / (1e6)
	signStr := fmt.Sprintf("%s&%d&%s&%s", token, t, APPKey, dataStr)
	sign := fmt.Sprintf("%x", md5.Sum([]byte(signStr)))

	return ghttp.Params{
		"t":    t,
		"sign": sign,
		"data": dataStr,
	}
}
