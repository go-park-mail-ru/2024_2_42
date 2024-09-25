package config

type internalParams struct {
	MainServerPort string
}

func GetInternalParams() internalParams {
	InternalParams := internalParams{}

	InternalParams.MainServerPort = ":8080"

	return InternalParams
}
