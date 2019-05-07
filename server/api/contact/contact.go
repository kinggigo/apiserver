package contact

import (
	"github.com/kinggigo/secret/server/db/schema"
	"github.com/kinggigo/secret/server/db/service"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"strconv"
)

//공지사항, 뉴스, 코인, 투자 정보 를 type에 따라 가져옴noti
/*
GET
param : type
notice(공지), news(뉴스), coin(코인정보), invest(투자 정보)
*/
func Board(e echo.Context) error {
	typeparam := e.QueryParam("type")
	nPagestr := e.QueryParam("nPage")
	nPage, err := strconv.Atoi(nPagestr)
	log.Info(nPage)
	if err != nil {
		nPage = 1
		log.Info("AtoI error " + err.Error())
	}
	if nPage < 1 {
		nPage = 1
	}

	nPageStartNum := (nPage - 1) * 10
	req := struct {
		BoardList []schema.Board
		NTotal    int
	}{}
	param := schema.BoardParm{Lang_type: "KOREAN", StartNum: nPageStartNum, Limit: 10}

	if typeparam == "notice" {
		param.Code = "B001"
	} else if typeparam == "news" {
		param.Code = "B002"
	} else if typeparam == "coin" {
		param.Code = "COIN"
	} else if typeparam == "invest" {
		param.Code = "B003"
	} else {
		return e.String(500, "없는 게시판 타입입니다.")
	}
	BorderList, err := service.BoardList(param)
	if err != nil {
		return e.String(400, err.Error())
	}
	NTotal, err := service.GetRowNum()
	if err != nil {
		log.Info(err.Error())
	}
	req.NTotal = NTotal
	req.BoardList = BorderList
	return e.JSON(200, req)
}

/**
공지사항/뉴스/ 투자정보 / 상세정보
return  게시판 타이틀,type은 미포함 시킴( controller/Contact.php(183))
조회수 증가, 게시글 정보 리턴
*/
func Board_View(e echo.Context) error {
	num := e.QueryParam("num")
	if num == "" || len(num) == 0 {
		e.String(500, "게시판 번호가 없다.")
	}
	//게시글 정보
	board, err := service.Board_View(num)
	if err != nil {
		return e.String(500, err.Error())
	}

	//조회수 올리기
	//실제에 반영은 올려야함
	//err = service.Board_View_Cnt(num)
	//if err != nil{
	//	log.Error("errr!!!!!! ", err)
	//}

	return e.JSON(200, board)
}
