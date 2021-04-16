- 点赞设计

Redis中保存点赞排行榜总表，数据结构为ZSET，KEY为: like_rank_set

ZSET中Key=file_id_like, Score=INT

Redis中也保存有用户点赞的文件，数据结构为Hash，KEY为: user_id_likes

HASH中Key=file_id, Value=Proto中的UserFile的JSON数据

总共需要的API列表如下
- ModifyLike(fileid) Need Login
    - 判断fileid是否在Hash中
    - 获取fileid的信息
    - 插入到user_id_likes中
    - 更改like_rank_set
- GetHotestLikesImage() []UserFile
    - 直接获取hot_likes_rank_list中的数据
- GetMyLikes(offset, limit) []UserFile Need Login
    - 

排行榜信息每半小时刷新一次，将会使用定时任务完成。每次都在like_rank_set中选取score最高的10张照片保存在Redis中

数据结构为List，Key=like_hot_rank_list, Value={Data:UserFile, Likes: INT}