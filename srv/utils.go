package srv

const myMagicNum = 7

func trunc(s string) string {
	if len(s) > myMagicNum {
		return s[:myMagicNum]
	}

	return s
}
