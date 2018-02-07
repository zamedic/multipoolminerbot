package monitor

import (
	"github.com/zamedic/multipoolminerbot/user"
	"github.com/zamedic/multipoolminerbot/multipoolminer"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"time"
	"strconv"
	"github.com/zamedic/telegram"
	"fmt"
)

type Service interface {
}
type service struct {
	userStore        user.Store
	multiPoolService multipoolminer.Service
	telegramService  telegram.Service
}

func NewService(userStore user.Store,
	multiPoolService multipoolminer.Service,
	telegramService telegram.Service) Service {

	s := &service{telegramService:telegramService,multiPoolService:multiPoolService,userStore:userStore}
	go func() {
		s.startMonitor()
	}()
	return s
}

func (s *service) startMonitor() {
	for {
		s.userStore.Iterate(s.processRecord)
		time.Sleep(2 * time.Minute)
	}
}

func (s service) processRecord(output *dynamodb.ScanOutput, last bool) bool {
	items := []user.Multipooluser{}
	err := dynamodbattribute.UnmarshalListOfMaps(output.Items, &items)
	log.Println(err)
	for _, item := range items {
		resp, err := s.multiPoolService.CheckAPIToken(item.Apikey)
		if err != nil {
			log.Println(err)
			continue
		}
		err = s.processRecords(resp, item)
		if err != nil {
			log.Println(err)
		}
	}
	return last
}

func (s service) processRecords(records multipoolminer.MiningPoolHubRespose, user user.Multipooluser) error {
	for _, record := range records {
		t, err := strconv.ParseInt(*record.LastSeen, 10, 64)
		if err != nil {
			return err
		}
		tm := time.Unix(t, 0)
		duration := time.Since(tm)
		log.Println(duration.Seconds())
		if duration.Minutes() > 5 {
			id, err := strconv.ParseInt(user.Userid, 10, 64)
			if err != nil {
				return err
			}
			s.telegramService.SendMessage(id, fmt.Sprintf("Worker %v has not been seen in %v minutes", record.Workername, duration.Minutes()), 0)
		}
	}
	return nil
}
