package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"

	store "mint-token-app/contracts"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type MainController struct {
	beego.Controller
}

type TokenRequest struct {
	Receiver string `json:"receiver" valid:"Required;Match(/^0x[0-9a-fA-F]{40}$/)"`
	Amount   string `json:"amount" valid:"Required"`
}

func validateReq(c *MainController, req TokenRequest) {
	valid := validation.Validation{}
	b, validationErr := valid.Valid(&req)
	if validationErr != nil {
		logs.Error(validationErr)
		c.Data["json"] = map[string]interface{}{
			"error": validationErr,
		}
		c.ServeJSON()
	}
	if !b {
		for _, err := range valid.Errors {
			logs.Error(err.Key + err.Message)
			c.Data["json"] = map[string]interface{}{
				"error": err.Key + err.Message,
			}
			c.ServeJSON()
		}
	}
}

func (c *MainController) MintToken() {
	// read req body
	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"error": "request body reading error",
		}
		logs.Error("request body reading error")
		c.ServeJSON()
		return
	}

	var req TokenRequest
	// parse req body
	err = json.Unmarshal(body, &req)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"error": "request body parcing error",
		}
		logs.Error("request body parcing error")
		c.ServeJSON()
		return
	}

	// receiver and amount validation
	validateReq(c, req)

	mintAddress := common.HexToAddress(req.Receiver)

	client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		logs.Error(err)
	}

	address := common.HexToAddress("0x3e1d0ced18a4454ba390b8f540682c718748b0e5")
	amount := new(big.Int)
	_, success := amount.SetString(req.Amount, 0)
	if !success {
		fmt.Println("Error converting string to big integer")
		return
	}

	instance, instanceErr := store.NewTokenTransactor(address, client)
	if instanceErr != nil {
		logs.Error(instanceErr)
	}

	transactionOption := bind.TransactOpts{
		From:      mintAddress,            // Ethereum account to send the transaction from
		Nonce:     nil,                    // Nonce to use for the transaction execution (nil = use pending state)
		Signer:    nil,                    // Method to use for signing the transaction (mandatory)
		Value:     nil,                    // Funds to transfer along the transaction (nil = 0 = no funds)
		GasPrice:  big.NewInt(1000000000), // Gas price to use for the transaction execution (nil = gas price oracle)
		GasFeeCap: nil,                    // Gas fee cap to use for the 1559 transaction execution (nil = gas price oracle)
		GasTipCap: nil,                    // Gas priority fee cap to use for the 1559 transaction execution (nil = gas price oracle)
		GasLimit:  uint64(200000),         // Gas limit to set for the transaction execution (0 = estimate)
		Context:   nil,                    // Network context to support cancellation and timeouts (nil = no timeout)
		NoSend:    false,                  // Do all transact steps but do not send the transaction
	}

	response, mintErr := instance.Mint(&transactionOption, mintAddress, amount)

	if mintErr != nil {
		logs.Error(mintErr)
		c.Data["json"] = map[string]interface{}{
			"error": mintErr,
		}
		c.ServeJSON()
		return
	}

	fmt.Println("response", response)

	// send response
	c.Data["json"] = map[string]interface{}{
		"response": "response",
	}
	c.ServeJSON()
}
