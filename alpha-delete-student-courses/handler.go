package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/cast"
)

const (
	apiKey = "P-CgtoLQxIfoXuz081FuZCenSJbUSMCjOf"

	getListURL   = "https://api.admin.u-code.io/v1/object/get-list/"
	getSingleURL = "https://api.admin.u-code.io/v1/object/"
	//multipleUpdateUrl    = "https://api.admin.u-code.io/v1/object/multiple-update/"
	// getListObjectBuilder = "https://api.admin.u-code.io/v1/object/get-list/"
)

// Handle a serverless request
func Handle(req []byte) string {

	var reqBody RequestBody

	if err := json.Unmarshal(req, &reqBody); err != nil {
		return Handler("error", err.Error())
	}

	data := cast.ToStringMap(reqBody.Data["object_data"])

	if cast.ToStringSlice(data["status"])[0] == "partly" {
		return Handler("OK", "Not changed")
	}

	var (
		getCoursesUrl = getListURL + "sutdent_course"
		getCoursesReq = Request{
			Data: map[string]interface{}{
				"director_course_id": data["id"],
			},
		}
		getCoursesResp = GetListClientApiResponse{}
	)

	body, err := DoRequest(getCoursesUrl, "POST", getCoursesReq)
	if err != nil {
		return Handler("error", err.Error())
	}
	if err := json.Unmarshal(body, &getCoursesResp); err != nil {
		return Handler("error", err.Error())
	}

	for _, val := range getCoursesResp.Data.Data.Response {
		var (
			deleteCourseUrl = getSingleURL + "sutdent_course/" + cast.ToString(val["guid"])
		)

		_, err := DoRequest(deleteCourseUrl, "DELETE", Request{Data: map[string]interface{}{}})
		if err != nil {
			return Handler("error", err.Error())
		}
	}

	return Handler("OK", "Success")
}

func DoRequest(url string, method string, body interface{}) ([]byte, error) {
	data, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	request.Header.Add("authorization", "API-KEY")
	request.Header.Add("X-API-KEY", apiKey)

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respByte, nil
}

func Handler(status, message string) string {
	var (
		response Response
		Message  = make(map[string]interface{})
	)

	sendMessage("delete-student-courses", status, message)
	response.Status = status
	data := Request{
		Data: map[string]interface{}{
			"data": message,
		},
	}
	response.Data = data.Data
	Message["message"] = message
	respByte, _ := json.Marshal(response)
	return string(respByte)
}

func sendMessage(functionName, errorStatus string, message interface{}) {
	bot, err := tgbotapi.NewBotAPI("5625907982:AAGf-AKQCngObyXjpxQBWBiKhZhmmq-HP_k")
	if err != nil {
		log.Panic(err)
	}

	chatID := int64(-1001990127540)
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("message from %s function: %s\n%s", functionName, errorStatus, message))
	_, err = bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

type Response struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

type RequestBody struct {
	ObjectIDs []string               `json:"object_ids"`
	Data      map[string]interface{} `json:"data"`
}

type ClientApiResponse struct {
	Data ClientApiData `json:"data"`
}

type ClientApiData struct {
	Data ClientApiResp `json:"data"`
}

type ClientApiResp struct {
	Response map[string]interface{} `json:"response"`
}

type Request struct {
	Data map[string]interface{} `json:"data"`
}

type MultipleUpdateRequest struct {
	Data Data `json:"data"`
}

type Data struct {
	Objects []map[string]interface{} `json:"objects"`
}

type GetListClientApiResponse struct {
	Data GetListClientApiData `json:"data"`
}

type GetListClientApiData struct {
	Data GetListClientApiResp `json:"data"`
}

type GetListClientApiResp struct {
	Response []map[string]interface{} `json:"response"`
}

type ResponseModel struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

type NewRequestBody struct {
	Data map[string]interface{} `json:"data"`
}

type CreateResponseBody struct {
	Data CreateResponseModel `json:"data"`
}

type CreateResponseModel struct {
	Data CreateResponse `json:"data"`
}

type CreateResponse struct {
	Data map[string]interface{} `json:"data"`
}
