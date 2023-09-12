package config

type Environment int

func (e Environment) String() string {
	switch e {
	case Production:
		return "production"
	case Development:
		return "development"
	case Testing:
		return "testing"
	default:
		return "UNKNOWN"
	}
}

const (
	Production  Environment = 1
	Development Environment = 2
	Testing     Environment = 3
)
