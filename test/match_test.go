/*

 */

package models

import (
	"os"
	"reflect"
	"testing"

	"github.com/bitspawngg/tournament-bracket-manager/models"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/sirupsen/logrus"
)

func setupDB(t *testing.T) *models.DB {
	db_type, exists := os.LookupEnv("DB_TYPE")
	if !exists {
		t.Fatal("missing DB_TYPE environment variable")
	}
	db_path, exists := os.LookupEnv("DB_PATH")
	if !exists {
		t.Fatal("missing DB_PATH environment variable")
	}
	db := models.NewDB(db_type, db_path)
	if err := db.Connect(); err != nil {
		t.Fatal("db connection failed")
	}
	db.DB.Migrator().DropTable(&models.Match{})
	db.DB.AutoMigrate(&models.Match{})
	return db
}

func TestDB_CreateMatches(t *testing.T) {
	type args struct {
		matches []models.Match
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1.3 positive", args{
			[]models.Match{models.Match{"4f3d9be9-226f-47f0-94f4-399c163fcd23", 1, 3, "C", "D", "Ready", 0},
				models.Match{"4f3d9be9-226f-47f0-94f4-399c163fcd23", 1, 2, "E", "F", "Ready", 0}},
		}, false},
		{"1.4 negative-unique", args{
			[]models.Match{models.Match{"4f3d9be9-226f-47f0-94f4-399c163fcd23", 1, 4, "C", "D", "Ready", 0},
				models.Match{"4f3d9be9-226f-47f0-94f4-399c163fcd23", 1, 4, "E", "F", "Ready", 0}},
		}, true},
	}
	db := setupDB(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.CreateMatches(tt.args.matches); (err != nil) != tt.wantErr {
				t.Errorf("CreateMatches() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_GetMatch(t *testing.T) {
	type args struct {
		tournamentId string
		round        int
		table        int
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Match
		wantErr bool
	}{
		{"2.1 positive", args{"4f3d9be9-226f-47f0-94f4-399c163fcd23", 2, 1},
			&models.Match{"4f3d9be9-226f-47f0-94f4-399c163fcd23", 2, 1, "", "", "", 0},
			false},
		{"2.2 negative", args{"4f3d9be9-226f-47f0-94f4-399c163fcd23", 2, 2}, nil, true},
	}

	db := setupDB(t)
	db.DB.Create(&models.Match{TournamentID: "4f3d9be9-226f-47f0-94f4-399c163fcd23", Round: 2, Table: 1})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.GetMatch(tt.args.tournamentId, tt.args.round, tt.args.table)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_GetMatchesByStatus(t *testing.T) {
	type args struct {
		status string
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Match
		wantErr bool
	}{
		{"2.3 positive", args{status: "Ready"},
			[]models.Match{models.Match{"4f3d9be9-226f-47f0-94f4-399c163fcd23", 1, 1, "C", "D", "Ready", 0},
				models.Match{"4f3d9be9-226f-47f0-94f4-399c163fcd23", 1, 2, "E", "F", "Ready", 0}}, false},
	}
	db := setupDB(t)
	db.DB.Create(&models.Match{TournamentID: "4f3d9be9-226f-47f0-94f4-399c163fcd23", Round: 1, Table: 1, TeamOne: "C", TeamTwo: "D", Status: "Ready", Result: 0})
	db.DB.Create(&models.Match{TournamentID: "4f3d9be9-226f-47f0-94f4-399c163fcd23", Round: 1, Table: 2, TeamOne: "E", TeamTwo: "F", Status: "Ready", Result: 0})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.GetMatchesByStatus(tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.GetMatchesByStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.GetMatchesByStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_DeleteMatch(t *testing.T) {
	type args struct {
		tournamentId string
		round        int
		table        int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"2.1 positive", args{"4f3d9be9-226f-47f0-94f4-399c163fcd23", 4, 1}, false},
		{"2.2 negative", args{"4f3d9be9-226f-47f0-94f4-399c163fcd23", 4, 1}, true},
	}

	db := setupDB(t)
	db.DB.Create(&models.Match{TournamentID: "4f3d9be9-226f-47f0-94f4-399c163fcd23", Round: 4, Table: 1})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.DeleteMatch(tt.args.tournamentId, tt.args.round, tt.args.table); (err != nil) != tt.wantErr {
				t.Errorf("DB.DeleteMatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_GetMatchesByTournament(t *testing.T) {
	matchOne := models.Match{TournamentID: "4f3d9be9-226f-47f0-94f4-399c163fcd23", Round: 1, Table: 1, TeamOne: "A", TeamTwo: "B", Status: "Finished", Result: 1}
	matchTwo := models.Match{TournamentID: "4f3d9be9-226f-47f0-94f4-399c163fcd23", Round: 1, Table: 2, TeamOne: "A", TeamTwo: "B", Status: "Finished", Result: 1}
	type args struct {
		tournamentId string
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Match
		wantErr bool
	}{
		{"3.1", args{"4f3d9be9-226f-47f0-94f4-399c163fcd23"}, []models.Match{matchOne, matchTwo}, false},
	}
	db := setupDB(t)
	db.DB.Create(&matchOne)
	db.DB.Create(&matchTwo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.GetMatchesByTournament(tt.args.tournamentId)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.GetMatchesByTournament() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.GetMatchesByTournament() = %v, want %v", got, tt.want)
			}
		})
	}
}
