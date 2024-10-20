package config

import "os"

var (
	PORT         = os.Getenv("PORT")
	DATABASE_URL = os.Getenv("DATABASE_URL")
)
