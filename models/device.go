package models

import (
	"github.com/astaxie/beego/orm"
	"strings"
	"reflect"
	"fmt"
	"errors"
)

type Device struct {
	Id int `json:"id, omitempty" orm:"column(id);pk;unique"`
	DeviceName string `json:"device_name" orm:"column(device_name);unique;size(32)"`
	Address string `json:"address" orm:"column(address);size(50)"`
	Status int `json:"status" orm:"column(status);size(1)"`// 0: enabled, 1:disabled
	CreatedAt int `json:"created_at, omitempty" orm:"column(create_at);size(11)"`
	UpdatedAt int `json:"updated_at, omitempty" orm:"column(updated_at);size(11)"`
	Latitude string `json:"latitude, omitempty" orm:"column(latitude);size(12)"`
	Longitude string `json:"longitude, omitempty" orm:"column(longitude);size(12)"`
	User *User `orm:"rel(fk)"`
	AirAd []*AirAd `orm:"reverse(many)"` // 设置一对多的反向关系
}

func init() {
	orm.RegisterModel(new(Device))
}

// AddDevice insert a new Device into database and returns
// last inserted Id on success.
func AddDevice(m *Device) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetDeviceById retrieves Device by Id. Returns error if
// Id doesn't exist
func GetDeviceById(id int) (v *Device, err error) {
	o := orm.NewOrm()
	v = &Device{Id: id}
	if err = o.QueryTable(new(Device)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetDeviceByUser retrieves Device by User. Returns error if
// Id doesn't exist
func GetDevicesByUser(user *User, limit int) []*Device {
	o := orm.NewOrm()
	var device Device
	var devices []*Device
	o.QueryTable(device).Filter("User", user).OrderBy("-CreatedAt").Limit(limit).All(&devices)
	return devices
}

// GetAllDevices retrieves all Device matches certain condition. Returns empty list if
// no records exist
func GetAllDevices(query map[string]string, fields []string, sortby []string, order []string,
	offset int, limit int) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Device))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
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
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
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
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Device
	qs = qs.OrderBy(sortFields...).RelatedSel()
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
		return ml, nil
	}
	return nil, err
}

// UpdateDevice updates Device by Id and returns error if
// the record to be updated doesn't exist
func UpdateDeviceById(m *Device) (err error) {
	o := orm.NewOrm()
	v := Device{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteDevice deletes Device by Id and returns error if
// the record to be deleted doesn't exist
func DeleteDevice(id int) (err error) {
	o := orm.NewOrm()
	v := Device{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Device{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}