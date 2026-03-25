package go_chipay

func getHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"charset":      "utf-8",
	}
}
