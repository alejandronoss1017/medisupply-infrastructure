// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package binding

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// SLAEnforcerCustomerContract is an auto generated low-level Go binding around an user-defined struct.
type SLAEnforcerCustomerContract struct {
	Id         string
	Path       string
	CustomerId string
	Slas       []SLAEnforcerSLA
}

// SLAEnforcerSLA is an auto generated low-level Go binding around an user-defined struct.
type SLAEnforcerSLA struct {
	Id          string
	Name        string
	Description string
	Target      *big.Int
	Comparator  uint8
	Status      uint8
}

// SLAEnforcerMetaData contains all meta data concerning the SLAEnforcer contract.
var SLAEnforcerMetaData = &bind.MetaData{
	ABI: "[{\"name\":\"ContractAdded\",\"type\":\"event\",\"inputs\":[{\"name\":\"contractId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"customerId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"name\":\"SLAAdded\",\"type\":\"event\",\"inputs\":[{\"name\":\"contractId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"slaId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"name\":\"SLAStatusUpdated\",\"type\":\"event\",\"inputs\":[{\"name\":\"contractId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"slaId\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"newStatus\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumSLAEnforcer.SLAStatus\"}],\"anonymous\":false},{\"name\":\"addContract\",\"type\":\"function\",\"inputs\":[{\"name\":\"_id\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_path\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_customerId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"name\":\"addSLA\",\"type\":\"function\",\"inputs\":[{\"name\":\"_contractId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_slaId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_target\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_comparator\",\"type\":\"uint8\",\"internalType\":\"enumSLAEnforcer.Comparator\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"name\":\"checkSLA\",\"type\":\"function\",\"inputs\":[{\"name\":\"_contractId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_slaId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_actualValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"name\":\"contracts\",\"type\":\"function\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"path\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"customerId\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"name\":\"getAllContracts\",\"type\":\"function\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"components\":[{\"name\":\"id\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"path\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"customerId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"slas\",\"type\":\"tuple[]\",\"components\":[{\"name\":\"id\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"target\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"comparator\",\"type\":\"uint8\",\"internalType\":\"enumSLAEnforcer.Comparator\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumSLAEnforcer.SLAStatus\"}],\"internalType\":\"structSLAEnforcer.SLA[]\"}],\"internalType\":\"structSLAEnforcer.CustomerContract[]\"}],\"stateMutability\":\"view\"},{\"name\":\"getContract\",\"type\":\"function\",\"inputs\":[{\"name\":\"_id\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"path\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"customerId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"slas\",\"type\":\"tuple[]\",\"components\":[{\"name\":\"id\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"target\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"comparator\",\"type\":\"uint8\",\"internalType\":\"enumSLAEnforcer.Comparator\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumSLAEnforcer.SLAStatus\"}],\"internalType\":\"structSLAEnforcer.SLA[]\"}],\"stateMutability\":\"view\"},{\"name\":\"getContractCount\",\"type\":\"function\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"name\":\"getSLA\",\"type\":\"function\",\"inputs\":[{\"name\":\"_contractId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_slaIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"target\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"comparator\",\"type\":\"uint8\",\"internalType\":\"enumSLAEnforcer.Comparator\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumSLAEnforcer.SLAStatus\"}],\"stateMutability\":\"view\"},{\"name\":\"getSLAById\",\"type\":\"function\",\"inputs\":[{\"name\":\"_contractId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_slaId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"id\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"target\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"comparator\",\"type\":\"uint8\",\"internalType\":\"enumSLAEnforcer.Comparator\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumSLAEnforcer.SLAStatus\"}],\"stateMutability\":\"view\"},{\"name\":\"getSLAs\",\"type\":\"function\",\"inputs\":[{\"name\":\"_contractId\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"components\":[{\"name\":\"id\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"target\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"comparator\",\"type\":\"uint8\",\"internalType\":\"enumSLAEnforcer.Comparator\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumSLAEnforcer.SLAStatus\"}],\"internalType\":\"structSLAEnforcer.SLA[]\"}],\"stateMutability\":\"view\"},{\"name\":\"setSLAStatus\",\"type\":\"function\",\"inputs\":[{\"name\":\"_contractId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_slaIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_status\",\"type\":\"uint8\",\"internalType\":\"enumSLAEnforcer.SLAStatus\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"name\":\"setSLAStatusById\",\"type\":\"function\",\"inputs\":[{\"name\":\"_contractId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_slaId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_status\",\"type\":\"uint8\",\"internalType\":\"enumSLAEnforcer.SLAStatus\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
}

// SLAEnforcerABI is the input ABI used to generate the binding from.
// Deprecated: Use SLAEnforcerMetaData.ABI instead.
var SLAEnforcerABI = SLAEnforcerMetaData.ABI

// SLAEnforcer is an auto generated Go binding around an Ethereum contract.
type SLAEnforcer struct {
	SLAEnforcerCaller     // Read-only binding to the contract
	SLAEnforcerTransactor // Write-only binding to the contract
	SLAEnforcerFilterer   // Log filterer for contract events
}

// SLAEnforcerCaller is an auto generated read-only Go binding around an Ethereum contract.
type SLAEnforcerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SLAEnforcerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SLAEnforcerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SLAEnforcerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SLAEnforcerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SLAEnforcerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SLAEnforcerSession struct {
	Contract     *SLAEnforcer      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SLAEnforcerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SLAEnforcerCallerSession struct {
	Contract *SLAEnforcerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// SLAEnforcerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SLAEnforcerTransactorSession struct {
	Contract     *SLAEnforcerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// SLAEnforcerRaw is an auto generated low-level Go binding around an Ethereum contract.
type SLAEnforcerRaw struct {
	Contract *SLAEnforcer // Generic contract binding to access the raw methods on
}

// SLAEnforcerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SLAEnforcerCallerRaw struct {
	Contract *SLAEnforcerCaller // Generic read-only contract binding to access the raw methods on
}

// SLAEnforcerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SLAEnforcerTransactorRaw struct {
	Contract *SLAEnforcerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSLAEnforcer creates a new instance of SLAEnforcer, bound to a specific deployed contract.
func NewSLAEnforcer(address common.Address, backend bind.ContractBackend) (*SLAEnforcer, error) {
	contract, err := bindSLAEnforcer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SLAEnforcer{SLAEnforcerCaller: SLAEnforcerCaller{contract: contract}, SLAEnforcerTransactor: SLAEnforcerTransactor{contract: contract}, SLAEnforcerFilterer: SLAEnforcerFilterer{contract: contract}}, nil
}

// NewSLAEnforcerCaller creates a new read-only instance of SLAEnforcer, bound to a specific deployed contract.
func NewSLAEnforcerCaller(address common.Address, caller bind.ContractCaller) (*SLAEnforcerCaller, error) {
	contract, err := bindSLAEnforcer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SLAEnforcerCaller{contract: contract}, nil
}

// NewSLAEnforcerTransactor creates a new write-only instance of SLAEnforcer, bound to a specific deployed contract.
func NewSLAEnforcerTransactor(address common.Address, transactor bind.ContractTransactor) (*SLAEnforcerTransactor, error) {
	contract, err := bindSLAEnforcer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SLAEnforcerTransactor{contract: contract}, nil
}

// NewSLAEnforcerFilterer creates a new log filterer instance of SLAEnforcer, bound to a specific deployed contract.
func NewSLAEnforcerFilterer(address common.Address, filterer bind.ContractFilterer) (*SLAEnforcerFilterer, error) {
	contract, err := bindSLAEnforcer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SLAEnforcerFilterer{contract: contract}, nil
}

// bindSLAEnforcer binds a generic wrapper to an already deployed contract.
func bindSLAEnforcer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SLAEnforcerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SLAEnforcer *SLAEnforcerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SLAEnforcer.Contract.SLAEnforcerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SLAEnforcer *SLAEnforcerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SLAEnforcer.Contract.SLAEnforcerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SLAEnforcer *SLAEnforcerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SLAEnforcer.Contract.SLAEnforcerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SLAEnforcer *SLAEnforcerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SLAEnforcer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SLAEnforcer *SLAEnforcerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SLAEnforcer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SLAEnforcer *SLAEnforcerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SLAEnforcer.Contract.contract.Transact(opts, method, params...)
}

// Contracts is a free data retrieval call binding the contract method 0x474da79a.
//
// Solidity: function contracts(uint256 ) view returns(string id, string path, string customerId)
func (_SLAEnforcer *SLAEnforcerCaller) Contracts(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Id         string
	Path       string
	CustomerId string
}, error) {
	var out []interface{}
	err := _SLAEnforcer.contract.Call(opts, &out, "contracts", arg0)

	outstruct := new(struct {
		Id         string
		Path       string
		CustomerId string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Path = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.CustomerId = *abi.ConvertType(out[2], new(string)).(*string)

	return *outstruct, err

}

// Contracts is a free data retrieval call binding the contract method 0x474da79a.
//
// Solidity: function contracts(uint256 ) view returns(string id, string path, string customerId)
func (_SLAEnforcer *SLAEnforcerSession) Contracts(arg0 *big.Int) (struct {
	Id         string
	Path       string
	CustomerId string
}, error) {
	return _SLAEnforcer.Contract.Contracts(&_SLAEnforcer.CallOpts, arg0)
}

// Contracts is a free data retrieval call binding the contract method 0x474da79a.
//
// Solidity: function contracts(uint256 ) view returns(string id, string path, string customerId)
func (_SLAEnforcer *SLAEnforcerCallerSession) Contracts(arg0 *big.Int) (struct {
	Id         string
	Path       string
	CustomerId string
}, error) {
	return _SLAEnforcer.Contract.Contracts(&_SLAEnforcer.CallOpts, arg0)
}

// GetAllContracts is a free data retrieval call binding the contract method 0x18d3ce96.
//
// Solidity: function getAllContracts() view returns((string,string,string,(string,string,string,uint256,uint8,uint8)[])[])
func (_SLAEnforcer *SLAEnforcerCaller) GetAllContracts(opts *bind.CallOpts) ([]SLAEnforcerCustomerContract, error) {
	var out []interface{}
	err := _SLAEnforcer.contract.Call(opts, &out, "getAllContracts")

	if err != nil {
		return *new([]SLAEnforcerCustomerContract), err
	}

	out0 := *abi.ConvertType(out[0], new([]SLAEnforcerCustomerContract)).(*[]SLAEnforcerCustomerContract)

	return out0, err

}

// GetAllContracts is a free data retrieval call binding the contract method 0x18d3ce96.
//
// Solidity: function getAllContracts() view returns((string,string,string,(string,string,string,uint256,uint8,uint8)[])[])
func (_SLAEnforcer *SLAEnforcerSession) GetAllContracts() ([]SLAEnforcerCustomerContract, error) {
	return _SLAEnforcer.Contract.GetAllContracts(&_SLAEnforcer.CallOpts)
}

// GetAllContracts is a free data retrieval call binding the contract method 0x18d3ce96.
//
// Solidity: function getAllContracts() view returns((string,string,string,(string,string,string,uint256,uint8,uint8)[])[])
func (_SLAEnforcer *SLAEnforcerCallerSession) GetAllContracts() ([]SLAEnforcerCustomerContract, error) {
	return _SLAEnforcer.Contract.GetAllContracts(&_SLAEnforcer.CallOpts)
}

// GetContract is a free data retrieval call binding the contract method 0x35817773.
//
// Solidity: function getContract(string _id) view returns(string id, string path, string customerId, (string,string,string,uint256,uint8,uint8)[] slas)
func (_SLAEnforcer *SLAEnforcerCaller) GetContract(opts *bind.CallOpts, _id string) (struct {
	Id         string
	Path       string
	CustomerId string
	Slas       []SLAEnforcerSLA
}, error) {
	var out []interface{}
	err := _SLAEnforcer.contract.Call(opts, &out, "getContract", _id)

	outstruct := new(struct {
		Id         string
		Path       string
		CustomerId string
		Slas       []SLAEnforcerSLA
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Path = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.CustomerId = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.Slas = *abi.ConvertType(out[3], new([]SLAEnforcerSLA)).(*[]SLAEnforcerSLA)

	return *outstruct, err

}

// GetContract is a free data retrieval call binding the contract method 0x35817773.
//
// Solidity: function getContract(string _id) view returns(string id, string path, string customerId, (string,string,string,uint256,uint8,uint8)[] slas)
func (_SLAEnforcer *SLAEnforcerSession) GetContract(_id string) (struct {
	Id         string
	Path       string
	CustomerId string
	Slas       []SLAEnforcerSLA
}, error) {
	return _SLAEnforcer.Contract.GetContract(&_SLAEnforcer.CallOpts, _id)
}

// GetContract is a free data retrieval call binding the contract method 0x35817773.
//
// Solidity: function getContract(string _id) view returns(string id, string path, string customerId, (string,string,string,uint256,uint8,uint8)[] slas)
func (_SLAEnforcer *SLAEnforcerCallerSession) GetContract(_id string) (struct {
	Id         string
	Path       string
	CustomerId string
	Slas       []SLAEnforcerSLA
}, error) {
	return _SLAEnforcer.Contract.GetContract(&_SLAEnforcer.CallOpts, _id)
}

// GetContractCount is a free data retrieval call binding the contract method 0x9399869d.
//
// Solidity: function getContractCount() view returns(uint256)
func (_SLAEnforcer *SLAEnforcerCaller) GetContractCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SLAEnforcer.contract.Call(opts, &out, "getContractCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetContractCount is a free data retrieval call binding the contract method 0x9399869d.
//
// Solidity: function getContractCount() view returns(uint256)
func (_SLAEnforcer *SLAEnforcerSession) GetContractCount() (*big.Int, error) {
	return _SLAEnforcer.Contract.GetContractCount(&_SLAEnforcer.CallOpts)
}

// GetContractCount is a free data retrieval call binding the contract method 0x9399869d.
//
// Solidity: function getContractCount() view returns(uint256)
func (_SLAEnforcer *SLAEnforcerCallerSession) GetContractCount() (*big.Int, error) {
	return _SLAEnforcer.Contract.GetContractCount(&_SLAEnforcer.CallOpts)
}

// GetSLA is a free data retrieval call binding the contract method 0xadbfd8ed.
//
// Solidity: function getSLA(string _contractId, uint256 _slaIndex) view returns(string id, string name, string description, uint256 target, uint8 comparator, uint8 status)
func (_SLAEnforcer *SLAEnforcerCaller) GetSLA(opts *bind.CallOpts, _contractId string, _slaIndex *big.Int) (struct {
	Id          string
	Name        string
	Description string
	Target      *big.Int
	Comparator  uint8
	Status      uint8
}, error) {
	var out []interface{}
	err := _SLAEnforcer.contract.Call(opts, &out, "getSLA", _contractId, _slaIndex)

	outstruct := new(struct {
		Id          string
		Name        string
		Description string
		Target      *big.Int
		Comparator  uint8
		Status      uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Description = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.Target = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Comparator = *abi.ConvertType(out[4], new(uint8)).(*uint8)
	outstruct.Status = *abi.ConvertType(out[5], new(uint8)).(*uint8)

	return *outstruct, err

}

// GetSLA is a free data retrieval call binding the contract method 0xadbfd8ed.
//
// Solidity: function getSLA(string _contractId, uint256 _slaIndex) view returns(string id, string name, string description, uint256 target, uint8 comparator, uint8 status)
func (_SLAEnforcer *SLAEnforcerSession) GetSLA(_contractId string, _slaIndex *big.Int) (struct {
	Id          string
	Name        string
	Description string
	Target      *big.Int
	Comparator  uint8
	Status      uint8
}, error) {
	return _SLAEnforcer.Contract.GetSLA(&_SLAEnforcer.CallOpts, _contractId, _slaIndex)
}

// GetSLA is a free data retrieval call binding the contract method 0xadbfd8ed.
//
// Solidity: function getSLA(string _contractId, uint256 _slaIndex) view returns(string id, string name, string description, uint256 target, uint8 comparator, uint8 status)
func (_SLAEnforcer *SLAEnforcerCallerSession) GetSLA(_contractId string, _slaIndex *big.Int) (struct {
	Id          string
	Name        string
	Description string
	Target      *big.Int
	Comparator  uint8
	Status      uint8
}, error) {
	return _SLAEnforcer.Contract.GetSLA(&_SLAEnforcer.CallOpts, _contractId, _slaIndex)
}

// GetSLAById is a free data retrieval call binding the contract method 0x962189e3.
//
// Solidity: function getSLAById(string _contractId, string _slaId) view returns(string id, string name, string description, uint256 target, uint8 comparator, uint8 status)
func (_SLAEnforcer *SLAEnforcerCaller) GetSLAById(opts *bind.CallOpts, _contractId string, _slaId string) (struct {
	Id          string
	Name        string
	Description string
	Target      *big.Int
	Comparator  uint8
	Status      uint8
}, error) {
	var out []interface{}
	err := _SLAEnforcer.contract.Call(opts, &out, "getSLAById", _contractId, _slaId)

	outstruct := new(struct {
		Id          string
		Name        string
		Description string
		Target      *big.Int
		Comparator  uint8
		Status      uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Description = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.Target = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Comparator = *abi.ConvertType(out[4], new(uint8)).(*uint8)
	outstruct.Status = *abi.ConvertType(out[5], new(uint8)).(*uint8)

	return *outstruct, err

}

// GetSLAById is a free data retrieval call binding the contract method 0x962189e3.
//
// Solidity: function getSLAById(string _contractId, string _slaId) view returns(string id, string name, string description, uint256 target, uint8 comparator, uint8 status)
func (_SLAEnforcer *SLAEnforcerSession) GetSLAById(_contractId string, _slaId string) (struct {
	Id          string
	Name        string
	Description string
	Target      *big.Int
	Comparator  uint8
	Status      uint8
}, error) {
	return _SLAEnforcer.Contract.GetSLAById(&_SLAEnforcer.CallOpts, _contractId, _slaId)
}

// GetSLAById is a free data retrieval call binding the contract method 0x962189e3.
//
// Solidity: function getSLAById(string _contractId, string _slaId) view returns(string id, string name, string description, uint256 target, uint8 comparator, uint8 status)
func (_SLAEnforcer *SLAEnforcerCallerSession) GetSLAById(_contractId string, _slaId string) (struct {
	Id          string
	Name        string
	Description string
	Target      *big.Int
	Comparator  uint8
	Status      uint8
}, error) {
	return _SLAEnforcer.Contract.GetSLAById(&_SLAEnforcer.CallOpts, _contractId, _slaId)
}

// GetSLAs is a free data retrieval call binding the contract method 0x37e67a79.
//
// Solidity: function getSLAs(string _contractId) view returns((string,string,string,uint256,uint8,uint8)[])
func (_SLAEnforcer *SLAEnforcerCaller) GetSLAs(opts *bind.CallOpts, _contractId string) ([]SLAEnforcerSLA, error) {
	var out []interface{}
	err := _SLAEnforcer.contract.Call(opts, &out, "getSLAs", _contractId)

	if err != nil {
		return *new([]SLAEnforcerSLA), err
	}

	out0 := *abi.ConvertType(out[0], new([]SLAEnforcerSLA)).(*[]SLAEnforcerSLA)

	return out0, err

}

// GetSLAs is a free data retrieval call binding the contract method 0x37e67a79.
//
// Solidity: function getSLAs(string _contractId) view returns((string,string,string,uint256,uint8,uint8)[])
func (_SLAEnforcer *SLAEnforcerSession) GetSLAs(_contractId string) ([]SLAEnforcerSLA, error) {
	return _SLAEnforcer.Contract.GetSLAs(&_SLAEnforcer.CallOpts, _contractId)
}

// GetSLAs is a free data retrieval call binding the contract method 0x37e67a79.
//
// Solidity: function getSLAs(string _contractId) view returns((string,string,string,uint256,uint8,uint8)[])
func (_SLAEnforcer *SLAEnforcerCallerSession) GetSLAs(_contractId string) ([]SLAEnforcerSLA, error) {
	return _SLAEnforcer.Contract.GetSLAs(&_SLAEnforcer.CallOpts, _contractId)
}

// AddContract is a paid mutator transaction binding the contract method 0x6c2bc72d.
//
// Solidity: function addContract(string _id, string _path, string _customerId) returns()
func (_SLAEnforcer *SLAEnforcerTransactor) AddContract(opts *bind.TransactOpts, _id string, _path string, _customerId string) (*types.Transaction, error) {
	return _SLAEnforcer.contract.Transact(opts, "addContract", _id, _path, _customerId)
}

// AddContract is a paid mutator transaction binding the contract method 0x6c2bc72d.
//
// Solidity: function addContract(string _id, string _path, string _customerId) returns()
func (_SLAEnforcer *SLAEnforcerSession) AddContract(_id string, _path string, _customerId string) (*types.Transaction, error) {
	return _SLAEnforcer.Contract.AddContract(&_SLAEnforcer.TransactOpts, _id, _path, _customerId)
}

// AddContract is a paid mutator transaction binding the contract method 0x6c2bc72d.
//
// Solidity: function addContract(string _id, string _path, string _customerId) returns()
func (_SLAEnforcer *SLAEnforcerTransactorSession) AddContract(_id string, _path string, _customerId string) (*types.Transaction, error) {
	return _SLAEnforcer.Contract.AddContract(&_SLAEnforcer.TransactOpts, _id, _path, _customerId)
}

// AddSLA is a paid mutator transaction binding the contract method 0x082b9f5b.
//
// Solidity: function addSLA(string _contractId, string _slaId, string _name, string _description, uint256 _target, uint8 _comparator) returns()
func (_SLAEnforcer *SLAEnforcerTransactor) AddSLA(opts *bind.TransactOpts, _contractId string, _slaId string, _name string, _description string, _target *big.Int, _comparator uint8) (*types.Transaction, error) {
	return _SLAEnforcer.contract.Transact(opts, "addSLA", _contractId, _slaId, _name, _description, _target, _comparator)
}

// AddSLA is a paid mutator transaction binding the contract method 0x082b9f5b.
//
// Solidity: function addSLA(string _contractId, string _slaId, string _name, string _description, uint256 _target, uint8 _comparator) returns()
func (_SLAEnforcer *SLAEnforcerSession) AddSLA(_contractId string, _slaId string, _name string, _description string, _target *big.Int, _comparator uint8) (*types.Transaction, error) {
	return _SLAEnforcer.Contract.AddSLA(&_SLAEnforcer.TransactOpts, _contractId, _slaId, _name, _description, _target, _comparator)
}

// AddSLA is a paid mutator transaction binding the contract method 0x082b9f5b.
//
// Solidity: function addSLA(string _contractId, string _slaId, string _name, string _description, uint256 _target, uint8 _comparator) returns()
func (_SLAEnforcer *SLAEnforcerTransactorSession) AddSLA(_contractId string, _slaId string, _name string, _description string, _target *big.Int, _comparator uint8) (*types.Transaction, error) {
	return _SLAEnforcer.Contract.AddSLA(&_SLAEnforcer.TransactOpts, _contractId, _slaId, _name, _description, _target, _comparator)
}

// CheckSLA is a paid mutator transaction binding the contract method 0x16cd2417.
//
// Solidity: function checkSLA(string _contractId, string _slaId, uint256 _actualValue) returns()
func (_SLAEnforcer *SLAEnforcerTransactor) CheckSLA(opts *bind.TransactOpts, _contractId string, _slaId string, _actualValue *big.Int) (*types.Transaction, error) {
	return _SLAEnforcer.contract.Transact(opts, "checkSLA", _contractId, _slaId, _actualValue)
}

// CheckSLA is a paid mutator transaction binding the contract method 0x16cd2417.
//
// Solidity: function checkSLA(string _contractId, string _slaId, uint256 _actualValue) returns()
func (_SLAEnforcer *SLAEnforcerSession) CheckSLA(_contractId string, _slaId string, _actualValue *big.Int) (*types.Transaction, error) {
	return _SLAEnforcer.Contract.CheckSLA(&_SLAEnforcer.TransactOpts, _contractId, _slaId, _actualValue)
}

// CheckSLA is a paid mutator transaction binding the contract method 0x16cd2417.
//
// Solidity: function checkSLA(string _contractId, string _slaId, uint256 _actualValue) returns()
func (_SLAEnforcer *SLAEnforcerTransactorSession) CheckSLA(_contractId string, _slaId string, _actualValue *big.Int) (*types.Transaction, error) {
	return _SLAEnforcer.Contract.CheckSLA(&_SLAEnforcer.TransactOpts, _contractId, _slaId, _actualValue)
}

// SetSLAStatus is a paid mutator transaction binding the contract method 0xa09a51b3.
//
// Solidity: function setSLAStatus(string _contractId, uint256 _slaIndex, uint8 _status) returns()
func (_SLAEnforcer *SLAEnforcerTransactor) SetSLAStatus(opts *bind.TransactOpts, _contractId string, _slaIndex *big.Int, _status uint8) (*types.Transaction, error) {
	return _SLAEnforcer.contract.Transact(opts, "setSLAStatus", _contractId, _slaIndex, _status)
}

// SetSLAStatus is a paid mutator transaction binding the contract method 0xa09a51b3.
//
// Solidity: function setSLAStatus(string _contractId, uint256 _slaIndex, uint8 _status) returns()
func (_SLAEnforcer *SLAEnforcerSession) SetSLAStatus(_contractId string, _slaIndex *big.Int, _status uint8) (*types.Transaction, error) {
	return _SLAEnforcer.Contract.SetSLAStatus(&_SLAEnforcer.TransactOpts, _contractId, _slaIndex, _status)
}

// SetSLAStatus is a paid mutator transaction binding the contract method 0xa09a51b3.
//
// Solidity: function setSLAStatus(string _contractId, uint256 _slaIndex, uint8 _status) returns()
func (_SLAEnforcer *SLAEnforcerTransactorSession) SetSLAStatus(_contractId string, _slaIndex *big.Int, _status uint8) (*types.Transaction, error) {
	return _SLAEnforcer.Contract.SetSLAStatus(&_SLAEnforcer.TransactOpts, _contractId, _slaIndex, _status)
}

// SetSLAStatusById is a paid mutator transaction binding the contract method 0x67c63392.
//
// Solidity: function setSLAStatusById(string _contractId, string _slaId, uint8 _status) returns()
func (_SLAEnforcer *SLAEnforcerTransactor) SetSLAStatusById(opts *bind.TransactOpts, _contractId string, _slaId string, _status uint8) (*types.Transaction, error) {
	return _SLAEnforcer.contract.Transact(opts, "setSLAStatusById", _contractId, _slaId, _status)
}

// SetSLAStatusById is a paid mutator transaction binding the contract method 0x67c63392.
//
// Solidity: function setSLAStatusById(string _contractId, string _slaId, uint8 _status) returns()
func (_SLAEnforcer *SLAEnforcerSession) SetSLAStatusById(_contractId string, _slaId string, _status uint8) (*types.Transaction, error) {
	return _SLAEnforcer.Contract.SetSLAStatusById(&_SLAEnforcer.TransactOpts, _contractId, _slaId, _status)
}

// SetSLAStatusById is a paid mutator transaction binding the contract method 0x67c63392.
//
// Solidity: function setSLAStatusById(string _contractId, string _slaId, uint8 _status) returns()
func (_SLAEnforcer *SLAEnforcerTransactorSession) SetSLAStatusById(_contractId string, _slaId string, _status uint8) (*types.Transaction, error) {
	return _SLAEnforcer.Contract.SetSLAStatusById(&_SLAEnforcer.TransactOpts, _contractId, _slaId, _status)
}

// SLAEnforcerContractAddedIterator is returned from FilterContractAdded and is used to iterate over the raw logs and unpacked data for ContractAdded events raised by the SLAEnforcer contract.
type SLAEnforcerContractAddedIterator struct {
	Event *SLAEnforcerContractAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SLAEnforcerContractAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SLAEnforcerContractAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SLAEnforcerContractAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SLAEnforcerContractAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SLAEnforcerContractAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SLAEnforcerContractAdded represents a ContractAdded event raised by the SLAEnforcer contract.
type SLAEnforcerContractAdded struct {
	ContractId common.Hash
	CustomerId string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterContractAdded is a free log retrieval operation binding the contract event 0x0fdbad36d02c7fa3bf7d92d489351b4122cefe532c727c18ff0cc22a45ad01ee.
//
// Solidity: event ContractAdded(string indexed contractId, string customerId)
func (_SLAEnforcer *SLAEnforcerFilterer) FilterContractAdded(opts *bind.FilterOpts, contractId []string) (*SLAEnforcerContractAddedIterator, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}

	logs, sub, err := _SLAEnforcer.contract.FilterLogs(opts, "ContractAdded", contractIdRule)
	if err != nil {
		return nil, err
	}
	return &SLAEnforcerContractAddedIterator{contract: _SLAEnforcer.contract, event: "ContractAdded", logs: logs, sub: sub}, nil
}

// WatchContractAdded is a free log subscription operation binding the contract event 0x0fdbad36d02c7fa3bf7d92d489351b4122cefe532c727c18ff0cc22a45ad01ee.
//
// Solidity: event ContractAdded(string indexed contractId, string customerId)
func (_SLAEnforcer *SLAEnforcerFilterer) WatchContractAdded(opts *bind.WatchOpts, sink chan<- *SLAEnforcerContractAdded, contractId []string) (event.Subscription, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}

	logs, sub, err := _SLAEnforcer.contract.WatchLogs(opts, "ContractAdded", contractIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SLAEnforcerContractAdded)
				if err := _SLAEnforcer.contract.UnpackLog(event, "ContractAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseContractAdded is a log parse operation binding the contract event 0x0fdbad36d02c7fa3bf7d92d489351b4122cefe532c727c18ff0cc22a45ad01ee.
//
// Solidity: event ContractAdded(string indexed contractId, string customerId)
func (_SLAEnforcer *SLAEnforcerFilterer) ParseContractAdded(log types.Log) (*SLAEnforcerContractAdded, error) {
	event := new(SLAEnforcerContractAdded)
	if err := _SLAEnforcer.contract.UnpackLog(event, "ContractAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SLAEnforcerSLAAddedIterator is returned from FilterSLAAdded and is used to iterate over the raw logs and unpacked data for SLAAdded events raised by the SLAEnforcer contract.
type SLAEnforcerSLAAddedIterator struct {
	Event *SLAEnforcerSLAAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SLAEnforcerSLAAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SLAEnforcerSLAAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SLAEnforcerSLAAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SLAEnforcerSLAAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SLAEnforcerSLAAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SLAEnforcerSLAAdded represents a SLAAdded event raised by the SLAEnforcer contract.
type SLAEnforcerSLAAdded struct {
	ContractId common.Hash
	SlaId      string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSLAAdded is a free log retrieval operation binding the contract event 0x95ca7048e73f7804cc6ca11e710e9a73c5dcf38f15666d206f81a00c5819f683.
//
// Solidity: event SLAAdded(string indexed contractId, string slaId)
func (_SLAEnforcer *SLAEnforcerFilterer) FilterSLAAdded(opts *bind.FilterOpts, contractId []string) (*SLAEnforcerSLAAddedIterator, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}

	logs, sub, err := _SLAEnforcer.contract.FilterLogs(opts, "SLAAdded", contractIdRule)
	if err != nil {
		return nil, err
	}
	return &SLAEnforcerSLAAddedIterator{contract: _SLAEnforcer.contract, event: "SLAAdded", logs: logs, sub: sub}, nil
}

// WatchSLAAdded is a free log subscription operation binding the contract event 0x95ca7048e73f7804cc6ca11e710e9a73c5dcf38f15666d206f81a00c5819f683.
//
// Solidity: event SLAAdded(string indexed contractId, string slaId)
func (_SLAEnforcer *SLAEnforcerFilterer) WatchSLAAdded(opts *bind.WatchOpts, sink chan<- *SLAEnforcerSLAAdded, contractId []string) (event.Subscription, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}

	logs, sub, err := _SLAEnforcer.contract.WatchLogs(opts, "SLAAdded", contractIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SLAEnforcerSLAAdded)
				if err := _SLAEnforcer.contract.UnpackLog(event, "SLAAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSLAAdded is a log parse operation binding the contract event 0x95ca7048e73f7804cc6ca11e710e9a73c5dcf38f15666d206f81a00c5819f683.
//
// Solidity: event SLAAdded(string indexed contractId, string slaId)
func (_SLAEnforcer *SLAEnforcerFilterer) ParseSLAAdded(log types.Log) (*SLAEnforcerSLAAdded, error) {
	event := new(SLAEnforcerSLAAdded)
	if err := _SLAEnforcer.contract.UnpackLog(event, "SLAAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SLAEnforcerSLAStatusUpdatedIterator is returned from FilterSLAStatusUpdated and is used to iterate over the raw logs and unpacked data for SLAStatusUpdated events raised by the SLAEnforcer contract.
type SLAEnforcerSLAStatusUpdatedIterator struct {
	Event *SLAEnforcerSLAStatusUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SLAEnforcerSLAStatusUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SLAEnforcerSLAStatusUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SLAEnforcerSLAStatusUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SLAEnforcerSLAStatusUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SLAEnforcerSLAStatusUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SLAEnforcerSLAStatusUpdated represents a SLAStatusUpdated event raised by the SLAEnforcer contract.
type SLAEnforcerSLAStatusUpdated struct {
	ContractId common.Hash
	SlaId      common.Hash
	NewStatus  uint8
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSLAStatusUpdated is a free log retrieval operation binding the contract event 0x55c708525bffba84acc408ffd134ca7094662d70ad28cc7dbf8dda4b009bcf62.
//
// Solidity: event SLAStatusUpdated(string indexed contractId, string indexed slaId, uint8 newStatus)
func (_SLAEnforcer *SLAEnforcerFilterer) FilterSLAStatusUpdated(opts *bind.FilterOpts, contractId []string, slaId []string) (*SLAEnforcerSLAStatusUpdatedIterator, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}
	var slaIdRule []interface{}
	for _, slaIdItem := range slaId {
		slaIdRule = append(slaIdRule, slaIdItem)
	}

	logs, sub, err := _SLAEnforcer.contract.FilterLogs(opts, "SLAStatusUpdated", contractIdRule, slaIdRule)
	if err != nil {
		return nil, err
	}
	return &SLAEnforcerSLAStatusUpdatedIterator{contract: _SLAEnforcer.contract, event: "SLAStatusUpdated", logs: logs, sub: sub}, nil
}

// WatchSLAStatusUpdated is a free log subscription operation binding the contract event 0x55c708525bffba84acc408ffd134ca7094662d70ad28cc7dbf8dda4b009bcf62.
//
// Solidity: event SLAStatusUpdated(string indexed contractId, string indexed slaId, uint8 newStatus)
func (_SLAEnforcer *SLAEnforcerFilterer) WatchSLAStatusUpdated(opts *bind.WatchOpts, sink chan<- *SLAEnforcerSLAStatusUpdated, contractId []string, slaId []string) (event.Subscription, error) {

	var contractIdRule []interface{}
	for _, contractIdItem := range contractId {
		contractIdRule = append(contractIdRule, contractIdItem)
	}
	var slaIdRule []interface{}
	for _, slaIdItem := range slaId {
		slaIdRule = append(slaIdRule, slaIdItem)
	}

	logs, sub, err := _SLAEnforcer.contract.WatchLogs(opts, "SLAStatusUpdated", contractIdRule, slaIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SLAEnforcerSLAStatusUpdated)
				if err := _SLAEnforcer.contract.UnpackLog(event, "SLAStatusUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSLAStatusUpdated is a log parse operation binding the contract event 0x55c708525bffba84acc408ffd134ca7094662d70ad28cc7dbf8dda4b009bcf62.
//
// Solidity: event SLAStatusUpdated(string indexed contractId, string indexed slaId, uint8 newStatus)
func (_SLAEnforcer *SLAEnforcerFilterer) ParseSLAStatusUpdated(log types.Log) (*SLAEnforcerSLAStatusUpdated, error) {
	event := new(SLAEnforcerSLAStatusUpdated)
	if err := _SLAEnforcer.contract.UnpackLog(event, "SLAStatusUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
