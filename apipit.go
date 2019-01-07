package apipit

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"google.golang.org/appengine"
)

type OnlineAction struct {
	UID           string           `json:"uid,omitempty"`
	Contract      string           `json:"contract,omitempty"`
	Card          string           `json:"card,omitempty"`
	Type          string           `json:"type,omitempty"`
	Refundable    int              `json:"refundable,omitempty"`
	NonRefundable int              `json:"nonRefundable,omitempty"`
	NonAccounting int              `json:"nonAccounting,omitempty"`
	PaymentMet    []PaymentMethods `json:"paymentMethods,omitempty"`
	Articles      []byte           `json:"articles,omitempty"`
}

type PaymentMethods struct {
	UID    string `json:"uid,omitempty"`
	Amount int    `json:"amount,omitempty"`
}

type AccountCashless struct {
	UID        string `json:"uid,omitempty"`
	Status     int    `json:"status,omitempty"`
	Email      string `json:"email,omitempty"`
	FirstName  string `json:"firstname,omitempty"`
	LastName   string `json:"lastname,omitempty"`
	Language   string `json:"language,omitempty"`
	Timezone   string `json:"timezone,omitempty"`
	CreatedAt  string `json:"createdAt,omitempty"`
	LastUpdate string `json:"lastUpdate,omitempty"`
}

type Log struct {
	Method  string
	Object  string
	Url     string
	Time    time.Time
	Status  int
	Message string
}

type ReturnBody struct {
	Type    string
	Message string
	Extra   string
	Error   string
}

func GetMethod(w http.ResponseWriter, req *http.Request) {

	ctx := appengine.NewContext(req)

	returnData := GetAPI("f154f0c1-9670-4372-b60e-4ade397678", "payment-methods", "", ctx)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	//json.NewEncoder(w).Encode(dados)
	w.Write(returnData)
}

func GetAPI(contract string, object string, params string, ctx context.Context) []byte {

	//client := urlfetch.Client(ctx)
	client := &http.Client{}

	// GET IN ANOTHER API
	url := "https://api.com" + contract + "/" + object

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer 76yjhXAiOiJKV1QiLCJhbGciOiJSUzUxMiJ9.eyJhcHAiOiJpby5wYXlpbnRlY2guYW5kcm9pZC5zZGsuc2FtcGxlIiwiYXVkIjoiaW8ucGF5aW50ZWNoLnBpdCIsImZlYXR1cmVzIjpbIioiXSwiaXNzIjoiaW8ucGF5aW50ZWNoLnBpdCIsImNlcnQiOiIzOTpDODpGOTpDNzoyOTpGRDozMDpEMTo1MzpCNzo5QjpERTo5RTowNzozMTpGMDpCRToxODpFNjpFMTo4MzoxNDo1QTo5OTpDNTpBMTpGNTozMDoxMDowQjo4MzoyRCIsImp0aSI6ImZkNzE3OTRlZTJhMzQ5ZDM5ZjkzYmI1OWFiOGY5YzBjIn0.H7rsB0IMP8sGKl9roQuImuHQnAUpjcnuKfEoJLO9n3wivWIIjtl-H8CigGek5DuYwCfFv1b3-3FvnMvLSk2uxUpOEp3MrAA5zxNJskFUDDlNdK-vRPRYc1PxRtRdiMQMDxNOfBiXKL-w7VKBB6rv_yusYiYDAG0hmXk3QqaD7Q1jtEWMNrqpLkVeUFQRGKYNc9H8UAXXLiXPZq2v88-Uw00iHP4AfDNtFoGnS4grVeE_WK6CKpObcPRdR8A__1U0Pp4FpHIaFAbGXzx-fuA3eGey_NsByT54fkHT_rXvYTvl7pAmvX7fI0ne6ZAzevfQK-TnGd7V9mowwt4CXhA2ww")

	msg := "OK"
	response, err := client.Do(req)
	if err != nil {
		msg = err.Error()
	}

	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)

	if msg == "OK" && response.StatusCode != 200 {

		var bodymsg ReturnBody
		err = json.Unmarshal([]byte(data), &bodymsg)

		if err == nil {

			msg = ""

			if bodymsg.Type != "" {
				msg = bodymsg.Type + ","
			}
			if bodymsg.Message != "" {
				msg = msg + bodymsg.Message + ","
			}
			if bodymsg.Extra != "" {
				msg = msg + bodymsg.Extra + ","
			}
			if bodymsg.Error != "" {
				msg = msg + bodymsg.Error
			}

		} else {
			msg = "Não foi possivel captar a mensagem de erro"
		}

	}

	RecordLog("GET", object, url, response.StatusCode, msg, ctx)

	return data
}

func PostMethod(w http.ResponseWriter, req *http.Request) {

	var onlineAction OnlineAction
	var msg string
	ctx := appengine.NewContext(req)

	_ = json.NewDecoder(req.Body).Decode(&onlineAction)

	//RecordLog(onlineAction.Contract, onlineAction.Card, onlineAction.Type, 1, onlineAction.PaymentMet[0].UID, ctx)

	oa := []OnlineAction{}
	oa = append(oa, onlineAction)

	param, _ := json.Marshal(oa)
	requestReader := bytes.NewReader(param)

	// POST in another APO
	url := "https://api.com"

	req, err := http.NewRequest("POST", url, requestReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer uy5reXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJpby5wYXlpbnRlY2gucGl0IiwiYXVkIjoiaW8ucGF5aW50ZWNoLnBpdCIsInNjb3BlcyI6eyIqIjpbIioiXX0sIm9iaiI6ImJjOWVlNWZlLWIzMzktNDA4YS1hZjg0LTE2NzE5NjcwYmRlZSIsImp0aSI6IjdmNWRhZmI5NDk0YzRmZTZiNDM1OWE2Y2NlZjg4YWQ5In0.4J1tu_h-JGjsAcsWtG5lg7CuP3cJHeXznHYCCOnjpcg")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		msg = err.Error()
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		msg = "OK"
	}
	data, _ := ioutil.ReadAll(response.Body)

	if msg == "OK" && response.StatusCode != 200 {

		var bodymsg ReturnBody
		err = json.Unmarshal([]byte(data), &bodymsg)

		if err == nil {

			msg = ""

			if bodymsg.Type != "" {
				msg = bodymsg.Type + ","
			}
			if bodymsg.Message != "" {
				msg = msg + bodymsg.Message + ","
			}
			if bodymsg.Extra != "" {
				msg = msg + bodymsg.Extra + ","
			}
			if bodymsg.Error != "" {
				msg = msg + bodymsg.Error
			}

		} else {
			msg = "Não foi possivel captar a mensagem de erro"
		}

	}

	RecordLog("POST", "online-action", url, response.StatusCode, msg, ctx)

}

func RecordLog(method string, object string, url string, status int, message string, ctx context.Context) {

	IDprojeto := "Your ID project in GCP"

	current_time := time.Now().Local()

	log := &Log{
		Method:  method,
		Object:  object,
		Url:     url,
		Time:    current_time,
		Status:  status,
		Message: message,
	}

	datastoreClient, err := datastore.NewClient(ctx, IDprojeto)
	if err != nil {
		//fmt.Errorf(err.Error())
	}

	// Sets the kind for the new entity.
	tipo := "logAPI"

	// Creates a Key instance.
	logKey := datastore.IncompleteKey(tipo, nil)

	// Saves the new entity.
	_, erro := datastoreClient.Put(ctx, logKey, log)
	if erro != nil {
		//fmt.Printf("Failed to save task: %v", logKey, erro)
	}

	
}
