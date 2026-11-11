package errors

type NotFoundError struct {
	Message string
}

type MethodNotAllowedError struct {
	Message string
}

type BadRequestError struct {
	Message string
}

type InternalServerError struct {
	Message string
}

type UnauthenticatedError struct {
	Message string
}

type UnauthorizededError struct {
	Message string
}

type ValidationError struct {
	Message string
}

// ServiceUnavailableError: layanan sementara tidak bisa dipakai (mis. belum diatur atau
// layanan lain tidak dapat dihubungi). Berbeda dari InternalServerError, pesannya
// diteruskan ke pengguna, jadi jangan isi dengan detail internal.
type ServiceUnavailableError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (e *MethodNotAllowedError) Error() string {
	return e.Message
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func (e *InternalServerError) Error() string {
	return e.Message
}

func (e *UnauthenticatedError) Error() string {
	return e.Message
}

func (e *UnauthorizededError) Error() string {
	return e.Message
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *ServiceUnavailableError) Error() string {
	return e.Message
}
