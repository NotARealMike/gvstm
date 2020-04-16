package interfaces

type OpFailedError interface {
    error
}

type opFailedErrorImpl struct {
    message string
}

func NewOpFailedError(message string) OpFailedError {
    return &opFailedErrorImpl{message: message}
}

func (e *opFailedErrorImpl) Error() string {
    return e.message
}
