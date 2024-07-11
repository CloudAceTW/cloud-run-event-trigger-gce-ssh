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
)

const (
	SshKeyIndex  = "sshKey"
	HostKeyIndex = "hostKey"
)
