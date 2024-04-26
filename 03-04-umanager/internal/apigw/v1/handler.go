package v1

import (
	"github.com/EfimVelichkin/3rd_module_GO/03-03-umanager/pkg/api/apiv1"
)

type serverInterface interface {
	apiv1.ServerInterface
}

var _ serverInterface = (*Handler)(nil)

func New(usersRepository usersClient, linksRepository linksClient) *Handler {
	return &Handler{
		usersHandler: newUsersHandler(usersRepository),
		linksHandler: newLinksHandler(linksRepository),
	}
}

type Handler struct {
	*usersHandler
	*linksHandler
}
