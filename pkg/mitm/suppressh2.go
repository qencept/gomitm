package mitm

func suppressH2(original []string) []string {
	for _, a := range original {
		if a == "h2" {
			return []string{"http/1.1"}
		}
	}
	return original
}
