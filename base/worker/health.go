package worker

type healthToken string

func (t healthToken) String() string {
	return "worker-init"
}
