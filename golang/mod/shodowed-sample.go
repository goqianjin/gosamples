package main

import (
	"fmt"
)

func main() {
	c := Client{}
	en, err := c.GetItblEntry()
	fmt.Printf("entry: %v, err: %v \n", en, err)
}


func (c Client) GetItblEntry() (entry Entry, err error) {
	//err = c.Get(&entry)
	//getEntryInternal(&entry)
	err = getItblEntryInn(c.GetI, &entry)
	//err = c.GetI.Get(&entry)
	if err != nil {
		fmt.Errorf("%v\n", err)
	}
	return
}

func getItblEntryInn(g GetI, ret interface{}) error {
	return g.Get(ret)
}
func GetItblEntry2() (entry Entry, err error) {
	getEntryInternal(&entry)
	if err != nil {
		fmt.Errorf("%v\n", err)
	}
	return
}

func getEntryInternal(ret interface{}) (err error) {
	ret.(*Entry).Uid = 1234
	return
}

type GetI interface {
	Get(ret interface{}) (err error)
}

type Client struct {
	GetI
}

func (c Client) Get(ret interface{}) (err error) {
	ret.(*Entry).Uid = 1234
	return
}

type Entry struct {
	Uid int
	Name string
	address []string
}