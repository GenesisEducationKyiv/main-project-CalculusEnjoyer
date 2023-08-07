package dtm

import (
	"api/config"
	"api/domain"
	"email/transport/proto"
	emailProto "email/transport/proto"
	storageProto "storage/transport/proto"
	"strconv"

	dtmgrpc "github.com/dtm-labs/client/dtmgrpc"
	"github.com/lithammer/shortuuid/v3"
	"github.com/pkg/errors"
)

type DTMClient struct {
	conf config.Config
}

func NewDTMClient(conf config.Config) *DTMClient {
	return &DTMClient{
		conf: conf,
	}
}

func (d *DTMClient) SubmitAddEmailWithGreetingMessage(request domain.AddEmailRequest, email domain.SendEmailsRequest) error {
	gid := shortuuid.New()

	addEmailURL := d.conf.StorageNetwork + ":" + strconv.Itoa(d.conf.StoragePort) + storageProto.StorageService_AddEmail_FullMethodName
	addEmailRevertURL := d.conf.StorageNetwork + ":" + strconv.Itoa(d.conf.StoragePort) + storageProto.StorageService_AddEmailRevert_FullMethodName

	sendEmailURL := d.conf.EmailNetwork + ":" + strconv.Itoa(d.conf.EmailPort) + emailProto.EmailService_SendEmail_FullMethodName

	saga := dtmgrpc.NewSagaGrpc(d.conf.DTMAddress, gid).
		Add(addEmailURL, addEmailRevertURL, modelAddEmailToProto(request)).
		Add(sendEmailURL, "", modelSendEmailsToProto(email))

	saga.RetryCount = 0
	saga.WaitResult = true
	saga.TimeoutToFail = 8

	err := saga.Submit()

	return errors.Wrap(err, "failed to submit add email with greeting message")
}

func modelAddEmailToProto(request domain.AddEmailRequest) *storageProto.AddEmailRequest {
	return &storageProto.AddEmailRequest{
		Email: request.Email.Value,
	}
}

func modelSendEmailsToProto(request domain.SendEmailsRequest) *proto.SendEmailRequest {
	return &proto.SendEmailRequest{
		Subject: request.Template.Subject,
		Body:    request.Template.Body,
		To:      request.Interceptor.Value,
	}
}
