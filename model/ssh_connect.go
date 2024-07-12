package model

import (
	"context"
	"log"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/CloudAceTW/go-ssh-restart/config"
	"golang.org/x/crypto/ssh"
)

type SshConnect struct {
	Client  *ssh.Client
	Session *ssh.Session
	User    string
	Host    string
}

func NewSshConnect(user, host string) *SshConnect {
	return &SshConnect{
		User: user,
		Host: host,
	}
}

func (sc *SshConnect) CreateConnect() error {
	keys, err := getSshAndHostKey()
	if err != nil {
		log.Printf("getSshKey err: %+v", err)
		return err
	}

	s, err := ssh.ParsePrivateKey(keys[config.SshKeyIndex])
	if err != nil {
		log.Printf("ssh.ParsePrivateKey err: %+v", err)
		return err
	}
	ps, _, _, _, err := ssh.ParseAuthorizedKey(keys[config.HostKeyIndex])
	if err != nil {
		log.Printf("ssh.ParsePublicKey err: %+v", err)
		return err
	}

	clientConfig := &ssh.ClientConfig{
		User: sc.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(s),
		},
		HostKeyCallback:   ssh.FixedHostKey(ps),
		HostKeyAlgorithms: config.HostKeyAlgorithms,
		Timeout:           time.Second * time.Duration(config.SshConnectTimeout),
	}
	var client *ssh.Client
	for i := 1; i <= config.SshRetryTimes; i++ {
		// Connect to ssh server
		client, err = ssh.Dial("tcp", sc.Host, clientConfig)
		if err != nil {
			log.Printf("ssh.Dial err: %+v", err)
			time.Sleep(time.Second * time.Duration(i))
			continue
		} else {
			break
		}
	}
	if err != nil {
		log.Printf("ssh.Dial err: %+v", err)
		return err
	}

	sc.Client = client
	return nil
}

func (sc *SshConnect) NewSession() error {
	session, err := sc.Client.NewSession()
	if err != nil {
		log.Printf("sc.Client.NewSession err: %+v", err)
		return err
	}
	sc.Session = session
	return nil
}

func (sc *SshConnect) Close() {
	clientErr := sc.Client.Close()
	if clientErr != nil {
		log.Printf("sc.Client.Close err: %+v", clientErr)
	}
}

func getSshAndHostKey() (map[string][]byte, error) {
	client, err := NewSecretClient()
	if err != nil {
		return nil, err
	}
	var secretChannelCount = 2
	secretChannel := make(chan ChannelObj, secretChannelCount)
	go getSecretManagerResult(client, config.SshKeyIndex, config.SecretManagerSshKey, secretChannel)
	go getSecretManagerResult(client, config.HostKeyIndex, config.SecretManagerHostKey, secretChannel)

	keys := map[string][]byte{
		config.SshKeyIndex:  nil,
		config.HostKeyIndex: nil,
	}
	for i := 0; i < secretChannelCount; i++ {
		c := <-secretChannel
		if !c.Status {
			log.Printf("getSecretManagerResult for %s err: %+v", c.Key, c.Error)
		}
		keys[c.Key] = c.Result
	}
	err = client.Close()
	if err != nil {
		log.Printf("client.Close err: %+v", err)
	}

	return keys, nil
}

func getSecretManagerResult(client *secretmanager.Client, key, secretName string, c chan ChannelObj) {
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}
	result, err := client.AccessSecretVersion(context.Background(), accessRequest)
	if err != nil {
		c <- ChannelObj{Key: key, Status: false, Error: err}
		return
	}
	c <- ChannelObj{Key: key, Status: true, Result: result.Payload.Data}
}
