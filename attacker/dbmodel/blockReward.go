package dbmodel

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type BlockReward struct {
	ID             int64 `orm:"column(id)" db:"id" json:"id" form:"id"`                                                     //  任务类型id
	Epoch          int64 `orm:"column(epoch)" db:"epoch" json:"epoch" form:"epoch"`                                         // epoch
	ValidatorIndex int   `orm:"column(validator_index)" db:"validator_index" json:"validator_index" form:"validator_index"` // 验证者索引
	HeadAmount     int64 `orm:"column(head_amount)" db:"head_amount" json:"head_amount" form:"head_amount"`                 // Head 奖励数量
	TargetAmount   int64 `orm:"column(target_amount)" db:"target_amount" json:"target_amount" form:"target_amount"`         // Target 奖励数量
	//Head	Target	Source	Inclusion Delay	Inactivity
}

func (BlockReward) TableName() string {
	return "t_block_reward"
}

type BlockRewardRepository interface {
	Create(reward *BlockReward) error
	GetListByFilter(filters ...interface{}) []*BlockReward
}

type blockRewardRepositoryImpl struct {
	o orm.Ormer
}

func NewBlockRewardRepository(o orm.Ormer) BlockRewardRepository {
	return &blockRewardRepositoryImpl{o}
}

func (repo *blockRewardRepositoryImpl) Create(reward *BlockReward) error {
	_, err := repo.o.Insert(reward)
	return err
}

func (repo *blockRewardRepositoryImpl) GetListByFilter(filters ...interface{}) []*BlockReward {
	list := make([]*BlockReward, 0)
	query := repo.o.QueryTable(new(BlockReward).TableName())
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	query.OrderBy("-epoch").All(&list)
	return list
}

func GetRewardListByEpoch(epoch int64) []*BlockReward {
	filters := make([]interface{}, 0)
	filters = append(filters, "epoch", epoch)
	return NewBlockRewardRepository(orm.NewOrm()).GetListByFilter(filters...)
}

func GetRewardListByValidatorIndex(index int) []*BlockReward {
	filters := make([]interface{}, 0)
	filters = append(filters, "validator_index", index)
	return NewBlockRewardRepository(orm.NewOrm()).GetListByFilter(filters...)
}

func GetRewardByValidatorAndEpoch(epoch int64, index int) *BlockReward {
	filters := make([]interface{}, 0)
	filters = append(filters, "epoch", epoch)
	filters = append(filters, "validator_index", index)

	list := NewBlockRewardRepository(orm.NewOrm()).GetListByFilter(filters...)
	if len(list) >= 0 {
		return list[0]
	}
	return nil
}

func GetMaxEpoch() int64 {
	var max int64
	sql := fmt.Sprintf("select max(epoch) from %s", new(BlockReward).TableName())
	if err := orm.NewOrm().Raw(sql).QueryRow(&max); err == orm.ErrNoRows {
		return -1
	}
	return max
}
