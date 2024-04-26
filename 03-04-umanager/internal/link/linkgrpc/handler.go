package linkgrpc

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/EfimVelichkin/3rd_module_GO/03-03-umanager/internal/database"
	"github.com/EfimVelichkin/3rd_module_GO/03-03-umanager/internal/link/models"
	"github.com/EfimVelichkin/3rd_module_GO/03-03-umanager/pkg/pb"
)

const ContentTypeJSON = "application/json"

var _ pb.LinkServiceServer = (*Handler)(nil)

func New(linksRepository linksRepository, timeout time.Duration, publisher amqpPublisher, queueName string) *Handler {
	return &Handler{
		linksRepository: linksRepository,
		pub:             publisher,
		queueName:       queueName,
		timeout:         timeout,
	}
}

type Handler struct {
	pb.UnimplementedLinkServiceServer
	linksRepository linksRepository
	pub             amqpPublisher
	queueName       string
	timeout         time.Duration
}

func (h Handler) GetLinkByUserID(ctx context.Context, id *pb.GetLinksByUserId) (*pb.ListLinkResponse, error) {
	links, err := h.linksRepository.FindByUserID(ctx, id.UserId)
	if err != nil {
		return nil, err
	}

	res := make([]*pb.Link, len(links))
	for i, l := range links {
		res[i] = &pb.Link{
			Id:        l.ID.Hex(),
			Title:     l.Title,
			Url:       l.URL,
			Images:    l.Images,
			Tags:      l.Tags,
			UserId:    l.UserID,
			CreatedAt: l.CreatedAt.String(),
			UpdatedAt: l.UpdatedAt.String(),
		}
	}
	return &pb.ListLinkResponse{Links: res}, err
}

func (h Handler) CreateLink(ctx context.Context, request *pb.CreateLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	var (
		id  primitive.ObjectID
		err error
	)
	if request.Id == "" {
		id = primitive.NewObjectID()
	} else {
		id, err = primitive.ObjectIDFromHex(request.Id)
	}

	if err != nil {
		return &pb.Empty{}, err
	}

	req := database.CreateLinkReq{
		ID:     id,
		Title:  request.Title,
		URL:    request.Url,
		Images: request.Images,
		Tags:   request.Tags,
		UserID: request.UserId,
	}

	link, err := h.linksRepository.Create(ctx, req)
	if err != nil {
		return &pb.Empty{}, err
	}

	data, err := json.Marshal(models.Message{ID: link.ID.Hex()})
	if err != nil {
		return &pb.Empty{}, err
	}

	err = h.pub.Publish("", h.queueName, false, false, amqp.Publishing{
		ContentType: ContentTypeJSON,
		Body:        data,
		Timestamp:   time.Now(),
	})

	return &pb.Empty{}, err
}

func (h Handler) GetLink(ctx context.Context, request *pb.GetLinkRequest) (*pb.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	l, err := h.linksRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.Link{
		Id:        l.ID.Hex(),
		Title:     l.Title,
		Url:       l.URL,
		Images:    l.Images,
		Tags:      l.Tags,
		UserId:    l.UserID,
		CreatedAt: l.CreatedAt.String(),
		UpdatedAt: l.UpdatedAt.String(),
	}, nil
}

func (h Handler) UpdateLink(ctx context.Context, request *pb.UpdateLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}

	req := database.UpdateLinkReq{
		ID:     id,
		Title:  request.Title,
		URL:    request.Url,
		Images: request.Images,
		Tags:   request.Tags,
		UserID: request.UserId,
	}
	_, err = h.linksRepository.Update(ctx, req)
	return &pb.Empty{}, err
}

func (h Handler) DeleteLink(ctx context.Context, request *pb.DeleteLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, h.linksRepository.Delete(ctx, id)
}

func (h Handler) ListLinks(ctx context.Context, request *pb.Empty) (*pb.ListLinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	links, err := h.linksRepository.FindAll(ctx)
	if err != nil {
		return &pb.ListLinkResponse{}, err
	}

	res := make([]*pb.Link, len(links))
	for i, l := range links {
		res[i] = &pb.Link{
			Id:        l.ID.Hex(),
			Title:     l.Title,
			Url:       l.URL,
			Images:    l.Images,
			Tags:      l.Tags,
			UserId:    l.UserID,
			CreatedAt: l.CreatedAt.String(),
			UpdatedAt: l.UpdatedAt.String(),
		}
	}
	return &pb.ListLinkResponse{Links: res}, err
}
