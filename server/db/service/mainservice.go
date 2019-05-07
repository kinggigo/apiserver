package service

import (
	"database/sql"
	"errors"
	"github.com/kinggigo/secret/server/db"
	"github.com/kinggigo/secret/server/db/schema"
	"time"
)

//메인에서 마켓을 가져오는 로직
func MarketList(market string) ([]schema.MarketList, error) {
	db := db.GormDB
	var coinInfo schema.MarketList
	var coinType string
	maketlist := []schema.MarketList{}
	//var symbol string
	amtdate := time.Now().Format("2006-01-02")
	closedate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	timeclosedate := time.Now().Add(-time.Minute).Format("2006-01-02 15:04:05")

	if db == nil {
		return nil, errors.New("DB가 없다!!")
	}
	rows, err := db.Raw(`SELECT A.mc_symbol AS symbol, A.mc_market_symbol AS market,
					(SELECT B.price FROM us_exchange_last_ticker B FORCE INDEX(SEARCH_IDX) WHERE B.coin_type = A.mc_symbol AND B.market_type=A.mc_market_symbol AND B.app_date = ? ) AS closing_price,
					(SELECT E.price FROM us_exchange_last_ticker_m E FORCE INDEX(SEARCH_IDX) WHERE E.coin_type = A.mc_symbol AND E.market_type=A.mc_market_symbol AND E.app_date = ? ) AS time_closing_price,
					(SELECT D.mc_kor_name FROM ma_crypto_currency D WHERE D.mc_ticker_symbol = A.mc_symbol ) AS kor_name,
					(SELECT F.mc_eng_name FROM ma_crypto_currency F WHERE F.mc_ticker_symbol = A.mc_symbol ) AS eng_name,
					(SELECT H.board_id FROM ma_crypto_currency H WHERE H.mc_ticker_symbol = A.mc_symbol ) AS board_id
               FROM ma_coin_marketsub A WHERE A.mc_use='Y' AND A.mc_market_symbol = ? ORDER BY mc_display_seq ASC`, closedate, timeclosedate, market).Rows()

	//rows , err := db.Raw(`SELECT * FROM tradedb.ma_coin_marketsub `).Rows()    , timeclosedate  //
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&coinInfo.Symbol, &coinInfo.Market, &coinInfo.Closing_price, &coinInfo.Time_closing_price, &coinInfo.Kor_name, &coinInfo.Eng_name, &coinInfo.Board_id)
		//log.Info(symbol)
		coinType = coinInfo.Symbol
		amtrows, err := db.Raw(`SELECT
 				(SELECT C.price FROM us_exchange_latest_price C WHERE C.COIN_TYPE= ? AND C.MARKET_TYPE=?) AS latest_price,
               (SELECT SUM(COIN_AMT) AS COIN_AMT FROM us_exchange_contract G WHERE G.COIN_TYPE=? AND G.MARKET_TYPE=? AND G.CONT_DT > ?) AS amt
                from dual`, market, coinType, amtdate, market, coinType).Rows()
		if err != nil {
			return nil, err
		}
		// , , coinType,amtdate
		var amt sql.NullInt64
		var latest_price sql.NullString
		for amtrows.Next() {
			amtrows.Scan(&latest_price, &amt)
		}
		coinInfo.Amt = amt
		coinInfo.Latest_price = latest_price
		maketlist = append(maketlist, coinInfo)
	}
	//TODO marketStatusWidthTimeLimit2 > 코인별 last_price와 amt를 넣어줘야함.

	if len(maketlist) > 0 {
		return maketlist, nil
	} else {
		return nil, errors.New("등록된 코인이 없습니다.")
	}
}
