package log

import "../log/base"

type LoggerCreator func(
    level base.LogLevel)
