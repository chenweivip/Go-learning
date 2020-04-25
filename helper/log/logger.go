package log

import "github.com/zbrechave/Go-learning/helper/log/base"

type LoggerCreator func(
	level base.LogLevel)
