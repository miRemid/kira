package repository

import (
	"github.com/gomodule/redigo/redis"
	"github.com/miRemid/kira/common"
	"github.com/miRemid/kira/proto/pb"
	"github.com/miRemid/kira/services/file/config"
)

func getUserFileLikeStatus(conn redis.Conn, fileid, userid string) (bool, error) {
	userKey := common.UserLikeKey(userid)
	_, err := redis.Int64(conn.Do("ZRANK", userKey, fileid))
	if err != nil {
		if err == redis.ErrNil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func getFileLikes(conn redis.Conn, fileid string) (int64, error) {
	index, err := redis.Int64(conn.Do("ZRANK", common.LikeRankKey, fileid))
	if err != nil {
		return 0, err
	}
	likes, err := redis.Int64Map(conn.Do("ZRANGE", common.LikeRankKey, index, index, "withscores"))
	return likes[fileid], err
}

func setFileIntoHash(conn redis.Conn, file *pb.UserFile) {
	conn.Do("HSET", file.FileID, "file_id", file.FileID)
	conn.Do("HSET", file.FileID, "user_name", file.UserName)
	conn.Do("HSET", file.FileID, "file_width", file.Width)
	conn.Do("HSET", file.FileID, "file_height", file.Height)
	conn.Do("HSET", file.FileID, "file_name", file.FileName)
}

func getFileFromHash(conn redis.Conn, userid, fileid string) *pb.UserFile {
	var file = new(pb.UserFile)
	res, _ := redis.StringMap(conn.Do("HGETALL", fileid))
	file.FileID = fileid
	file.FileName = res["file_name"]
	file.UserName = res["user_name"]
	file.Width = res["file_width"]
	file.Height = res["file_height"]
	file.FileURL = config.Path(fileid)
	if exist, err := getUserFileLikeStatus(conn, fileid, userid); err != nil {
		file.Liked = false
	} else {
		file.Liked = exist
		file.Likes, _ = getFileLikes(conn, fileid)
	}
	return file
}
