package http

import (
	"gtask/internal/session"
	"gtask/pkg/restful"
	filter "gtask/pkg/restful/filters"
	"gtask/pkg/restful/parser"
	"net/http"
)

var GetToken = restful.NewHandler(
	func(w http.ResponseWriter, r *http.Request, params map[string]interface{}, resp *restful.Response) {
		token := session.Ins().GetToken(params["secretKey"].(string))
		resp.Success(map[string]interface{}{"token": token})
	},
	restful.HandlerOpts{
		ParseFunc: parser.JsonParser,
		Filters: []restful.Filter{
			&filter.CheckParams{Params: map[string]interface{}{
				"secretKey": filter.FieldString().SetLength(6, 40),
			},
			},
		},
	},
)
