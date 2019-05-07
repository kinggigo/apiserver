package mainpage

import (
	"github.com/kinggigo/secret/server/db/schema"
	"github.com/kinggigo/secret/server/db/service"
	"github.com/labstack/echo"
	"net/http"
)

/**
main에서 보내는 코인 리스트 ( market, coinType 에따라 지금 사용되는 코인의 정보들을 리턴)
GET
param : market , coinType
return : coinList
*/
type Res struct {
	Maket    string `json:"m"`
	Cointype string `json:"c"`
}

func MarketList(e echo.Context) error {
	market := e.QueryParam("market")
	//coinType := e.QueryParam("coinType")
	//r := Res{market, coinType}
	if market == "" {
		return e.String(500, "Market에 대한 값이 없습니다")
	}
	list, err := service.MarketList(market)
	if err != nil {
		return e.String(500, "DB조회 에러"+err.Error())
	}
	//if err := e.Bind(r); err != nil {
	//	return nil
	//}

	return e.JSON(http.StatusOK, list)
}

/**
메인 공지사항, 뉴스 3개 씩만 뿌려주는것
GET
param : type
return : notice 3개, news 3개
*/
//TODO contact/board 를 불러오는데 type으로 main을 넘겨서 하는데 가능한가? context 를 넘겨서 확인해봐야함
func MainBoard(e echo.Context) error {
	var Mainboard = struct {
		Notice []schema.Board
		News   []schema.Board
	}{}

	param := schema.BoardParm{
		Code:      "B001",
		Lang_type: "KOREAN",
		StartNum:  0,
		Limit:     3,
	}
	//post로 파라메터를 넘길때 써야함.
	//if err := e.Bind(&param); err != nil {
	// 	return err
	//}

	BoardList, err := service.BoardList(param)
	if err != nil {
		return err
	}
	Mainboard.Notice = BoardList
	param = schema.BoardParm{
		Code:      "B002",
		Lang_type: "KOREAN",
		StartNum:  0,
		Limit:     3,
	}
	NewsList, err := service.BoardList(param)
	if err != nil {
		return err
	}
	Mainboard.News = NewsList
	return e.JSON(200, Mainboard)
}

/**
로그인이 되어 있으면 로그인된 계정 정보를 넘겨줌
POST
param : 없음
return : 로그인된 계정 정보
*/
func MemberInfo(e echo.Context) error {

	return nil
}
