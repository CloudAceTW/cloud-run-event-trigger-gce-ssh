package controller

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/CloudAceTW/go-ssh-restart/config"
	"github.com/CloudAceTW/go-ssh-restart/model"
)

func Restart(w http.ResponseWriter, r *http.Request) {
	log.Printf("Req: %s %s %s", r.Host, r.Method, r.URL.Path)

	if r.Header.Get("Authorization") != config.AuthToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sshConnect := model.NewSshConnect(config.VmUser, config.VmIp)
	err := sshConnect.CreateConnect()
	if err != nil {
		log.Printf("cannot ssh vm err: %+v", err)
		if config.EnableRestart {
			VmRestart(w, r)
			return
		}
		http.Error(w, "Gateway Timeout", http.StatusGatewayTimeout)
		return
	}

	err = sshConnect.NewSession()
	if err != nil {
		log.Printf("sshConnect.NewSession err %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	defer sshConnect.Close()

	log.Printf("start to service restart")
	var b bytes.Buffer
	sshConnect.Session.Stdout = &b
	err = sshConnect.Session.Run(config.SshCommand)
	if err != nil {
		log.Printf("session.Run err: %+v", err)
	}
	log.Printf("%s", b.String())
	log.Printf("restart finished")

	SuccessResponse(w, r)
}

func VmRestart(w http.ResponseWriter, r *http.Request) {
	log.Printf("start to VM reset")
	gceInstance := model.NewGceInstance(config.Project, config.Zone, config.Instance)
	restartErr := gceInstance.RestartVM()
	if restartErr != nil {
		log.Printf("restart vm err: %+v", restartErr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if gceInstance.RestartStatus {
		log.Printf("reset finished")
		SuccessResponse(w, r)
		return
	}
}

func SuccessResponse(w http.ResponseWriter, r *http.Request) {
	resp := model.Resp{
		Result: "Status OK",
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&resp)
	if err != nil {
		http.Error(w, "Encoding result to json fail", http.StatusBadRequest)
		return
	}
}
