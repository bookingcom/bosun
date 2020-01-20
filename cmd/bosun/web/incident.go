package web

import (
	"fmt"
	"net/http"

	"bosun.org/cmd/bosun/sched"

	"github.com/MiniProfiler/go/miniprofiler"
	"github.com/kylebrandt/boolq"
)

func ListOpenIncidents(t miniprofiler.Timer, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// TODO: Retune this when we no longer store email bodies with incidents
	list, err := schedule.DataAccess.State().GetAllOpenIncidents()
	if err != nil {
		return nil, err
	}
	suppressor, err := schedule.Silenced()
	if err != nil {
		return nil, fmt.Errorf("failed to get silences: %v", err)
	}
	summaries := []*sched.IncidentSummaryView{}
	filterText := r.FormValue("filter")
	var parsedExpr *boolq.Tree
	parsedExpr, err = boolq.Parse(filterText)
	if err != nil {
		return nil, fmt.Errorf("bad filter: %v", err)
	}
	for _, iState := range list {
		is, err := sched.MakeIncidentSummary(schedule.RuleConf, suppressor, iState)
		if err != nil {
			return nil, err
		}
		match, err := boolq.AskParsedExpr(parsedExpr, is)
		if err != nil {
			return nil, err
		}
		if match {
			summaries = append(summaries, is)
		}
	}
	return summaries, nil
}
