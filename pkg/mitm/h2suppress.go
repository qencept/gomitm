package mitm

func h2suppress(original []string) []string {
	for i, a := range original {
		if a == "h2" {
			original[i] = "http/1.1"
		}
	}
	return original
}
