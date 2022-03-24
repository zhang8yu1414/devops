package service

import (
	"context"
	"fmt"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
)

type sAnsible struct{}

var insAnsible = sAnsible{}

func Ansible() *sAnsible {
	return &insAnsible
}

func (s *sAnsible) ExecutePlaybook(path string) {
	ansiblePlaybookConnectOptions := &options.AnsibleConnectionOptions{
		Connection: "local",
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Check:     false,
		Forks:     "3",
		Inventory: fmt.Sprintf("%s/hosts", path),
	}

	privilegeEscalationOptions := &options.AnsiblePrivilegeEscalationOptions{
		Become:       true,
		BecomeMethod: "sudo",
	}

	cmd := &playbook.AnsiblePlaybookCmd{
		Playbooks:                  []string{fmt.Sprintf("%s/site.yml", path)},
		ConnectionOptions:          ansiblePlaybookConnectOptions,
		Options:                    ansiblePlaybookOptions,
		PrivilegeEscalationOptions: privilegeEscalationOptions,
	}

	err := cmd.Run(context.TODO())
	if err != nil {
		panic(err)
	}
}
