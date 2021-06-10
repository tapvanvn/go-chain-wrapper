package campain

type IABI interface {
	GetMethod(input string) (string, []interface{}, error)
}
