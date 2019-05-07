package exchange

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kinggigo/secret/server/db/schema"
	"github.com/kinggigo/secret/server/db/service"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

// TODO
//1. 코인 마켓 거래정보 (getCoinMarketSub 'coinType'=>$symbol, 'marketType'=>$market)
//2. 마켓리스트 가져오기(market_service->getMarketList)
//3. 마켓에 등록된 코인들의 현재가 종가 이름을 가져온다(거래소 좌측 코인 리스트)(exchange_service -> getSymbolsPriceWithMarket)
//4. 특정 심볼의 설정된 정보를 가져온다 (기본설정 :  coinTitle = mc_kor_name )(currency_service->get_one_info)
//5. 특정 코인 종가 현재가, 고가, 저가, 거래량, 거래대금 등의 정보를 가져온다. (exchange_service->marketStatus)
//6. 주문 목록 가져오기 -> controller/param.go 320번째 줄
//7. 주문 완료 목록 가져오기 -> controller/param.go 313번째 줄
//8. 로그인되어 있다면 자신의 심볼과 마켓 코인의 자산 정보를 가져온다.exchange_service->getAssets (coinType를 symbol과 market의 파라미터를 넘겨 각각 가져온다. )\
//8-1 거래 수수료 (exchange_service->getTradeFeeRate(array('memId'=> $mem_id,'coinType' => $symbol));)
//8-2 미체결 리스트(ㅒexchange_service->getTransListProceeding(array('memId'=> $mem_id,'coinType' => $symbol,'marketType'=>$market));)
//8-3 체결 리스트(exchange_service->getTransListCompleted(array('memId'=> $mem_id,'coinType' => $symbol,'marketType'=>$market));)
//9 로그인상태가 아니면 (tradeFeeRate 에 0.0015를 넣는다 => front에서 처리)

var (
	upgrader = websocket.Upgrader{}
)

//1. 코인 마켓 거래정보 (getCoinMarketSub 'coinType'=>$symbol, 'marketType'=>$market)
func GetCoinMarketSub(e echo.Context) error {

	return nil
}

//2. 코인리스트
/**
param : market
return :
*/
func GetCoinList(e echo.Context) error {
	market := e.QueryParam("market")
	CoinInfoList := []schema.CoinListInfo{}
	if market == "" || len(market) == 0 {
		market = "KRW"
	}
	CoinInfoList, err := service.CoinList(market)
	if err != nil {
		log.Info(err)
		e.Error(err)
		e.String(500, err.Error())
	}

	return e.JSON(200, CoinInfoList)
}

/**
한 코인 정보 가져오기
*/
func GetCoinInfo(e echo.Context) error {
	market := e.QueryParam("market")
	cointype := e.QueryParam("coinType")
	if market == "" || len(market) == 0 {
		//e.Error(errors.New("Market을 지정하세요"))
		return e.String(400, "Market을 지정하세요")
	}
	if cointype == "" || len(cointype) == 0 {
		//e.Error(errors.New("CoinType을 지정하세요"))
		return e.String(400, "CoinType을 지정하세요")
	}
	coinInfo, err := service.CoinInfo(market, cointype)
	if err != nil {
		e.Error(err)
		e.String(500, err.Error())
	}

	return e.JSON(200, coinInfo)
}

/**
체결완료 리스트를 가져온다
*/
func GetContract(e echo.Context) error {
	market := e.QueryParam("market")
	cointype := e.QueryParam("coinType")
	if market == "" || len(market) == 0 {
		//e.Error(errors.New("Market을 지정하세요"))
		return e.String(400, "Market을 지정하세요")
	}
	if cointype == "" || len(cointype) == 0 {
		//e.Error(errors.New("CoinType을 지정하세요"))
		return e.String(400, "CoinType을 지정하세요")
	}
	contractList, err := service.ContractList(market, cointype)
	if err != nil {
		return e.String(500, err.Error())
	}
	return e.JSON(200, contractList)
}

/**
/exchange/order ?coinType,market
오더북을 가져온다.
구매리스트,
*/
func GetOrderBook(e echo.Context) error {
	BuyList := []schema.BuyOrder{}
	SellList := []schema.BuyOrder{}

	market := e.QueryParam("market")
	coinType := e.QueryParam("coinType")
	BuyList, err := service.Order_BuyList(market, coinType)
	if err != nil {
		e.String(500, err.Error())
	}
	SellList, err = service.Order_SellList(market, coinType)
	if err != nil {
		e.String(500, err.Error())
	}
	ClosePriz, err := service.Closing_Price(market, coinType)
	if err != nil {
		e.String(500, err.Error())
	}
	zeroBuy := schema.BuyOrder{"0", "0"}
	//8개를 풀로 주기 위함
	if len(BuyList) < 8 {
		lenbuy := len(BuyList)
		for i := 0; i < (8 - lenbuy); i++ {
			BuyList = append(BuyList, zeroBuy)
		}

	}
	if len(SellList) < 8 {
		lensell := len(SellList)
		for i := 0; i < (8 - lensell); i++ {
			SellList = append(SellList, zeroBuy)
		}
	}
	for i, j := 0, len(SellList)-1; i < j; i, j = i+1, j-1 {
		SellList[i], SellList[j] = SellList[j], SellList[i]
	}
	Res := struct {
		BUY   []schema.BuyOrder
		SELL  []schema.BuyOrder
		PRICE string
	}{BuyList, SellList, ClosePriz}

	return e.JSON(200, Res)
}

/**
/exchange/order ?coinType,market
오더북을 가져온다.
websocket으로
*/
func GetOrderBook_WS(e echo.Context) error {
	//BuyList := []schema.BuyOrder{}
	//SellList := []schema.BuyOrder{}

	//market := e.QueryParam("market")
	//coinType := e.QueryParam("coinType")
	//BuyList, err := service.Order_BuyList(market, coinType)
	//if err != nil {
	//	e.String(500, err.Error())
	//}
	//SellList, err = service.Order_SellList(market, coinType)
	//if err != nil {
	//	e.String(500, err.Error())
	//}
	//ClosePriz, err := service.Closing_Price(market, coinType)
	//if err != nil {
	//	e.String(500, err.Error())
	//}
	//zeroBuy := schema.BuyOrder{"0", "0"}
	////8개를 풀로 주기 위함
	//if len(BuyList) < 8 {
	//	lenbuy := len(BuyList)
	//	for i := 0; i < (8 - lenbuy); i++ {
	//		BuyList = append(BuyList, zeroBuy)
	//	}
	//
	//}
	//if len(SellList) < 8 {
	//	lensell := len(SellList)
	//	for i := 0; i < (8 - lensell); i++ {
	//		SellList = append(SellList, zeroBuy)
	//	}
	//}
	//for i, j := 0, len(SellList)-1; i < j; i, j = i+1, j-1 {
	//	SellList[i], SellList[j] = SellList[j], SellList[i]
	//}
	//Res := struct {
	//	BUY   []schema.BuyOrder
	//	SELL  []schema.BuyOrder
	//	PRICE string
	//}{BuyList, SellList, ClosePriz}
	//
	//return e.JSON(200, Res)

	ws, err := upgrader.Upgrade(e.Response(), e.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("hello~"))
		if err != nil {
			e.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			e.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}

/*
미체결 리스트 - 로그인 되어있어야함
*/
func GetTransList(e echo.Context) error {
	market := e.QueryParam("market")
	coinType := e.QueryParam("coinType")
	//TODO : JWT로 토큰을 발행 한 후 ID에 대한 값을 memID로 넣어 service를 불러와야함.
	memID := e.QueryParam("memID")
	TransList, err := service.TransList(coinType, market, memID)
	if err != nil {
		return e.String(500, "에러!!!"+err.Error())
	}

	return e.JSON(200, TransList)
}

/*
체결 리스트 - 로그인 되어있어야함
*/
func GetTransListCom(e echo.Context) error {
	market := e.QueryParam("market")
	coinType := e.QueryParam("coinType")
	//TODO : JWT로 토큰을 발행 한 후 ID에 대한 값을 memID로 넣어 service를 불러와야함.
	memID := e.QueryParam("memID")

	//최종 체결리스트
	TransListFinal := []schema.TransComRet{}
	temp := schema.TransComRet{}
	TransList, err := service.TransListCom(coinType, market, memID)
	if err != nil {
		return e.String(500, "에러!!!"+err.Error())
	}
	for _, trans := range TransList {
		//판매, 구매가 같은 경우
		if trans.BUY_MEM_ID == memID && trans.SELL_MEM_ID == memID {
			temp.CONT_NO = trans.CONT_NO
			temp.SEARCH_GB = "1F"
			temp.TRAN_GN = "매수"
			temp.TRANSFER_DATE = trans.TRAN_DT
			temp.COIN_AMT = trans.COIN_AMT
			temp.KRW_AMT = trans.KRW_AMT
			temp.CONTRACT_AMT = trans.CONT_AMT
			temp.FEE = trans.COIN_FEE
			TransListFinal = append(TransListFinal, temp)
			temp.CONT_NO = trans.CONT_NO
			temp.SEARCH_GB = "2F"
			temp.TRAN_GN = "매도"
			temp.TRANSFER_DATE = trans.TRAN_DT
			temp.COIN_AMT = trans.COIN_AMT
			temp.KRW_AMT = trans.KRW_AMT
			temp.CONTRACT_AMT = trans.CONT_AMT
			temp.FEE = trans.KRW_FEE
			TransListFinal = append(TransListFinal, temp)
			//구매완료
		} else if trans.BUY_MEM_ID == memID && trans.SELL_MEM_ID != memID {
			temp.CONT_NO = trans.CONT_NO
			temp.SEARCH_GB = "1F"
			temp.TRAN_GN = "매수"
			temp.TRANSFER_DATE = trans.TRAN_DT
			temp.COIN_AMT = trans.COIN_AMT
			temp.KRW_AMT = trans.KRW_AMT
			temp.CONTRACT_AMT = trans.CONT_AMT
			temp.FEE = trans.COIN_FEE
			TransListFinal = append(TransListFinal, temp)
			//판매 완료
		} else if trans.BUY_MEM_ID != memID && trans.SELL_MEM_ID == memID {
			temp.CONT_NO = trans.CONT_NO
			temp.SEARCH_GB = "2F"
			temp.TRAN_GN = "매도"
			temp.TRANSFER_DATE = trans.TRAN_DT
			temp.COIN_AMT = trans.COIN_AMT
			temp.KRW_AMT = trans.KRW_AMT
			temp.CONTRACT_AMT = trans.CONT_AMT
			temp.FEE = trans.KRW_FEE
			TransListFinal = append(TransListFinal, temp)
		}
	}
	return e.JSON(200, TransListFinal)
}
