package process

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

func (svc *LinuxProcessService) StartLinuxProcess() error {

	svc.Logger.Debug("[StartLinuxProcess]")
	tmp := strings.Split(svc.Parameters.URL, "/")
	processName := tmp[len(tmp)-1]
	script := fmt.Sprintf("wget %s; chmod 774 %s; ./%s", svc.Parameters.URL, processName, processName)

	err := runOnRemote(svc.Parameters.Username, svc.Parameters.Password, svc.Parameters.Host, script)
	if err != nil {
		svc.Logger.Info(fmt.Sprintf("[StartLinuxProcess]: error while running remote command: %s", err.Error()))
		return err
	}
	return nil

}

func (svc *LinuxProcessService) ObserveLinuxProcess() (int64, error) {

	script := fmt.Sprintf("pgrep -f %s | head -n 1", svc.Parameters.Name)

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

func (svc *LinuxProcessService) KillLinuxProcess(process_pid int64) error {

	script := fmt.Sprintf("kill %d", process_pid)
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
