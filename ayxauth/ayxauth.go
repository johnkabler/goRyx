package ayxauth

type AyxSigner struct {
  ConsumerKey string
  ConsumerSecret string
  GalleryUrl string
  Endpoint string
}

func (ayxSigner *AyxSigner) BuildRequest() (string, error) {
    signer := *ayxSigner
    consumerKey := signer.ConsumerKey
    consumerSecret :=
  }


// PercentEncode percent encodes a string according to RFC 3986 2.1.
func PercentEncode(input string) string {
	var buf bytes.Buffer
	for _, b := range []byte(input) {
		// if in unreserved set
		if shouldEscape(b) {
			buf.Write([]byte(fmt.Sprintf("%%%02X", b)))
		} else {
			// do not escape, write byte as-is
			buf.WriteByte(b)
		}
	}
	return buf.String()
}

// shouldEscape returns false if the byte is an unreserved character that
// should not be escaped and true otherwise, according to RFC 3986 2.1.
func shouldEscape(c byte) bool {
	// RFC3986 2.3 unreserved characters
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}
	switch c {
	case '-', '.', '_', '~':
		return false
	}
	// all other bytes must be escaped
	return true
}

func generateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func generateTimestamp() string {
	unixEpochSeconds := time.Now().Unix()
	return strconv.FormatInt(unixEpochSeconds, 10)
}
