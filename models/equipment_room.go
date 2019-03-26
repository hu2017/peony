package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type EquipmentRoom struct {
	Id         int       `orm:"column(id);auto" description:"序号"`
	RoomNo     string    `orm:"column(room_no);size(20)" description:"电房编号"`
	RoomName   string    `orm:"column(room_name);size(30)" description:"电房名称"`
	CustomerNo string    `orm:"column(customer_no);size(20);null" description:"客户编号"`
	Tag        int       `orm:"column(tag);null" description:"状态:0.正常1.删除"`
	Createuser string    `orm:"column(createuser);size(20);null" description:"创建人"`
	Createdate time.Time `orm:"column(createdate);type(datetime);null;auto_now_add" description:"创建时间"`
	Changeuser string    `orm:"column(changeuser);size(20);null" description:"修改人"`
	Changedate time.Time `orm:"column(changedate);type(datetime);null;auto_now_add" description:"修改时间"`
}

func (t *EquipmentRoom) TableName() string {
	return "equipment_room"
}

func init() {
	orm.RegisterModel(new(EquipmentRoom))
}

// AddEquipmentRoom insert a new EquipmentRoom into database and returns
// last inserted Id on success.
func AddEquipmentRoom(m *EquipmentRoom) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetEquipmentRoomById retrieves EquipmentRoom by Id. Returns error if
// Id doesn't exist
func GetEquipmentRoomById(id int) (v *EquipmentRoom, err error) {
	o := orm.NewOrm()
	v = &EquipmentRoom{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllEquipmentRoom retrieves all EquipmentRoom matches certain condition. Returns empty list if
// no records exist
func GetAllEquipmentRoom(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64, paging bool) (ml []interface{}, total int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(EquipmentRoom))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, 0, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, 0, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, 0, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, 0, errors.New("Error: unused 'order' fields")
		}
	}

	var l []EquipmentRoom
	qs = qs.OrderBy(sortFields...)
	if paging{
		total, _ = qs.Count()
	}

	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, total, nil
	}
	return nil, 0, err
}

// UpdateEquipmentRoom updates EquipmentRoom by Id and returns error if
// the record to be updated doesn't exist
func UpdateEquipmentRoomById(m *EquipmentRoom) (err error) {
	o := orm.NewOrm()
	v := EquipmentRoom{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteEquipmentRoom deletes EquipmentRoom by Id and returns error if
// the record to be deleted doesn't exist
func DeleteEquipmentRoom(id int) (err error) {
	o := orm.NewOrm()
	v := EquipmentRoom{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&EquipmentRoom{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
