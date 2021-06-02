package publisher

import (
	"github.com/miRemid/kira/common"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
)

var (
	MailPub micro.Event
	FilePub micro.Event
)

func Init(client client.Client) {
	MailPub = micro.NewPublisher(common.MailEvent, client)
	FilePub = micro.NewPublisher("kira.micro.service.user.delete", client)
}
