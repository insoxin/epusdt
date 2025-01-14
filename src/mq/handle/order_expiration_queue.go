package handle

import (
	"context"
	"fmt"
	"github.com/assimon/luuu/model/data"
	"github.com/assimon/luuu/model/mdb"
	"github.com/hibiken/asynq"
)

const QueueOrderExpiration = "order:expiration"

func NewOrderExpirationQueue(tradeId string) (*asynq.Task, error) {
	return asynq.NewTask(QueueOrderExpiration, []byte(tradeId)), nil
}

// OrderExpirationHandle 设置订单过期
func OrderExpirationHandle(ctx context.Context, t *asynq.Task) error {
	tradeId := string(t.Payload())
	orderInfo, err := data.GetOrderInfoByTradeId(tradeId)
	if err != nil {
		return err
	}
	if orderInfo.ID <= 0 || orderInfo.Status != mdb.StatusWaitPay {
		return nil
	}
	err = data.UpdateOrderIsExpirationById(orderInfo.ID)
	if err != nil {
		return err
	}
	actualAmount := fmt.Sprintf("%.4f", orderInfo.ActualAmount)
	err = data.ClearPayCache(orderInfo.Token, actualAmount)
	if err != nil {
		return err
	}
	return nil
}
