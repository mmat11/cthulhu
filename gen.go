package tgbot

//go:generate mockgen -destination mock/bot.go -package mock -mock_names Service=BotService tg.bot/bot Service
