package externals

import (
	"context"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/liuzemei/bot-manager/session"
)

func AssetList(ctx context.Context, token string) ([]*bot.Asset, error) {
	list, err := bot.AssetList(ctx, token)
	if err != nil {
		return nil, parseError(err.(bot.Error))
	}
	return list, nil
}

func AssetShow(ctx context.Context, assetId, token string) (*bot.Asset, error) {
	asset, err := bot.AssetShow(ctx, assetId, token)
	if err != nil {
		return nil, parseError(err.(bot.Error))
	}
	return asset, nil
}

func parseError(err bot.Error) *session.Error {
	if err.Code > 0 {
		switch err.Code {
		case 401:
			return session.AuthorizationError()
		case 403:
			return session.ForbiddenError()
		case 404:
			return session.NotFoundError()
		default:
			return session.ServerError()
		}
	}
	return nil
}
