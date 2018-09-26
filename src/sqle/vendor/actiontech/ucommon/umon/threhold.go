package umon

type UmonThreholds []UmonThrehold
type UmonThrehold struct {
	Id              string
	Expr            string
	DurationSeconds int
	Labels          map[string]string
	Summary         string
	Detail          string
	Enable          bool
}
