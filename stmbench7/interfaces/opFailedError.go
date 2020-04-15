package interfaces

type OpFailedError struct {
    Message string
}

func (e *OpFailedError) Error() string {
    return e.Message
}
