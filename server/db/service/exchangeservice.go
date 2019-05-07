package service

import (
	"errors"
	db2 "github.com/kinggigo/secret/server/db"
	"github.com/kinggigo/secret/server/db/schema"
	"github.com/labstack/gommon/log"
)

//마켓(KRW, BTC, ETH)에 따른 코인 리스트 가져오기
func CoinList(market string) ([]schema.CoinListInfo, error) {
	db := db2.GormDB
	coinList := []schema.CoinListInfo{}
	coininfo := schema.CoinListInfo{}
	rows, err := db.Raw(`SELECT A.mc_symbol AS symbol, A.mc_market_symbol AS market,
                (SELECT B.price FROM us_exchange_last_ticker B FORCE INDEX(SEARCH_IDX) WHERE B.coin_type = A.mc_symbol AND market_type=A.mc_market_symbol AND B.app_date < NOW() AND B.app_date >= date_add(now(), interval -2 day)) AS closing_price,
                (SELECT C.price FROM us_exchange_latest_price C WHERE C.COIN_TYPE=A.mc_symbol AND C.MARKET_TYPE=A.mc_market_symbol) AS latest_price,
                (SELECT D.mc_kor_name FROM ma_crypto_currency D WHERE D.mc_ticker_symbol = A.mc_symbol ) AS kor_name,
                (SELECT F.mc_eng_name FROM ma_crypto_currency F WHERE F.mc_ticker_symbol = A.mc_symbol ) AS eng_name
                FROM ma_coin_marketsub A WHERE A.mc_use ='Y' AND A.mc_market_symbol = ? ORDER BY mc_display_seq ASC`, market).Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&coininfo.SYMBOL, &coininfo.MARKET, &coininfo.CLOSING_PRICE, &coininfo.LASTEST_PRICE, &coininfo.KOR_NAME, &coininfo.ENG_NAME)
		if err != nil {
			log.Info("*****************", err)
		}
		log.Info(coininfo)
		coinList = append(coinList, coininfo)
	}
	return coinList, nil
}

//내가 선택한 코인 정보 가져오기
func CoinInfo(market, cointype string) (schema.CoinInfo, error) {
	db := db2.GormDB
	coinInfo := schema.CoinInfo{}
	rows, err := db.Raw(`SELECT SUM(A.COIN_AMT) AS COIN_AMT ,SUM(A.KRW_AMT) AS KRW_AMT, MIN(A.CONT_AMT) AS KRW_MIN, MAX(A.CONT_AMT) AS KRW_MAX,
                (SELECT B.price FROM us_exchange_last_ticker B FORCE INDEX(SEARCH_IDX) WHERE B.coin_type = A.COIN_TYPE AND market_type=A.MARKET_TYPE AND B.app_date < NOW() AND B.app_date >= date_add(now(), interval -2 day)) AS closing_price,
                (SELECT C.price FROM us_exchange_latest_price C WHERE C.COIN_TYPE=? AND C.MARKET_TYPE=?) AS latest_price
                FROM us_exchange_contract A FORCE INDEX(SEARCH_IDX)
                WHERE A.COIN_TYPE=? AND A.MARKET_TYPE=? AND A.CONT_DT > DATE_ADD(NOW(3), INTERVAL -24 HOUR)`, cointype, market, cointype, market).Rows()

	if rows.Next() {
		rows.Scan(&coinInfo.COIN_AMT, &coinInfo.KRW_AMT, &coinInfo.KRW_MIN, &coinInfo.KRW_MAX, &coinInfo.CLOSING_PRICE, &coinInfo.LATEST_PRICE)
	}
	log.Info("coinInfo : ", coinInfo)
	if err != nil {
		return schema.CoinInfo{}, err
	}
	return coinInfo, nil
}

//채결완료 리스트 가져오기
func ContractList(market, cointype string) ([]schema.Contract, error) {
	db := db2.GormDB
	contractList := []schema.Contract{}
	rows, err := db.Raw(`SELECT CONT_AMT, COIN_AMT, DATE_FORMAT(CONT_DT, '%m-%d %H:%i:%s') AS CONTRACT_DT, CONT_GB
        FROM us_exchange_contract FORCE INDEX(SEARCH_IDX2)
        WHERE coin_type=? AND market_type=?
        ORDER BY CONT_NO DESC
        LIMIT 0, 20`, cointype, market).Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var row schema.Contract
		rows.Scan(&row.COIN_AMT, &row.CONT_AMT, &row.CONTRACT_DT, &row.CONT_GB)
		contractList = append(contractList, row)
	}
	return contractList, nil
}

//오더북 Buy리스트 가져오기
func Order_BuyList(market, cointype string) ([]schema.BuyOrder, error) {
	db := db2.GormDB
	buylist := []schema.BuyOrder{}
	rows, err := db.Raw(` SELECT BUY_KRW, fn_altcoin_trim(ROUND(SUM(BUY_COIN), 8)) AS BUY_COIN FROM (
        SELECT A.COIN1_KRW_AMT     AS BUY_KRW,
        (A.TOT_BUY_COIN_AMT - A.TOT_TRADE_COIN) AS BUY_COIN
        FROM us_exchange_buy A FORCE INDEX(SEARCH_IDX4)
        WHERE A.BUY_STA_CD = 'P'  AND A.coin_type=? AND A.market_type=?
        ) T
        WHERE BUY_COIN > 0
        GROUP BY T.BUY_KRW
        ORDER BY T.BUY_KRW DESC
        LIMIT 0, 8`, cointype, market).Rows()
	if err != nil {
		return nil, errors.New("BUY리스트 DB조회 에러 " + err.Error())
	}
	for rows.Next() {
		var buy schema.BuyOrder
		rows.Scan(&buy.BUY_KRW, &buy.BUY_COIN)
		buylist = append(buylist, buy)
	}
	return buylist, nil
}

//오더북 Sell리스트 가져오기
func Order_SellList(market, cointype string) ([]schema.BuyOrder, error) {
	db := db2.GormDB
	selllist := []schema.BuyOrder{}
	rows, err := db.Raw(`SELECT BUY_KRW, fn_altcoin_trim(ROUND(BUY_COIN, 8)) AS  BUY_COIN FROM (
        SELECT BUY_KRW, SUM(BUY_COIN) AS BUY_COIN FROM (
        SELECT A.COIN1_KRW_AMT     AS BUY_KRW,
        (A.TOT_SELL_COIN_AMT - A.TOT_TRADE_COIN) AS BUY_COIN
        FROM us_exchange_sell A FORCE INDEX(SEARCH_IDX3)
        WHERE A.SELL_STA_CD = 'P' AND A.coin_type=? AND A.market_type=?
        ) T
        WHERE BUY_COIN > 0
        GROUP BY T.BUY_KRW
        ORDER BY T.BUY_KRW ASC
        LIMIT 0, 8
        ) T ORDER BY BUY_KRW DESC`, cointype, market).Rows()
	if err != nil {
		return nil, errors.New("SELL리스트 DB조회 에러 " + err.Error())
	}
	for rows.Next() {
		var buy schema.BuyOrder
		err := rows.Scan(&buy.BUY_KRW, &buy.BUY_COIN)
		if err != nil {
			log.Info("**************** err : ", err)
		}
		log.Info(buy)
		selllist = append(selllist, buy)
	}
	return selllist, nil
}

//해당코인 종가 가져오기 (closing prise)
func Closing_Price(market, cointype string) (string, error) {
	db := db2.GormDB
	price := ""
	rows, err := db.Raw(`SELECT price 
								FROM us_exchange_last_ticker FORCE INDEX(SEARCH_IDX) 
								WHERE coin_type = ? 
									AND market_type=? 
									AND app_date <= date_format(CURDATE() - INTERVAL 1 DAY, '%Y/%m/%d') 
								ORDER BY app_date DESC LIMIT 1`, cointype, market).Rows()
	if err != nil {
		return "", err
	}
	if rows.Next() {
		rows.Scan(&price)
	}

	return price, nil
}

//미체결 리스트 가져오기
func TransList(coinType, marketType, memID string) ([]schema.TransInfo, error) {
	db := db2.GormDB
	translist := []schema.TransInfo{}
	rows, err := db.Raw(`SELECT 'SELL' AS TRAN_TYPE, COIN_TYPE, SELL_REG_DT AS REG_DT,COIN1_KRW_AMT AS PRICE, TOT_SELL_COIN_AMT AS AMT, TOT_TRADE_COIN AS TRADE_COIN, TOT_SELL_COIN_AMT-TOT_TRADE_COIN AS COIN_REMAIN   
								FROM us_exchange_sell 
								WHERE coin_type=? AND market_type =? AND mem_id = ?
            			UNION ALL
								SELECT 'BUY' AS TRAN_TYPE, COIN_TYPE, BUY_REG_DT AS REG_DT,COIN1_KRW_AMT AS PRICE, TOT_BUY_COIN_AMT AS AMT, TOT_TRADE_COIN AS TRADE_COIN, TOT_BUY_COIN_AMT-TOT_TRADE_COIN AS COIN_REMAIN  
								FROM us_exchange_buy
								WHERE coin_type=? AND market_type=? AND mem_id = ? ORDER BY REG_DT DESC`, coinType, marketType, memID, coinType, marketType, memID).Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		transinfo := schema.TransInfo{}
		err := rows.Scan(&transinfo.TRAN_TYPE, &transinfo.COIN_TYPE, &transinfo.REG_DT, &transinfo.PRICE, &transinfo.AMT, &transinfo.TRADE_COIN, &transinfo.COIN_REMAIN)
		if err != nil {
			return nil, err
		}
		translist = append(translist, transinfo)
	}
	return translist, nil
}

//미체결 리스트 가져오기
func TransListCom(coinType, marketType, memID string) ([]schema.TransComInfo, error) {
	db := db2.GormDB
	translist := []schema.TransComInfo{}
	rows, err := db.Raw(`SELECT * FROM
								(
									(
										SELECT *,DATE_FORMAT(CONT_DT, '%Y-%m-%d %H:%i:%s') AS TRAN_DT
	        					FROM us_exchange_contract FORCE INDEX(tran2_idx) WHERE BUY_MEM_ID =? AND MARKET_TYPE=? AND coin_type= ?  order by cont_no desc limit 100
	        				)
	        				UNION
	        				(
	        						SELECT *,DATE_FORMAT(CONT_DT, '%Y-%m-%d %H:%i:%s') AS TRAN_DT
	        						FROM us_exchange_contract FORCE INDEX(tran_idx) WHERE SELL_MEM_ID =? AND MARKET_TYPE=? AND coin_type= ? order by cont_no desc limit 100
	        				)
								) A ORDER BY A.TRAN_DT DESC LIMIT 100`, memID, marketType, coinType, memID, marketType, coinType).Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		transinfo := schema.TransComInfo{}
		err := rows.Scan(&transinfo.TRAN_DT, &transinfo.CONT_NO, &transinfo.COIN_TYPE, &transinfo.SELL_MEM_ID, &transinfo.SELL_REG_DT, &transinfo.BUY_MEM_ID, &transinfo.BUY_REG_DT, &transinfo.CONT_AMT, &transinfo.COIN_AMT,
			&transinfo.KRW_AMT, &transinfo.COIN_FEE, &transinfo.KRW_FEE, &transinfo.SELL_RATE, &transinfo.BUY_RATE, &transinfo.SELL_BEF_COIN_AMT, &transinfo.SELL_BEF_KRW_AMT,
			&transinfo.BUY_BEF_COIN_AMT, &transinfo.BUY_BEF_KRW_AMT, &transinfo.CONT_DT, &transinfo.SELL_REM, &transinfo.BUY_REM, &transinfo.CONT_GB, &transinfo.MARKET_TYPE)
		if err != nil {
			return nil, err
		}
		translist = append(translist, transinfo)
	}
	return translist, nil
}
