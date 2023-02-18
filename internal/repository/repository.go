package repository

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"makves-task/internal/entity"
)

type GormDB struct {
	db *gorm.DB
}

func NewGormLiteSQL() Databaser {
	r := &GormDB{}
	r.Init()
	return r
}

func (g *GormDB) Init() {
	file, err := os.Open("../makves-task/pkg/ueba.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var entries []entity.NodeGorm
	err = gocsv.Unmarshal(file, &entries)
	if err != nil {
		log.Fatal().Err(err).Msg("something goes wrong via unmarshal csv.")
	}
	g.db, err = gorm.Open(sqlite.Open("gorm.db"))
	if err != nil {
		log.Fatal().Err(err).Msg("something goes wrong via connection db")
	}
	if ok := g.db.Migrator().HasTable(&entity.NodeGorm{}); ok {
		err = g.db.Migrator().DropTable(&entity.NodeGorm{})
		if err != nil {
			log.Fatal().Err(err).Msg("something goes wrong via drop table")
		}
	}
	err = g.db.Migrator().CreateTable(&entity.NodeGorm{})
	if err != nil {
		log.Fatal().Err(err).Msg("something goes wrong migrate db")
	}

	result := g.db.CreateInBatches(entries, 100)
	if result.Error != nil {
		log.Fatal().Err(result.Error).Msg("something goes wrong via push csv to db")
	}
}

func (g *GormDB) GetItems(slice []string) []*entity.Node {
	result := make([]*entity.Node, 0)
	for _, v := range slice {
		node := &entity.Node{}
		g.db.Table("node_gorms").Select("id", "uid", "domain", "cn", "department",
			"title", "who", "logon_count", "num_logons7", "num_share7", "num_file7", "num_ad7", "num_n7",
			"num_logons14", "num_share14", "num_file14", "num_ad14", "num_n14", "num_logons30", "num_share30",
			"num_file30", "num_ad30", "num_n30", "num_logons150", "num_share150", "num_file150", "num_ad150",
			"num_n150", "num_logons365", "num_share365", "num_file365", "num_ad365", "num_n365", "has_user_principal_name",
			"has_mail", "has_phone", "flag_disabled", "flag_lockout", "flag_password_not_required", "flag_password_cant_change",
			"flag_dont_expire_password", "owned_files", "num_mailboxes", "num_member_of_groups", "num_member_of_indirect_groups",
			"member_of_indirect_groups_ids", "member_of_groups_id_s", "is_admin", "is_service").
			Where("id = ?", v).Scan(&node)
		result = append(result, node)
	}
	return result
}
