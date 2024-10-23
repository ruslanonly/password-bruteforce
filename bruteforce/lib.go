package bruteforce

var dictionary = []rune("abcdefghijklmnopqrstuvwxyz")
var dictionaryLength = len(dictionary);

func intToPassword(n int) string {
	password := make([]rune, 5)
	for i := 4; i >= 0; i-- {
		password[i] = dictionary[n % dictionaryLength]
		n /= dictionaryLength
	}
	return string(password)
}
