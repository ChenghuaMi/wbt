package main

import (
	"flag"
	"fmt"

	"wbt/server/run"
	"wbt/server/internal/config"
	"wbt/server/internal/handler"
	"wbt/server/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/server-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	runTask := run.NewRunTask(ctx)
	//tk := time.NewTicker(1000)
	//chs := make([]map[string]string,10)
	done := make(chan bool,1)

	go func() {
		for  {
			go func() {
				//url := "http://210.51.167.173:2130/culture-frontend-trade/webHttpServlet"
				//data := []byte(`{"U": "0000000000006037","SI": "4128568194094858240","name": "hold_query_mobile"}`)
				//rs,_ := run.HttpPostJson(url,data)
				//fmt.Println(rs)

				data := make(map[string]string)
				//data["U"] = "0000000000006087"
				data["U"] = "0000000000006103"
				data["P"] = "SWeo+ccDMfb3JdB5XnIQ/g=="
				ss := runTask.Login(run.BaseUrl + run.CultureFrountendTrade,data)
				//ss := runTask.All(runTask.CheckLogin(run.BaseUrl + run.CultureFrountendTrade,runTask.Login(run.BaseUrl + run.CultureFrountendTrade,data)))
				fmt.Println("allllllllllllllllllllllllllll：",<-ss)

				//data := make([]map[string]string,4)
				//for i := 0;i<len(data);i++ {
				//	data[i] = make(map[string]string)
				//}
				//data[0]["U"] = " 0000000000006103"
				//data[0]["P"] = "SWeo+ccDMfb3JdB5XnIQ/g=="
				//data[1]["U"] = "0000000000006090"
				//data[1]["P"] = "SWeo+ccDMfb3JdB5XnIQ/g=="
				//go func() {
				//	for i := 0;i<len(data);i++ {
				//		ss := runTask.All(runTask.CheckLogin(run.BaseUrl + run.CultureFrountendTrade,runTask.Login(run.BaseUrl + run.CultureFrountendTrade,data[i])))
				//		fmt.Println("allllllllllllllllllllllllllll：====>",i,<-ss)
				//		time.Sleep(5 * time.Second)
				//	}
				//
				//}()
				//
				<-done
			}()
			done <- true
		}
	}()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}


