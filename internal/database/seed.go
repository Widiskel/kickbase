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
		{Username: "admin", Password: string(hashedPassword), Name: "System Administrator", Role: "admin"},
		{Username: "staff", Password: string(hashedPassword), Name: "Club Operator Staff", Role: "staff"},
		{Username: "viewer", Password: string(hashedPassword), Name: "Public Viewer", Role: "viewer"},
	}

	for _, u := range users {
		db.Where("username = ?", u.Username).FirstOrCreate(&u)
	}

	// 2. Seed 2 Main Teams
	t1 := domain.Team{
		Name:        "Persija Jakarta",
		FoundedYear: 1928,
		Address:     "Jl. Pintu Satu Senayan, Gelora Bung Karno",
		City:        "Jakarta",
		LogoURL:     "https://upload.wikimedia.org/wikipedia/id/5/5e/Logo_Persija.png",
		Version:     1,
	}
	t2 := domain.Team{
		Name:        "Persib Bandung",
		FoundedYear: 1933,
		Address:     "Jl. Sulanjana No. 17, Tamansari",
		City:        "Bandung",
		LogoURL:     "https://upload.wikimedia.org/wikipedia/id/8/80/Persib_Bandung_logo.svg",
		Version:     1,
	}
	db.Create(&t1)
	db.Create(&t2)

	// 3. Seed 22 Players for Persija Jakarta (Complete eFootball Positions & Playstyles)
	persijaPlayers := []domain.Player{
		// Goalkeepers
		{TeamID: t1.ID, Name: "Andritany Ardhiyasa", Height: 180, Weight: 76, Position: "GK", Playstyle: strPtr("Defensive Goalkeeper"), JerseyNumber: 26, Version: 1},
		{TeamID: t1.ID, Name: "Cahya Supriadi", Height: 179, Weight: 72, Position: "GK", Playstyle: strPtr("Offensive Goalkeeper"), JerseyNumber: 99, Version: 1},
		// Center Backs
		{TeamID: t1.ID, Name: "Ondrej Kudela", Height: 182, Weight: 78, Position: "CB", Playstyle: strPtr("Build Up"), JerseyNumber: 17, Version: 1},
		{TeamID: t1.ID, Name: "Rizky Ridho", Height: 183, Weight: 74, Position: "CB", Playstyle: strPtr("Build Up"), JerseyNumber: 74, Version: 1},
		{TeamID: t1.ID, Name: "Muhammad Ferarri", Height: 181, Weight: 73, Position: "CB", Playstyle: strPtr("The Destroyer"), JerseyNumber: 41, Version: 1},
		{TeamID: t1.ID, Name: "Hansamu Yama", Height: 181, Weight: 75, Position: "CB", Playstyle: strPtr("Extra Frontman"), JerseyNumber: 23, Version: 1},
		// Full Backs (LB / RB)
		{TeamID: t1.ID, Name: "Firza Andika", Height: 169, Weight: 64, Position: "LB", Playstyle: strPtr("Offensive Fullback"), JerseyNumber: 11, Version: 1},
		{TeamID: t1.ID, Name: "Rezaldi Hehanussa", Height: 178, Weight: 70, Position: "LB", Playstyle: strPtr("Fullback Finisher"), JerseyNumber: 28, Version: 1},
		{TeamID: t1.ID, Name: "Rio Fahmi", Height: 170, Weight: 65, Position: "RB", Playstyle: strPtr("Offensive Fullback"), JerseyNumber: 2, Version: 1},
		{TeamID: t1.ID, Name: "Oliver Bias", Height: 164, Weight: 60, Position: "RB", Playstyle: strPtr("Defensive Fullback"), JerseyNumber: 22, Version: 1},
		// Defensive & Central Midfielders (DMF / CMF)
		{TeamID: t1.ID, Name: "Resky Fandi", Height: 172, Weight: 68, Position: "DMF", Playstyle: strPtr("Anchor Man"), JerseyNumber: 24, Version: 1},
		{TeamID: t1.ID, Name: "Hanif Sjahbandi", Height: 181, Weight: 75, Position: "DMF", Playstyle: strPtr("The Destroyer"), JerseyNumber: 19, Version: 1},
		{TeamID: t1.ID, Name: "Syahrian Abimanyu", Height: 171, Weight: 66, Position: "CMF", Playstyle: strPtr("The Orchestrator"), JerseyNumber: 8, Version: 1},
		{TeamID: t1.ID, Name: "Ramon Bueno", Height: 176, Weight: 71, Position: "CMF", Playstyle: strPtr("Box-to-Box"), JerseyNumber: 6, Version: 1},
		// Side & Attacking Midfielders (LMF / RMF / AMF)
		{TeamID: t1.ID, Name: "Maciej Gajos", Height: 175, Weight: 70, Position: "AMF", Playstyle: strPtr("Hole Player"), JerseyNumber: 10, Version: 1},
		{TeamID: t1.ID, Name: "Dony Tri Pamungkas", Height: 177, Weight: 68, Position: "LMF", Playstyle: strPtr("Cross Specialist"), JerseyNumber: 77, Version: 1},
		{TeamID: t1.ID, Name: "Riko Simanjuntak", Height: 158, Weight: 55, Position: "RMF", Playstyle: strPtr("Roaming Flank"), JerseyNumber: 25, Version: 1},
		{TeamID: t1.ID, Name: "Rayhan Hannan", Height: 167, Weight: 61, Position: "AMF", Playstyle: strPtr("Creative Playmaker"), JerseyNumber: 58, Version: 1},
		// Wingers & Second Strikers (LWF / RWF / SS)
		{TeamID: t1.ID, Name: "Witan Sulaeman", Height: 170, Weight: 64, Position: "LWF", Playstyle: strPtr("Prolific Winger"), JerseyNumber: 78, Version: 1},
		{TeamID: t1.ID, Name: "Ryo Matsumura", Height: 168, Weight: 62, Position: "SS", Playstyle: strPtr("Deep-Lying Forward"), JerseyNumber: 7, Version: 1},
		// Center Forwards (CF)
		{TeamID: t1.ID, Name: "Bambang Pamungkas", Height: 178, Weight: 72, Position: "CF", Playstyle: strPtr("Goal Poacher"), JerseyNumber: 20, Version: 1},
		{TeamID: t1.ID, Name: "Gustavo Almeida", Height: 180, Weight: 77, Position: "CF", Playstyle: strPtr("Fox in the Box"), JerseyNumber: 70, Version: 1},
	}

	for _, p := range persijaPlayers {
		db.Create(&p)
	}

	// 4. Seed 22 Players for Persib Bandung (Complete eFootball Positions & Playstyles)
	persibPlayers := []domain.Player{
		// Goalkeepers
		{TeamID: t2.ID, Name: "Kevin Ray Mendoza", Height: 187, Weight: 82, Position: "GK", Playstyle: strPtr("Offensive Goalkeeper"), JerseyNumber: 29, Version: 1},
		{TeamID: t2.ID, Name: "Teja Paku Alam", Height: 177, Weight: 70, Position: "GK", Playstyle: strPtr("Defensive Goalkeeper"), JerseyNumber: 14, Version: 1},
		// Center Backs
		{TeamID: t2.ID, Name: "Nick Kuipers", Height: 193, Weight: 88, Position: "CB", Playstyle: strPtr("Build Up"), JerseyNumber: 2, Version: 1},
		{TeamID: t2.ID, Name: "Gustavo Franca", Height: 188, Weight: 83, Position: "CB", Playstyle: strPtr("The Destroyer"), JerseyNumber: 4, Version: 1},
		{TeamID: t2.ID, Name: "Victor Igbonefo", Height: 183, Weight: 82, Position: "CB", Playstyle: strPtr("Build Up"), JerseyNumber: 32, Version: 1},
		{TeamID: t2.ID, Name: "Kakang Rudianto", Height: 183, Weight: 74, Position: "CB", Playstyle: strPtr("Extra Frontman"), JerseyNumber: 5, Version: 1},
		// Full Backs (LB / RB)
		{TeamID: t2.ID, Name: "Edo Febriansah", Height: 173, Weight: 67, Position: "LB", Playstyle: strPtr("Offensive Fullback"), JerseyNumber: 97, Version: 1},
		{TeamID: t2.ID, Name: "Zalnando", Height: 176, Weight: 69, Position: "LB", Playstyle: strPtr("Defensive Fullback"), JerseyNumber: 27, Version: 1},
		{TeamID: t2.ID, Name: "Henhen Herdiana", Height: 170, Weight: 66, Position: "RB", Playstyle: strPtr("Offensive Fullback"), JerseyNumber: 12, Version: 1},
		{TeamID: t2.ID, Name: "Rachmat Irianto", Height: 175, Weight: 72, Position: "RB", Playstyle: strPtr("Defensive Fullback"), JerseyNumber: 53, Version: 1},
		// Defensive & Central Midfielders (DMF / CMF)
		{TeamID: t2.ID, Name: "Dedi Kusnandar", Height: 175, Weight: 71, Position: "DMF", Playstyle: strPtr("Anchor Man"), JerseyNumber: 11, Version: 1},
		{TeamID: t2.ID, Name: "Mateo Kocijan", Height: 189, Weight: 81, Position: "DMF", Playstyle: strPtr("The Destroyer"), JerseyNumber: 17, Version: 1},
		{TeamID: t2.ID, Name: "Marc Klok", Height: 177, Weight: 73, Position: "CMF", Playstyle: strPtr("Box-to-Box"), JerseyNumber: 23, Version: 1},
		{TeamID: t2.ID, Name: "Adam Alis", Height: 172, Weight: 67, Position: "CMF", Playstyle: strPtr("The Orchestrator"), JerseyNumber: 18, Version: 1},
		// Side & Attacking Midfielders (LMF / RMF / AMF)
		{TeamID: t2.ID, Name: "Tyronne del Pino", Height: 180, Weight: 74, Position: "AMF", Playstyle: strPtr("Hole Player"), JerseyNumber: 10, Version: 1},
		{TeamID: t2.ID, Name: "Beckham Putra", Height: 173, Weight: 63, Position: "AMF", Playstyle: strPtr("Creative Playmaker"), JerseyNumber: 7, Version: 1},
		{TeamID: t2.ID, Name: "Febri Hariyadi", Height: 168, Weight: 62, Position: "RMF", Playstyle: strPtr("Roaming Flank"), JerseyNumber: 13, Version: 1},
		{TeamID: t2.ID, Name: "Ryan Kurnia", Height: 177, Weight: 69, Position: "LMF", Playstyle: strPtr("Cross Specialist"), JerseyNumber: 96, Version: 1},
		// Wingers & Second Strikers (LWF / RWF / SS)
		{TeamID: t2.ID, Name: "Atep", Height: 170, Weight: 65, Position: "LWF", Playstyle: strPtr("Prolific Winger"), JerseyNumber: 77, Version: 1},
		{TeamID: t2.ID, Name: "Ciro Alves", Height: 175, Weight: 72, Position: "RWF", Playstyle: strPtr("Prolific Winger"), JerseyNumber: 77, Version: 1},
		{TeamID: t2.ID, Name: "Mailson Lima", Height: 178, Weight: 75, Position: "SS", Playstyle: strPtr("Deep-Lying Forward"), JerseyNumber: 94, Version: 1},
		// Center Forwards (CF)
		{TeamID: t2.ID, Name: "David da Silva", Height: 185, Weight: 84, Position: "CF", Playstyle: strPtr("Goal Poacher"), JerseyNumber: 19, Version: 1},
	}

	// Adjust Ciro jersey number to avoid duplicate with Atep
	persibPlayers[19].JerseyNumber = 78

	for _, p := range persibPlayers {
		db.Create(&p)
	}

	// 5. Seed Initial Scheduled Match (El Clasico Indonesia)
	m := domain.Match{
		MatchDate:  "2026-09-01",
		MatchTime:  "15:30:00",
		HomeTeamID: t1.ID,
		AwayTeamID: t2.ID,
		Status:     "scheduled",
		Version:    1,
	}
	db.Create(&m)

	return nil
}

func strPtr(s string) *string {
	return &s
}
