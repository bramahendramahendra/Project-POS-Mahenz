package repo

import dto "pos_api/domain/payment_method/dto"

const getAllPaymentMethodQuery = `SELECT id, code, label, is_active, sort_order FROM payment_methods WHERE is_active = 1 ORDER BY sort_order ASC`

func (r *paymentMethodRepo) GetAll() ([]*dto.PaymentMethodResponse, error) {
	var dataDB []*dto.PaymentMethodResponse
	if err := r.db.Raw(getAllPaymentMethodQuery).Scan(&dataDB).Error; err != nil {
		return nil, err
	}
	if dataDB == nil {
		dataDB = []*dto.PaymentMethodResponse{}
	}
	return dataDB, nil
}
