package transport

import (
	"context"
	"errors"

	"github.com/laidingqing/dabanshan-go/pb"
	"github.com/laidingqing/dabanshan-go/svcs/order/model"
	"github.com/laidingqing/dabanshan-go/utils"
)

// CreateOrder encode/decode
func decodeGRPCCreateOrderRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateOrderRequest)
	logger := utils.NewLogger()
	logger.Log("amount", req.Amount, "userId", req.Userid)
	return model.CreateOrderRequest{
		Invoice: model.Invoice{
			Amount:     req.Amount,
			UserID:     req.Userid,
			OrdereItem: pbInvoice2Model(req.Items),
		},
	}, nil
}

func encodeGRPCCreateOrderResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(model.CreatedOrderResponse)
	return &pb.CreatedOrderResponse{
		Id:  resp.ID,
		Err: err2str(resp.Err),
	}, nil
}

// GetOrders encode/decode

func decodeGRPCGetOrdersRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetOrdersRequest)
	return model.GetOrdersRequest{
		UserID:   req.Userid,
		TenantID: req.Tenantid,
	}, nil
}

func encodeGRPCGetOrdersResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(model.GetOrdersResponse)
	logger := utils.NewLogger()
	logger.Log("userID", resp.UserID)

	invoices := resp.Orders.Data.([]model.Invoice)
	return &pb.GetOrdersResponse{
		Userid:    resp.UserID,
		Tenantid:  resp.TenantID,
		PageIndex: int32(resp.Orders.PageIndex),
		PageSize:  int32(resp.Orders.PageSize),
		Invoices:  modelOrder2Pb(invoices),
		Err:       err2str(resp.Err),
	}, nil
}

// GetOrder encode/decode

func decodeGRPCGetOrderRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetOrderRequest)
	return model.GetOrderRequest{
		OrderID: req.Orderid,
	}, nil
}

func encodeGRPCGetOrderResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(model.GetOrderResponse)
	return &pb.GetOrderResponse{
		Invoice: &pb.InvoiceRecord{}, //todo
		Err:     err2str(resp.Err),
	}, nil
}

// addCart encode/decode func

func decodeGRPCAddCartRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateCartRequest)
	return model.CreateCartRequest{
		UserID:    req.Item.Userid,
		Price:     req.Item.Price,
		ProductID: req.Item.Productid,
	}, nil
}

func encodeGRPCAddCartResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(model.CreatedCartResponse)
	return &pb.CreatedCartResponse{
		Id:  resp.ID,
		Err: err2str(resp.Err),
	}, nil
}

// GetCartItems encode/decode

func decodeGRPCGetCartItemsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetCartItemsRequest)
	return model.GetCartItemsRequest{
		UserID: req.Userid,
	}, nil
}

func encodeGRPCGetCartItemsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(model.GetCartItemsResponse)
	return &pb.GetCartItemsResponse{
		Items: modelCartItem2Pb(resp.Items),
		Err:   err2str(resp.Err),
	}, nil
}

// removeCartItem
func decodeGRPCRemoveCartItemRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.RemoveCartItemRequest)
	return model.RemoveCartItemRequest{
		CartID: req.Cartid,
	}, nil
}

func encodeGRPCRemoveCartItemResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(model.RemoveCartItemResponse)
	return &pb.RemoveCartItemResponse{
		Err: err2str(resp.Err),
	}, nil
}

func decodeGRPCUpdateQuantityRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UpdateQuantityRequest)
	return model.UpdateQuantityRequest{
		CartID:   req.Cartid,
		Quantity: req.Quantity,
	}, nil
}

func encodeGRPCUpdateQuantityResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(model.UpdateQuantityResponse)
	return &pb.UpdateQuantityResponse{
		Err: err2str(resp.Err),
	}, nil
}

// client encode and decode

func encodeGRPCCreateOrderRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(model.CreateOrderRequest)
	logger := utils.NewLogger()
	logger.Log("amount", req.Invoice.Amount, "userId", req.Invoice.UserID)
	return &pb.CreateOrderRequest{
		Amount: req.Invoice.Amount,
		Userid: req.Invoice.UserID,
		Items:  modelInvoice2Pb(req.Invoice.OrdereItem),
	}, nil
}

func decodeGRPCCreateOrderResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.CreatedOrderResponse)
	return model.CreatedOrderResponse{
		ID:  reply.Id,
		Err: str2err(reply.Err)}, nil
}

// getOrders encode/decode func

func encodeGRPCGetOrdersRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(model.GetOrdersRequest)
	return &pb.GetOrdersRequest{
		Userid:    req.UserID,
		PageIndex: int32(req.PageIndex),
		PageSize:  int32(req.PageSize),
	}, nil
}

func decodeGRPCGetOrdersResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GetOrdersResponse)
	return model.GetOrdersResponse{
		UserID:   reply.Userid,
		TenantID: reply.Tenantid,
		Orders: utils.Pagination{
			PageIndex: int(reply.PageIndex),
			PageSize:  int(reply.PageSize),
			Data:      pbOrder2Model(reply.Invoices),
		},
		Err: str2err(reply.Err)}, nil
}

// getOrder encode/decode func

func encodeGRPCGetOrderRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(model.GetOrderRequest)
	return &pb.GetOrderRequest{
		Orderid: req.OrderID,
	}, nil
}

func decodeGRPCGetOrderResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GetOrderResponse)
	return model.GetOrderResponse{
		Order: model.Invoice{}, // TODO
		Err:   str2err(reply.Err)}, nil
}

func encodeGRPCAddCartRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(model.CreateCartRequest)
	return &pb.CreateCartRequest{
		Item: &pb.OrderItemRecord{
			Price:     req.Price,
			Productid: req.ProductID,
			Userid:    req.UserID,
		},
	}, nil
}

func decodeGRPCAddCartResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.CreatedCartResponse)
	return model.CreatedCartResponse{
		ID:  reply.Id,
		Err: str2err(reply.Err)}, nil
}

func encodeGRPCCartItemsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(model.GetCartItemsRequest)
	return &pb.GetCartItemsRequest{
		Userid: req.UserID,
	}, nil
}

func decodeGRPCCartItemsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GetCartItemsResponse)
	return model.GetCartItemsResponse{
		Items: pbCartItem2Model(reply.Items),
		Err:   str2err(reply.Err)}, nil
}

func encodeGRPCRemoveCartItemRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(model.RemoveCartItemRequest)
	return &pb.RemoveCartItemRequest{
		Cartid: req.CartID,
	}, nil
}

func decodeGRPCRemoveCartItemResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.RemoveCartItemResponse)
	return model.RemoveCartItemResponse{
		Err: str2err(reply.Err)}, nil
}

// UpdateQuantity encode/decode

func encodeGRPCUpdateQuantityRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(model.UpdateQuantityRequest)
	return &pb.UpdateQuantityRequest{
		Cartid:   req.CartID,
		Quantity: req.Quantity,
	}, nil
}

func decodeGRPCUpdateQuantityResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UpdateQuantityResponse)
	return model.UpdateQuantityResponse{
		Err: str2err(reply.Err)}, nil
}

func str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func pbInvoice2Model(records []*pb.OrderItemRecord) []model.OrderItem {
	var models []model.OrderItem
	for _, record := range records {
		models = append(models, model.OrderItem{
			CartID:   record.Cartid,
			Quantity: record.Quantity,
			Price:    record.Price,
		})
	}
	return models
}

func modelInvoice2Pb(records []model.OrderItem) []*pb.OrderItemRecord {
	var models []*pb.OrderItemRecord
	for _, record := range records {
		models = append(models, &pb.OrderItemRecord{
			Cartid:   record.CartID,
			Quantity: record.Quantity,
			Price:    record.Price,
		})
	}
	return models
}

func pbCartItem2Model(records []*pb.OrderItemRecord) []model.Cart {
	var models []model.Cart
	for _, record := range records {
		models = append(models, model.Cart{
			UserID:    record.Userid,
			Price:     record.Price,
			ProductID: record.Productid,
			CartID:    record.Cartid,
			Quantity:  record.Quantity,
		})
	}
	return models
}

func modelCartItem2Pb(models []model.Cart) []*pb.OrderItemRecord {
	var records []*pb.OrderItemRecord
	for _, model := range models {
		records = append(records, &pb.OrderItemRecord{
			Price:     model.Price,
			Productid: model.ProductID,
			Userid:    model.UserID,
			Cartid:    model.CartID,
			Quantity:  model.Quantity,
		})
	}

	return records
}

func pbOrderItem2Model(records []*pb.OrderItemRecord) []model.OrderItem {
	var models []model.OrderItem
	for _, record := range records {
		models = append(models, model.OrderItem{
			Price:     record.Price,
			ProductID: record.Productid,
			Quantity:  record.Quantity,
		})
	}
	return models
}

func pbOrder2Model(records []*pb.InvoiceRecord) []model.Invoice {
	var models []model.Invoice
	for _, record := range records {
		models = append(models, model.Invoice{
			UserID:     record.Userid,
			Amount:     record.Amount,
			OrdereItem: pbOrderItem2Model(record.Items),
		})
	}
	return models
}

func modelOrder2Pb(models []model.Invoice) []*pb.InvoiceRecord {
	var records []*pb.InvoiceRecord
	for _, model := range models {
		records = append(records, &pb.InvoiceRecord{
			Amount: model.Amount,
			Userid: model.UserID,
			Items:  modelInvoice2Pb(model.OrdereItem),
		})
	}

	return records
}
