// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package localcachelayer

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/store"
)

type LocalCacheWebhookStore struct {
	store.WebhookStore
	rootStore *LocalCacheStore
}

func (s *LocalCacheWebhookStore) handleClusterInvalidateWebhook(msg *model.ClusterMessage) {
	if msg.Data == CLEAR_CACHE_MESSAGE_DATA {
		s.rootStore.webhookCache.Purge()
	} else {
		s.rootStore.webhookCache.Remove(msg.Data)
	}
}

func (s LocalCacheWebhookStore) ClearCaches() {
	s.rootStore.webhookCache.Purge()

	if s.rootStore.metrics != nil {
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter("Webhook - Purge")
	}
}

func (s LocalCacheWebhookStore) InvalidateWebhookCache(webhookId string) {
	s.rootStore.webhookCache.Remove(webhookId)
	if s.rootStore.metrics != nil {
		s.rootStore.metrics.IncrementMemCacheInvalidationCounter("Webhook - Remove by WebhookId")
	}
}

func (s LocalCacheWebhookStore) GetIncoming(id string, allowFromCache bool) (*model.IncomingWebhook, *model.AppError) {
	if !allowFromCache {
		return s.WebhookStore.GetIncoming(id, allowFromCache)
	}

	if incomingWebhook := s.rootStore.doStandardReadCache(s.rootStore.webhookCache, id); incomingWebhook != nil {
		if s.rootStore.metrics != nil {
			s.rootStore.metrics.IncrementMemCacheHitCounter("Webhook")
		}
		return incomingWebhook.(*model.IncomingWebhook), nil
	}

	if s.rootStore.metrics != nil {
		s.rootStore.metrics.IncrementMemCacheMissCounter("Webhook")
	}

	incomingWebhook, err := s.WebhookStore.GetIncoming(id, allowFromCache)
	if err != nil {
		return nil, err
	}

	s.rootStore.doStandardAddToCache(s.rootStore.webhookCache, id, incomingWebhook)

	return incomingWebhook, nil
}
