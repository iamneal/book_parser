package server

import(
	"io/ioutil"
	"encoding/json"
	"fmt"
)

var (
	TOKEN_KEY string
	GRPC_GATEWAY_TOKEN string
	HTML_LOCATION string
	SECRET_LOCATION_NAME string
)

type Globals struct{
	TOKEN_KEY string
	GRPC_GATEWAY_TOKEN string
	HTML_LOCATION string
	SECRET_LOCATION_NAME string
}

func init() {
	bytes, err := ioutil.ReadFile("./server/globals.json")
	if err != nil {
		panic(fmt.Sprintf("could not read globals file: %s", err))
	}

	var globals Globals
	err = json.Unmarshal(bytes, &globals)
	if err != nil {
		panic(fmt.Sprintf("could not unmarshal json: %s", err))
	}
	fmt.Printf("our globals %+v\n", globals)
	TOKEN_KEY = globals.TOKEN_KEY
	GRPC_GATEWAY_TOKEN  = globals.GRPC_GATEWAY_TOKEN
	HTML_LOCATION = globals.HTML_LOCATION
	SECRET_LOCATION_NAME = globals.SECRET_LOCATION_NAME
}
