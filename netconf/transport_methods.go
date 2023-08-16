package netconf

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

// CISCO IOS XR

// dialXR creates a new NETCONF session using an SSH Transport.
func dialXR(target string, user, pass string) (*Session, error) {
	config := &ssh.ClientConfig{User: user, HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{ssh.Password(pass)}}

	var t TransportSSH
	err := t.Dial(target, config)
	if err != nil {
		fmt.Println(err)
		err := t.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	return NewSession(&t), nil
}

// ARISTA EOS

// dialArista creates a new NETCONF session using an SSH Transport.
// Removes complexity of the device types as Arista requires interactive keyboard.
func dialArista(target string, user, pass string) (*Session, error) {
	config := &ssh.ClientConfig{User: user, HostKeyCallback: ssh.InsecureIgnoreHostKey(), Auth: []ssh.AuthMethod{
		ssh.KeyboardInteractive(
			func(user, instruction string, questions []string, echos []bool) ([]string, error) {
				answers := make([]string, len(questions))
				for i := range answers {
					answers[i] = pass
				}

				return answers, nil
			},
		),
	}}

	var t TransportSSH
	err := t.Dial(target, config)
	if err != nil {
		fmt.Println(err)
		err := t.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	return NewSession(&t), nil
}
