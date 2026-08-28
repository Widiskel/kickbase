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

	// 1. Seed Demo Users (RBAC)
	users := []domain.User{
		{Username: "admin", Password: string(hashedPassword), Name: "System Administrator", Role: "admin"},
		{Username: "staff", Password: string(hashedPassword), Name: "Club Operator Staff", Role: "staff"},
		{Username: "viewer", Password: string(hashedPassword), Name: "Public Viewer", Role: "viewer"},
	}
	for _, u := range users {
		db.Where("username = ?", u.Username).FirstOrCreate(&u)
	}

	// 2. Seed 5 Professional Football Clubs
	teams := []domain.Team{
		{Name: "Persija Jakarta", FoundedYear: 1928, Address: "Jl. Pintu Satu Senayan, GBK", City: "Jakarta", LogoURL: "https://upload.wikimedia.org/wikipedia/id/5/5e/Logo_Persija.png", Version: 1},
		{Name: "Persib Bandung", FoundedYear: 1933, Address: "Jl. Sulanjana No. 17, Tamansari", City: "Bandung", LogoURL: "https://upload.wikimedia.org/wikipedia/id/8/80/Persib_Bandung_logo.svg", Version: 1},
		{Name: "Persebaya Surabaya", FoundedYear: 1927, Address: "Jl. Hayam Wuruk No. 6", City: "Surabaya", LogoURL: "https://upload.wikimedia.org/wikipedia/id/b/bf/Persebaya_Surabaya_logo.svg", Version: 1},
		{Name: "Arema FC", FoundedYear: 1987, Address: "Jl. Mayjend Panjaitan No. 42", City: "Malang", LogoURL: "https://upload.wikimedia.org/wikipedia/id/4/44/Logo_Arema_FC_2017.svg", Version: 1},
		{Name: "Bali United", FoundedYear: 1989, Address: "Jl. Bypass Ngurah Rai No. 748", City: "Gianyar", LogoURL: "https://upload.wikimedia.org/wikipedia/id/0/09/Logo_Bali_United.svg", Version: 1},
	}
	for i := range teams {
		db.Create(&teams[i])
	}

	// 3. Seed Players for each team (Total 110 realistic players covering all 15 tactical positions)
	allPlayers := []domain.Player{
		// --- PERSIJA JAKARTA (Team 0) ---
		{TeamID: teams[0].ID, Name: "Andritany Ardhiyasa", Height: 180, Weight: 76, Position: "GK", Playstyle: strPtr("Defensive Goalkeeper"), JerseyNumber: 26, Version: 1},
		{TeamID: teams[0].ID, Name: "Cahya Supriadi", Height: 179, Weight: 72, Position: "GK", Playstyle: strPtr("Offensive Goalkeeper"), JerseyNumber: 99, Version: 1},
		{TeamID: teams[0].ID, Name: "Ondrej Kudela", Height: 182, Weight: 78, Position: "CB", Playstyle: strPtr("Build Up"), JerseyNumber: 17, Version: 1},
		{TeamID: teams[0].ID, Name: "Rizky Ridho", Height: 183, Weight: 74, Position: "CB", Playstyle: strPtr("Build Up"), JerseyNumber: 74, Version: 1},
		{TeamID: teams[0].ID, Name: "Muhammad Ferarri", Height: 181, Weight: 73, Position: "CB", Playstyle: strPtr("The Destroyer"), JerseyNumber: 41, Version: 1},
		{TeamID: teams[0].ID, Name: "Hansamu Yama", Height: 181, Weight: 75, Position: "CB", Playstyle: strPtr("Extra Frontman"), JerseyNumber: 23, Version: 1},
		{TeamID: teams[0].ID, Name: "Firza Andika", Height: 169, Weight: 64, Position: "LB", Playstyle: strPtr("Offensive Fullback"), JerseyNumber: 11, Version: 1},
		{TeamID: teams[0].ID, Name: "Rio Fahmi", Height: 170, Weight: 65, Position: "RB", Playstyle: strPtr("Offensive Fullback"), JerseyNumber: 2, Version: 1},
		{TeamID: teams[0].ID, Name: "Oliver Bias", Height: 164, Weight: 60, Position: "RB", Playstyle: strPtr("Defensive Fullback"), JerseyNumber: 22, Version: 1},
		{TeamID: teams[0].ID, Name: "Resky Fandi", Height: 172, Weight: 68, Position: "DMF", Playstyle: strPtr("Anchor Man"), JerseyNumber: 24, Version: 1},
		{TeamID: teams[0].ID, Name: "Syahrian Abimanyu", Height: 171, Weight: 65, Position: "CMF", Playstyle: strPtr("Orchestrator"), JerseyNumber: 8, Version: 1},
		{TeamID: teams[0].ID, Name: "Hanif Sjahbandi", Height: 181, Weight: 74, Position: "DMF", Playstyle: strPtr("The Destroyer"), JerseyNumber: 19, Version: 1},
		{TeamID: teams[0].ID, Name: "Maciej Gajos", Height: 174, Weight: 70, Position: "AMF", Playstyle: strPtr("Classic No. 10"), JerseyNumber: 10, Version: 1},
		{TeamID: teams[0].ID, Name: "Rayhan Hannan", Height: 167, Weight: 60, Position: "AMF", Playstyle: strPtr("Hole Player"), JerseyNumber: 58, Version: 1},
		{TeamID: teams[0].ID, Name: "Dony Tri Pamungkas", Height: 177, Weight: 68, Position: "LMF", Playstyle: strPtr("Cross Specialist"), JerseyNumber: 77, Version: 1},
		{TeamID: teams[0].ID, Name: "Riko Simanjuntak", Height: 158, Weight: 55, Position: "RMF", Playstyle: strPtr("Prolific Winger"), JerseyNumber: 25, Version: 1},
		{TeamID: teams[0].ID, Name: "Witan Sulaeman", Height: 170, Weight: 63, Position: "LWF", Playstyle: strPtr("Roaming Flank"), JerseyNumber: 78, Version: 1},
		{TeamID: teams[0].ID, Name: "Ryo Matsumura", Height: 168, Weight: 64, Position: "SS", Playstyle: strPtr("Creative Playmaker"), JerseyNumber: 7, Version: 1},
		{TeamID: teams[0].ID, Name: "Marko Simic", Height: 187, Weight: 84, Position: "CF", Playstyle: strPtr("Target Man"), JerseyNumber: 9, Version: 1},
		{TeamID: teams[0].ID, Name: "Gustavo Almeida", Height: 180, Weight: 75, Position: "CF", Playstyle: strPtr("Goal Poacher"), JerseyNumber: 70, Version: 1},
		{TeamID: teams[0].ID, Name: "Bambang Pamungkas", Height: 171, Weight: 68, Position: "CF", Playstyle: strPtr("Fox in the Box"), JerseyNumber: 20, Version: 1},
		{TeamID: teams[0].ID, Name: "Maman Abdurahman", Height: 177, Weight: 74, Position: "CB", Playstyle: strPtr("Build Up"), JerseyNumber: 56, Version: 1},

		// --- PERSIB BANDUNG (Team 1) ---
		{TeamID: teams[1].ID, Name: "Kevin Ray Mendoza", Height: 187, Weight: 82, Position: "GK", Playstyle: strPtr("Offensive Goalkeeper"), JerseyNumber: 29, Version: 1},
		{TeamID: teams[1].ID, Name: "Teja Paku Alam", Height: 177, Weight: 70, Position: "GK", Playstyle: strPtr("Defensive Goalkeeper"), JerseyNumber: 14, Version: 1},
		{TeamID: teams[1].ID, Name: "Nick Kuipers", Height: 193, Weight: 88, Position: "CB", Playstyle: strPtr("Build Up"), JerseyNumber: 2, Version: 1},
		{TeamID: teams[1].ID, Name: "Alberto Rodriguez", Height: 191, Weight: 85, Position: "CB", Playstyle: strPtr("The Destroyer"), JerseyNumber: 22, Version: 1},
		{TeamID: teams[1].ID, Name: "Victor Igbonefo", Height: 183, Weight: 80, Position: "CB", Playstyle: strPtr("The Destroyer"), JerseyNumber: 32, Version: 1},
		{TeamID: teams[1].ID, Name: "Kakang Rudianto", Height: 183, Weight: 73, Position: "CB", Playstyle: strPtr("Build Up"), JerseyNumber: 5, Version: 1},
		{TeamID: teams[1].ID, Name: "Edo Febriansah", Height: 173, Weight: 67, Position: "LB", Playstyle: strPtr("Offensive Fullback"), JerseyNumber: 97, Version: 1},
		{TeamID: teams[1].ID, Name: "Rezaldi Hehanusa", Height: 178, Weight: 70, Position: "LB", Playstyle: strPtr("Fullback Finisher"), JerseyNumber: 56, Version: 1},
		{TeamID: teams[1].ID, Name: "Henhen Herdiana", Height: 169, Weight: 65, Position: "RB", Playstyle: strPtr("Defensive Fullback"), JerseyNumber: 12, Version: 1},
		{TeamID: teams[1].ID, Name: "Dedi Kusnandar", Height: 175, Weight: 70, Position: "DMF", Playstyle: strPtr("Anchor Man"), JerseyNumber: 11, Version: 1},
		{TeamID: teams[1].ID, Name: "Rachmat Irianto", Height: 175, Weight: 69, Position: "DMF", Playstyle: strPtr("The Destroyer"), JerseyNumber: 53, Version: 1},
		{TeamID: teams[1].ID, Name: "Marc Klok", Height: 177, Weight: 72, Position: "CMF", Playstyle: strPtr("Box-to-Box"), JerseyNumber: 23, Version: 1},
		{TeamID: teams[1].ID, Name: "Abdul Aziz", Height: 172, Weight: 66, Position: "CMF", Playstyle: strPtr("Orchestrator"), JerseyNumber: 8, Version: 1},
		{TeamID: teams[1].ID, Name: "Stefano Beltrame", Height: 183, Weight: 76, Position: "AMF", Playstyle: strPtr("Hole Player"), JerseyNumber: 93, Version: 1},
		{TeamID: teams[1].ID, Name: "Beckham Putra", Height: 173, Weight: 62, Position: "AMF", Playstyle: strPtr("Creative Playmaker"), JerseyNumber: 7, Version: 1},
		{TeamID: teams[1].ID, Name: "Febri Hariyadi", Height: 170, Weight: 63, Position: "RMF", Playstyle: strPtr("Speedster"), JerseyNumber: 13, Version: 1},
		{TeamID: teams[1].ID, Name: "Ryan Kurnia", Height: 176, Weight: 68, Position: "RWF", Playstyle: strPtr("Prolific Winger"), JerseyNumber: 96, Version: 1},
		{TeamID: teams[1].ID, Name: "Ciro Alves", Height: 175, Weight: 73, Position: "LWF", Playstyle: strPtr("Roaming Flank"), JerseyNumber: 77, Version: 1},
		{TeamID: teams[1].ID, Name: "David da Silva", Height: 185, Weight: 80, Position: "CF", Playstyle: strPtr("Goal Poacher"), JerseyNumber: 19, Version: 1},
		{TeamID: teams[1].ID, Name: "Ferdiansyah", Height: 172, Weight: 64, Position: "SS", Playstyle: strPtr("Dummy Runner"), JerseyNumber: 18, Version: 1},
		{TeamID: teams[1].ID, Name: "Achmad Jufriyanto", Height: 182, Weight: 78, Position: "CB", Playstyle: strPtr("Build Up"), JerseyNumber: 16, Version: 1},
		{TeamID: teams[1].ID, Name: "Ezra Walian", Height: 185, Weight: 81, Position: "CF", Playstyle: strPtr("Deep-Lying Forward"), JerseyNumber: 9, Version: 1},

		// --- PERSEBAYA SURABAYA (Team 2) ---
		{TeamID: teams[2].ID, Name: "Ernando Ari", Height: 178, Weight: 73, Position: "GK", Playstyle: strPtr("Offensive Goalkeeper"), JerseyNumber: 21, Version: 1},
		{TeamID: teams[2].ID, Name: "Andhika Ramadhani", Height: 182, Weight: 75, Position: "GK", Playstyle: strPtr("Defensive Goalkeeper"), JerseyNumber: 52, Version: 1},
		{TeamID: teams[2].ID, Name: "Dusan Stevanovic", Height: 188, Weight: 82, Position: "CB", Playstyle: strPtr("Build Up"), JerseyNumber: 5, Version: 1},
		{TeamID: teams[2].ID, Name: "Yan Victor", Height: 189, Weight: 84, Position: "CB", Playstyle: strPtr("The Destroyer"), JerseyNumber: 26, Version: 1},
		{TeamID: teams[2].ID, Name: "Kadek Raditya", Height: 178, Weight: 72, Position: "CB", Playstyle: strPtr("Build Up"), JerseyNumber: 23, Version: 1},
		{TeamID: teams[2].ID, Name: "Riswan Lauhin", Height: 179, Weight: 71, Position: "CB", Playstyle: strPtr("The Destroyer"), JerseyNumber: 44, Version: 1},
		{TeamID: teams[2].ID, Name: "Reva Adi Utama", Height: 172, Weight: 66, Position: "LB", Playstyle: strPtr("Offensive Fullback"), JerseyNumber: 3, Version: 1},
		{TeamID: teams[2].ID, Name: "Mikael Tata", Height: 174, Weight: 65, Position: "LB", Playstyle: strPtr("Defensive Fullback"), JerseyNumber: 25, Version: 1},
		{TeamID: teams[2].ID, Name: "Arief Catur", Height: 169, Weight: 63, Position: "RB", Playstyle: strPtr("Offensive Fullback"), JerseyNumber: 2, Version: 1},
		{TeamID: teams[2].ID, Name: "Andre Oktaviansyah", Height: 160, Weight: 58, Position: "DMF", Playstyle: strPtr("Anchor Man"), JerseyNumber: 8, Version: 1},
		{TeamID: teams[2].ID, Name: "Ripal Wahyudi", Height: 168, Weight: 64, Position: "DMF", Playstyle: strPtr("The Destroyer"), JerseyNumber: 36, Version: 1},
		{TeamID: teams[2].ID, Name: "Song Ui-young", Height: 171, Weight: 67, Position: "CMF", Playstyle: strPtr("Box-to-Box"), JerseyNumber: 7, Version: 1},
		{TeamID: teams[2].ID, Name: "Muhammad Hidayat", Height: 170, Weight: 66, Position: "CMF", Playstyle: strPtr("Orchestrator"), JerseyNumber: 96, Version: 1},
		{TeamID: teams[2].ID, Name: "Ze Valente", Height: 180, Weight: 74, Position: "AMF", Playstyle: strPtr("Classic No. 10"), JerseyNumber: 10, Version: 1},
		{TeamID: teams[2].ID, Name: "Robson Duarte", Height: 174, Weight: 68, Position: "RWF", Playstyle: strPtr("Roaming Flank"), JerseyNumber: 30, Version: 1},
		{TeamID: teams[2].ID, Name: "Bruno Moreira", Height: 178, Weight: 72, Position: "LWF", Playstyle: strPtr("Prolific Winger"), JerseyNumber: 99, Version: 1},
		{TeamID: teams[2].ID, Name: "Kasim Botan", Height: 170, Weight: 64, Position: "RMF", Playstyle: strPtr("Cross Specialist"), JerseyNumber: 77, Version: 1},
		{TeamID: teams[2].ID, Name: "Toni Firmansyah", Height: 167, Weight: 60, Position: "AMF", Playstyle: strPtr("Hole Player"), JerseyNumber: 68, Version: 1},
		{TeamID: teams[2].ID, Name: "Paulo Henrique", Height: 187, Weight: 83, Position: "CF", Playstyle: strPtr("Goal Poacher"), JerseyNumber: 9, Version: 1},
		{TeamID: teams[2].ID, Name: "Wildan Ramdhani", Height: 170, Weight: 65, Position: "SS", Playstyle: strPtr("Creative Playmaker"), JerseyNumber: 20, Version: 1},
		{TeamID: teams[2].ID, Name: "Oktafianus Fernando", Height: 168, Weight: 62, Position: "LMF", Playstyle: strPtr("Speedster"), JerseyNumber: 27, Version: 1},
		{TeamID: teams[2].ID, Name: "Chandra Waskito", Height: 176, Weight: 71, Position: "CF", Playstyle: strPtr("Fox in the Box"), JerseyNumber: 11, Version: 1},

		// --- AREMA FC (Team 3) ---
		{TeamID: teams[3].ID, Name: "Julian Schwarzer", Height: 181, Weight: 76, Position: "GK", Playstyle: strPtr("Offensive Goalkeeper"), JerseyNumber: 18, Version: 1},
		{TeamID: teams[3].ID, Name: "Teguh Amiruddin", Height: 183, Weight: 77, Position: "GK", Playstyle: strPtr("Defensive Goalkeeper"), JerseyNumber: 23, Version: 1},
		{TeamID: teams[3].ID, Name: "Bagas Adi Nugroho", Height: 176, Weight: 70, Position: "CB", Playstyle: strPtr("Build Up"), JerseyNumber: 5, Version: 1},
		{TeamID: teams[3].ID, Name: "Charles Almeida", Height: 185, Weight: 80, Position: "CB", Playstyle: strPtr("The Destroyer"), JerseyNumber: 32, Version: 1},
		{TeamID: teams[3].ID, Name: "Syaeful Anwar", Height: 184, Weight: 79, Position: "CB", Playstyle: strPtr("Extra Frontman"), JerseyNumber: 4, Version: 1},
		{TeamID: teams[3].ID, Name: "Johan Alfarizi", Height: 173, Weight: 68, Position: "LB", Playstyle: strPtr("Fullback Finisher"), JerseyNumber: 87, Version: 1},
		{TeamID: teams[3].ID, Name: "Achmad Syarif", Height: 177, Weight: 69, Position: "RB", Playstyle: strPtr("Defensive Fullback"), JerseyNumber: 26, Version: 1},
		{TeamID: teams[3].ID, Name: "Rifad Marasabessy", Height: 171, Weight: 65, Position: "RB", Playstyle: strPtr("Offensive Fullback"), JerseyNumber: 12, Version: 1},
		{TeamID: teams[3].ID, Name: "Jayus Hariono", Height: 174, Weight: 71, Position: "DMF", Playstyle: strPtr("The Destroyer"), JerseyNumber: 14, Version: 1},
		{TeamID: teams[3].ID, Name: "Julian Guevara", Height: 193, Weight: 86, Position: "DMF", Playstyle: strPtr("Anchor Man"), JerseyNumber: 6, Version: 1},
		{TeamID: teams[3].ID, Name: "Arkhan Fikri", Height: 165, Weight: 58, Position: "CMF", Playstyle: strPtr("Orchestrator"), JerseyNumber: 8, Version: 1},
		{TeamID: teams[3].ID, Name: "Samuel Balinsa", Height: 168, Weight: 62, Position: "CMF", Playstyle: strPtr("Box-to-Box"), JerseyNumber: 11, Version: 1},
		{TeamID: teams[3].ID, Name: "Ariel Lucero", Height: 170, Weight: 66, Position: "AMF", Playstyle: strPtr("Creative Playmaker"), JerseyNumber: 7, Version: 1},
		{TeamID: teams[3].ID, Name: "Muhammad Rafli", Height: 180, Weight: 73, Position: "AMF", Playstyle: strPtr("Hole Player"), JerseyNumber: 10, Version: 1},
		{TeamID: teams[3].ID, Name: "Dendi Santoso", Height: 172, Weight: 67, Position: "RMF", Playstyle: strPtr("Cross Specialist"), JerseyNumber: 41, Version: 1},
		{TeamID: teams[3].ID, Name: "Flabio Soares", Height: 171, Weight: 64, Position: "LMF", Playstyle: strPtr("Speedster"), JerseyNumber: 17, Version: 1},
		{TeamID: teams[3].ID, Name: "Ginanjar Wahyu", Height: 178, Weight: 70, Position: "LWF", Playstyle: strPtr("Prolific Winger"), JerseyNumber: 27, Version: 1},
		{TeamID: teams[3].ID, Name: "Charles Lokolingoy", Height: 188, Weight: 84, Position: "CF", Playstyle: strPtr("Target Man"), JerseyNumber: 19, Version: 1},
		{TeamID: teams[3].ID, Name: "Gilbert Alvarez", Height: 185, Weight: 82, Position: "CF", Playstyle: strPtr("Goal Poacher"), JerseyNumber: 91, Version: 1},
		{TeamID: teams[3].ID, Name: "Dedik Setiawan", Height: 177, Weight: 71, Position: "CF", Playstyle: strPtr("Fox in the Box"), JerseyNumber: 27, Version: 1},
		{TeamID: teams[3].ID, Name: "Achmad Maulana", Height: 182, Weight: 72, Position: "SS", Playstyle: strPtr("Dummy Runner"), JerseyNumber: 20, Version: 1},
		{TeamID: teams[3].ID, Name: "Bayu Aji", Height: 175, Weight: 68, Position: "CB", Playstyle: strPtr("Build Up"), JerseyNumber: 3, Version: 1},

		// --- BALI UNITED (Team 4) ---
		{TeamID: teams[4].ID, Name: "Adilson Maringa", Height: 194, Weight: 89, Position: "GK", Playstyle: strPtr("Defensive Goalkeeper"), JerseyNumber: 1, Version: 1},
		{TeamID: teams[4].ID, Name: "Muhammad Ridho", Height: 179, Weight: 74, Position: "GK", Playstyle: strPtr("Offensive Goalkeeper"), JerseyNumber: 88, Version: 1},
		{TeamID: teams[4].ID, Name: "Elias Dolah", Height: 196, Weight: 92, Position: "CB", Playstyle: strPtr("The Destroyer"), JerseyNumber: 4, Version: 1},
		{TeamID: teams[4].ID, Name: "Jajang Mulyana", Height: 182, Weight: 80, Position: "CB", Playstyle: strPtr("Extra Frontman"), JerseyNumber: 73, Version: 1},
		{TeamID: teams[4].ID, Name: "Haudi Abdillah", Height: 180, Weight: 74, Position: "CB", Playstyle: strPtr("Build Up"), JerseyNumber: 5, Version: 1},
		{TeamID: teams[4].ID, Name: "Ricky Fajrin", Height: 175, Weight: 68, Position: "LB", Playstyle: strPtr("Fullback Finisher"), JerseyNumber: 24, Version: 1},
		{TeamID: teams[4].ID, Name: "Ardi Idrus", Height: 166, Weight: 62, Position: "LB", Playstyle: strPtr("Defensive Fullback"), JerseyNumber: 2, Version: 1},
		{TeamID: teams[4].ID, Name: "Novri Setiawan", Height: 173, Weight: 66, Position: "RB", Playstyle: strPtr("Offensive Fullback"), JerseyNumber: 22, Version: 1},
		{TeamID: teams[4].ID, Name: "Made Andhika", Height: 172, Weight: 65, Position: "RB", Playstyle: strPtr("Defensive Fullback"), JerseyNumber: 33, Version: 1},
		{TeamID: teams[4].ID, Name: "Luthfi Kamal", Height: 170, Weight: 64, Position: "DMF", Playstyle: strPtr("Anchor Man"), JerseyNumber: 71, Version: 1},
		{TeamID: teams[4].ID, Name: "Mohammed Rashid", Height: 183, Weight: 78, Position: "CMF", Playstyle: strPtr("Box-to-Box"), JerseyNumber: 74, Version: 1},
		{TeamID: teams[4].ID, Name: "Kadek Agung", Height: 170, Weight: 65, Position: "CMF", Playstyle: strPtr("Orchestrator"), JerseyNumber: 18, Version: 1},
		{TeamID: teams[4].ID, Name: "Eber Bessa", Height: 167, Weight: 64, Position: "AMF", Playstyle: strPtr("Creative Playmaker"), JerseyNumber: 10, Version: 1},
		{TeamID: teams[4].ID, Name: "Ramdani Lestaluhu", Height: 168, Weight: 62, Position: "AMF", Playstyle: strPtr("Hole Player"), JerseyNumber: 77, Version: 1},
		{TeamID: teams[4].ID, Name: "Privat Mbarga", Height: 175, Weight: 70, Position: "RWF", Playstyle: strPtr("Prolific Winger"), JerseyNumber: 37, Version: 1},
		{TeamID: teams[4].ID, Name: "Yabes Roni", Height: 168, Weight: 63, Position: "RMF", Playstyle: strPtr("Speedster"), JerseyNumber: 11, Version: 1},
		{TeamID: teams[4].ID, Name: "Irfan Jaya", Height: 162, Weight: 58, Position: "LWF", Playstyle: strPtr("Roaming Flank"), JerseyNumber: 41, Version: 1},
		{TeamID: teams[4].ID, Name: "Rahmat", Height: 167, Weight: 61, Position: "LMF", Playstyle: strPtr("Cross Specialist"), JerseyNumber: 91, Version: 1},
		{TeamID: teams[4].ID, Name: "Ilija Spasojevic", Height: 187, Weight: 85, Position: "CF", Playstyle: strPtr("Goal Poacher"), JerseyNumber: 9, Version: 1},
		{TeamID: teams[4].ID, Name: "Jefferson Assis", Height: 185, Weight: 82, Position: "CF", Playstyle: strPtr("Target Man"), JerseyNumber: 94, Version: 1},
		{TeamID: teams[4].ID, Name: "Taufik Hidayat", Height: 179, Weight: 72, Position: "CF", Playstyle: strPtr("Fox in the Box"), JerseyNumber: 17, Version: 1},
		{TeamID: teams[4].ID, Name: "Ryuji Utomo", Height: 185, Weight: 80, Position: "CB", Playstyle: strPtr("Build Up"), JerseyNumber: 23, Version: 1},
	}

	for _, p := range allPlayers {
		db.Create(&p)
	}

	// 4. Seed Matches
	m1 := domain.Match{
		HomeTeamID: teams[0].ID,
		AwayTeamID: teams[1].ID,
		MatchDate:  "2026-09-01",
		MatchTime:  "15:30:00",
		Status:     "scheduled",
		Version:    1,
	}
	m2 := domain.Match{
		HomeTeamID: teams[2].ID,
		AwayTeamID: teams[3].ID,
		MatchDate:  "2026-09-02",
		MatchTime:  "19:00:00",
		Status:     "scheduled",
		Version:    1,
	}
	m3 := domain.Match{
		HomeTeamID: teams[4].ID,
		AwayTeamID: teams[0].ID,
		MatchDate:  "2026-09-08",
		MatchTime:  "18:30:00",
		Status:     "scheduled",
		Version:    1,
	}
	db.Create(&m1)
	db.Create(&m2)
	db.Create(&m3)

	return nil
}

func strPtr(s string) *string {
	return &s
}
