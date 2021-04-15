package test

import (
	"log"
	"testing"

	"github.com/miRemid/kira/common"
)

func TestGormTest(t *testing.T) {
	db, _ := common.DBConnect()
	var res = make(map[interface{}]interface{})
	if err := db.Raw(`select tf.owner, tf.file_name, tf.file_id, tf.file_width, tf.file_height, tf.anony 
	from tbl_file tf left join tbl_token_user ttu on tf.owner = ttu.user_id 
	where tf.anony != 1 
	and tf.id >= ((SELECT MAX(tf2.id) from tbl_file tf2) - (select MIN(tf3.id) from tbl_file tf3)) * RAND() + (select MIN(tu.id) from tbl_user tu)  
	limit 20`).Scan(&res).Error; err != nil {
		log.Println(err)
	} else {
		log.Println(res)
	}

}
