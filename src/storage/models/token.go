package models

type Token struct {
	AccessToken  string `bson:"access_token" json:"access_token"`
	RefreshToken string `bson:"refresh_token,omitempty" json:"-"`
}
