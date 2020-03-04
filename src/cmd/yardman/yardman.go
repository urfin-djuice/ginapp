package main

import (
	"oko/pkg/account"
	"oko/pkg/db"
	"oko/pkg/log"
	"oko/pkg/worker"
)

func main() {
	worker.StartScheduler(handler, "24h")
}

func handler() {
	log.Println("Start cleaning")
	log.Println("\tDelete from sign_up_token...")
	if err := db.GetDB().Where("is_used or expire_at < now()").Delete(&account.SignUpToken{}).Error; err != nil {
		log.Errorln("Fail", err)
	} else {
		log.Println("Done!")
	}
	log.Println("\tDelete from resover_token...")
	if err := db.GetDB().Where("is_used or expire_at < now()").Delete(&account.RecoverToken{}).Error; err != nil {
		log.Errorln("Fail", err)
	} else {
		log.Println("Done!")
	}
	log.Println("Finish it!")
}
