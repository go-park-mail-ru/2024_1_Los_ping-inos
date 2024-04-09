package delivery

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"main.go/internal/types"
	"strconv"
	"strings"
	"time"
)

func CreateCSRFToken(SID string, UID types.UserID, tokenExpTime int64) (string, error) {
	h := hmac.New(sha256.New, []byte("8Ke9h5ZGsrNK9BghcWjGphMD9Zy79QM7"))
	data := fmt.Sprintf("%s:%d:%d", SID, UID, tokenExpTime)
	h.Write([]byte(data))
	token := hex.EncodeToString(h.Sum(nil)) + ":" + strconv.FormatInt(tokenExpTime, 10)
	return token, nil
}

func CheckCSRFToken(SID string, UID types.UserID, inputToken string) (bool, error) {
	tokenData := strings.Split(inputToken, ":")
	if len(tokenData) != 2 {
		return false, fmt.Errorf("bad token data")
	}

	tokenExp, err := strconv.ParseInt(tokenData[1], 10, 64)
	if err != nil {
		return false, fmt.Errorf("bad token time")
	}

	if tokenExp < time.Now().Unix() {
		return false, fmt.Errorf("token expired")
	}

	h := hmac.New(sha256.New, []byte("8Ke9h5ZGsrNK9BghcWjGphMD9Zy79QM7"))
	data := fmt.Sprintf("%s:%d:%d", SID, UID, tokenExp)
	h.Write([]byte(data))
	expectedMAC := h.Sum(nil)
	messageMAC, err := hex.DecodeString(tokenData[0])
	if err != nil {
		return false, fmt.Errorf("cand hex decode token")
	}

	return hmac.Equal(messageMAC, expectedMAC), nil
}
