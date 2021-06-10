package campain

type IABI interface {
	Info()
	GetMethod(input string) (string, []interface{}, error)
}
