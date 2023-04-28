package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	gosdk "github.com/okx/okbchain-go-sdk"
	"github.com/okx/okbchain-go-sdk/utils"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/query"
	"log"
)

const (
	rpcURL = "tcp://127.0.0.1:26657"
	// user's name
	name = "alice"
	// user's mnemonic
	mnemonic  = "giggle sibling fun arrow elevator spoon blood grocery laugh tortoise culture tool"
	mnemonic2 = "antique onion adult slot sad dizzy sure among cement demise submit scare"
	// user's password
	passWd = "12345678"
	// target address
	addr     = "ex1qj5c07sm6jetjz8f509qtrxgh4psxkv3ddyq7u"
	baseCoin = "okt"

	addr2   = "ex1fsfwwvl93qv6r56jpu084hxxzn9zphnyxhske5"
	chainId = "okbchain-67"
)

func main() {
	//-------------------- 1. preparation --------------------//
	// NOTE: either of the both ways below to pay fees is available

	// WAY 1: create a client config with fixed fees
	config, err := gosdk.NewClientConfig(rpcURL, "okbchain-67", gosdk.BroadcastBlock, "0.01okb", 100000000,
		0, "")
	if err != nil {
		log.Fatal(err)
	}

	// WAY 2: alternative client config with the fees by auto gas calculation
	//config, err = gosdk.NewClientConfig(rpcURL, "exchain-64", gosdk.BroadcastBlock, "", 200000,
	//	1.1, "0.00000000000001okt")
	//if err != nil {
	//	log.Fatal(err)
	//}

	cli := gosdk.NewClient(config)

	// create an account with your own mnemonic，name and password
	fromInfo, _, err := utils.CreateAccountWithMnemo(mnemonic, name, passWd)
	if err != nil {
		log.Fatal(err)
	}

	fromInfo2, _, err := utils.CreateAccountWithMnemo(mnemonic2, "bob", passWd)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("account 2 address: ", fromInfo2.GetAddress().String())

	//-------------------- 2. query for the information of your address --------------------//

	accInfo, err := cli.Auth().QueryAccount(fromInfo.GetAddress().String())
	if err != nil {
		log.Fatal(err)
	}

	log.Println(accInfo)

	//-------------------- 3. transfer to other address --------------------//
	accountNum, sequenceNum := accInfo.GetAccountNumber(), accInfo.GetSequence()

	wasmFile := "sample/wasm/hackatom.wasm"
	codeId, err := cli.Wasm().StoreCode(fromInfo, passWd, accountNum, sequenceNum, "memo", wasmFile, "", false, false)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("=============================================================StoreCode1===============================================================")
	log.Println(codeId)

	// instantiate a wasm contract
	log.Println("=========================================================InstantiateContract==========================================================")
	sequenceNum++
	address := common.Address{}
	address.SetBytes(fromInfo.GetAddress().Bytes())
	initMsg := fmt.Sprintf(`{"verifier": "%s", "beneficiary": "%s"}`, address.String(), address.String())
	contractAddr, err := cli.Wasm().InstantiateContract(fromInfo, passWd, accountNum, sequenceNum, "memo", uint64(codeId), initMsg, "1okb", "local0.1.0", "ex1qj5c07sm6jetjz8f509qtrxgh4psxkv3ddyq7u", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("contract address: ", contractAddr)

	// execute a wasm contract
	sequenceNum++
	execMsg := `{"release":{}}`
	executeRes, err := cli.Wasm().ExecuteContract(fromInfo, passWd, accountNum, sequenceNum, "memo", contractAddr, execMsg, "2okb")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("=========================================================ExecuteContract==========================================================")
	log.Println(executeRes)

	// set new admin for the contract
	sequenceNum++
	updateAdminRes, err := cli.Wasm().UpdateContractAdmin(fromInfo, passWd, accountNum, sequenceNum, "memo", contractAddr, addr2)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("=========================================================UpdateContractAdmin==========================================================")
	log.Println(updateAdminRes)

	// query contract state all
	contractStateAll, err := cli.Wasm().QueryContractStateAll(contractAddr, &query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      50,
		CountTotal: false,
	})
	log.Println("=========================================================QueryContractStateAll==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(contractStateAll)
	}

	// query contract state raw
	// 0006636f6e666967 is hex of the key "config"
	fmt.Println(hex.EncodeToString([]byte(`"config"`)))
	contractStateRaw, err := cli.Wasm().QueryContractStateRaw(contractAddr, "22636f6e66696722")
	log.Println("=========================================================QueryContractStateRaw==========================================================")
	if err != nil {
		log.Println("err:", err)
	} else {
		log.Println(contractStateRaw)
	}

	// query contract state smart
	contractStateSmart, err := cli.Wasm().QueryContractStateSmart(contractAddr, `{"verifier":{}}`)
	log.Println("=========================================================QueryContractStateSmart==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(contractStateSmart)
	}

	// migrate contract to the new code
	accInfo2, err := cli.Auth().QueryAccount(fromInfo2.GetAddress().String())
	if err != nil {
		log.Fatal(err)
	}

	log.Println(accInfo2)

	// store new code
	sequenceNum++
	migrateWasmFile := "sample/wasm/burner.wasm"
	codeId, err = cli.Wasm().StoreCode(fromInfo, passWd, accountNum, sequenceNum, "memo", migrateWasmFile, "", false, false)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("=============================================================StoreCode2===============================================================")
	log.Println(codeId)

	migrateMsg := `{"payout": "ex1fsfwwvl93qv6r56jpu084hxxzn9zphnyxhske5"}`
	migrateRes, err := cli.Wasm().MigrateContract(fromInfo2, passWd, accInfo2.GetAccountNumber(), accInfo2.GetSequence(), "memo", uint64(codeId), contractAddr, migrateMsg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("=========================================================MigrateContract==========================================================")
	log.Println(migrateRes)

	// clear admin for the contract
	clearAdminRes, err := cli.Wasm().ClearContractAdmin(fromInfo2, passWd, accInfo2.GetAccountNumber(), accInfo2.GetSequence()+1, "memo", contractAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("=========================================================ClearContractAdmin==========================================================")
	log.Println(clearAdminRes)

	// query code
	queryCodeRes, err := cli.Wasm().QueryCode(uint64(codeId))
	log.Println("=========================================================QueryCode==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(queryCodeRes.DataHash)
	}

	// 	query listCode
	listCodeRes, err := cli.Wasm().QueryListCode(&query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      50,
		CountTotal: false,
	})

	log.Println("=========================================================QueryListCode==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(len(listCodeRes.CodeInfos))
	}

	// query ListContractByCode
	listContract, err := cli.Wasm().QueryListContract(uint64(codeId), &query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      50,
		CountTotal: false,
	})

	log.Println("=========================================================QueryListContract==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(listContract)
	}

	// query code info
	codeInfo, err := cli.Wasm().QueryCodeInfo(uint64(codeId))
	log.Println("=========================================================QueryCodeInfo==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(codeInfo)
	}

	// query contract info
	contractInfo, err := cli.Wasm().QueryContractInfo(contractAddr)
	log.Println("=========================================================QueryContractInfo==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(contractInfo)
	}

	// query contract history
	contractHistory, err := cli.Wasm().QueryContractHistory(contractAddr, &query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      50,
		CountTotal: false,
	})
	log.Println("=========================================================QueryContractHistory==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(contractHistory)
	}

	// query contract ListPinnedCode
	pinnedCode, err := cli.Wasm().QueryListPinnedCode(&query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      50,
		CountTotal: false,
	})
	log.Println("=========================================================QueryListPinnedCode==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(pinnedCode)
	}
}
