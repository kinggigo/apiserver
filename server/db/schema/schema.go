package schema

import "database/sql"

//메인 마켓 리스트
type MarketList struct {
	Symbol             string
	Market             string
	Closing_price      string
	Kor_name           string
	Time_closing_price sql.NullString `json:"time_closing_price"`
	Eng_name           string
	Board_id           string
	Latest_price       sql.NullString `json:"lastest_price"`
	Amt                sql.NullInt64  `json:"amt"`
}

// NO ID MA_CD MS_CD TITLE CONTENT MGN_ID MOD_MGN_ID NAME VIEW_CNT FILE_NAME FILE_PATH REG_DT MOD_DT USE_YN LANG_TYPE
//공지사항, 뉴스 투자정보 사항
type Board struct {
	NO         string
	ID         string
	MA_CD      string
	MS_CD      string
	TITLE      string
	CONTENT    string
	MGN_ID     sql.NullString
	MOD_MGN_ID sql.NullString
	NAME       sql.NullString
	VIEW_CNT   sql.NullInt64
	FILE_NAME  sql.NullString
	FILE_PATH  sql.NullString
	REG_DT     sql.NullString
	MOD_DT     sql.NullString
	USE_YN     sql.NullString
	LANG_TYPE  sql.NullString
}

//coinList용 struct
type CoinListInfo struct {
	SYMBOL        string
	MARKET        string
	CLOSING_PRICE sql.NullString
	LASTEST_PRICE sql.NullString
	KOR_NAME      string
	ENG_NAME      string
}

//코인하나당 정보
type CoinInfo struct {
	COIN_AMT      string
	KRW_AMT       string
	KRW_MIN       string
	KRW_MAX       string
	CLOSING_PRICE sql.NullString
	LATEST_PRICE  sql.NullString
	DIF           string
}

//체결완료 정보
type Contract struct {
	CONT_AMT    string
	COIN_AMT    string
	CONTRACT_DT string
	CONT_GB     string
}

//오더북 정보
type BuyOrder struct {
	BUY_KRW  string
	BUY_COIN string
}

//게시판 정보
type BoardView struct {
	ID       string
	MA_CD    string
	MS_CD    string
	TITLE    string
	CONTENT  string
	NAME     string
	VIEW_CNT string
	REG_DT   string
}

//미체결 정보
type TransInfo struct {
	TRAN_TYPE   string
	COIN_TYPE   string
	REG_DT      string
	PRICE       string
	AMT         string
	TRADE_COIN  string
	COIN_REMAIN string
}

//체결 정보
type TransComInfo struct {
	TRAN_DT           string
	CONT_NO           string
	COIN_TYPE         string
	SELL_MEM_ID       string
	SELL_REG_DT       string
	BUY_MEM_ID        string
	BUY_REG_DT        string
	CONT_AMT          string
	COIN_AMT          string
	KRW_AMT           string
	COIN_FEE          string
	KRW_FEE           string
	SELL_RATE         string
	BUY_RATE          string
	SELL_BEF_COIN_AMT string
	SELL_BEF_KRW_AMT  string
	BUY_BEF_COIN_AMT  string
	BUY_BEF_KRW_AMT   string
	CONT_DT           string
	SELL_REM          sql.NullString
	BUY_REM           sql.NullString
	CONT_GB           string
	MARKET_TYPE       string
}

//체결 리턴 목록
type TransComRet struct {
	BUY_MEM_ID    string
	SELL_MEM_ID   string
	CONT_NO       string
	SEARCH_GB     string
	TRAN_GN       string
	TRANSFER_DATE string
	COIN_AMT      string
	KRW_AMT       string
	CONTRACT_AMT  string
	FEE           string
}
