package repository

import (
	fcli "github.com/miRemid/kira/services/file/client"
	"github.com/micro/go-micro/v2/client"
)

type SiteRepository interface{}

type SiteRepositoryImpl struct {
	fileCli fcli.FileClient
}

func NewSiteRepository(cli client.Client) SiteRepository {
	var res SiteRepositoryImpl
	var ff = fcli.NewFileClient(cli)
	res.fileCli = ff
	return res
}

func (repo SiteRepositoryImpl) GetImage(fileID string) {

}
