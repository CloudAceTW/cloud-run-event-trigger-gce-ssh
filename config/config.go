package config

import "os"

var (
	SecretManagerSshKey  = os.Getenv("SECRET_MANAGER_SSH_KEY")
	SecretManagerHostKey = os.Getenv("SECRET_MANAGER_HOST_KEY")
	AuthToken            = os.Getenv("AUTH_TOKEN")
	VmIp                 = os.Getenv("VM_IP")
	VmUser               = os.Getenv("VM_USER")

	Project  = os.Getenv("PROJECT")
	Zone     = os.Getenv("ZONE")
	Instance = os.Getenv("INSTANCE")

	SshCommand = os.Getenv("SSH_COMMAND")

	EnableRestart = os.Getenv("ENABLE_RESTART") == "true"
)

const (
	SshKeyIndex       = "sshKey"
	HostKeyIndex      = "hostKey"
	SshRetryTimes     = 4
	SshConnectTimeout = 4
)

var (
	HostKeyAlgorithms = []string{
		"ssh-ed25519",
		"ssh-rsa",
		"ecdsa-sha2-nistp256"}
)
