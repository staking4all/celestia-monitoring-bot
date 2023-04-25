package telegram

import (
	"fmt"

	"github.com/staking4all/celestia-monitoring-bot/services/models"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

func (t *telegramNotificationService) addValidatorHandler(c telebot.Context) error {
	if len(c.Args()) != 2 {
		t.Send(c.Chat(), "Please try `/add ValidatorName celestiavalcons1XXXXXXX`")
		return nil
	}

	zap.L().Info("add validator", zap.Int64("userID", c.Sender().ID), zap.String("userName", c.Sender().Username), zap.String("name", c.Args()[0]), zap.String("address", c.Args()[1]))
	err := t.mm.Add(c.Sender().ID, models.NewValidator(c.Args()[0], c.Args()[1]))
	if err != nil {
		t.Send(c.Chat(), err.Error())
		return err
	}

	t.Send(c.Chat(), fmt.Sprintf("validator added to monitor list *%s*", c.Args()[1]))

	return nil
}

func (t *telegramNotificationService) removeValidatorHandler(c telebot.Context) error {
	if len(c.Args()) != 1 {
		// TODO: list validators regitered to user
		t.Send(c.Chat(), "Please try `/remove celestiavalcons1XXXXXXX`")
		return nil
	}

	zap.L().Info("remove validator", zap.Int64("userID", c.Sender().ID), zap.String("userName", c.Sender().Username), zap.String("address", c.Args()[0]))
	err := t.mm.Remove(c.Sender().ID, c.Args()[0])
	if err != nil {
		t.Send(c.Chat(), err.Error())
		return err
	}

	t.Send(c.Chat(), fmt.Sprintf("validator removed from monitor list *%s*", c.Args()[0]))

	return nil
}

func (t *telegramNotificationService) statusHandler(c telebot.Context) error {

	return nil
}

func (t *telegramNotificationService) listValidatorHandler(c telebot.Context) error {
	return nil
}
