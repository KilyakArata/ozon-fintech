package graph

import (
	"context"
	"sync"

	"github.com/KilyakArata/gqlgen-todos/models"
)

var commentChannels = make(map[string][]chan *models.Comment)
var mu sync.Mutex

func notifySubscribers(postID string, comment *models.Comment) {
	mu.Lock()
	defer mu.Unlock()

	if channels, found := commentChannels[postID]; found {
		for _, ch := range channels {
			ch <- comment
		}
	}

}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *models.Comment, error) {
	ch := make(chan *models.Comment, 1)

	mu.Lock()
	commentChannels[postID] = append(commentChannels[postID], ch)
	mu.Unlock()

	go func() {
		<-ctx.Done()
		mu.Lock()
		defer mu.Unlock()
		for i, c := range commentChannels[postID] {
			if c == ch {
				commentChannels[postID] = append(commentChannels[postID][:i], commentChannels[postID][i+1:]...)
				break
			}
		}
	}()

	return ch, nil
}

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
