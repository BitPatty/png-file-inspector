package pngscanner

type IENDValidationResult struct {
	HasValidLength bool
}

type IENDInspectionResult struct {
	Validation IENDValidationResult
}

func inspectIENDData(c *ChunkData) IENDInspectionResult {
	return IENDInspectionResult{
		Validation: IENDValidationResult{
			HasValidLength: validateIENDLength(len(*c)),
		},
	}
}

func validateIENDLength(length int) bool {
	return length == 0
}
