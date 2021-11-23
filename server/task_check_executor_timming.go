package server

import (
	"context"
	"fmt"
	"time"

	"github.com/bytebase/bytebase/api"
	"github.com/bytebase/bytebase/common"
	"go.uber.org/zap"
)

func NewTaskCheckTimingExecutor(logger *zap.Logger) TaskCheckExecutor {
	return &TaskCheckTimingExecutor{
		l: logger,
	}
}

type TaskCheckTimingExecutor struct {
	l *zap.Logger
}

func (exec *TaskCheckTimingExecutor) Run(ctx context.Context, server *Server, taskCheckRun *api.TaskCheckRun) (result []api.TaskCheckResult, err error) {
	taskFind := &api.TaskFind{
		ID: &taskCheckRun.TaskID,
	}
	task, err := server.TaskService.FindTask(ctx, taskFind)
	if err != nil {
		return []api.TaskCheckResult{}, common.Errorf(common.Internal, err)
	}
	// Test if it has already passed the earliest allowed time specified by the task
	if time.Now().Before(time.Unix(task.NotBeforeTs, 0)) {
		return []api.TaskCheckResult{
			{
				Status:  api.TaskCheckStatusError,
				Code:    common.TimingTaskNotAllowed,
				Title:   "Execution not allowed",
				Content: fmt.Sprintf("Eearlier than the earliest expeted execution timing: ,%q", time.Unix(task.NotBeforeTs, 0).Format("2006-01-02 15:04")),
			},
		}, nil
	}

	// The following codes is a little awkward, for we only have 'TaskID' here
	pipelineFind := &api.PipelineFind{
		ID: &task.ID,
	}
	pipeline, err := server.PipelineService.FindPipeline(ctx, pipelineFind)
	if err != nil {
		return []api.TaskCheckResult{}, common.Errorf(common.Internal, err)
	}
	env := pipeline.StageList[0].Environment
	if env == nil {
		return []api.TaskCheckResult{}, common.Errorf(common.Internal, err)
	}
	policy, err := server.PolicyService.GetAllowedWindowPolicy(ctx, env.ID)
	if err != nil {
		return []api.TaskCheckResult{}, common.Errorf(common.Internal, err)
	}
	cron, err := api.GetAllowedWindowCronParser().Parse(string(policy.Cron))
	if err != nil {
		return []api.TaskCheckResult{}, common.Errorf(common.Internal, err)
	}
	nextTime := cron.Next(time.Now().Truncate(time.Second))
	// Test if it is within the allowed window specified by the environment
	if nextTime.Before(time.Now()) {
		return []api.TaskCheckResult{
			{
				Status:  api.TaskCheckStatusSuccess,
				Code:    common.TimingEnvironmentNotAllowed,
				Title:   "Execution not allowed",
				Content: fmt.Sprintf("Not in Allowed Window, %q", policy.Cron),
			},
		}, nil
	}

	return []api.TaskCheckResult{
		{
			Status:  api.TaskCheckStatusSuccess,
			Code:    common.Ok,
			Title:   "OK",
			Content: fmt.Sprintf("Good to be executed"),
		},
	}, nil
}
