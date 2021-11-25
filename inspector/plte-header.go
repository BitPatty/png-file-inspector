package pngscanner

type PLTEValidationResult struct {
	HasValidLength bool
}

type PLTEInspectionResult struct {
	Validation PLTEValidationResult
}

func inspectPLTEData(c *ChunkData) PLTEInspectionResult {
	return PLTEInspectionResult{
		Validation: PLTEValidationResult{
			HasValidLength: validatePLTELength(len(*c)),
		},
	}
}

func validatePLTELength(length int) bool {
	return (length > 0) && (length%3 == 0) && (length < (256 * 3))
}
