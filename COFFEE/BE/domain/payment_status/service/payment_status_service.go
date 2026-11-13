package service

import dto "coffee_api/domain/payment_status/dto"

func (s *paymentStatusService) GetAll() ([]*dto.PaymentStatusResponse, error) {
	return s.repo.GetAll()
}
