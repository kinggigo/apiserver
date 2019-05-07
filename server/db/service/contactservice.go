package service

import (
	"errors"
	"github.com/kinggigo/secret/server/db"
	"github.com/kinggigo/secret/server/db/schema"
	"github.com/labstack/gommon/log"
)

//공지사항, 보도자료 리스트
func BoardList(param schema.BoardParm) ([]schema.Board, error) {
	db := db.GormDB

	boardList := []schema.Board{}
	if db == nil {
		return nil, errors.New("DB가 없다!!")
	}
	//if param.Limit_str == 0 {
	//	param.Limit_str = 3
	//}
	if param.Limit == 0 {
		param.Limit = 3
	}
	if param.Lang_type == "" || len(param.Lang_type) == 0 {
		param.Lang_type = "KOREAN"
	}

	//limitnum , err := strconv.Atoi(param.Limit_str)
	//if err != nil {
	//	log.Info("#####################", err)
	//}

	rows, err := db.Raw(`SELECT SQL_CALC_FOUND_ROWS sub.* FROM(
					SELECT @rownum:=@rownum+1 AS 'NO', a.* FROM us_board a ,(SELECT @rownum:='', @rownum:=0 FROM DUAL) AS b WHERE a.MA_CD = 'BRD' AND a.MS_CD = ? and LANG_TYPE= ? ORDER BY a.ID ASC
				)sub ORDER BY sub.NO DESC LIMIT ?,?`, param.Code, param.Lang_type, param.StartNum, param.Limit).Rows()
	if err != nil {
		return nil, errors.New("DB 조회 에러 : " + err.Error())
	}
	for rows.Next() {
		var board schema.Board
		rows.Scan(&board.NO, &board.ID, &board.MA_CD, &board.MS_CD, &board.TITLE, &board.CONTENT, &board.CONTENT, &board.MGN_ID, &board.MOD_MGN_ID, &board.VIEW_CNT, &board.FILE_NAME, &board.FILE_PATH, &board.REG_DT, &board.MOD_DT, &board.USE_YN, &board.LANG_TYPE)
		//err := rows.Scan(board)
		if err != nil {
			log.Error(err)
		}
		log.Info(board)
		boardList = append(boardList, board)
	}

	// 서비스중 rownum 가져오는 테스트
	return boardList, nil
}

//게시판 번호에 대해 게시판 내용 가져오기
func Board_View(num string) (schema.BoardView, error) {
	db := db.GormDB
	board := schema.BoardView{}
	rows, err := db.Raw(`SELECT a.ID,a.MA_CD,a.MS_CD,a.TITLE,a.CONTENT , a.NAME , a.VIEW_CNT , LEFT(a.REG_DT,10) AS REG_DT FROM us_board a WHERE a.ID = ?`, num).Rows()
	if err != nil {
		log.Error("errrr!!!!!! : ", err)
		return schema.BoardView{}, err
	}
	for rows.Next() {
		rows.Scan(&board.ID, &board.MA_CD, &board.MS_CD, &board.TITLE, &board.CONTENT, &board.NAME, &board.VIEW_CNT, &board.REG_DT)
		log.Info(board)
	}
	return board, nil
}

//게시판 번호 조회수 올리기
func Board_View_Cnt(num string) error {
	db := db.GormDB
	rows, err := db.Raw(`UPDATE us_board b SET b.VIEW_CNT = b.VIEW_CNT + 1  WHERE   b.ID = ?`, num).Rows()
	if err != nil {
		log.Error("errrr!!!!!! : ", err)
		return err
	}
	if rows.Next() {
		log.Info("count up ! num : ", num)
	}
	return nil
}
