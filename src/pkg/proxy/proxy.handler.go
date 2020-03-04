package proxy

import (
	"context"
	"net/http"
	"oko/pkg/e"
	"oko/pkg/env"
	"oko/pkg/ginapp"
	"oko/pkg/ginapp/types"
	"oko/pkg/log"
	"oko/pkg/util"
	pb "oko/srv/proxy/proto"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
)

type handler struct {
	srv pb.ProxyService
}

func NewHandler() *handler { //nolint
	reg := etcd.NewRegistry(
		registry.Addrs(env.GetEnvOrPanic("ETCD_ADDRESS")),
	)
	service := micro.NewService(micro.Registry(reg))

	service.Init()

	cl := service.Client()

	_ = cl.Init(client.RequestTimeout(time.Second * 30))

	proxyService := pb.NewProxyService("go.micro.srv.proxy", cl)

	return &handler{srv: proxyService}
}

// Delete godoc
// @Summary Delete proxy item
// @Description Delete proxy item
// @ID delete-proxy
// @Tags Proxy
// @Accept json
// @Produce json
// @Param id path int true "Delete proxy item id"
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /proxy/{id} [delete]
func (p handler) Delete(c *gin.Context) {
	proxyID, err := ginapp.GetUint32PathParam(c, "id")
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}

	request := &pb.ProxyRequest{Id: proxyID}

	res, err := p.srv.Delete(context.Background(), request)
	if err != nil || res == nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")

		return
	}

	types.SuccessEmptyResponse(c)
}

// Create godoc
// @Summary Create proxy item
// @Description Create proxy item
// @ID create-proxy
// @Tags Proxy
// @Accept json
// @Produce json
// @Param object body proxy.CreateProxyRequest true "Create proxy item id"
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /proxy/ [post]
func (p handler) Create(c *gin.Context) {
	request := &CreateProxyRequest{}
	if err := c.BindJSON(request); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}

	createPrxReq := &pb.ProxyCreateRequest{
		Host: request.Host,
	}

	if request.HTTPS {
		createPrxReq.ProtocolHttps = &pb.ProxyCreateRequest_Https{Https: true}
	}
	if request.HTTP {
		createPrxReq.ProtocolHttp = &pb.ProxyCreateRequest_Http{Http: true}
	}

	_, err := p.srv.Create(context.Background(), createPrxReq)
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	types.SuccessEmptyResponse(c)
}

// Get godoc
// @Summary Get proxy item
// @Description Get proxy item
// @ID get-proxy
// @Tags Proxy
// @Accept json
// @Produce json
// @Param id path int true "Get proxy item id"
// @Success 200 {object} types.Response
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /proxy/{id} [get]
func (p handler) Get(c *gin.Context) {
	param := c.Param("id")
	proxyID, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}

	res, err := p.srv.Get(context.Background(), &pb.ProxyRequest{Id: uint32(proxyID)})
	if err != nil || res == nil || res.Data == nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	types.SuccessResponse(c, p.toDTO(res.Data))
}

// Update godoc
// @Summary Update proxy item
// @Description Update proxy item
// @ID update-proxy
// @Tags Proxy
// @Accept json
// @Produce json
// @Param id path int true "Update proxy item id"
// @Param object body proxy.CreateProxyRequest true "Update proxy item id"
// @Success 200 {object} types.Response
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /proxy/{id} [put]
func (p handler) Update(c *gin.Context) {
	param := c.Param("id")
	var proxyID uint64
	var err error
	if proxyID, err = strconv.ParseUint(param, 10, 32); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}
	request := &CreateProxyRequest{}
	if err = c.BindJSON(request); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}
	var HTTP *pb.ProxyCreateRequest_Http = nil
	if request.HTTP {
		HTTP = &pb.ProxyCreateRequest_Http{Http: true}
	}
	var https *pb.ProxyCreateRequest_Https = nil
	if request.HTTPS {
		https = &pb.ProxyCreateRequest_Https{Https: true}
	}

	_, err = p.srv.Update(context.Background(), &pb.ProxyCreateRequest{
		Id:            uint32(proxyID),
		Host:          request.Host,
		ProtocolHttp:  HTTP,
		ProtocolHttps: https,
	})

	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	types.SuccessEmptyResponse(c)
}

// CheckProxyExist godoc
// @Summary Check proxy with id exist
// @Description Update proxy item
// @ID exist-proxy
// @Tags Proxy
// @Accept json
// @Produce json
// @Param id path int true "Update proxy item id"
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /proxy/{id} [head]
func (p handler) Exist(c *gin.Context) {
	param := c.Param("id")
	proxyID, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}

	exist, err := p.srv.Exist(context.Background(), &pb.ProxyRequest{Id: uint32(proxyID)})
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if exist.Data {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}

// List godoc
// @Summary Proxies list
// @Description Proxies list
// @ID list-proxy
// @Tags Proxy
// @Accept json
// @Produce json
// @Success 200 {object} types.Response
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /proxy/ [get]
func (p handler) List(c *gin.Context) {
	list, err := p.srv.List(context.Background(), &pb.ProxyListRequest{})
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		log.Error(err)
		return
	}

	res := make([]*Response, 0, len(list.Data))

	for _, data := range list.Data {
		if data != nil {
			res = append(res, p.toDTO(data))
		}
	}
	types.SuccessResponse(c, res)
}

func (p handler) toDTO(data *pb.Proxy) (res *Response) {
	res = &Response{}

	if data == nil {
		return
	}

	res.ID = data.Id
	res.HTTPS = data.Https
	res.HTTP = data.Http
	res.Host = data.Host
	if len(data.State) > 0 {
		for _, st := range data.State {
			newState := State{Host: st.Host}
			if st.StartBan != nil {
				ts := util.TimestampToTime(*st.StartBan)
				newState.StartBan = &ts
			}
			if st.EndBan != nil {
				ts := util.TimestampToTime(*st.EndBan)
				newState.EndBan = &ts
			}
			res.State = append(res.State, newState)
		}
	}
	if data.CreatedAt != nil {
		res.CreatedAt = util.TimestampToTime(*data.CreatedAt)
	}
	if data.DeletedAt != nil {
		toTime := util.TimestampToTime(*data.DeletedAt)
		res.DeletedAt = &toTime
	}
	if data.UpdatedAt != nil {
		toTime := util.TimestampToTime(*data.UpdatedAt)
		res.UpdatedAt = &toTime
	}

	return
}
