package database

import (
	"context"
	"gorm.io/gorm"
)

const IteratorLimit = 1000

type Iterator struct {
	query  *gorm.DB
	branch interface{}
	count  int
	offset int
	limit  int
	index  int
}

func NewIterator(query *gorm.DB, branch interface{}, limit ...int) (it *Iterator, err error) {
	it = &Iterator{
		branch: branch,
		query:  query,
		limit:  IteratorLimit,
		offset: 0,
		index:  -1,
	}
	if len(limit) > 0 {
		it.limit = limit[0]
	}
	err = it.nextBranch()
	return
}

func (it *Iterator) nextBranch() error {

	if err := it.query.WithContext(context.Background()).Limit(it.limit).Offset(it.offset).Find(it.branch).Error; err != nil {
		return err
	}
	var count int64
	if err := it.query.WithContext(context.Background()).Model(it.branch).Count(&count).Error; err != nil {
		return err
	}
	if c := int(count) - it.offset; c < it.limit {
		it.count = c
	} else {
		it.count = it.limit
	}
	it.offset += it.limit
	it.index = -1
	return nil
}

func (it *Iterator) Next() (ok bool, err error) {
	if it.index+1 == it.limit {
		if err = it.saveBranch(); err != nil {
			return
		}
		if err = it.nextBranch(); err != nil {
			return
		}
	}
	if it.index+1 == it.count {
		if it.count > 0 {
			if err = it.saveBranch(); err != nil {
				return
			}
		}
		return false, nil
	}
	it.index++
	return true, nil
}

func (it *Iterator) Index() interface{} {
	return it.index
}

func (it *Iterator) saveBranch() error {
	return it.query.Save(it.branch).Error
}
