// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: ride/driver/v1alpha1/driver_service.proto

package v1alpha1connect

import (
	context "context"
	errors "errors"
	http "net/http"
	strings "strings"

	connect "connectrpc.com/connect"
	v1alpha1 "github.com/ride-app/driver-service/api/ride/driver/v1alpha1"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// DriverServiceName is the fully-qualified name of the DriverService service.
	DriverServiceName = "ride.driver.v1alpha1.DriverService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// DriverServiceCreateDriverProcedure is the fully-qualified name of the DriverService's
	// CreateDriver RPC.
	DriverServiceCreateDriverProcedure = "/ride.driver.v1alpha1.DriverService/CreateDriver"
	// DriverServiceGetDriverProcedure is the fully-qualified name of the DriverService's GetDriver RPC.
	DriverServiceGetDriverProcedure = "/ride.driver.v1alpha1.DriverService/GetDriver"
	// DriverServiceUpdateDriverProcedure is the fully-qualified name of the DriverService's
	// UpdateDriver RPC.
	DriverServiceUpdateDriverProcedure = "/ride.driver.v1alpha1.DriverService/UpdateDriver"
	// DriverServiceDeleteDriverProcedure is the fully-qualified name of the DriverService's
	// DeleteDriver RPC.
	DriverServiceDeleteDriverProcedure = "/ride.driver.v1alpha1.DriverService/DeleteDriver"
	// DriverServiceGetVehicleProcedure is the fully-qualified name of the DriverService's GetVehicle
	// RPC.
	DriverServiceGetVehicleProcedure = "/ride.driver.v1alpha1.DriverService/GetVehicle"
	// DriverServiceUpdateVehicleProcedure is the fully-qualified name of the DriverService's
	// UpdateVehicle RPC.
	DriverServiceUpdateVehicleProcedure = "/ride.driver.v1alpha1.DriverService/UpdateVehicle"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	driverServiceServiceDescriptor             = v1alpha1.File_ride_driver_v1alpha1_driver_service_proto.Services().ByName("DriverService")
	driverServiceCreateDriverMethodDescriptor  = driverServiceServiceDescriptor.Methods().ByName("CreateDriver")
	driverServiceGetDriverMethodDescriptor     = driverServiceServiceDescriptor.Methods().ByName("GetDriver")
	driverServiceUpdateDriverMethodDescriptor  = driverServiceServiceDescriptor.Methods().ByName("UpdateDriver")
	driverServiceDeleteDriverMethodDescriptor  = driverServiceServiceDescriptor.Methods().ByName("DeleteDriver")
	driverServiceGetVehicleMethodDescriptor    = driverServiceServiceDescriptor.Methods().ByName("GetVehicle")
	driverServiceUpdateVehicleMethodDescriptor = driverServiceServiceDescriptor.Methods().ByName("UpdateVehicle")
)

// DriverServiceClient is a client for the ride.driver.v1alpha1.DriverService service.
type DriverServiceClient interface {
	CreateDriver(context.Context, *connect.Request[v1alpha1.CreateDriverRequest]) (*connect.Response[v1alpha1.CreateDriverResponse], error)
	GetDriver(context.Context, *connect.Request[v1alpha1.GetDriverRequest]) (*connect.Response[v1alpha1.GetDriverResponse], error)
	UpdateDriver(context.Context, *connect.Request[v1alpha1.UpdateDriverRequest]) (*connect.Response[v1alpha1.UpdateDriverResponse], error)
	DeleteDriver(context.Context, *connect.Request[v1alpha1.DeleteDriverRequest]) (*connect.Response[v1alpha1.DeleteDriverResponse], error)
	GetVehicle(context.Context, *connect.Request[v1alpha1.GetVehicleRequest]) (*connect.Response[v1alpha1.GetVehicleResponse], error)
	UpdateVehicle(context.Context, *connect.Request[v1alpha1.UpdateVehicleRequest]) (*connect.Response[v1alpha1.UpdateVehicleResponse], error)
}

// NewDriverServiceClient constructs a client for the ride.driver.v1alpha1.DriverService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewDriverServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) DriverServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &driverServiceClient{
		createDriver: connect.NewClient[v1alpha1.CreateDriverRequest, v1alpha1.CreateDriverResponse](
			httpClient,
			baseURL+DriverServiceCreateDriverProcedure,
			connect.WithSchema(driverServiceCreateDriverMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getDriver: connect.NewClient[v1alpha1.GetDriverRequest, v1alpha1.GetDriverResponse](
			httpClient,
			baseURL+DriverServiceGetDriverProcedure,
			connect.WithSchema(driverServiceGetDriverMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updateDriver: connect.NewClient[v1alpha1.UpdateDriverRequest, v1alpha1.UpdateDriverResponse](
			httpClient,
			baseURL+DriverServiceUpdateDriverProcedure,
			connect.WithSchema(driverServiceUpdateDriverMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		deleteDriver: connect.NewClient[v1alpha1.DeleteDriverRequest, v1alpha1.DeleteDriverResponse](
			httpClient,
			baseURL+DriverServiceDeleteDriverProcedure,
			connect.WithSchema(driverServiceDeleteDriverMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getVehicle: connect.NewClient[v1alpha1.GetVehicleRequest, v1alpha1.GetVehicleResponse](
			httpClient,
			baseURL+DriverServiceGetVehicleProcedure,
			connect.WithSchema(driverServiceGetVehicleMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updateVehicle: connect.NewClient[v1alpha1.UpdateVehicleRequest, v1alpha1.UpdateVehicleResponse](
			httpClient,
			baseURL+DriverServiceUpdateVehicleProcedure,
			connect.WithSchema(driverServiceUpdateVehicleMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// driverServiceClient implements DriverServiceClient.
type driverServiceClient struct {
	createDriver  *connect.Client[v1alpha1.CreateDriverRequest, v1alpha1.CreateDriverResponse]
	getDriver     *connect.Client[v1alpha1.GetDriverRequest, v1alpha1.GetDriverResponse]
	updateDriver  *connect.Client[v1alpha1.UpdateDriverRequest, v1alpha1.UpdateDriverResponse]
	deleteDriver  *connect.Client[v1alpha1.DeleteDriverRequest, v1alpha1.DeleteDriverResponse]
	getVehicle    *connect.Client[v1alpha1.GetVehicleRequest, v1alpha1.GetVehicleResponse]
	updateVehicle *connect.Client[v1alpha1.UpdateVehicleRequest, v1alpha1.UpdateVehicleResponse]
}

// CreateDriver calls ride.driver.v1alpha1.DriverService.CreateDriver.
func (c *driverServiceClient) CreateDriver(ctx context.Context, req *connect.Request[v1alpha1.CreateDriverRequest]) (*connect.Response[v1alpha1.CreateDriverResponse], error) {
	return c.createDriver.CallUnary(ctx, req)
}

// GetDriver calls ride.driver.v1alpha1.DriverService.GetDriver.
func (c *driverServiceClient) GetDriver(ctx context.Context, req *connect.Request[v1alpha1.GetDriverRequest]) (*connect.Response[v1alpha1.GetDriverResponse], error) {
	return c.getDriver.CallUnary(ctx, req)
}

// UpdateDriver calls ride.driver.v1alpha1.DriverService.UpdateDriver.
func (c *driverServiceClient) UpdateDriver(ctx context.Context, req *connect.Request[v1alpha1.UpdateDriverRequest]) (*connect.Response[v1alpha1.UpdateDriverResponse], error) {
	return c.updateDriver.CallUnary(ctx, req)
}

// DeleteDriver calls ride.driver.v1alpha1.DriverService.DeleteDriver.
func (c *driverServiceClient) DeleteDriver(ctx context.Context, req *connect.Request[v1alpha1.DeleteDriverRequest]) (*connect.Response[v1alpha1.DeleteDriverResponse], error) {
	return c.deleteDriver.CallUnary(ctx, req)
}

// GetVehicle calls ride.driver.v1alpha1.DriverService.GetVehicle.
func (c *driverServiceClient) GetVehicle(ctx context.Context, req *connect.Request[v1alpha1.GetVehicleRequest]) (*connect.Response[v1alpha1.GetVehicleResponse], error) {
	return c.getVehicle.CallUnary(ctx, req)
}

// UpdateVehicle calls ride.driver.v1alpha1.DriverService.UpdateVehicle.
func (c *driverServiceClient) UpdateVehicle(ctx context.Context, req *connect.Request[v1alpha1.UpdateVehicleRequest]) (*connect.Response[v1alpha1.UpdateVehicleResponse], error) {
	return c.updateVehicle.CallUnary(ctx, req)
}

// DriverServiceHandler is an implementation of the ride.driver.v1alpha1.DriverService service.
type DriverServiceHandler interface {
	CreateDriver(context.Context, *connect.Request[v1alpha1.CreateDriverRequest]) (*connect.Response[v1alpha1.CreateDriverResponse], error)
	GetDriver(context.Context, *connect.Request[v1alpha1.GetDriverRequest]) (*connect.Response[v1alpha1.GetDriverResponse], error)
	UpdateDriver(context.Context, *connect.Request[v1alpha1.UpdateDriverRequest]) (*connect.Response[v1alpha1.UpdateDriverResponse], error)
	DeleteDriver(context.Context, *connect.Request[v1alpha1.DeleteDriverRequest]) (*connect.Response[v1alpha1.DeleteDriverResponse], error)
	GetVehicle(context.Context, *connect.Request[v1alpha1.GetVehicleRequest]) (*connect.Response[v1alpha1.GetVehicleResponse], error)
	UpdateVehicle(context.Context, *connect.Request[v1alpha1.UpdateVehicleRequest]) (*connect.Response[v1alpha1.UpdateVehicleResponse], error)
}

// NewDriverServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewDriverServiceHandler(svc DriverServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	driverServiceCreateDriverHandler := connect.NewUnaryHandler(
		DriverServiceCreateDriverProcedure,
		svc.CreateDriver,
		connect.WithSchema(driverServiceCreateDriverMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	driverServiceGetDriverHandler := connect.NewUnaryHandler(
		DriverServiceGetDriverProcedure,
		svc.GetDriver,
		connect.WithSchema(driverServiceGetDriverMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	driverServiceUpdateDriverHandler := connect.NewUnaryHandler(
		DriverServiceUpdateDriverProcedure,
		svc.UpdateDriver,
		connect.WithSchema(driverServiceUpdateDriverMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	driverServiceDeleteDriverHandler := connect.NewUnaryHandler(
		DriverServiceDeleteDriverProcedure,
		svc.DeleteDriver,
		connect.WithSchema(driverServiceDeleteDriverMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	driverServiceGetVehicleHandler := connect.NewUnaryHandler(
		DriverServiceGetVehicleProcedure,
		svc.GetVehicle,
		connect.WithSchema(driverServiceGetVehicleMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	driverServiceUpdateVehicleHandler := connect.NewUnaryHandler(
		DriverServiceUpdateVehicleProcedure,
		svc.UpdateVehicle,
		connect.WithSchema(driverServiceUpdateVehicleMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/ride.driver.v1alpha1.DriverService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case DriverServiceCreateDriverProcedure:
			driverServiceCreateDriverHandler.ServeHTTP(w, r)
		case DriverServiceGetDriverProcedure:
			driverServiceGetDriverHandler.ServeHTTP(w, r)
		case DriverServiceUpdateDriverProcedure:
			driverServiceUpdateDriverHandler.ServeHTTP(w, r)
		case DriverServiceDeleteDriverProcedure:
			driverServiceDeleteDriverHandler.ServeHTTP(w, r)
		case DriverServiceGetVehicleProcedure:
			driverServiceGetVehicleHandler.ServeHTTP(w, r)
		case DriverServiceUpdateVehicleProcedure:
			driverServiceUpdateVehicleHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedDriverServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedDriverServiceHandler struct{}

func (UnimplementedDriverServiceHandler) CreateDriver(context.Context, *connect.Request[v1alpha1.CreateDriverRequest]) (*connect.Response[v1alpha1.CreateDriverResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("ride.driver.v1alpha1.DriverService.CreateDriver is not implemented"))
}

func (UnimplementedDriverServiceHandler) GetDriver(context.Context, *connect.Request[v1alpha1.GetDriverRequest]) (*connect.Response[v1alpha1.GetDriverResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("ride.driver.v1alpha1.DriverService.GetDriver is not implemented"))
}

func (UnimplementedDriverServiceHandler) UpdateDriver(context.Context, *connect.Request[v1alpha1.UpdateDriverRequest]) (*connect.Response[v1alpha1.UpdateDriverResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("ride.driver.v1alpha1.DriverService.UpdateDriver is not implemented"))
}

func (UnimplementedDriverServiceHandler) DeleteDriver(context.Context, *connect.Request[v1alpha1.DeleteDriverRequest]) (*connect.Response[v1alpha1.DeleteDriverResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("ride.driver.v1alpha1.DriverService.DeleteDriver is not implemented"))
}

func (UnimplementedDriverServiceHandler) GetVehicle(context.Context, *connect.Request[v1alpha1.GetVehicleRequest]) (*connect.Response[v1alpha1.GetVehicleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("ride.driver.v1alpha1.DriverService.GetVehicle is not implemented"))
}

func (UnimplementedDriverServiceHandler) UpdateVehicle(context.Context, *connect.Request[v1alpha1.UpdateVehicleRequest]) (*connect.Response[v1alpha1.UpdateVehicleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("ride.driver.v1alpha1.DriverService.UpdateVehicle is not implemented"))
}
