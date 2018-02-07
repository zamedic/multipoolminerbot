package multipoolminer

import (
	"net/http"
	"fmt"
	"encoding/json"
)

//https://multipoolminer.io/monitor/stats.php?address=3GdDe78CzXVWwVB4JaXZViohmtHT1KvpUR



type Service interface {
	CheckAPIToken(token string) (MiningPoolHubRespose, error)
}

type service struct {

}

func NewService() Service{
	return &service{}
}


type MiningPoolHubRespose []struct {
	Address *string
	Workername *string
	LastSeen *string
	Miners *[]Miner
	Profit *string

}

type Miner struct {
	Name *string
	Path *string
	MinerType *[]string
	Active *string
	Algorithm *[]string
	Pool *[]string
	CurrenySpeed *[]float64
	EstimatedSpeed*[]float64
	PID *float64
	BTCDay *float64

}

func (s service)CheckAPIToken(token string) (response MiningPoolHubRespose, err error){

	resp, err := http.Get(fmt.Sprintf("https://multipoolminer.io/monitor/stats.php?address=%v",token))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err =  json.NewDecoder(resp.Body).Decode(&response)
	return response,err
}