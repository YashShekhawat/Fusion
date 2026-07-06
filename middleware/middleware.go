package middleware

import "github.com/YashShekhawat/fusion/drivers"

type Middleware func(drivers.Driver) drivers.Driver
