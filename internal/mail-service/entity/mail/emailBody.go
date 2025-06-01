package mail

import (
	"path/filepath"

	"github.com/aymerick/raymond"
)

const (
	REGISTER       = 1
	RESET_PASSWORD = 2
	EDIT_PASSWORD  = 3
	TOPUP_BALANCE  = 4
	LOW_BALANCE    = 5
	TEMPLATES_DIR  = "./internal/mail-service/entity/mail/templates" // TODO: Вынести в конфиг это
)

var emailTemplates = map[int]*raymond.Template{}

func getTemplate(name string) (*raymond.Template, error) {
	filePath := filepath.Join(TEMPLATES_DIR, name+".hbs")
	tmpl, err := raymond.ParseFile(filePath)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func init() {
	templates := map[int]string{
		REGISTER:       "registerEmail",
		RESET_PASSWORD: "resetPasswordEmail",
		EDIT_PASSWORD:  "editPasswordEmail",
		TOPUP_BALANCE:  "topUpedBalanceEmail",
		LOW_BALANCE:    "lowBalanceEmail",
	}

	for operation, name := range templates {
		tmpl, err := getTemplate(name)
		if err != nil {
			panic(err)
		}
		emailTemplates[operation] = tmpl
	}
}

func GetEmailBody(operation int, code string) (string, error) {
	tmpl, ok := emailTemplates[operation]
	if !ok {
		return "", nil
	}

	if tmpl == nil {
		return "", nil
	}

	onInsert := map[string]interface{}{
		"Code": code,
	}

	result, err := tmpl.Exec(onInsert)
	if err != nil {
		return "", err
	}

	return result, nil
}

func GetEmailLowBalanceBody(operation int, username, balance, href string) (string, error) {
	tmpl, ok := emailTemplates[operation]
	if !ok {
		return "", nil
	}

	if tmpl == nil {
		return "", nil
	}

	onInsert := map[string]interface{}{
		"Username": username,
		"Href":     href,
		"Balance":  balance,
	}

	result, err := tmpl.Exec(onInsert)
	if err != nil {
		return "", err
	}

	return result, nil
}

func GetEmailTopUpBody(operation int, username, amount string) (string, error) {
	tmpl, ok := emailTemplates[operation]
	if !ok {
		return "", nil
	}

	if tmpl == nil {
		return "", nil
	}

	onInsert := map[string]interface{}{
		"Username": username,
		"Amount":   amount,
	}

	result, err := tmpl.Exec(onInsert)
	if err != nil {
		return "", err
	}

	return result, nil
}
