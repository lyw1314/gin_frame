package mailx

import (
	"context"
	"crypto/tls"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gopkg.in/gomail.v2"
)

var (
	externalEmailAddr   = "mta.ucmail.360.cn"
	internalEmailAddr   = "127.0.0.1"
	instrumentationName = "adgit.src.corp.qihoo.net/plat-arch/xkits/netx/mailx"
)

func SendInternal(ctx context.Context, msg *Message) (err error) {
	return sendmail(ctx, msg, internalEmailAddr)
}

func SendExternal(ctx context.Context, msg *Message) (err error) {
	return sendmail(ctx, msg, externalEmailAddr)
}

func sendmail(ctx context.Context, msg *Message, server string) (err error) {
	tr := otel.Tracer(instrumentationName)
	_, span := tr.Start(ctx, instrumentationName, oteltrace.WithSpanKind(oteltrace.SpanKindClient))
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	for i := 0; i < 3; i++ {
		dialer := gomail.NewDialer(server, 25, "", "")
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		if err = dialer.DialAndSend(msg.Message); err == nil {
			return
		}

		if server == externalEmailAddr {
			time.Sleep(3 * time.Second)
		}
	}

	return

}
