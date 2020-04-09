package interfaces

type OpFailedError struct {
    message string
}

func (e *OpFailedError) Error() string {
    return e.message
}
