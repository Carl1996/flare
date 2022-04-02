package appearance

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	FlareData "github.com/soulteary/flare/data"
	FlareModel "github.com/soulteary/flare/model"
	FlareAuth "github.com/soulteary/flare/pkg/auth"
	FlareState "github.com/soulteary/flare/state"
)

func RegisterRouting(router *gin.Engine) {

	router.GET(FlareState.SettingPages.Appearance.Path, FlareAuth.AuthRequired, pageHome)
	router.POST(FlareState.SettingPages.Appearance.Path, FlareAuth.AuthRequired, updateAppearanceOptions)

}

func pageHome(c *gin.Context) {

	render(c)

}

func updateAppearanceOptions(c *gin.Context) {

	type UpdateBody struct {
		OptionTitle  string `form:"title"`
		OptionFooter string `form:"footer"`

		OptionOpenAppNewTab      bool `form:"open-app-newtab"`
		OptionOpenBookmarkNewTab bool `form:"open-bookmark-newtab"`

		OptionShowTitle     bool   `form:"show-title"`
		OptionGreetings     string `form:"greetings"`
		OptionShowDateTime  bool   `form:"show-datetime"`
		OptionShowApps      bool   `form:"show-apps"`
		OptionShowBookmarks bool   `form:"show-bookmarks"`
		HideSettingsButton  bool   `form:"hide-settings-button"`
		HideHelpButton      bool   `form:"hide-help-button"`
		EnableEncryptedLink bool   `form:"enable-encrypted-link"`
		IconMode            string `form:"icon-mode"`
		KeepLetterCase      bool   `form:"keep-letter-case"`

		OptionCustomDay   string `form:"custom-day"`
		OptionCustomMonth string `form:"custom-month"`
	}

	var body UpdateBody
	if c.ShouldBind(&body) != nil {
		c.PureJSON(http.StatusForbidden, "提交数据缺失")
		return
	}

	var update FlareModel.Application
	update.Title = body.OptionTitle
	update.Footer = body.OptionFooter
	update.OpenAppNewTab = body.OptionOpenAppNewTab
	update.OpenBookmarkNewTab = body.OptionOpenBookmarkNewTab
	update.ShowTitle = body.OptionShowTitle
	update.Greetings = body.OptionGreetings
	update.ShowDateTime = body.OptionShowDateTime
	update.ShowApps = body.OptionShowApps
	update.ShowBookmarks = body.OptionShowBookmarks
	update.HideSettingsButton = body.HideSettingsButton
	update.HideHelpButton = body.HideHelpButton
	update.EnableEncryptedLink = body.EnableEncryptedLink
	update.KeepLetterCase = body.KeepLetterCase

	requestIconMode := strings.ToUpper(body.IconMode)
	if requestIconMode != "DEFAULT" && requestIconMode != "FILLING" {
		update.IconMode = "DEFAULT"
	} else {
		update.IconMode = requestIconMode
	}

	FlareData.UpdateAppearance(update)

	render(c)
}

func render(c *gin.Context) {
	options := FlareData.GetAllSettingsOptions()

	IconModeDefault := options.IconMode == "DEFAULT"
	IconModeFilling := options.IconMode == "FILLING"

	c.HTML(
		http.StatusOK,
		"settings.html",
		gin.H{
			"DebugMode":       FlareState.AppFlags.DebugMode,
			"PageInlineStyle": FlareState.GetPageInlineStyle(),

			"PageName":       "Appearance",
			"PageAppearance": FlareState.GetAppBodyStyle(),
			"SettingPages":   FlareState.SettingPages,
			"SettingsURI":    FlareState.RegularPages.Settings.Path,

			"OptionTitle":               options.Title,
			"OptionFooter":              template.HTML(options.Footer),
			"OptionOpenAppNewTab":       options.OpenAppNewTab,
			"OptionOpenBookmarkNewTab":  options.OpenBookmarkNewTab,
			"OptionShowTitle":           options.ShowTitle,
			"OptionGreetings":           options.Greetings,
			"OptionShowDateTime":        options.ShowDateTime,
			"OptionShowApps":            options.ShowApps,
			"OptionShowBookmarks":       options.ShowBookmarks,
			"OptionHideSettingsButton":  options.HideSettingsButton,
			"OptionHideHelpButton":      options.HideHelpButton,
			"OptionEnableEncryptedLink": options.EnableEncryptedLink,
			"OptionIconModeDefault":     IconModeDefault,
			"OptionIconModeFilling":     IconModeFilling,
		},
	)
}
