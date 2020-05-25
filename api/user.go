package api

import (
	"github.com/antihax/optional"
	openapi "github.com/sapphi-red/go-traq"
)

var (
	allUsersCache     []openapi.User
	currentUsersCache []openapi.User
)

type NameUserMap map[string]*openapi.User

func GetNameUserMap(includeSuspended bool, canUseCache bool) (NameUserMap, error) {
	users, err := GetUsers(includeSuspended, canUseCache)
	if err != nil {
		return nil, err
	}

	ret := make(NameUserMap, len(users))
	for _, u := range users {
		user := u
		ret[user.Name] = &user
	}
	return ret, nil
}

func GetUsers(includeSuspended bool, canUseCache bool) ([]openapi.User, error) {
	if canUseCache {
		if includeSuspended && allUsersCache != nil {
			return allUsersCache, nil
		} else if !includeSuspended && currentUsersCache != nil {
			return currentUsersCache, nil
		}
	}

	users, _, err := client.UserApi.GetUsers(auth, &openapi.UserApiGetUsersOpts{
		IncludeSuspended: optional.NewBool(includeSuspended),
	})
	if err != nil {
		return nil, err
	}

	if includeSuspended {
		allUsersCache = users
	} else {
		currentUsersCache = users
	}

	return users, err
}
