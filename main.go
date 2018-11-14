package main

import (
	"net/http"
	"strconv"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gorilla/mux"
	"time"
	"log"
)

type ConsumeInfo struct {
	Result      bool 				`json:result`
	Msg 		string 				`json:msg`
	List 		[]ConsumeInfoMonth	`json:list`
}

type ConsumeInfoMonth struct {
	Title		string 				`json:title`
	Data 		[]ConsumeMeta 		`json:data`
}

type ConsumeMeta struct {
	SmtDealName 		string 		`json:smtDealName`
	SmtTransMoney 		string 		`json:smtTransMoney`
	SmtDealDateTimeTxt 	string		`json:smtDealDateTimeTxt`
	Date 				string 		`json:date`
	Time  				string		`json:time`
	SmtOrgName			string 		`json:smtOrgName`
	SmtInMoney			string 		`json:smtInMoney`
	SmtOutMoney			string		`json:smtOutMoney`
}

type Message struct {
	Msg  				string 		`json:msg`
}

func call(sid string, page int, endtimestring string, res *ConsumeInfo) error {
	const BASEURL = "http://weixin.ccnu.edu.cn/App/weixin/queryTrans"
	const PAGESIZE = 20
	URL := BASEURL + "?page=" + strconv.Itoa(page) + "&pageSize=" + strconv.Itoa(PAGESIZE) + "&startTime=2018-01-01" +
			"&endTime=" + endtimestring

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/16A404 MicroMessenger/6.7.3(0x16070321) NetType/WIFI Language/zh_CN")
	req.Header.Set("Cookie", "wxqyuserid=" + sid)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(1)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, res)
	if err != nil {
		fmt.Println(2)
		return err
	}

	return nil
}

func ConsumeDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page, err := strconv.Atoi(vars["page"])
	if err != nil {
		msg := new(Message)
		msg.Msg = err.Error()
		json.NewEncoder(w).Encode(msg)
		return
	}

	sid := r.Header.Get("sid")
	y, m, d := time.Now().Date()
	endtime := strconv.Itoa(y) + "-" + strconv.Itoa(int(m)) + "-" + strconv.Itoa(d)
	//fmt.Println(endtime)
	res := ConsumeInfo{}
	fmt.Println(sid, page, endtime)
	err = call(sid, page, endtime, &res)

	if err != nil {
		msg := new(Message)
		msg.Msg = err.Error()
		json.NewEncoder(w).Encode(msg)
		return
	}
	w.Header().Set("Content-type","application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/consume/details/{page}/", ConsumeDetails)
	log.Fatal(http.ListenAndServe(":8080", router))
}
