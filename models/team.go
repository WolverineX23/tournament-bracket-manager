package models

import	(
	//"errors"
)

type Team struct {
	ID 			 uint   `json:"id" gorm:"primary_key"`
	TeamName  	 string `json:"teamName"`
	ChampionGame int   	`json:"championGame"`
	TotalGame 	 int 	`json:"totalGame"`
}

func (db DB) CreateTeams(teams []Team) error {
	return db.DB.Create(teams).Error
}