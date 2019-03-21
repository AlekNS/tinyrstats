package config

import (
	"testing"

	"github.com/spf13/viper"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDefaultSettingsSpec(t *testing.T) {
	Convey("Given empty viper", t, func(c C) {
		v := viper.New()

		c.Convey("When gettings settings", func(c C) {
			c.So(func() {
				GetSettings(v)
			}, ShouldNotPanic)

			settings := GetSettings(v)

			c.Convey("Then all configs should be filled up by default values", func(c C) {
				c.So(settings, ShouldNotBeNil)
				c.So(settings.Endpoints, ShouldNotBeNil)
				c.So(settings.Logger, ShouldNotBeNil)
				c.So(settings.Tasks, ShouldNotBeNil)
				c.So(settings.Scheduler, ShouldNotBeNil)
			})
		})
	})
}
