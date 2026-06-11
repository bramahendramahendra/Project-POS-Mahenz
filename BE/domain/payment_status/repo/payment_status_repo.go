package repo

import dto "pos_api/domain/payment_status/dto"

const getAllPaymentStatusQuery = `SELECT id, code, label, is_active, sort_order FROM payment_statuses WHERE is_active = 1 ORDER BY sort_order ASC`

func (r *paymentStatusRepo) GetAll() ([]*dto.PaymentStatusResponse, error) {
	var dataDB []*dto.PaymentStatusResponse
	if err := r.db.Raw(getAllPaymentStatusQuery).Scan(&dataDB).Error; err != nil {
		return nil, err
	}
	if dataDB == nil {
		dataDB = []*dto.PaymentStatusResponse{}
	}
	return dataDB, nil
}
