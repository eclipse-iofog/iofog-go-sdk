package util

var pkg struct {
	errorVariableNotInteger string
	errorVariableNotBool    string
}

func init() {
	pkg.errorVariableNotInteger = "Variable (%s) is not of type integer"
	pkg.errorVariableNotBool = "Variable (%s) is not of type bool"
}
