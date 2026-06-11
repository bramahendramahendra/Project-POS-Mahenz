package service

import dto "pos_api/domain/payment_method/dto"

func (s *paymentMethodService) GetAll() ([]*dto.PaymentMethodResponse, error) {
	return s.repo.GetAll()
}
