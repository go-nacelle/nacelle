package http

type healthToken string

func (t healthToken) String() string {
	return "http-init"
}
