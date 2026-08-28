package database

import (
	"kickbase/internal/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed populates initial demo users and teams if database is empty
func Seed(db *gorm.DB) error {
	var count int64
	db.Model(&domain.User{}).Count(&count)
	if count > 0 {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 1. Seed Demo Users
	users := []domain.User{
		{
			Username: "admin",
			Password: string(hashedPassword),
			Name:     "System Administrator",
			Role:     "admin",
		},
		{
			Username: "staff",
			Password: string(hashedPassword),
			Name:     "Club Operator Staff",
			Role:     "staff",
		},
		{
			Username: "viewer",
			Password: string(hashedPassword),
			Name:     "Public Viewer",
			Role:     "viewer",
		},
	}

	for _, u := range users {
		db.Where("username = ?", u.Username).FirstOrCreate(&u)
	}

	// 2. Seed Demo Teams
	var teamCount int64
	db.Model(&domain.Team{}).Count(&teamCount)
	if teamCount == 0 {
		t1 := domain.Team{
			Name:        "Persija Jakarta",
			FoundedYear: 1928,
			Address:     "Jl. Pintu Satu Senayan",
			City:        "Jakarta",
			LogoURL:     "https://upload.wikimedia.org/wikipedia/id/5/5e/Logo_Persija.png",
			Version:     1,
		}
		t2 := domain.Team{
			Name:        "Persib Bandung",
			FoundedYear: 1933,
			Address:     "Jl. Sulanjana No. 17",
			City:        "Bandung",
			LogoURL:     "https://upload.wikimedia.org/wikipedia/id/8/80/Persib_Bandung_logo.svg",
			Version:     1,
		}
		db.Create(&t1)
		db.Create(&t2)

		// 3. Seed Demo Players
		p1 := domain.Player{
			TeamID:       t1.ID,
			Name:         "Bambang Pamungkas",
			Height:       178.0,
			Weight:       72.0,
			Position:     "CF",
			Playstyle:    stringPtr("Goal Poacher"),
			JerseyNumber: 20,
			Version:      1,
		}
		p2 := domain.Player{
			TeamID:       t2.ID,
			Name:         "Atep",
			Height:       170.0,
			Weight:       65.0,
			Position:     "LWF",
			Playstyle:    stringPtr("Prolific Winger"),
			JerseyNumber: 7,
			Version:      1,
		}
		db.Create(&p1)
		db.Create(&p2)

		// 4. Seed Demo Match
		m := domain.Match{
			MatchDate:  "2026-09-01",
			MatchTime:  "15:30:00",
			HomeTeamID: t1.ID,
			AwayTeamID: t2.ID,
			Status:     "scheduled",
			Version:    1,
		}
		db.Create(&m)
	}

	return nil
}

func stringPtr(s string) *string {
	return &s
}
