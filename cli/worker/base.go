package worker

import (
	"os"
	"os/user"
	"pmon2/pmond/model"
	"strconv"
)

func GetProcUser(a *model.ExecFlags) (*user.User, error) {
	runUser := a.User
	var curUser *user.User
	var err error

	if len(runUser) <= 0 {
		curUser, err = user.LookupId(strconv.Itoa(os.Getuid()))
	} else {
		curUser, err = user.Lookup(runUser)
	}

	if err != nil {
		return nil, err
	}

	return curUser, nil
}
