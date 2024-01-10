package main

type User struct {
	id   int32
	name string
}

func (u User) clone() *User {
	println("u addr = ", &u)
	newUser := new(User)
	*newUser = u
	return newUser
}

func (u *User) clone2() *User {
	println("u addr = ", u)
	newUser := new(User)
	*newUser = *u
	return newUser
}

func main() {
	user := User{id: 3, name: "zjy"}

	newUser := user.clone()
	newUser2 := user.clone2()

	println("addr = ", &user)
	println("addr = ", newUser)
	println("addr = ", newUser2)

	println("id =", newUser.id)
	println("name =", newUser.name)
}
