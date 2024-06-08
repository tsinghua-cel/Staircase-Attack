package dbmodel

import (
	"github.com/astaxie/beego/orm"
	"github.com/tsinghua-cel/attacker-service/types"
)

type ChainReorg struct {
	ID                    int64  `orm:"column(id)" db:"id" json:"id" form:"id"`                                                                                         //  任务类型id
	Epoch                 int64  `orm:"column(epoch)" db:"epoch" json:"epoch" form:"epoch"`                                                                             // epoch
	Slot                  int64  `orm:"column(slot)" db:"slot" json:"slot" form:"slot"`                                                                                 // slot
	Depth                 int    `orm:"column(depth)" db:"depth" json:"depth" form:"depth"`                                                                             // depth
	OldBlockSlot          int64  `orm:"column(old_block_slot)" db:"old_block_slot" json:"old_block_slot" form:"old_block_slot"`                                         // old_block_slot
	NewBlockSlot          int64  `orm:"column(new_block_slot)" db:"new_block_slot" json:"new_block_slot" form:"new_block_slot"`                                         // new_block_slot
	OldBlockProposerIndex int64  `orm:"column(old_block_proposer_index)" db:"old_block_proposer_index" json:"old_block_proposer_index" form:"old_block_proposer_index"` // old_block_proposer_index
	NewBlockProposerIndex int64  `orm:"column(new_block_proposer_index)" db:"new_block_proposer_index" json:"new_block_proposer_index" form:"new_block_proposer_index"` // new_block_proposer_index
	OldHeadState          string `orm:"column(old_head_state)" db:"old_head_state" json:"old_head_state" form:"old_head_state"`                                         // old_head_state
	NewHeadState          string `orm:"column(new_head_state)" db:"new_head_state" json:"new_head_state" form:"new_head_state"`                                         // new_head_state
}

func (ChainReorg) TableName() string {
	return "t_chain_reorg"
}

type ChainReorgRepository interface {
	Create(reorg *ChainReorg) error
	GetListByFilter(filters ...interface{}) []*ChainReorg
}

type chainReorgRepositoryImpl struct {
	o orm.Ormer
}

func NewChainReorgRepository(o orm.Ormer) ChainReorgRepository {
	return &chainReorgRepositoryImpl{o}
}

func (repo *chainReorgRepositoryImpl) Create(reorg *ChainReorg) error {
	_, err := repo.o.Insert(reorg)
	return err
}

func (repo *chainReorgRepositoryImpl) GetListByFilter(filters ...interface{}) []*ChainReorg {
	list := make([]*ChainReorg, 0)
	query := repo.o.QueryTable(new(ChainReorg).TableName())
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	query.OrderBy("-slot").All(&list)
	return list
}

func InsertNewReorg(ev types.ReorgEvent) {
	NewChainReorgRepository(orm.NewOrm()).Create(&ChainReorg{
		Epoch:                 int64(ev.Epoch),
		Slot:                  int64(ev.Slot),
		Depth:                 int(ev.Depth),
		OldBlockSlot:          ev.OldBlockSlot,
		NewBlockSlot:          ev.NewBlockSlot,
		OldBlockProposerIndex: ev.OldBlockProposerIndex,
		NewBlockProposerIndex: ev.NewBlockProposerIndex,
		OldHeadState:          ev.OldHeadState,
		NewHeadState:          ev.NewHeadState,
	})
}

func GetAllReorgList() []*ChainReorg {
	return NewChainReorgRepository(orm.NewOrm()).GetListByFilter()
}
