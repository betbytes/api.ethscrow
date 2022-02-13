package user

import (
	"api.ethscrow/controllers/broker"
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
	user := r.Context().Value("user").(models.User)

	pools, err := user.GetAllPools()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	var resp poolResponse

	resp.Active = make([]models.Pool, 0)
	resp.Completed = make([]models.Pool, 0)
	resp.Resolve = make([]models.Pool, 0)
	resp.Sent = make([]models.Pool, 0)
	resp.Inbox = make([]models.Pool, 0)

	for _, pool := range pools {
		if pool.Mediator != user.Username {
			pool.ConflictTempData = nil
		}

		if pool.CallerState == broker.LostState || pool.BetterState == broker.LostState {
			resp.Completed = append(resp.Completed, pool)
		} else if pool.Accepted {
			resp.Active = append(resp.Active, pool)
		} else if pool.Bettor == user.Username {
			resp.Sent = append(resp.Sent, pool)
		} else {
			resp.Inbox = append(resp.Inbox, pool)
		}

		if pool.Mediator == user.Username && ((pool.BetterState == broker.WonState && pool.CallerState == broker.WonState) || pool.CallerState == broker.ConflictState || pool.BetterState == broker.ConflictState) {
			resp.Resolve = append(resp.Resolve, pool)
		}
	}

	utils.JSON(w, http.StatusOK, resp)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	if err := utils.ParseRequestBody(r, user); err != nil || user.Username == "" || user.PublicKey == "" || user.EncPublicKey == "" {
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
	session.Values["enc_public_key"] = user.EncPublicKey
	if err := session.Save(r, w); err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, &utils.BasicData{Data: true})
}
