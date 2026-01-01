package main

import (
	"log"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type App struct {
	DB *badger.DB
	
}	

var defaultTTL = 60 * time.Second

func newApp(db *badger.DB) *App {
	return &App{DB: db}
}	

func (a *App) handleGet(key string) (itemVal []byte, err error){
	err =a.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		itemVal, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		log.Printf("Value: %s", string(itemVal))
		return nil
	})

	return itemVal, err
}

func (a *App) handlePut(key string, ttlInSeconds *int, value []byte) error {
	return a.DB.Update(func(txn *badger.Txn) error {
		item := badger.NewEntry([]byte(key), value)
		if ttlInSeconds != nil {
			item.WithTTL(time.Duration(*ttlInSeconds) * time.Second)
		} else {
			item.WithTTL(defaultTTL)
		}
		return txn.SetEntry(item)
	})
}	
