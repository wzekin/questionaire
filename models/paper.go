package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

var (
	ErrTimeEarly error = errors.New("问卷还未发布")
	ErrTimeLate  error = errors.New("问卷已经结束")
)

type Paper struct {
	Id        int            `orm:"pk"`
	Title     string         `orm:"size(255)"`
	StartTime time.Time      `orm:"type(datetime)"`
	EndTime   time.Time      `orm:"type(datetime)"`
	Questions *[3][]Question `orm:"-"`
}

func (p *Paper) Read() error {
	o := orm.NewOrm()
	//_, err := o.Raw("select id,title,type from question where paper_id = ? order by id", p.Id).QueryRows(&p.Questions)
	err := p.RandomQuestion(1, 1, 1)
	if err != nil {
		return err
	}
	var options []*Option
	_, err = o.Raw("select `option`.id,`option`.context,`option`.question_id  from `option` "+
		"inner join question inner join paper "+
		"on paper.id=question.paper_id and question.id=`option`.question_id"+
		" where paper.id = ? order by `option`.question_id", p.Id).QueryRows(&options)
	i := 0
	j := 0
	for i < len(p.Questions) && j < len(options) {
		if p.Questions[i].QId == options[j].QuestionId {
			p.Questions[i].Options = append(p.Questions[i].Options, options[j])
			j++
		} else {
			i++
		}
	}
	return nil
}

func GetPaper(id int) (p Paper, err error) {
	p = Paper{}
	o := orm.NewOrm()
	err = o.QueryTable("paper").Filter("id", id).One(&p)
	if err != nil {
		return Paper{}, err
	}
	if p.StartTime.After(time.Now()) {
		return Paper{}, ErrTimeEarly
	} else if p.EndTime.Before(time.Now()) {
		return Paper{}, ErrTimeLate
	} else {
		return
	}
}

func (p *Paper) RandomQuestion(select_, judge, text int) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("select * from question where type = 0 and paper_id = ? order by rand() limit ?", p.Id, select_).QueryRows(&p.Questions[0])
	if err != nil {
		return err
	}
	_, err = o.Raw("select * from question where type = 1 and paper_id = ? order by rand() limit ?", p.Id, judge).QueryRows(&p.Questions[0])
	if err != nil {
		return err
	}
	_, err = o.Raw("select * from question where type = 2 and paper_id = ? order by rand() limit ?", p.Id, text).QueryRows(&p.Questions[0])
	if err != nil {
		return err
	}
	return
}
