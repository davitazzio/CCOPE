package emqxprocess

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

func (svc *EmqxLinuxProcessService) StartEmqxProcess() error {

	svc.Logger.Debug("[StartLinuxProcess]")

	command := fmt.Sprintf(`rm emqx-5.7.0*
wget https://www.emqx.com/en/downloads/broker/5.7.0/emqx-5.7.0-ubuntu22.04-amd64.tar.gz
mkdir -p emqx && tar -zxvf emqx-5.7.0-ubuntu22.04-amd64.tar.gz -C emqx
echo '
api_key = {
  bootstrap_file = "etc/default_api_key.conf"
}'>> ~/emqx/etc/emqx.conf
echo %s>> ~/emqx/etc/default_api_key.conf
./emqx/bin/emqx start`, svc.Parameters.BrokerAPIKey)

	err := runOnRemote(svc.Parameters.Username, svc.Parameters.Password, svc.Parameters.Host, command)
	if err != nil {
		svc.Logger.Info(fmt.Sprintf("[StartLinuxProcess]: error while running remote command: %s", err.Error()))
		return err
	}
	return nil
}

func (svc *EmqxLinuxProcessService) ObserveEmqxProcess() (int64, error) {

	script := "~/emqx/bin/emqx pid"

	output, err := combineOutputOnRemote(svc.Parameters.Username, svc.Parameters.Password, svc.Parameters.Host, script)
	if err != nil {
		svc.Logger.Info(fmt.Sprintf("[ObserveLinuxProcess]: error while obtaining process PID: %s", err.Error()))

		return -1, err
	}
	output_str := strings.Split(output, "\n")
	process_pid, err := strconv.ParseInt(output_str[0], 10, 64)
	if err != nil {
		svc.Logger.Info(fmt.Sprintf("[ObserveLinuxProcess]: error while converting String to Integer: %s", err.Error()))
		return -1, err
	}
	svc.Logger.Debug(fmt.Sprint("[ObserveLinuxProcess]: Process PID: %d", process_pid))

	return process_pid, nil

}

func (svc *EmqxLinuxProcessService) StopEmqxProcess(process_pid int64) error {

	script := "./emqx/bin/emqx stop"
	err := runOnRemote(svc.Parameters.Username, svc.Parameters.Password, svc.Parameters.Host, script)
	if err != nil {
		svc.Logger.Info(fmt.Sprintf("[KillLinuxProcess]: error sending KILL signal to remote process: %s", err.Error()))
		return err
	}

	return nil

}

func runOnRemote(user, passwd, host, command string) error {

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(passwd)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", host+":22", sshConfig)
	if err != nil {
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return err
	}
	err = session.Run(command)
	if err != nil {
		client.Close()
		return err
	}

	return nil
}

func combineOutputOnRemote(user, passwd, host, command string) (string, error) {

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(passwd)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", host+":22", sshConfig)
	if err != nil {
		return "", err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return "", err
	}
	output, err := session.CombinedOutput(command)
	if err != nil {
		client.Close()
		return "", err
	}

	return string(output), nil
}
