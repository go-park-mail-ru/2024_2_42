package configs

type internalParams struct {
	MainServerPort string
}

func NewInternalParams() internalParams {
	internalParams := internalParams{}

	internalParams.MainServerPort = ":8080"

	return internalParams
}
