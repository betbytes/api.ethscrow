package user

import (
	"api.ethscrow/models"
	"api.ethscrow/utils"
	session2 "api.ethscrow/utils/session"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func PublicKey(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	user.Username = chi.URLParam(r, "Username")
	if user.Username == "" {
		return
	}

	if err := user.GetPublicKey(); err != nil {
		utils.Error(w, http.StatusNotFound, "User not found.")
		return
	}
	utils.JSON(w, http.StatusOK, user)
}

func AllPools(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("claims").(models.User)

	pools, err := user.GetAllPools()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	var resp poolResponse

	for _, pool := range pools {
		if pool.Accepted {
			resp.Active = append(resp.Active, pool)
			continue
		}
		if pool.Bettor == user.Username {
			resp.Active = append(resp.Sent, pool)
		} else {
			resp.Active = append(resp.Inbox, pool)
		}
	}

	utils.JSON(w, http.StatusOK, resp)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	if err := utils.ParseRequestBody(r, user); err != nil || user.Username == "" || user.PublicKey == "" {
		utils.Error(w, http.StatusBadRequest, "Invalid request.")
		return
	}

	if err := user.Create(); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Could not create user.")
		return
	}

	session, _ := session2.Store.Get(r, "session.id")
	session.Values["authenticated"] = true
	session.Values["username"] = user.Username
	session.Values["public_key"] = user.PublicKey
	if err := session.Save(r, w); err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, &utils.BasicData{Data: true})
}
