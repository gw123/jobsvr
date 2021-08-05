package jobsvr

import (
	context "context"
	"encoding/json"

	"github.com/gw123/gworker"
)

func (j *Job) UUID() string {
	return j.Uuid
}

func (j *Job) Queue() string {
	return j.Name
}

func (j *Job) Delay() int {
	return int(j.DelaySecond)
}

func (j *Job) Marshal() ([]byte, error) {
	return json.Marshal(j)
}

func (j *Job) JobHandler(ctx context.Context, job gworker.Jobber) error {
	return nil
}
