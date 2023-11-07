/**
 * bug595988 【外部反馈 111279-月度版】golang驱动中的sql.Out参数没有正确地被实现
 */

package main

import (
	"database/sql"
	_ "dm"
	"log"
)

func main() {
	//url := "dm://" + os.Getenv("dm_username") + ":" + os.Getenv("dm_password") + "@" + os.Getenv("dm_host") + "?noConvertToHex=true"
	conn, err := sql.Open("dm", "dm://SYSDBA:SYSDBA@192.168.100.168:7777")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	for _, sqlstr := range []string{
		"DROP TABLE IF EXISTS tmp_testabc",
		`CREATE TABLE  tmp_testabc (
			id bigint primary key,
			name varchar(12)
		)`,
	} {
		_, err := conn.Exec(sqlstr)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	var id int64
	sqlstr := "insert into tmp_testabc(id, name) values(1, 'a') returning id into :id"
	_, err = conn.Exec(sqlstr, sql.Out{Dest: &id})
	if err != nil {
		log.Fatal(err)
		return
	}
	if id != 1 {
		log.Fatal("Error cmp")
	}
	sqlstr = "insert into tmp_testabc(id, name) values(2, 'b') returning id into :id"
	_, err = conn.Exec(sqlstr, sql.Named("id123", sql.Out{Dest: &id}))
	if err != nil {
		log.Fatal(err)
		return
	}
	if id != 2 {
		log.Fatal("Error cmp")
	}
}
