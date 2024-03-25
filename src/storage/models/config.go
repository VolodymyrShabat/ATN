package models

type Config struct {
	Port               string
	Domain             string
	Env                string
	EmailFrom          string
	SMTPPass           string
	SMTPUser           string
	SMTPHost           string
	SMTPPort           string
	DatabaseConnectURL string
}
