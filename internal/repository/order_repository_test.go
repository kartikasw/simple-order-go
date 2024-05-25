package repository

import (
	"simple-order-go/common"
	"simple-order-go/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomOrder(t *testing.T) entity.Order {
	n := common.RandomInt(1, 10)
	items := make([]entity.Item, n)

	for i := 1; 1 <= n; i++ {
		items = append(items, createRandomItem())
	}

	arg := entity.Order{
		CustomerName: common.RandomName(),
		OrderedAt:    time.Now(),
		Items:        items,
	}

	order, err := testOrderRepo.CreateOrder(arg)

	require.NoError(t, err)
	require.Equal(t, arg.CustomerName, order.CustomerName)
	require.Equal(t, n, len(order.Items))

	return order
}

func createRandomItem() (item entity.Item) {
	item = entity.Item{
		Name:        common.RandomName(),
		Description: common.RandomString(10),
		Quantity:    int32(common.RandomInt(1, 100)),
	}
	return
}

func TestCreateOrder(t *testing.T) {
	createRandomOrder(t)
}

func TestGetOrder(t *testing.T) {
	order1 := createRandomOrder(t)

	order2, err := testOrderRepo.GetOrder(order1.ID)

	require.NoError(t, err)
	require.Equal(t, order1.CustomerName, order2.CustomerName)
	require.Equal(t, order1.OrderedAt, order2.OrderedAt)
	require.Equal(t, len(order1.Items), len(order2.Items))
}

func TestGetAllOrders(t *testing.T) {
	for i := 1; i <= 10; i++ {
		createRandomOrder(t)
	}

	orders, err := testOrderRepo.GetAllOrders()

	require.NoError(t, err)
	require.Equal(t, 10, len(orders))
}

func TestUpdateOrder(t *testing.T) {
	n := 3

	updErrs := make(chan error, n)
	getErrs := make(chan error, n)
	custNames := make(chan string, n)
	firstItemNames := make(chan string, n)
	results := make(chan entity.Order, n)

	for i := 0; i < n; i++ {
		order := createRandomOrder(t)

		go func(entity.Order) {
			name := common.RandomName()
			firstItemName := common.RandomName()
			order.CustomerName = name
			order.Items[0].Name = firstItemName

			err := testOrderRepo.UpdateOrder(order)

			updErrs <- err
			custNames <- name
			firstItemNames <- firstItemName
		}(order)

		go func(orderID int64) {
			time.Sleep(1 * time.Millisecond)

			order, err := testOrderRepo.GetOrder(orderID)

			getErrs <- err
			results <- order
		}(order.ID)
	}

	for i := 0; i < n; i++ {
		updErr := <-updErrs
		getErr := <-getErrs
		require.NoError(t, updErr)
		require.NoError(t, getErr)

		order := <-results
		custName := <-custNames
		firstItemName := <-firstItemNames
		require.Equal(t, order.CustomerName, custName)
		require.Equal(t, order.Items[0].Name, firstItemName)
	}
}

func TestDeleteOrder(t *testing.T) {
	order := createRandomOrder(t)

	err := testOrderRepo.DeleteOrder(order.ID)
	require.NoError(t, err)

	delOrder, err := testOrderRepo.GetOrder(order.ID)
	require.Error(t, err)
	require.Equal(t, delOrder.ID, 0)
}
