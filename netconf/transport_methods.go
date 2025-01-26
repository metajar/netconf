package netconf

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

// NewNetconfConnection creates a new NETCONF session using an SSH Transport.
func NewNetconfConnection(target string, user, pass string, keyboardAuth bool) (*Session, error) {
	config := &ssh.ClientConfig{}
	if keyboardAuth {
		config = &ssh.ClientConfig{User: user, HostKeyCallback: ssh.InsecureIgnoreHostKey(), Auth: []ssh.AuthMethod{
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
	} else {
		config = &ssh.ClientConfig{User: user, HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Auth: []ssh.AuthMethod{ssh.Password(pass)}}
	}

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
