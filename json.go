package main

import	(
	"fmt"
	"github.com/dancannon/gorethink"
	"encoding/json"
	"io/ioutil"
)

func main() {

	//Connectivity part----------------
	s, err := gorethink.Connect(gorethink.ConnectOpts{
		Address : "localhost:28015",
		//Database: "New",
	})
	if err != nil {
		fmt.Println("The Error is  :----- ", err)
		return
	}

	//------------------------------***********************-----------------------------------------------------------

	//Structure----------------------------------------------------------------------------------------------------

	type Children struct {
		Acquired_year   float64 `json:"acquired_year"`
		Acquired_month float64  `json:"acquired_month"`
		Name string `json:"name"`
		Company  struct{
			Name string `json:"name"`
			Permalink string `json:"permalink"`
			 } `json:"children,omitempty"`
	}

	type Format struct {
		Name           string   `json:"name"`
		Children       []Children `json:"children"`
			}
	//---------------------------------***************************-------------------------------------------------



	//Query---------------------------------------------------------------------------------------------------
	q, e := gorethink.DB("New").Table("companies").ConcatMap(func(x gorethink.Term) gorethink.Term {
		return x.Field("acquisitions")
	}).Run(s)
//	-------------------------**************************------------------------------------------------------





	var m  []map[string]interface{}
	if e != nil {
		fmt.Println(err)
	}
	q.All(&m)
	//fmt.Println(m)
	o := Format{}
	o.Name = "Acquisitions"

	//f.Name = "Company"
	for _,i := range m{
		//fmt.Println(i["acquired_year"])
		//
		f := Children{}


		if i["acquired_year"]  == nil {
			f.Acquired_year = 0
		}else{
			x := i["acquired_year"].(float64)
			f.Acquired_year = x
		}


		if i["acquired_month"]  == nil {
			f.Acquired_month = 0
		}else{
			//fmt.Println(i["acquired_month"])

			x1 := i["acquired_month"].(float64)
			f.Acquired_month = x1
		}
		f.Name = "Company"

		if i["company"]  == nil {
			f.Company.Permalink = ""
			f.Company.Name = ""
		}else{
			f.Company.Name = i["company"].(map[string]interface{})["name"].(string)  //type cast the values of company
			f.Company.Permalink = i["company"].(map[string]interface{})["permalink"].(string)


		}


		o.Children = append(o.Children,f)



	}
	b, err := json.Marshal(o)
	if err  != nil{
		fmt.Println(err)
	}

	//fmt.Println(string(b))

	//fmt.Println(o)
	err = ioutil.WriteFile("output.txt", b, 0644)
	if err != nil {
		panic(err)
	}
}