package server

import (
	"context"
	"net/http"
	"net/url"

	normanapi "github.com/rancher/norman/api"
	"github.com/rancher/norman/parse"
	"github.com/rancher/norman/types"
	"github.com/rancher/rancher/pkg/api/user/api/setup"
	"github.com/rancher/rancher/pkg/api/user/store"
	"github.com/rancher/rancher/pkg/rbac"
	"github.com/rancher/types/config"
)

func New(ctx context.Context, cluster *config.UserContext) (http.Handler, error) {
	if err := setup.Schemas(ctx, cluster, cluster.Schemas); err != nil {
		return nil, err
	}

	server := normanapi.NewAPIServer()
	server.AccessControl = rbac.NewAccessControl(cluster.RBAC)
	server.URLParser = func(schemas *types.Schemas, url *url.URL) (parse.ParsedURL, error) {
		return URLParser(cluster.ClusterName, schemas, url)
	}
	server.StoreWrapper = store.ProjectSetter(server.StoreWrapper)

	if err := server.AddSchemas(cluster.Schemas); err != nil {
		return nil, err
	}

	return server, nil
}
