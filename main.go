package main

import (
	"github.com/wagfog/payment/domain/repository"
	service2 "github.com/wagfog/payment/domain/service"
	handler "github.com/wagfog/payment/handler"
	payment "github.com/wagfog/payment/proto"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	prometheus "github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	ratelimiter "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	common "github.com/wagfog/mycommon"
)

func main() {
	//
	consulConfig, err := common.GetConsulConfig("localhost", 8500, "/micro/config")
	if err != nil {
		common.Debug(err)
	}

	consul := consul.NewRegistry(func(o *registry.Options) {
		o.Addrs = []string{
			"localhost:8500",
		}
	})
	//jaeger
	t, io, err := common.NewTracer("go.micro.service.payment", "localhost:6831")
	if err != nil {
		common.Debug(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//mysql set
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	//init database
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		common.Debug(err)
	}
	defer db.Close()
	//禁止复数表
	db.SingularTable(true)

	//创建表
	tableInit := repository.NewPaymentRepository(db)
	tableInit.InitTable()

	//health
	common.PrometheusBoot(9089)

	service := micro.NewService(
		micro.Name("go.micro.service.payment"),
		micro.Version("latest"),
		micro.Address("0.0.0.0:8089"),
		micro.Registry(consul),
		//tracer
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		//hystrix
		micro.WrapHandler(ratelimiter.NewHandlerWrapper(1000)),
		//jiazai
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)

	service.Init()

	PaymentDataService := service2.NewPaymentDataService(repository.NewPaymentRepository(db))

	payment.RegisterPaymentHandler(service.Server(), &handler.Payment{PaymentDataService: PaymentDataService})

	if err := service.Run(); err != nil {
		common.Debug(err)
	}

}
