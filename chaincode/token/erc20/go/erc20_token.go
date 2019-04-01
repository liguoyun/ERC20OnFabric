package main

import (
	"fmt"
)

// ERC20Token   implementation
type ERC20Token struct {
	TokenName   string           `json:"TokenName"`
	TokenSymbol string           `json:"TokenSymbol"`
	TokenOwner  string           `json:"TokenOwner"`
	TotalSupply int64            `json:"TotalSupply"`
	Balances    map[string]int64 `json:"Balances"`
}

func NewERC20Token(name string, symbol string, owner string, totalSuplly int64) ERC20Token {
	token := ERC20Token{TokenName: name,
		TokenSymbol: symbol,
		TokenOwner:  owner,
		TotalSupply: totalSuplly,
		Balances:    map[string]int64{owner: totalSuplly}}
	return token
}

func (s *ERC20Token) transfer(from string, to string, amount int64) error {
	//check if from exist
	if _, ok := s.Balances[from]; !ok {
		return fmt.Errorf("No balance of %s exists", from)
	}
	//check if from equal to to
	if from == to {
		return fmt.Errorf("Don't transfer to same person, from:%s to:%s", from, to)
	}
	if _, ok := s.Balances[from]; !ok {
		return fmt.Errorf("No balance of %s exists", from)
	}
	//check if to exist,assign 0 if not yet
	if _, ok := s.Balances[to]; !ok {
		fmt.Println("Balances[", to, "] no exist, init one")
		s.Balances[to] = 0
	}
	//check if amountl less or eq to balance of from
	if s.Balances[from] >= amount {
		s.Balances[from] -= amount
		s.Balances[to] += amount
	} else {
		return fmt.Errorf("The balance of %s is %d, less than the tranferring amount: %d", from, s.Balances[from], amount)
	}
	return nil
}
