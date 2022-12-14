package run

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"wbt/server/internal/svc"
)

const (
	Timeout = 2000
	BaseUrl = "http://210.51.167.173:2130/"
	CultureFrountendTrade = "culture-frontend-trade/webHttpServlet"
)
type currentResult struct {
	name string
	result interface{}
}

type RunTask struct {
	c *svc.ServiceContext
}

func NewRunTask(c *svc.ServiceContext) *RunTask {
	return &RunTask{
		c: c,
	}
}

//{
//"name":"logon",
//"U":"13638451711",
//"P":"rLuT6jXY2EIDV+7LFwrh7w==",
//"LT":"web"
//}
// data[U] = "13638451711"
// data[P] = "rLuT6jXY2EIDV+7LFwrh7w=="
//return data
//name: "xxx",
//result: xxxx

func(r *RunTask) Login(url string,data map[string]string) <-chan string {
	valStream := make(chan string)
	u := data["U"]
	retVal,_ := r.c.Redis.C.Get(r.c.Redis.Ctx,u).Result()
	//if err != nil {
	//	fmt.Println(err)
	//	panic(fmt.Sprintf("redis data error %s \n",err.Error()))
	//}
	if len(retVal) > 0 {
		result := make(map[string]interface{})
		result["U"] = u
		result["RETCODE"] = retVal
		rs := make(map[string]interface{})
		rs["name"] = "logon"
		rs["result"] = result

		by,err := json.Marshal(rs)
		if err != nil {
			fmt.Println("map 转成json 格式错误")
			close(valStream)
			return valStream
		}
		jsonStr := string(by)
		fmt.Println("???????????????????????",by,jsonStr)
		go func() {
			valStream <- string(by)
		}()

		return valStream

	}else {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>to login....")

		var req map[string]string
		req = make(map[string]string,4)
		req["name"] = "logon"
		req["LT"] = "web"
		for k,v := range data {
			req[k] = v
		}
		bt,err := json.Marshal(req)
		if err != nil {
			fmt.Println("map 转成json 格式错误")
			close(valStream)
			return valStream
		}
		jsonStr := string(bt)

		go func() {
			defer close(valStream)
			for  {
				select {
				case <-time.After(Timeout * time.Second):
					fmt.Println("login time out........................",time.Now())
					return
				case valStream <- HttpPostJson(url, []byte(jsonStr)):
					fmt.Println("get " + req["name"])

				}
			}
		}()


	}

	return valStream

}

//data login 返回的数据
func(r *RunTask) CheckLogin(url string,data <-chan string) <-chan string {
	//var dataMap map[string]interface{}
	valStream := make(chan string)
	go func() {
		defer close(valStream)
		for  {
			select {
			case <-time.After(Timeout * time.Second):
				fmt.Println("checkLogin timeout")
				return
			case ss := <-data:
				fmt.Println("Sss=============================================:",ss)
				<-r.CheckUserInfo(ss)
				fmt.Println("okokok")
				go func() {
					valStream <- ss // 验证 登录信息的有效，过期，删除redis,重新请求接口
				}()



			}
		}

	}()
	return valStream
}
//验证用户信息是否过期
//data => map[string]string 的json字符串   data[0000000000006087]=SI
func(r *RunTask) CheckUserInfo(data string) <-chan string {
	valStream := make(chan string)
	var dataMap map[string]interface{}
	err := json.Unmarshal([]byte(data),&dataMap)
	if err != nil {
		panic(fmt.Sprintf("数据解析失败 %s\n",err.Error()))
	}
	result := dataMap["result"]

	if val,ok := result.(map[string]interface{});ok {
		vvv,_ := val["RETCODE"].(string)
		u,_ := val["U"].(string)
		fmt.Println("uuuuu:",u,"vvvvvv:",vvv)
		checkReq := make(map[string]string)
		checkReq["U"] = u
		checkReq["SI"] = vvv
		checkReq["name"] = "check_userinfo"
		checkReq["FLT"]="web"
		checkReq["LT"]="web"
		reqby,_ := json.Marshal(checkReq)
		res := HttpPostJson(BaseUrl + CultureFrountendTrade,reqby)
		checkInfo := make(map[string]interface{})
		json.Unmarshal([]byte(res),&checkInfo)
		fmt.Println("checkInfo:",checkInfo)
		rs := checkInfo["result"]
		if rr,ok := rs.(map[string]interface{});ok {
			retcode,_ := rr["RETCODE"].(string)

			fmt.Println("retcode:",retcode)
			code,_ := strconv.Atoi(retcode)
			if code < 0 {
				//清除 redis缓存
				r.c.Redis.C.Del(r.c.Redis.Ctx,u)
				fmt.Println("del redis >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
				//todo 重新登录
				dd := make(map[string]string)
				dd["U"] = "0000000000006087"
				//data["U"] = "0000000000006090"
				dd["P"] = "SWeo+ccDMfb3JdB5XnIQ/g=="
				valStream <- <-r.Login(BaseUrl + CultureFrountendTrade,dd)
			}else {
				r.c.Redis.C.Set(r.c.Redis.Ctx,u,vvv,0)
				fmt.Println("write redis data.... success................................................")
				go func() {
					valStream <- data
				}()
			}
		}


	}

	return valStream
}
func(r *RunTask) CommonReqParam(data string,addParam map[string]interface{}) map[string]interface{} {
	reqData := make(map[string]interface{})
	var dataMap map[string]interface{}
	err := json.Unmarshal([]byte(data),&dataMap)
	if err != nil {
		panic(fmt.Sprintf("数据解析失败 %s\n",err.Error()))
	}
	result := dataMap["result"]
	if val,ok := result.(map[string]interface{});ok{
		u,_ := val["U"].(string)
		si := r.c.Redis.C.Get(r.c.Redis.Ctx,u).Val() //session
		reqData["U"] = u
		reqData["SI"] = si
	}
	for k,v := range addParam {
		reqData[k] = v
	}
	return reqData
}
//交易商总资金查询
func(r *RunTask) FirmSumInfo(data  string) <-chan string{
	valStream := make(chan string)
	reqData := make(map[string]interface{})
	reqData["name"] = "firm_sum_info"
	requestParam := r.CommonReqParam(data,reqData)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx:",requestParam)
	bt,_ := json.Marshal(requestParam)
	go func() {
		defer close(valStream)
		for  {
			select {
			case valStream <- HttpPostJson(BaseUrl + CultureFrountendTrade,bt):
			}
		}
	}()
	return valStream
}
//交易员资金信息查询(firm_info)
func(r *RunTask) FirmInfo(data  string) <-chan string{
	valStream := make(chan string)
	reqData := make(map[string]interface{})
	reqData["name"] = "firm_info"
	requestParam := r.CommonReqParam(data,reqData)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!FirmInfo:",requestParam)
	bt,_ := json.Marshal(requestParam)

	go func() {
		defer close(valStream)
		for  {
			select {
			case valStream <- HttpPostJson(BaseUrl + CultureFrountendTrade,bt):
			}
		}
	}()
	return valStream
}
//充值提现记录查询（不分当前历史）(invest_draw_query)
func(r *RunTask) InvestDrawQuery(data  string) <-chan string{
	valStream := make(chan string)
	reqData := make(map[string]interface{})
	reqData["name"] = "invest_draw_query"
	reqData["RECCNT"] = 100
	requestParam := r.CommonReqParam(data,reqData)
	bt,_ := json.Marshal(requestParam)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!invest_draw_query:",reqData)

	go func() {
		defer close(valStream)
		for  {
			select {
			case valStream <- HttpPostJson(BaseUrl + CultureFrountendTrade,bt):
			}
		}
	}()
	return valStream
}
//持仓汇总查询-修改 hold_query_mobile  =>存货汇总
func(r *RunTask) HoldQueryMobile(data  string) <-chan string{
	valStream := make(chan string)
	reqData := make(map[string]interface{})
	reqData["name"] = "hold_query_mobile"
	requestParam := r.CommonReqParam(data,reqData)
	bt,_ := json.Marshal(requestParam)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!hold_query_mobile:",reqData)

	go func() {
		defer close(valStream)
		for  {
			select {
			case valStream <- HttpPostJson(BaseUrl + CultureFrountendTrade,bt):
			}
		}
	}()
	return valStream
}
func(r *RunTask) All(check <-chan string) <-chan string {
	valStream := make(chan string)
	go func() {
		defer close(valStream)
		for  {
			select {
			case <-time.After(Timeout * time.Second):
				fmt.Println("all timeout........................",time.Now())
			case valStream <- <-r.FirmSumInfo(<-check):
			case valStream <- <-r.FirmInfo(<-check):
			case valStream <- <-r.InvestDrawQuery(<-check):
			case valStream <- <-r.HoldQueryMobile(<-check):

			}
		}
	}()

	return valStream

}
//map 转换成byte
func(r *RunTask) jsonToByte(data map[string]string) ([]byte,error) {
	return json.Marshal(data)
}
//json字符串 解析成 map
func(r *RunTask) PassStrToJson(data string) (map[string]interface{},error) {
	mp := make(map[string]interface{})
	err := json.Unmarshal([]byte(data),&mp)

	return mp,err
}

//组装 接口请求数据
func(r *RunTask) makeData(data map[string]interface{}) map[string]string {
	mp := make(map[string]string)
	for k,v := range data {
		vv,ok := v.(string)
		if ok {
			mp[k] = vv
		}
	}

	return mp

}

func HttpPostJson(url string,data []byte) string {
	//jsonStr := []byte(`{"username": "mch","age": 10}`)
	//url := "http://baidu.com"
	req,err := http.NewRequest("POST",url,bytes.NewBuffer(data))

	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type","application/json")

	client := &http.Client{}
	resp,err := client.Do(req)


	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(ioutil.Discard,resp.Body)
	statusCode := resp.StatusCode
	header := resp.Header
	fmt.Println(statusCode,header)

	body,_ := ioutil.ReadAll(resp.Body)
	fmt.Println("body:",string(body))
	return string(body)
}
func HttpPostJson2(url string,data []byte) string {
	//jsonStr := []byte(`{"username": "mch","age": 10}`)
	//url := "http://baidu.com"
	req,err := http.NewRequest("POST",url,bytes.NewBuffer(data))

	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type","application/json")

	client := &http.Client{}
	resp,err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//statusCode := resp.StatusCode
	//header := resp.Header
	//fmt.Println(statusCode,header)


	body,_ := ioutil.ReadAll(resp.Body)
	return string(body)
}

