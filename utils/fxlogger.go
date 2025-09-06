package utils

import (
	"strings"

	"github.com/rs/zerolog"
	"go.uber.org/fx/fxevent"
)

type FXZerologLogger struct {
	Logger zerolog.Logger
}

func strToZerologArray(arr []string) *zerolog.Array {
	zerologArr := zerolog.Arr()
	for _, v := range arr {
		zerologArr.Str(v)
	}

	return zerologArr
}

func (l *FXZerologLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Logger.Trace().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStart hook executing")
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Logger.Error().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("error", e.Err.Error()).
				Msg("OnStart hook failed")
		} else {
			l.Logger.Trace().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStart hook executed")
		}
	case *fxevent.OnStopExecuting:
		l.Logger.Trace().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStop hook executing")
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Logger.Error().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("error", e.Err.Error()).
				Msg("OnStop hook failed")
		} else {
			l.Logger.Trace().
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStop hook executed")
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.Logger.Error().
				Str("type", e.TypeName).
				Array("stacktrace", strToZerologArray(e.StackTrace)).
				Array("moduletrace", strToZerologArray(e.ModuleTrace)).
				Str("modulename", e.ModuleName).
				Str("error", e.Err.Error()).
				Msg("Error encountered while applying options")
		} else {
			l.Logger.Trace().
				Str("type", e.TypeName).
				Array("stacktrace", strToZerologArray(e.StackTrace)).
				Array("moduletrace", strToZerologArray(e.ModuleTrace)).
				Str("modulename", e.ModuleName).
				Msg("Supplied")
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.Trace().
				Str("constructor", e.ConstructorName).
				Array("stacktrace", strToZerologArray(e.StackTrace)).
				Array("moduletrace", strToZerologArray(e.ModuleTrace)).
				Str("modulename", e.ModuleName).
				Str("type", rtype).
				Msg("Provided")
		}

		if e.Err != nil {
			l.Logger.Error().
				Array("stacktrace", strToZerologArray(e.StackTrace)).
				Array("moduletrace", strToZerologArray(e.ModuleTrace)).
				Str("modulename", e.ModuleName).
				Str("error", e.Err.Error()).
				Msg("Error encountered while applying options")
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.Trace().
				Str("decorator", e.DecoratorName).
				Array("stacktrace", strToZerologArray(e.StackTrace)).
				Array("moduletrace", strToZerologArray(e.ModuleTrace)).
				Str("modulename", e.ModuleName).
				Str("type", rtype).
				Msg("Decorated")
		}

		if e.Err != nil {
			l.Logger.Error().
				Array("stacktrace", strToZerologArray(e.StackTrace)).
				Array("moduletrace", strToZerologArray(e.ModuleTrace)).
				Str("modulename", e.ModuleName).
				Str("error", e.Err.Error()).
				Msg("Error encountered while applying options")
		}

	case *fxevent.BeforeRun:
		l.Logger.Trace().
			Str("name", e.Name).
			Str("kind", e.Kind).
			Str("modulename", e.ModuleName).
			Msg("Before run")

	case *fxevent.Run:
		if e.Err != nil {
			l.Logger.Error().
				Str("name", e.Name).
				Str("kind", e.Kind).
				Str("modulename", e.ModuleName).
				Str("error", e.Err.Error()).
				Msg("Error returned")
		} else {
			l.Logger.Trace().
				Str("name", e.Name).
				Str("kind", e.Kind).
				Str("runtime", e.Runtime.String()).
				Str("modulename", e.ModuleName).
				Msg("Run")
		}

	case *fxevent.Invoking:
		l.Logger.Trace().
			Str("function", e.FunctionName).
			Str("modulename", e.ModuleName).
			Msg("Invoking")

	case *fxevent.Invoked:
		if e.Err != nil {
			l.Logger.Error().
				Str("stack", e.Trace).
				Str("function", e.FunctionName).
				Str("modulename", e.ModuleName).
				Str("error", e.Err.Error()).
				Msg("Invoke failed")
		}

	case *fxevent.Stopping:
		l.Logger.Info().
			Str("signal", strings.ToUpper(e.Signal.String())).
			Msg("Stopping")

	case *fxevent.Stopped:
		if e.Err != nil {
			l.Logger.Error().
				Str("error", e.Err.Error()).
				Msg("Stop failed")
		}

	case *fxevent.RollingBack:
		l.Logger.Error().
			Str("error", e.StartErr.Error()).
			Msg("Start failed, rolling back")

	case *fxevent.RolledBack:
		if e.Err != nil {
			l.Logger.Error().
				Str("error", e.Err.Error()).
				Msg("Rollback failed")
		}

	case *fxevent.Started:
		if e.Err != nil {
			l.Logger.Error().
				Str("error", e.Err.Error()).
				Msg("Start failed")
		} else {
			l.Logger.Info().
				Msg("Started pharmacy finder")
		}

	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.Logger.Error().
				Str("error", e.Err.Error()).
				Msg("Custom logger initialization failed")
		} else {
			l.Logger.Trace().
				Str("constructor", e.ConstructorName).
				Msg("Initialized custom fxevent.Logger")
		}
	}
}
