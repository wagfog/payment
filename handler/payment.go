package handler

import (
	"context"

	"github.com/wagfog/payment/domain/model"
	"github.com/wagfog/payment/domain/service"
	payment "github.com/wagfog/payment/proto"

	common "github.com/wagfog/mycommon"
)

type Payment struct {
	PaymentDataService service.IPaymentDataService
}

func (e *Payment) AddPayment(ctx context.Context, request *payment.PaymentInfo, response *payment.PaymentID) error {
	payment := &model.Payment{}
	if err := common.SwapTo(request, payment); err != nil {
		common.Debug(err)
	}
	paymentID, err := e.PaymentDataService.AddPayment(payment)
	if err != nil {
		common.Debug(err)
	}
	response.PaymentId = paymentID
	return nil
}

func (e *Payment) UpdatePayment(ctx context.Context, request *payment.PaymentInfo, response *payment.Response) error {
	payment := &model.Payment{}
	if err := common.SwapTo(request, payment); err != nil {
		common.Debug(err)
	}
	return e.PaymentDataService.UpdatePayment(payment)
}

func (e *Payment) DeletePaymentByID(ctx context.Context, request *payment.PaymentID, response *payment.Response) error {
	return e.PaymentDataService.DeletePayment(request.PaymentId)
}

func (e *Payment) FindPaymentByID(ctx context.Context, request *payment.PaymentID, response *payment.PaymentInfo) error {
	payment, err := e.PaymentDataService.FindPaymentByID(request.PaymentId)
	if err != nil {
		common.Debug(err)
	}
	return common.SwapTo(payment, response)
}

func (e *Payment) FindAllPayment(ctx context.Context, request *payment.All, response *payment.PaymentAll) error {
	allPayment, err := e.PaymentDataService.FindAllPayment()
	if err != nil {
		common.Debug(err)
	}

	for _, v := range allPayment {
		paymentInfo := &payment.PaymentInfo{}
		if err := common.SwapTo(v, paymentInfo); err != nil {
			common.Debug(err)
		}
		response.PaymentInfo = append(response.PaymentInfo, paymentInfo)
	}
	return nil
}
