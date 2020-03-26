package http

import (
	"gtask/internal/task"
	"gtask/pkg/restful"
	filter "gtask/pkg/restful/filters"
	"gtask/pkg/restful/parser"
	"net/http"
)

var StartTaskHandler = restful.NewHandler(
	func(w http.ResponseWriter, r *http.Request, params map[string]interface{}, resp *restful.Response) {
		task.Ins().Loop()
	},
	restful.HandlerOpts{
		ParseFunc: parser.JsonParser,
		Filters: []restful.Filter{
			&filter.CheckAdmin{},
		},
	},
)

var StopTaskHandler = restful.NewHandler(
	func(w http.ResponseWriter, r *http.Request, params map[string]interface{}, resp *restful.Response) {
		task.Ins().Stop()
	},
	restful.HandlerOpts{
		ParseFunc: parser.JsonParser,
		Filters: []restful.Filter{
			&filter.CheckAdmin{},
		},
	},
)

var CreateJobHandler = restful.NewHandler(
	func(w http.ResponseWriter, r *http.Request, params map[string]interface{}, resp *restful.Response) {
		if task.Ins().AddJob(params["key"].(string)) == nil {
			resp.UseError("TASK.10000")
		}
	},
	restful.HandlerOpts{
		MakeErrorFunc: func(err *restful.Errors) {
			err.NewError("TASK.10000", "create job error")
		},
		ParseFunc: parser.JsonParser,
		Filters: []restful.Filter{
			&filter.CheckAdmin{},
			&filter.CheckParams{Params: map[string]interface{}{
				"key": filter.FieldString().SetLength(1, 40),
			},
			},
		},
	},
)

var RunJobHandler = restful.NewHandler(
	func(w http.ResponseWriter, r *http.Request, params map[string]interface{}, resp *restful.Response) {
		job := task.Ins().GetJob(params["key"].(string))
		if job == nil {
			resp.UseError("TASK.10001")
			return
		}
		job.Run()
	},
	restful.HandlerOpts{
		MakeErrorFunc: func(err *restful.Errors) {
			err.NewError("TASK.10001", "run job error")
		},
		ParseFunc: parser.JsonParser,
		Filters: []restful.Filter{
			&filter.CheckAdmin{},
			&filter.CheckParams{Params: map[string]interface{}{
				"key": filter.FieldString().SetLength(1, 40),
			},
			},
		},
	},
)

var CreateProcessorHandler = restful.NewHandler(
	func(w http.ResponseWriter, r *http.Request, params map[string]interface{}, resp *restful.Response) {
		job := task.Ins().GetJob(params["key"].(string))
		if job != nil {
			job.AddProcessor(
				params["code"].(string),
				int64(params["trigger"].(float64)),
				params["bReset"].(bool),
				params["bLoop"].(bool),
				params["bExit"].(bool),
			)
		}
	},
	restful.HandlerOpts{
		ParseFunc: parser.JsonParser,
		Filters: []restful.Filter{
			&filter.CheckAdmin{},
			&filter.CheckParams{Params: map[string]interface{}{
				"key":     filter.FieldString().SetLength(1, 40),
				"code":    filter.FieldString().SetLength(1, 1024*1024),
				"trigger": filter.FieldFloat64(),
				"bReset":  filter.FieldBool(),
				"bLoop":   filter.FieldBool(),
				"bExit":   filter.FieldBool(),
			},
			},
		},
	},
)

var DeleteJobHandler = restful.NewHandler(
	func(w http.ResponseWriter, r *http.Request, params map[string]interface{}, resp *restful.Response) {
		task.Ins().RemoveJob(params["key"].(string))
	},
	restful.HandlerOpts{
		ParseFunc: parser.JsonParser,
		Filters: []restful.Filter{
			&filter.CheckAdmin{},
			&filter.CheckParams{Params: map[string]interface{}{
				"key": filter.FieldString().SetLength(1, 40),
			},
			},
		},
	},
)

var UpdateTaskHandler = restful.NewHandler(
	func(w http.ResponseWriter, r *http.Request, params map[string]interface{}, resp *restful.Response) {
		task.Ins().Reset()
	},
	restful.HandlerOpts{
		ParseFunc: parser.JsonParser,
		Filters: []restful.Filter{
			&filter.CheckAdmin{},
		},
	},
)

var UpdateJobHandler = restful.NewHandler(
	func(w http.ResponseWriter, r *http.Request, params map[string]interface{}, resp *restful.Response) {
		job := task.Ins().GetJob(params["key"].(string))
		if job != nil {
			job.Reset()
		}
	},
	restful.HandlerOpts{
		ParseFunc: parser.JsonParser,
		Filters: []restful.Filter{
			&filter.CheckAdmin{},
			&filter.CheckParams{Params: map[string]interface{}{
				"key": filter.FieldString().SetLength(1, 40),
			},
			},
		},
	},
)

var GetJobHandler = restful.NewHandler(
	func(w http.ResponseWriter, r *http.Request, params map[string]interface{}, resp *restful.Response) {
		job := task.Ins().GetJob(params["key"].(string))
		if job != nil {
			resp.Success(map[string]interface{}{"processorsCount": job.GetProcessorCount(), "state": job.GetState()})
			return
		}
		resp.UseError("TASK.10006")
	},
	restful.HandlerOpts{
		MakeErrorFunc: func(err *restful.Errors) {
			err.NewError("TASK.10006", "get job error")
		},
		ParseFunc: parser.JsonParser,
		Filters: []restful.Filter{
			&filter.CheckAdmin{},
			&filter.CheckParams{Params: map[string]interface{}{
				"key": filter.FieldString().SetLength(1, 40),
			},
			},
		},
	},
)
