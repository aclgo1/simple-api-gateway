package auth

type ErrInvalidToken struct {
}

func (e ErrInvalidToken) Error() string {
	return "invalid token"
}

type ErrRevogedToken struct {
}

func (e ErrRevogedToken) Error() string {
	return "revoged token"
}

type ErrSessionExpiredOrLoginNewDisp struct {
}

func (e ErrSessionExpiredOrLoginNewDisp) Error() string {
	return "session expired or login in new dispositivy"
}
