package service

type sAnsible struct{}

var insAnsible = sAnsible{}

func Ansible() *sAnsible {
	return &insAnsible
}
