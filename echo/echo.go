package echo

import (
	"regexp"

	"github.com/labstack/echo/v4"
)

func Skipper(pathToSkips []string) func(c echo.Context) bool {

	return func(c echo.Context) bool {

		path := c.Request().URL.Path

		for _, pattern := range pathToSkips {
			re := regexp.MustCompile(pattern)

			isMatch := re.MatchString(path)

			if isMatch {
				return true
			}
		}

		return false
	}
}
