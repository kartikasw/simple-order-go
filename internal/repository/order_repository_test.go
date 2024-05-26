package repository

import (
	"simple-order-go/common"
	"simple-order-go/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomOrder(t *testing.T) entity.Order {
	n := common.RandomInt(1, 3)
	items := make([]entity.Item, n)

	for i := 0; i < int(n); i++ {
		items[i] = createRandomItem()
	}

	arg := entity.Order{
		CustomerName: common.RandomName(),
		OrderedAt:    time.Now(),
		Items:        items,
	}

	order, err := testOrderRepo.CreateOrder(arg)

	require.NoError(t, err)
	require.Equal(t, arg.CustomerName, order.CustomerName)
	require.Equal(t, int(n), len(order.Items))

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
	defer tearDown()

	createRandomOrder(t)
}

func TestGetOrder(t *testing.T) {
	defer tearDown()

	order1 := createRandomOrder(t)

	order2, err := testOrderRepo.GetOrder(order1.ID)

	require.NoError(t, err)
	require.Equal(t, order1.CustomerName, order2.CustomerName)
	require.WithinDuration(t, order1.OrderedAt, order2.OrderedAt, time.Second)
	require.Equal(t, len(order1.Items), len(order2.Items))
}

func TestGetAllOrders(t *testing.T) {
	defer tearDown()

	for i := 1; i <= 10; i++ {
		createRandomOrder(t)
	}

	orders, err := testOrderRepo.GetAllOrders()

	require.NoError(t, err)
	require.Equal(t, 10, len(orders))
}

func TestUpdateToCreateOrder(t *testing.T) {
	defer tearDown()

	n := 3

	updErrs := make(chan error, n)
	getErrs := make(chan error, n)
	itemsLen := make(chan int, n)
	result := make(chan entity.Order, n)

	for i := 1; i <= n; i++ {
		order := createRandomOrder(t)

		start := make(chan bool, 1)
		g1 := make(chan bool, 1)
		g2 := make(chan bool, 1)

		go func(entity.Order) {
			order.Items = append(order.Items, entity.Item{
				Name:        common.RandomName(),
				Description: common.RandomString(10),
				Quantity:    int32(common.RandomInt(1, 10)),
				OrderID:     order.ID,
			})

			start <- true

			err := testOrderRepo.UpdateOrder(order)

			updErrs <- err
			itemsLen <- len(order.Items)
			g1 <- true
		}(order)

		go func(orderID int64) {
			<-start

			time.Sleep(2 * time.Millisecond)

			o, err := testOrderRepo.GetOrder(orderID)

			getErrs <- err
			result <- o
			g2 <- true
		}(order.ID)

		<-g1
		<-g2
	}

	for i := 0; i < n; i++ {
		updErr := <-updErrs
		getErr := <-getErrs
		require.NoError(t, updErr)
		require.NoError(t, getErr)

		orderResult := <-result
		itemsLen := <-itemsLen
		require.NotEqual(t, len(orderResult.Items), itemsLen)
	}
}

func TestUpdateOrder(t *testing.T) {
	defer tearDown()

	n := 3

	updErrs := make(chan error, n)
	getErrs := make(chan error, n)
	custNames := make(chan string, n)
	firstItemNames := make(chan string, n)
	results := make(chan entity.Order, n)

	for i := 1; i <= n; i++ {
		order := createRandomOrder(t)

		start := make(chan bool, 1)
		g1 := make(chan bool, 1)
		g2 := make(chan bool, 1)

		go func(entity.Order) {
			name := common.RandomName()
			firstItemName := common.RandomName()
			order.CustomerName = name
			order.Items[0].Name = firstItemName

			start <- true

			err := testOrderRepo.UpdateOrder(order)

			updErrs <- err
			custNames <- name
			firstItemNames <- firstItemName
			g1 <- true
		}(order)

		go func(orderID int64) {
			<-start

			time.Sleep(1 * time.Millisecond)

			order, err := testOrderRepo.GetOrder(orderID)

			getErrs <- err
			results <- order
			g2 <- true
		}(order.ID)

		<-g1
		<-g2
	}

	for i := 0; i < n; i++ {
		updErr := <-updErrs
		getErr := <-getErrs
		require.NoError(t, updErr)
		require.NoError(t, getErr)

		orderResult := <-results
		custName := <-custNames
		firstItemName := <-firstItemNames
		require.NotEqual(t, orderResult.CustomerName, custName)
		require.NotEqual(t, orderResult.Items[0].Name, firstItemName)
	}
}

func TestDeleteOrder(t *testing.T) {
	defer tearDown()

	order := createRandomOrder(t)

	err := testOrderRepo.DeleteOrder(order.ID)
	require.NoError(t, err)

	delOrder, err := testOrderRepo.GetOrder(order.ID)
	require.Error(t, err)
	require.Equal(t, delOrder.ID, int64(0))
}

func tearDown() {
	tx := testDB.Begin()
	defer tx.Rollback()

	tx.Exec("DELETE FROM orders")

	tx.Commit()
}
