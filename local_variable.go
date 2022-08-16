package gocc

type LocalVariable struct {
	next   *LocalVariable // 次の変数かnil
	name   []rune         // 変数の名前
	len    int            // 変数名の長さ...goではいらないかも
	offset int            // frame pointer(x29)からの距離
}
