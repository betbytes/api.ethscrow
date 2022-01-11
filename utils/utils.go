package utils

import (
	"github.com/joho/godotenv"
	"os"
)

var PORT string
var MAIL_USER string
var MAIL_PASS string
var AES_KEY string
var DATABASE_URL string
var ETHNET_URL string
var BASE_URl string
var FB_TYPE string
var FB_PROJ_ID string
var FB_PRIV_KEY_ID string
var FB_PRIV_KEY string
var FB_CLI_EMAIL string
var FB_CLI_ID string
var FB_AUTH_URI string
var FB_TOKEN_URI string
var FB_PROV_URL string
var FB_CLI_URL string

func SetParams() error {
	err := godotenv.Load()
	PORT = os.Getenv("PORT")
	MAIL_USER = os.Getenv("MAIL_USER")
	MAIL_PASS = os.Getenv("MAIL_PASS")
	AES_KEY = os.Getenv("AES_KEY")
	DATABASE_URL = os.Getenv("DATABASE_URL")
	ETHNET_URL = os.Getenv("ETHNET_URL")
	FB_TYPE = os.Getenv("FB_TYPE")
	FB_PROJ_ID = os.Getenv("FB_PROJ_ID")
	FB_PRIV_KEY_ID = os.Getenv("FB_PRIV_KEY_ID")
	FB_PRIV_KEY = os.Getenv("FB_PRIV_KEY")
	FB_CLI_EMAIL = os.Getenv("FB_CLI_EMAIL")
	FB_CLI_ID = os.Getenv("FB_CLI_ID")
	FB_AUTH_URI = os.Getenv("FB_AUTH_URI")
	FB_TOKEN_URI = os.Getenv("FB_TOKEN_URI")
	FB_PROV_URL = os.Getenv("FB_PROV_URL")
	FB_CLI_URL = os.Getenv("FB_CLI_URL")
	return err

}
