package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	jsonStr := `{"rc":0,"rt":4,"svr":2887257261,"lt":1,"full":1,"dlmkts":"","data":{"f43":12.43,"f44":12.56,"f45":12.21,"f46":12.28,"f47":308763,"f48":382298576.0,"f50":0.77,"f51":13.64,"f52":11.16,"f55":0.287524125,"f57":"601238","f58":"广汽集团","f60":12.4,"f71":12.38,"f92":9.2931902,"f105":3008640272.0,"f116":130066993676.51,"f117":90282023753.45999,"f162":10.81,"f167":1.34,"f168":0.43,"f169":0.03,"f170":0.24,"f173":3.21,"f183":23267750729.0,"f186":6.1190936268,"f187":12.9383267312,"f188":35.5373915625,"f191":-25.48,"f192":-1863,"f193":0.31,"f194":1.54,"f195":-1.23,"f196":-0.71,"f197":0.41}}`

	//json str 转map
	var res map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &res); err == nil {
		fmt.Println("==============json str 转map=======================")
		/*fmt.Println(reflect.TypeOf(res)) //map[string]interface {}
		fmt.Println(res)
		fmt.Println(reflect.TypeOf(res["data"]))
		fmt.Println(res["data"])*/
		data := res["data"].(map[string]interface{})["f43"].(float64)
		/*for key, value := range data {
			fmt.Printf("%s : %s", key, value)
		}*/
		fmt.Printf("%.2f", data)

	}

}
