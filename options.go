// Copyright confd. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE-confd file.

package libconfd

import (
	"text/template"
)

type Options func(*Config)

func (p *Config) applyOptions(opts ...Options) *Config {
	for _, fn := range opts {
		fn(p)
	}
	return p
}

func WithOnetimeMode() Options {
	return func(opt *Config) {
		opt.Onetime = true
	}
}

func WithIntervalMode() Options {
	return func(opt *Config) {
		opt.Onetime = false
		opt.Watch = false
	}
}

func WithInterval(interval int) Options {
	return func(opt *Config) {
		opt.Interval = interval
	}
}

func WithWatchMode() Options {
	return func(opt *Config) {
		opt.Onetime = false
		opt.Watch = true
	}
}

func WithFuncMap(maps ...template.FuncMap) Options {
	return func(opt *Config) {
		if opt.FuncMap == nil {
			opt.FuncMap = make(template.FuncMap)
		}
		for _, m := range maps {
			for k, fn := range m {
				opt.FuncMap[k] = fn
			}
		}
	}
}

func WithAbsKeyAdjuster(fn func(absKey string) (realKey string)) Options {
	return func(opt *Config) {
		opt.HookAbsKeyAdjuster = fn
	}
}

func WithFuncMapUpdater(funcMapUpdater ...func(m template.FuncMap)) Options {
	return func(opt *Config) {
		opt.FuncMapUpdater = append(opt.FuncMapUpdater, funcMapUpdater...)
	}
}

func WithHookBeforeCheckCmd(fn func(trName, cmd string, err error)) Options {
	return func(opt *Config) {
		opt.HookBeforeCheckCmd = fn
	}
}

func WithHookAfterCheckCmd(fn func(trName, cmd string, err error)) Options {
	return func(opt *Config) {
		opt.HookAfterCheckCmd = fn
	}
}

func WithHookBeforeReloadCmd(fn func(trName, cmd string, err error)) Options {
	return func(opt *Config) {
		opt.HookBeforeReloadCmd = fn
	}
}

func WithHookAfterReloadCmd(fn func(trName, cmd string, err error)) Options {
	return func(opt *Config) {
		opt.HookAfterReloadCmd = fn
	}
}

func WithHookError(fn func(trName string, err error)) Options {
	return func(opt *Config) {
		opt.HookError = fn
	}
}
