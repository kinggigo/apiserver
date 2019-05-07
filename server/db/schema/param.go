package schema

type BoardParm struct {
	Code      string
	Lang_type string //기본 korean - 현재 한글로만 되어있음
	StartNum  int    //가져올 게시판 위치 시작 부분
	Limit     int    //가져올 갯수
}
