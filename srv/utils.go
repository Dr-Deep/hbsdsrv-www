package srv

func trunc(s string) string {
	if len(s) > 7 {
		return s[:7]
	}

	return s
}
