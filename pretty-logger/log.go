package prettylogger

// NOTE: Can't import the "asserts" package because it creates a circular dependency
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"

	"github.com/wavly/surf/env"
)

const (
	reset = "\033[0m"

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97
	brightWhite  = 98
)

type Handler struct {
	handler slog.Handler
	buf     *bytes.Buffer
	mutex   *sync.Mutex
}

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

func (handler *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return handler.handler.Enabled(ctx, level)
}

func (handler *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{handler: handler.handler.WithAttrs(attrs), buf: handler.buf, mutex: handler.mutex}
}

func (handler *Handler) WithGroup(name string) slog.Handler {
	return &Handler{handler: handler.handler.WithGroup(name), buf: handler.buf, mutex: handler.mutex}
}

func (handler *Handler) Handle(ctx context.Context, r slog.Record) error {

	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = colorize(darkGray, level)
	case slog.LevelInfo:
		level = colorize(cyan, level)
	case slog.LevelWarn:
		level = colorize(lightYellow, level)
	case slog.LevelError:
		level = colorize(lightRed, level)
	}

	attrs, err := handler.computeAttrs(ctx, r)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return fmt.Errorf("error when marshaling attrs: %w", err)
	}

	// Skip 3 stack frames which is all the logging functions and get to the caller
	_, file, lineNo, ok := runtime.Caller(3)
	if !ok {
		log.Fatalln("failed to retrieve the runtime information")
		return nil
	}
	fileName := path.Base(file)

	fmt.Println(
		colorize(brightWhite, fmt.Sprintf("[%s:%d]", fileName, lineNo)),
		level,
		colorize(white, r.Message),
		colorize(lightGray, string(bytes)),
	)

	return nil
}

func (handler *Handler) computeAttrs(
	ctx context.Context,
	r slog.Record,
) (map[string]any, error) {

	handler.mutex.Lock()
	defer func() {
		handler.buf.Reset()
		handler.mutex.Unlock()
	}()
	if err := handler.handler.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(handler.buf.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}
	return attrs, nil
}

func suppressDefaults(
	next func([]string, slog.Attr) slog.Attr,
) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}

func NewHandler(opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	buf := &bytes.Buffer{}
	return &Handler{
		buf: buf,
		handler: slog.NewJSONHandler(buf, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		}),
		mutex: &sync.Mutex{},
	}
}

// TODO: save the `prod` logs in a file
func GetLogger(opts *slog.HandlerOptions) *slog.Logger {
	if env.MODE == "prod" {
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}
	return slog.New(NewHandler(opts))
}
