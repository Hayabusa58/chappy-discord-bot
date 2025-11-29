package util

// 文字列 s を n文字単位のチャンクに分割して返す
func SplitCharsByN(s string, n int) []string {
	chunkSize := n
	var chunks []string
	runes := []rune(s)
	for start := 0; start < len(runes); start += chunkSize {
		end := start + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[start:end]))
	}
	return chunks
}
