package maxbot

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"go.uber.org/zap"

	maxbotcli "github.com/max-messenger/max-bot-api-client-go"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/domain"
)

func (c *Client) UpdateHandler() (string, http.HandlerFunc) {
	maxHandler := c.client.GetHandler(c.updates)

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		secret := r.Header.Get("X-Max-Bot-Api-Secret")
		if secret != c.config.Secret {
			if r.Body != nil {
				_ = r.Body.Close()
			}
			host, err := os.Hostname()
			if err != nil {
				c.logger.Error("get host name", zap.Error(err))
			}
			c.logger.Error("secret is wrong, check subscription", zap.String("host", host))
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		maxHandler(w, r)
	}

	return c.config.Path, handlerFunc
}

func (c *Client) Updates() <-chan schemes.UpdateInterface {
	if c.config.SubscriptionType == subscribeTypePoll {
		return c.client.GetUpdates(c.ctx)
	}

	return c.updates
}

func (c *Client) unsubscribe(ctx context.Context) error {
	result, rErr := c.client.Subscriptions.GetSubscriptions(ctx)
	if rErr != nil {
		return fmt.Errorf("cannot get subscriptions: %w", rErr)
	}

	for _, s := range result.Subscriptions {
		c.logger.Info("unsubscribe", zap.Any("subscription", s))
		res, err := c.client.Subscriptions.Unsubscribe(ctx, s.Url)
		if err != nil {
			c.logger.Error("cannot unsubscribe", zap.Error(err))

			continue
		}
		c.logger.Info("unsubscribe", zap.Any("response", res))
	}

	return nil
}

func (c *Client) subscribe(ctx context.Context) error {
	// subscribe only for webhook
	if c.config.SubscriptionType != subscribeTypeWebhook {
		return nil
	}

	webhookURL, err := url.JoinPath(c.config.Url, c.config.Path)
	if err != nil {
		return fmt.Errorf("join webhook url: %w", err)
	}

	subscriptionResp, err := c.client.Subscriptions.Subscribe(ctx, webhookURL, subscriptions, c.config.Secret)
	if err != nil {
		return fmt.Errorf("cannot subscribe: %w", err)
	}

	c.logger.Info("subscribe", zap.Any("response", subscriptionResp))

	return nil
}

func makeMessage(m domain.Message) *maxbotcli.Message {
	result := maxbotcli.NewMessage().
		SetUser(m.UserID).SetChat(m.ChatID).SetText(m.Text)

	if m.Text != "" {
		result.SetFormat("markdown")
	}

	for _, kbd := range m.Keyboards {
		result.AddKeyboard(kbd)
	}

	for _, info := range m.Files {
		result.AddFile(info)
	}

	return result
}
