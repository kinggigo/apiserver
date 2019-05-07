package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kinggigo/secret/server/api/account"
	"github.com/kinggigo/secret/server/api/contact"
	"github.com/kinggigo/secret/server/api/exchange"
	"github.com/kinggigo/secret/server/api/main"
	"github.com/kinggigo/secret/server/config"
	"github.com/kinggigo/secret/server/db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/sirupsen/logrus"
)

func main() {
	//DB 생성 후 db에 넣어놓음
	dbUrl := db.GetDBConnURL()

	e := echo.New()
	e.Use(middleware.Logger())
	gormdb, err := gorm.Open("mysql", dbUrl)
	if err != nil {
		panic("FAIL TO CONNECT DATABASE     " + err.Error())
	}
	gormdb.LogMode(true)

	db.GormDB = gormdb
	//기본 react로 바인딩
	//e.Static("/", "../client/woori/build")

	//jwt 설정 시 필요함
	logined := e.Group("/logined")
	logined.Use(middleware.JWT([]byte(config.JWT_KEY)))
	logined.GET("/test", account.Happyy)

	notlogined := e.Group("/")
	notlogined.Use(middleware.Logger())

	//MAIN 마켓리스트, 보도자료, 뉴스 가져오기,
	notlogined.GET("main/marketlist", mainpage.MarketList)
	notlogined.POST("main/mainboard", mainpage.MainBoard)

	//CONTACT 게시판 (공지 / 보도자료 /투자정보 / 상장코인정보)
	notlogined.GET("contact/board", contact.Board)
	notlogined.GET("contact/board_view", contact.Board_View)

	//로그인부분
	notlogined.GET("login", account.Login)

	//거래소
	//코인 리스트 가져오기(거래소 코인리스트 ) get으로 market 을 보내서 각 리스트를 받는다. - market
	notlogined.GET("exchange/coinlist", exchange.GetCoinList)
	//특정 코인 거래정보 - market , coinType
	notlogined.GET("exchange/coininfo", exchange.GetCoinInfo)
	//차트
	//notlogined.GET("exchange/coininfo", exchange.GetCoinList)

	//오더북
	notlogined.GET("exchange/order", exchange.GetOrderBook)
	//TODO :  websocket 으로
	notlogined.GET("exchange/order_WS", exchange.GetOrderBook_WS)

	//매수,매도,
	//미체결, 체결 - 로그인
	logined.GET("exchange/transList", exchange.GetTransList)
	//체결가
	logined.GET("exchange/transListCom", exchange.GetTransListCom)
	//최근 체결 내역 -market , coinType
	notlogined.GET("exchange/contract", exchange.GetContract)

	//입출금부분 BANK
	//logined.GET("bank/account",exchange.GetTransList)

	e.POST("/", mainpage.MarketList)
	//e.GET("/main", )

	//TODO :: https://echo.labstack.com/guide/routing 로 URL 그룹핑/바인딩 하기

	e.Use(middleware.CORS())

	e.Logger.Fatal(e.Start(":1323"))

}
