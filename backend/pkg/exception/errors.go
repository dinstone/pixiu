package exception

type AppError interface {
	error
	Code() int
	Type() string
	Cause() error
}

type appError struct {
	code  int    // error code
	mesg  string // error message
	cause error  // error cause
	etype string // error type
}

const (
	ErrorTypeService  = "service"
	ErrorTypeBusiness = "business"
)

func (e *appError) Code() int {
	return e.code
}

func (e *appError) Type() string {
	return e.etype
}

func (e *appError) Error() string {
	if e.mesg != "" {
		return e.mesg
	}

	if e.cause != nil {
		return e.cause.Error()
	}

	return e.mesg
}

func (e *appError) Cause() error {
	return e.cause
}

// wrap a new service error
func WrapService(code int, message string, cause error) AppError {
	if message == "" && cause != nil {
		message = cause.Error()
	}

	return &appError{
		code:  code,
		cause: cause,
		mesg:  message,
		etype: ErrorTypeService,
	}
}

// wrap a new business error
func WrapBusiness(code int, message string, cause error) AppError {
	if message == "" && cause != nil {
		message = cause.Error()
	}

	return &appError{
		code:  code,
		mesg:  message,
		cause: cause,
		etype: ErrorTypeBusiness,
	}
}

// creates a new service error
func NewService(code int, message string) AppError {
	return &appError{
		code:  code,
		mesg:  message,
		cause: nil,
		etype: ErrorTypeService,
	}
}

// creates a new business error
func NewBusiness(code int, message string) AppError {
	return &appError{
		code:  code,
		mesg:  message,
		cause: nil,
		etype: ErrorTypeBusiness,
	}
}
