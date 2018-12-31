package main

import (
	"fmt"

	"github.com/mihirkelkar/lenslocked.com/models"
	"github.com/mihirkelkar/lenslocked.com/rand"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "lenslocked_dev"
)

const remembertokenbytes = 32

func rememberToken() (string, error) {
	str, err := rand.String(remembertokenbytes)
	if err != nil {
		return "", err
	}
	return str, nil
}

func main() {
	plsqlString := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable\n",
		host, port, user, dbname)

	userService, err := models.NewUserService(plsqlString)
	if err != nil {
		panic(err)
	}
	userService.DestructiveReset()

	//	userService.Create("Mihir Kelkar", "test@test.com")
	//	//user, err := userService.ById(1)
	//	//if err != nil {
	//	//	panic(err)
	//	//	}
	//	//fmt.Println(user)
	//	user, _ := userService.ByEmail("girish@notrelated.com")
	//	fmt.Println(user)
	//	user, _ = userService.ByAge(25)
	//	fmt.Println(user)
	//	users, _ := userService.InAgeRange(23, 25)
	//	for _, user := range users {
	//		fmt.Println(user.Name)
	//		fmt.Println(user.Age)
	//	}//

	//	str, err := rememberToken()
	//	fmt.Println(str)//

	//	hh := hash.NewHMAC("this-is-my-secret-token")
	//	fmt.Println(hh.Hash(str))
	//fmt.Println("Deleting the given user")
	//err = userService.DeleteUser(user.ID)
	user := models.User{
		Name:     "Michael Scott",
		Email:    "michael@dundermiflin.com",
		Password: "ThisIsPassword",
	}
	userService.Create(&user)
	user_new, err := userService.ByRememberToken(user.Remember)
	fmt.Println(user_new)
}
