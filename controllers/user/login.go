package user

import (
	"api.ethscrow/models"
	"api.ethscrow/utils"
	session2 "api.ethscrow/utils/session"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"math/big"
	"net/http"
)

var Logins = make(map[string][]byte)

func RequestChallenge(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	if err := utils.ParseRequestBody(r, user); err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	nonce := make([]byte, 10)
	if _, err := rand.Read(nonce); err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	h := sha256.New()
	sHash := h.Sum(nonce)

	Logins[user.Username] = sHash

	utils.JSON(w, http.StatusOK, nonce)
}

func SubmitChallenge(w http.ResponseWriter, r *http.Request) {
	var body map[string]string

	if err := utils.ParseRequestBody(r, &body); err != nil || body["r"] == "" || body["s"] == "" || body["username"] == "" {
		utils.Error(w, http.StatusInternalServerError, "Invalid")
		return
	}

	user := &models.User{
		Username: body["username"],
	}
	//user.GetData()

	rByte, _ := hex.DecodeString(body["r"])
	sByte, _ := hex.DecodeString(body["s"])
	rVal := new(big.Int).SetBytes(rByte)
	sVal := new(big.Int).SetBytes(sByte)

	key, err := x509.ParsePKIXPublicKey(user.PublicKey)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	valid := ecdsa.Verify(key.(*ecdsa.PublicKey), Logins[user.Username], rVal, sVal)
	delete(Logins, user.Username)

	if valid {

		session, _ := session2.Store.Get(r, "session.id")
		session.Values["authenticated"] = true
		session.Values["username"] = user.Username
		session.Values["public_key"] = user.PublicKey
		if err = session.Save(r, w); err != nil {
			utils.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.JSON(w, http.StatusOK, valid)
	} else {
		utils.Error(w, http.StatusForbidden, "Unauthorized")
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := session2.Store.Get(r, "session.id")
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "Logged out")
}
