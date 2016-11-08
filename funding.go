package main

import "fmt"
import (
	r "github.com/dancannon/gorethink"
	"encoding/json"
	//"github.com/golang/protobuf/jsonpb/jsonpb_test_proto"
	"io/ioutil"
)
func main() {
//	------------------------Structure-------------------------------------------------------------------------------
	type Naming struct{
	Name string `json:"name,omitempty" gorethink:"name"`
	First_Name string`json:"first_name,omitempty" gorethink:"first_name"`
	Last_Name string`json:"last_name,omitempty" gorethink:"last_name"`
	Permalink string `json:"permalink,omitempty" gorethink:"permalink"`
}

	type Investments struct{
	Name string `json:"name"`
	Children Naming  `json:"children"`
	}


	type Children struct
	{
	Funded_Year float64 `json:"funded_year,omitempty"`
	Id float64 `json:"id,omitempty"`
	Name string `json:"name"`
	Investments [][]Investments  `json:"children"`
	}

	type Format struct {
		Name           string   `json:"name"`
		Children       []Children `json:"children"`
	}
//	------------------------Structure-Ends--------------------------------------------------------------------------





	//connectivity--------------------------
	s,e := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	if e != nil {
		fmt.Println("error is -----",e)
		return
	}
//	----------------------------------------------------------------------------------------------------------------
	c,e := r.DB("New").Table("companies").ConcatMap(func(fund r.Term ) r.Term  {
		return fund.Field("funding_rounds")
	}).Run(s)
//	----------------------------------------------------------------------------------------------------------------

	var m   []map[string]interface{}
	c.All(&m)
	//fmt.Println(m)
	c1:=Format{}
	c1.Name = "Funding Rounds"
	for _,v := range m  {
		c := Children{}

		if v["funded_year"] == nil {
			c.Funded_Year = 0
		}else {
			x := v["funded_year"].(float64)
			c.Funded_Year=x
		}
		if v["id"] ==nil {
			c.Id = 0
		}else {
			x:=v["id"].(float64)
			c.Id=x
		}
		c.Name = "Investments"

		if v["investments"] ==nil {
			c.Investments = nil
		}else {

			x:=v["investments"].([]interface{})    //
			for _,val := range x {
				i2 := []Investments{}
				v1, _ := val.(map[string]interface{})["company"]
				if v1 != nil {
					i3 := Investments{}
					i3.Name = "company"
					i3.Children.Name = v1.(map[string]interface{})["name"].(string)
					i3.Children.Permalink = v1.(map[string]interface{})["permalink"].(string)
					i2 = append(i2, i3)
				}

				v2, _ := val.(map[string]interface{})["financial_org"]
				if v2 != nil {
					i3 := Investments{}
					i3.Name = "financial_org"
					i3.Children.Name = v2.(map[string]interface{})["name"].(string)
					i3.Children.Permalink = v2.(map[string]interface{})["permalink"].(string)
					i2 = append(i2, i3)
				}
				v3, _ := val.(map[string]interface{})["person"]
				if v3 != nil {
					i3 := Investments{}
					i3.Name = "person"
					i3.Children.First_Name = v3.(map[string]interface{})["first_name"].(string)
					i3.Children.Last_Name = v3.(map[string]interface{})["last_name"].(string)
					i3.Children.Permalink = v3.(map[string]interface{})["permalink"].(string)
					i2 = append(i2, i3)
				}
				c.Investments = append(c.Investments,i2)
			}
		}
		c1.Children = append(c1.Children,c)
		}



	j,_:=json.Marshal(c1)
	//fmt.Println(string(j))
	f := ioutil.WriteFile("new.json",j,0644)
	if f != nil {
		fmt.Println(f)
	}
}
