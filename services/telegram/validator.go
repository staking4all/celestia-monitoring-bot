package telegram

import (
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

func (t *telegramNotificationService) addValidatorHandler(c telebot.Context) error {
	// t.mm.Add(c.Sender().ID, models.NewValidator("name", "address"))
	zap.L().Debug("add validator", zap.Int64("userID", c.Sender().ID), zap.String("userName", c.Sender().Username), zap.String("address", "address"))
	return nil
}
