package telegram

import (
	"go.uber.org/zap"
	"gopkg.in/telebot.v3/middleware"
)

func (t *telegramNotificationService) RegisterHandlers() {
	// Global-scoped middleware:
	t.bot.Use(middleware.Logger())
	t.bot.Use(middleware.AutoRespond())

	t.bot.Handle("/addvalidator", t.addValidatorHandler, OnlyDM())

	zap.L().Info("Starting Telegram Controller", zap.String("user", t.bot.Me.Username))

	go t.bot.Start()
}

func (t *telegramNotificationService) Stop() {
	t.bot.Stop()
}
