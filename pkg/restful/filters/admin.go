package filter

import (
	"fmt"
	"gtask/internal/session"
	"net/http"
)

type CheckAdmin struct{}

func (c *CheckAdmin) Processor(r *http.Request, in map[string]interface{}) (out map[string]interface{}, err error) {
	token := r.Header.Get("token")
	if !session.Ins().CheckToken(token) {
		err = fmt.Errorf("filter.check_admin_error")
		return
	}
	out = in
	return
}
