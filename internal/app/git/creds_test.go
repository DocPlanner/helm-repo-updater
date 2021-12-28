package git

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/argoproj-labs/argocd-image-updater/ext/git"
)

const validGitCredentialsEmail = "test-user@docplanner.com"
const validGitCredentialsUsername = "test-user"
const validGitCredentialsPassword = "test-password"
const validGitCredentialsSSHKey = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAnYHxlioxx8fqeeqjoR8hhaV2akOEI83Pn6alEFge+Fh3OKswWFIj
a0NRj/M6ppHiKPABA5jhzgJUx2WTWINUyAMHHiGSsPR5gSxojLhDcxrYkfnt0Byn4aKxhG
Z0QB8vvj3e7x/7Nb6bL/AzbWtk8TNUC0IeIfTRxv7zdQrftMHgALSM1yeDKEOcOY3TOtY/
Xim3ZQ/dIyWfO6YbjF8jXyYINB4sQNI4bqWhwZiJ7shG2tHBo59BeCMDu2lfbRAYgIcs4T
iytZ60qhiBF7YWkK05a67DFHQp1pPtWq7HSygUEnjb+34I3ZfYobF94itKFuATsCRRY6I7
Ro/rRLyTqQAAA9iN09F3jdPRdwAAAAdzc2gtcnNhAAABAQCdgfGWKjHHx+p56qOhHyGFpX
ZqQ4Qjzc+fpqUQWB74WHc4qzBYUiNrQ1GP8zqmkeIo8AEDmOHOAlTHZZNYg1TIAwceIZKw
9HmBLGiMuENzGtiR+e3QHKfhorGEZnRAHy++Pd7vH/s1vpsv8DNta2TxM1QLQh4h9NHG/v
N1Ct+0weAAtIzXJ4MoQ5w5jdM61j9eKbdlD90jJZ87phuMXyNfJgg0HixA0jhupaHBmInu
yEba0cGjn0F4IwO7aV9tEBiAhyzhOLK1nrSqGIEXthaQrTlrrsMUdCnWk+1arsdLKBQSeN
v7fgjdl9ihsX3iK0oW4BOwJFFjojtGj+tEvJOpAAAAAwEAAQAAAQBicq0RAhCZYbCCQZHD
DJVEVracFtVKF8MVc/CqNZot+gWSyxVtrvFqgupBAnN/V6G3msPXfsBspnJdK3UclwHv/k
x9ndh1eGlVvu8ePbITCQ2iuEfXk4Gve6RfMDarOZL64ussJZ476oZPQWCznLO8Oyvl2Y7C
BKb2Lbb4SjKnZKe3OI8Dw4pKLh9kxWJNJKaIZWU6MwQPSzjTgJgwpF0QUpGqIx2IXC47hW
h+onKD7AvMcpg4SwZ0JAxuEKPXXvixrDXGNC8RufqUl73q3XZg4wt9W37hlH6g+QuuCgN3
iyc1dVEWvs5QfNW5d0SUThPDLfTs5PbEKGzJGpn4KHD9AAAAgDhyOL3o7CuKYX+TLwBZpO
WEaKxOnGhmqkBDmB7IjqaSV9Wdeej9PL0DL4kHZdnM1PuoM4Nx6cXNSAQl/3kKGULacK5D
fidfGVSi08oIH/RSxrIHMLK0/VNa5gfo0o2fhEOIDFot7d59dAjl/+NgGTy2+GvQzJFNsK
SV+0/CVKdGAAAAgQDPrk5JIJ4uXlvWcDyumO1MUu6Gr2E1XRyvSJPsumAvPPdr5kZidJUL
PyZ+6xr/euwuAa8SSWwvJOiD4DmfNl45KXLK0nQRhseboT4DyjA4uPeh2Ha3vsxmMXaXr0
S393FXrQpI/BMoYnravxGf8Ohn6XrbmTNr2nOioE74GN+aHwAAAIEAwidDbzBnwiHnEBTk
4b8e4PN+8z3iHMtYYxFaisF+YJ1zSW/d5wLuNbfQWuPFYnc375K4LRTSYVfhm89H3LYklR
XXXBY3wcmShQeDmVDTmMulJBS9o6Xg29cc6wh6Z8+QOLbsNJEWtsNs/z7q8tzAnhypLqXB
mJDUjAc+gxmMqTcAAAAdeG9hbkBBZG1pbnMtTWFjQm9vay1Qcm8ubG9jYWwBAgMEBQY=
-----END OPENSSH PRIVATE KEY-----`
const validGitRepoSSHURL = "git@github.com:kubernetes/kubernetes.git"
const validGitRepoHTTPSURL = "https://github.com/kubernetes/kubernetes.git"
const invalidGitRepoURL = "github.com/kubernetes/kubernetes.git"

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestNewCredsSSHURLSShPrivKey(t *testing.T) {

	g := Credentials{
		Email:      validGitCredentialsEmail,
		SSHPrivKey: validGitCredentialsSSHKey,
	}

	repoURL := validGitRepoSSHURL

	creds, err := g.NewCreds(repoURL)

	if err != nil {
		log.Fatal(err)
	}

	expectedCreds := git.NewSSHCreds(validGitCredentialsSSHKey, "", true)

	if !(reflect.DeepEqual(creds, expectedCreds)) {
		log.Fatalf("creds and expectedCreds are equal: %t\n", reflect.DeepEqual(creds, expectedCreds))
	}
}

func TestNewCredsHTPPSURLUsernamePassword(t *testing.T) {

	g := Credentials{
		Email:    validGitCredentialsEmail,
		Username: validGitCredentialsUsername,
		Password: validGitCredentialsPassword,
	}

	repoURL := validGitRepoHTTPSURL

	creds, err := g.NewCreds(repoURL)

	if err != nil {
		log.Fatal(err)
	}

	expectedCreds := git.NewHTTPSCreds(g.Username, g.Password, "", "", true, "")

	if !(reflect.DeepEqual(creds, expectedCreds)) {
		log.Fatalf("creds and expectedCreds are equal: %t\n", reflect.DeepEqual(creds, expectedCreds))
	}
}

func TestNewCredsSSHURLWithoutSShPrivKey(t *testing.T) {

	g := Credentials{
		Email:      validGitCredentialsEmail,
		SSHPrivKey: "",
	}

	repoURL := validGitRepoSSHURL

	_, err := g.NewCreds(repoURL)

	expectedError := fmt.Errorf(
		"sshPrivKey not provided for authenticatication to repository %s",
		repoURL,
	)

	if err.Error() != expectedError.Error() {
		t.Fatalf("The error obtained %s is not the expected %s", err.Error(), expectedError.Error())
	}
}

func TestNewCredsHTPPSURLWithoutUsernameWithPassword(t *testing.T) {

	g := Credentials{
		Email:    validGitCredentialsEmail,
		Username: "",
		Password: validGitCredentialsPassword,
	}

	repoURL := validGitRepoHTTPSURL

	_, err := g.NewCreds(repoURL)

	expectedError := fmt.Errorf(
		"no value provided for username and password for authentication to repository %s",
		repoURL,
	)

	if err.Error() != expectedError.Error() {
		t.Fatalf("The error obtained %s is not the expected %s", err.Error(), expectedError.Error())
	}
}

func TestNewCredsHTPPSURLWitUsernameWithoutPassword(t *testing.T) {

	g := Credentials{
		Email:    validGitCredentialsEmail,
		Username: validGitCredentialsUsername,
		Password: "",
	}

	repoURL := validGitRepoHTTPSURL

	_, err := g.NewCreds(repoURL)

	expectedError := fmt.Errorf(
		"no value provided for username and password for authentication to repository %s",
		repoURL,
	)

	if err.Error() != expectedError.Error() {
		t.Fatalf("The error obtained %s is not the expected %s", err.Error(), expectedError.Error())
	}
}

func TestNewCredsInvalidURL(t *testing.T) {

	g := Credentials{
		Email:    validGitCredentialsEmail,
		Username: validGitCredentialsUsername,
		Password: validGitCredentialsPassword,
	}

	repoURL := invalidGitRepoURL

	_, err := g.NewCreds(repoURL)

	expectedError := fmt.Errorf(
		"unknown repository type for git repository URL %s",
		repoURL,
	)

	if err.Error() != expectedError.Error() {
		t.Fatalf("The error obtained %s is not the expected %s", err.Error(), expectedError.Error())
	}
}
