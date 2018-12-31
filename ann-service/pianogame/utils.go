package pianogame

import "fmt"

func strConcate(s ...string) string {
	return fmt.Sprint(s) // why return array?
}
func strConcateF(str string, user string, pwd string, db string) string {
	return fmt.Sprintf(str, user, pwd, db)
}
