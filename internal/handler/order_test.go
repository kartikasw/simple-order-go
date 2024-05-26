package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"simple-order-go/common"
	"simple-order-go/internal/entity"
	mockService "simple-order-go/internal/service/mock"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	order := randomOrder(false)
	fmt.Println("order: ", order)

	testCases := []struct {
		name          string
		body          requiredOrderRequest
		buildStubs    func(service *mockService.MockIOrderService)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: requiredOrderRequest{
				CustomerName: order.CustomerName,
				OrderedAt:    common.ParseTimeToString(order.OrderedAt),
				Items: []itemRequest{
					{
						Name:     order.Items[0].Name,
						Desc:     order.Items[0].Description,
						Quantity: order.Items[0].Quantity,
					},
					{
						Name:     order.Items[1].Name,
						Desc:     order.Items[1].Description,
						Quantity: order.Items[1].Quantity,
					},
				},
			},
			buildStubs: func(service *mockService.MockIOrderService) {
				service.EXPECT().CreateOrder(order).Times(1).Return(order, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("recorder: ", recorder)
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchOrder(t, recorder.Body, order)
			},
		},
		{
			name: "MissingRequiredData",
			body: requiredOrderRequest{
				CustomerName: order.CustomerName,
				OrderedAt:    common.ParseTimeToString(order.OrderedAt),
			},
			buildStubs: func(service *mockService.MockIOrderService) {
				service.EXPECT().CreateOrder(order).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("recorder: ", recorder)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Request = &http.Request{Header: make(http.Header), Method: "POST"}
			mockRequest(ctx, tc.body, 0)

			handler, service := setUpHandler(t)
			tc.buildStubs(service)

			handler.CreateOrder(ctx)
			tc.checkResponse(w)
		})
	}
}

func TestGetOrderByID(t *testing.T) {
	order := randomOrder(true)

	testCases := []struct {
		name          string
		param         int64
		buildStubs    func(service *mockService.MockIOrderService)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			param: order.ID,
			buildStubs: func(service *mockService.MockIOrderService) {
				service.EXPECT().GetOrder(order.ID).Times(1).Return(order, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchOrder(t, recorder.Body, order)
			},
		},
		{
			name:  "NotFound",
			param: order.ID,
			buildStubs: func(service *mockService.MockIOrderService) {
				service.EXPECT().GetOrder(order.ID).Times(1).Return(entity.OrderViewModel{}, pgx.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "NoParam",
			buildStubs: func(service *mockService.MockIOrderService) {
				service.EXPECT().CreateOrder(order).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Request = &http.Request{Header: make(http.Header), Method: "GET"}
			mockRequest(ctx, nil, tc.param)

			handler, service := setUpHandler(t)
			tc.buildStubs(service)

			handler.GetOrderByID(ctx)
			tc.checkResponse(w)
		})
	}
}

func TestUpdateOrder(t *testing.T) {
	order := randomOrder(true)

	testCases := []struct {
		name          string
		body          orderRequest
		buildStubs    func(service *mockService.MockIOrderService)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: orderRequest{
				CustomerName: order.CustomerName,
				OrderedAt:    common.ParseTimeToString(order.OrderedAt),
				Items: []itemRequest{
					{
						Name:     order.Items[0].Name,
						Desc:     order.Items[0].Description,
						Quantity: order.Items[0].Quantity,
					},
					{
						Name:     order.Items[1].Name,
						Desc:     order.Items[1].Description,
						Quantity: order.Items[1].Quantity,
					},
				},
			},
			buildStubs: func(service *mockService.MockIOrderService) {
				service.EXPECT().UpdateOrder(order).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "WithoutItems",
			body: orderRequest{
				CustomerName: order.CustomerName,
				OrderedAt:    common.ParseTimeToString(order.OrderedAt),
			},
			buildStubs: func(service *mockService.MockIOrderService) {
				arg := entity.OrderViewModel{
					ID:           order.ID,
					CustomerName: order.CustomerName,
					OrderedAt:    order.OrderedAt,
					Items:        []entity.ItemViewModel{},
				}
				service.EXPECT().UpdateOrder(arg).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				fmt.Println("recorder: ", recorder)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "MissingRequiredData",
			body: orderRequest{
				CustomerName: order.CustomerName,
			},
			buildStubs: func(service *mockService.MockIOrderService) {
				service.EXPECT().UpdateOrder(gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Request = &http.Request{Header: make(http.Header), Method: "PUT"}
			mockRequest(ctx, tc.body, order.ID)

			handler, service := setUpHandler(t)
			tc.buildStubs(service)

			handler.UpdateOrder(ctx)
			tc.checkResponse(w)
		})
	}
}

func TestDeleteOrde(t *testing.T) {
	var orderID int64 = 1

	testCases := []struct {
		name          string
		param         int64
		buildStubs    func(service *mockService.MockIOrderService)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			param: orderID,
			buildStubs: func(service *mockService.MockIOrderService) {
				service.EXPECT().DeleteOrder(orderID).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "NotFound",
			param: orderID,
			buildStubs: func(service *mockService.MockIOrderService) {
				service.EXPECT().DeleteOrder(orderID).Times(1).Return(pgx.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "NoParam",
			buildStubs: func(service *mockService.MockIOrderService) {
				service.EXPECT().CreateOrder(orderID).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Request = &http.Request{Header: make(http.Header), Method: "DELETE"}
			mockRequest(ctx, nil, tc.param)

			handler, service := setUpHandler(t)
			tc.buildStubs(service)

			handler.DeleteOrder(ctx)
			tc.checkResponse(w)
		})
	}
}

func requireBodyMatchOrder(t *testing.T, body *bytes.Buffer, order entity.OrderViewModel) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotOrder entity.OrderViewModel
	err = json.Unmarshal(data, &gotOrder)

	require.NoError(t, err)
	require.NotEmpty(t, gotOrder)
	require.Equal(t, order.CustomerName, gotOrder.CustomerName)
	require.WithinDuration(t, order.OrderedAt, gotOrder.OrderedAt, time.Second)
	require.Equal(t, len(order.Items), len(gotOrder.Items))
}

func randomOrder(withID bool) entity.OrderViewModel {
	n := 2

	items := make([]entity.ItemViewModel, n)

	for i := 0; i < n; i++ {
		items[i] = randomItem()
	}

	time, err := common.ParseStringToTime(common.ParseTimeToString(time.Now()))
	if err != nil {
		log.Fatal("Couldn't parse string to time: ", err)
	}

	var id int64 = 0
	if withID {
		id = common.RandomInt(1, 99)
	}

	order := entity.OrderViewModel{
		ID:           id,
		CustomerName: common.RandomName(),
		OrderedAt:    time,
		Items:        items,
	}

	return order
}

func randomItem() entity.ItemViewModel {
	item := entity.ItemViewModel{
		Name:        common.RandomName(),
		Description: common.RandomString(10),
		Quantity:    int32(common.RandomInt(1, 100)),
	}

	return item
}

func mockRequest(ctx *gin.Context, content interface{}, param int64) {
	ctx.Request.Header.Set("Content-Type", "application/json")

	if content != nil {
		jsonbytes, err := json.Marshal(content)
		if err != nil {
			log.Fatal(err)
		}

		bytes := bytes.NewBuffer(jsonbytes)
		ctx.Request.Body = io.NopCloser(bytes)
	}

	if param != 0 {
		ctx.Params = []gin.Param{{Key: "id", Value: big.NewInt(param).String()}}
	}
}

func setUpHandler(t *testing.T) (*OrderHandler, *mockService.MockIOrderService) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderService := mockService.NewMockIOrderService(ctrl)
	orderHandler := NewOrderHandler(orderService)

	return orderHandler, orderService
}
