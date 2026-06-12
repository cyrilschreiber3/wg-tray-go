package notifications

import (
	"log/slog"

	"github.com/cyrilschreiber3/wg-tray-go/assets"
	"github.com/gen2brain/beeep"
)

var notificationEnabled bool

func init() {
	beeep.AppName = "Wg Tray Go"
}

func InitNotifications(enabled bool) {
	notificationEnabled = enabled
}

func Notify(title, message string) {
	if notificationEnabled {
		err := beeep.Notify(title, message, assets.IconByte)
		if err != nil {
			slog.Error("Error sending notification", slog.String("title", title), slog.String("message", message), slog.Any("error", err))
		}
	}
}
