package static_data

var Users = map[string]string{
	"sameer123@gamil.com": "123456",
	"zohaib123@gamil.com": "qwerty",
}

var IndexCollection = []string{
	"users",
	"token",
	"contactmessages",
	"pets",
	"activity",
	"participants",
	"messages",
}

var IndexCollectionAttribute = map[string]string{
	"users":           "token",
	"token":           "token",
	"contactmessages": "token",
	"pets":            "token",
	"activity":        "token",
	"participants":    "token",
	"messages":        "token",
}
