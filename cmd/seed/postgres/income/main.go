package main

import (
	"fmt"
	"golang.org/x/net/context"
	"incomster/backend/store/postgres"
	"incomster/config"
	"incomster/core"
	"log"
	"math/rand"
	"time"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load[config.Config]("incomster")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := postgres.Connect(ctx, cfg.Store.Postgres)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}

	store := postgres.NewIncomeStore(db)

	for i := 0; i < 100; i++ {
		amount := rand.Intn(100)
		comment := fmt.Sprintf("Comment %d", amount)
		createdAt := time.Now().Add(time.Duration(i) * 2 * time.Hour)

		income, fail := store.Create(ctx, &core.IncomeCreateInput{
			UserID:    1,
			Amount:    float64(amount),
			Comment:   &comment,
			CreatedAt: &createdAt,
			UpdatedAt: &createdAt,
		})

		if fail != nil {
			log.Printf("create: %v", fail)
		} else {
			log.Printf("create: %v", income)
		}
	}

}
