package cthulhu

//go:generate mockgen -destination mock/bot.go -package mock -mock_names Service=BotService cthulhu/bot Service
//go:generate mockgen -destination mock/store.go -package mock -mock_names Service=StoreService cthulhu/store Service
//go:generate mockgen -destination mock/telegram.go -package mock -mock_names Service=TelegramService cthulhu/telegram Service
//go:generate mockgen -destination mock/metrics.go -package mock -mock_names Service=MetricsService cthulhu/metrics Service
