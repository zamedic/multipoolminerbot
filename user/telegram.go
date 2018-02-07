package user

import (
	"gopkg.in/telegram-bot-api.v4"
	"github.com/zamedic/telegram"
	"fmt"
	"strconv"
)

type startCommand struct {
	telegramStore   telegram.Store
	telegramService telegram.Service
	userStore       Store
}

func NewStartCommand(telegramStore telegram.Store,
	telegramService telegram.Service,
	userStore Store) telegram.Command {
	return &startCommand{
		telegramStore:   telegramStore,
		telegramService: telegramService,
		userStore:       userStore,
	}

}

func (startCommand) CommandIdentifier() string {
	return "start"
}

func (startCommand) CommandDescription() string {
	return "process will check if you have a multipool token associated with your account, if not, it will start the registration flow"
}

func (s startCommand) Execute(update tgbotapi.Update) {

	response, err := s.userStore.getDetails(strconv.Itoa(update.Message.From.ID))
	if err != nil {
		s.telegramService.SendMessage(update.Message.Chat.ID, fmt.Sprintf("hmmm, we have encoutered an error. %v", err.Error()), update.Message.MessageID)
		return
	}

	if response.Apikey == "" {
		s.telegramStore.SetState(update.Message.From.ID, "START", nil)
		s.telegramService.SendMessage(update.Message.Chat.ID, "Ah, we see you are new. Please reply with you miningpoolhub token or btc ID so we can start monitoring your workers. See https://multipoolminer.io/monitor/ for more details.",
			update.Message.MessageID)
		return
	}
	s.telegramService.SendMessage(update.Message.Chat.ID, "We have you on our records and we are activly monitoring you miningpoolhub workers. please see /help for a list of commands", update.Message.MessageID)
	return

}

//--------------------------------------------------------------//

type addTokenCommandlet struct {
	service         Service
	telegramService telegram.Service
}

func NewAddTokenCommandlet(service Service, telegramService telegram.Service) telegram.Commandlet {
	return &addTokenCommandlet{service: service, telegramService: telegramService}
}

func (*addTokenCommandlet) CanExecute(update tgbotapi.Update, state telegram.State) bool {
	return state.State == "START"
}

func (s *addTokenCommandlet) Execute(update tgbotapi.Update, state telegram.State) {
	api := update.Message.Text
	err := s.service.setToken(strconv.Itoa(update.Message.From.ID), api)
	if err != nil {
		s.telegramService.SendMessage(update.Message.Chat.ID, fmt.Sprintf("We envountred an error adding your token. Error is %v",
			err.Error()), update.Message.MessageID)
	} else {
		s.telegramService.SendMessage(update.Message.Chat.ID, "API Token successfully added.",
			update.Message.MessageID)
	}
}

func (*addTokenCommandlet) NextState(update tgbotapi.Update, state telegram.State) string {
	return ""
}

func (*addTokenCommandlet) Fields(update tgbotapi.Update, state telegram.State) []string {
	return nil
}
