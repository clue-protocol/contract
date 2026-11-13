// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package code

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

// ClueCodeCodeMetadata is an auto generated low-level Go binding around an user-defined struct.
type ClueCodeCodeMetadata struct {
	Name        string
	Description string
	Image       string
	ExternalUrl string
	Traits      []ClueCodeTrait
}

// ClueCodeTrait is an auto generated low-level Go binding around an user-defined struct.
type ClueCodeTrait struct {
	TraitType string
	Value     string
}

// ClueCodeMetaData contains all meta data concerning the ClueCode contract.
var ClueCodeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"adminDelay\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"minter\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"initialURI\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint48\",\"name\":\"schedule\",\"type\":\"uint48\"}],\"name\":\"AccessControlEnforcedDefaultAdminDelay\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AccessControlEnforcedDefaultAdminRules\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"defaultAdmin\",\"type\":\"address\"}],\"name\":\"AccessControlInvalidDefaultAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"code\",\"type\":\"uint256\"}],\"name\":\"CodeMetadataIsFrozen\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721IncorrectOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721InsufficientApproval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOperator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC721InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721NonexistentToken\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"code\",\"type\":\"uint256\"}],\"name\":\"InvalidProviderCode\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LengthMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MetadataIsFrozen\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"bits\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"SafeCastOverflowedUintDowncast\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"newBaseURI\",\"type\":\"string\"}],\"name\":\"BaseURIUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"code\",\"type\":\"uint256\"}],\"name\":\"CodeMetadataFrozen\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"DefaultAdminDelayChangeCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"newDelay\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"effectSchedule\",\"type\":\"uint48\"}],\"name\":\"DefaultAdminDelayChangeScheduled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"DefaultAdminTransferCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"acceptSchedule\",\"type\":\"uint48\"}],\"name\":\"DefaultAdminTransferScheduled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"MetadataFrozen\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"code\",\"type\":\"uint256\"}],\"name\":\"MetadataSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"MetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_SUPPLY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"METADATA_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINTER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptDefaultAdminTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"baseURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"beginDefaultAdminTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cancelDefaultAdminTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"code\",\"type\":\"uint256\"}],\"name\":\"canonical\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint48\",\"name\":\"newDelay\",\"type\":\"uint48\"}],\"name\":\"changeDefaultAdminDelay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"code\",\"type\":\"uint256\"}],\"name\":\"clearMetadata\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"codeFrozen\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"defaultAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"defaultAdminDelay\",\"outputs\":[{\"internalType\":\"uint48\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"defaultAdminDelayIncreaseWait\",\"outputs\":[{\"internalType\":\"uint48\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"code\",\"type\":\"uint256\"}],\"name\":\"freezeCode\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"freezeMetadata\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"code\",\"type\":\"uint256\"}],\"name\":\"gradeOf\",\"outputs\":[{\"internalType\":\"enumClueCode.Grade\",\"name\":\"g\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"code\",\"type\":\"uint256\"}],\"name\":\"hasMetadata\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"code\",\"type\":\"uint256\"}],\"name\":\"isValidCode\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"code\",\"type\":\"uint256\"}],\"name\":\"metadata\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"metadataFrozen\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"code\",\"type\":\"uint256\"}],\"name\":\"metadataOf\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"image\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"externalUrl\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"traitType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"internalType\":\"structClueCode.Trait[]\",\"name\":\"traits\",\"type\":\"tuple[]\"}],\"internalType\":\"structClueCode.CodeMetadata\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"code\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"codes\",\"type\":\"uint256[]\"}],\"name\":\"mintBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"code\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"image\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"externalUrl\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"traitType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"internalType\":\"structClueCode.Trait[]\",\"name\":\"traits\",\"type\":\"tuple[]\"}],\"internalType\":\"structClueCode.CodeMetadata\",\"name\":\"md\",\"type\":\"tuple\"}],\"name\":\"mintWithMetadata\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingDefaultAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"schedule\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingDefaultAdminDelay\",\"outputs\":[{\"internalType\":\"uint48\",\"name\":\"newDelay\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"schedule\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollbackDefaultAdminDelay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newBaseURI\",\"type\":\"string\"}],\"name\":\"setBaseURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"code\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"image\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"externalUrl\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"traitType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"internalType\":\"structClueCode.Trait[]\",\"name\":\"traits\",\"type\":\"tuple[]\"}],\"internalType\":\"structClueCode.CodeMetadata\",\"name\":\"md\",\"type\":\"tuple\"}],\"name\":\"setMetadata\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"codes\",\"type\":\"uint256[]\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"image\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"externalUrl\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"traitType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"}],\"internalType\":\"structClueCode.Trait[]\",\"name\":\"traits\",\"type\":\"tuple[]\"}],\"internalType\":\"structClueCode.CodeMetadata[]\",\"name\":\"mds\",\"type\":\"tuple[]\"}],\"name\":\"setMetadataBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"id\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"code\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalMinted\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620045b0380380620045b08339810160408190526200003491620002d2565b828460405180604001604052806012815260200171436c75652050726f766964657220436f646560701b8152506040518060400160405280600381526020016243504360e81b81525081600090816200008e91906200047c565b5060016200009d82826200047c565b5050506001600160a01b038116620000cf57604051636116401160e11b81526000600482015260240160405180910390fd5b600780546001600160d01b0316600160d01b65ffffffffffff851602179055620000fb60008262000177565b505050620001307f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6836200017760201b60201c565b506200015d7f6bd6b5318a46e5fff572d5e4258a20774aab40cc35ac7680654b9081fcc82f808562000177565b50600a6200016c82826200047c565b505050505062000548565b600082620001d8576000620001946008546001600160a01b031690565b6001600160a01b031614620001bc57604051631fe1e13d60e11b815260040160405180910390fd5b600880546001600160a01b0319166001600160a01b0384161790555b620001e48383620001ed565b90505b92915050565b60008281526006602090815260408083206001600160a01b038516845290915281205460ff16620002965760008381526006602090815260408083206001600160a01b03861684529091529020805460ff191660011790556200024d3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a4506001620001e7565b506000620001e7565b80516001600160a01b0381168114620002b757600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b60008060008060808587031215620002e957600080fd5b620002f4856200029f565b935060208086015165ffffffffffff811681146200031157600080fd5b935062000321604087016200029f565b60608701519093506001600160401b03808211156200033f57600080fd5b818801915088601f8301126200035457600080fd5b815181811115620003695762000369620002bc565b604051601f8201601f19908116603f01168101908382118183101715620003945762000394620002bc565b816040528281528b86848701011115620003ad57600080fd5b600093505b82841015620003d15784840186015181850187015292850192620003b2565b600086848301015280965050505050505092959194509250565b600181811c908216806200040057607f821691505b6020821081036200042157634e487b7160e01b600052602260045260246000fd5b50919050565b601f82111562000477576000816000526020600020601f850160051c81016020861015620004525750805b601f850160051c820191505b8181101562000473578281556001016200045e565b5050505b505050565b81516001600160401b03811115620004985762000498620002bc565b620004b081620004a98454620003eb565b8462000427565b602080601f831160018114620004e85760008415620004cf5750858301515b600019600386901b1c1916600185901b17855562000473565b600085815260208120601f198616915b828110156200051957888601518255948401946001909101908401620004f8565b5085821015620005385787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61405880620005586000396000f3fe608060405234801561001057600080fd5b506004361061030c5760003560e01c806375ceb3411161019d578063c7c02bf5116100e9578063d5391393116100a2578063dca99a281161007c578063dca99a28146106ee578063e3684e3914610701578063e985e9c514610714578063fb3cc6c21461072757600080fd5b8063d5391393146106be578063d547741f146106d3578063d602b9fd146106e657600080fd5b8063c7c02bf514610642578063c87b56dd14610665578063cc8463c814610678578063cefc142914610680578063cf6eefb714610688578063d111515d146106b657600080fd5b806395d89b4111610156578063a22cb46511610130578063a22cb46514610600578063a2309ff814610613578063b88d4fde1461061c578063c15b34171461062f57600080fd5b806395d89b41146105c9578063a1eda53c146105d1578063a217fddf146105f857600080fd5b806375ceb3411461055457806384ef8ffc146105675780638a49dc59146105785780638da5cb5b1461059b578063918ba32a146105a357806391d14854146105b657600080fd5b80632fd718ac1161025c57806355f804b311610215578063649a5ec7116101ef578063649a5ec7146105135780636c0360eb146105265780636de17ce41461052e57806370a082311461054157600080fd5b806355f804b3146104da578063634e93da146104ed5780636352211e1461050057600080fd5b80632fd718ac1461046f57806332cb6b0c1461048257806336568abe1461048c578063388417821461049f57806340c10f19146104b457806342842e0e146104c757600080fd5b80630ef7cc8e116102c95780631cfaf219116102a35780631cfaf2191461040557806323b872dd14610418578063248a9ca31461042b5780632f2ff15d1461045c57600080fd5b80630ef7cc8e146103b25780631a3eb71e146103d25780631ac9a7d2146103e557600080fd5b806301ffc9a714610311578063022d63fb1461033957806306fdde0314610355578063081812fc1461036a578063095ea7b3146103955780630aa6220b146103aa575b600080fd5b61032461031f366004613120565b610734565b60405190151581526020015b60405180910390f35b620697805b60405165ffffffffffff9091168152602001610330565b61035d61075f565b604051610330919061318d565b61037d6103783660046131a0565b6107f1565b6040516001600160a01b039091168152602001610330565b6103a86103a33660046131d0565b61081a565b005b6103a8610829565b6103c56103c03660046131a0565b61083f565b60405161033091906131fa565b6103246103e03660046131a0565b610c64565b6103f86103f33660046131a0565b610c89565b604051610330919061330d565b6103a86104133660046131a0565b610cd5565b6103a8610426366004613335565b610dbe565b61044e6104393660046131a0565b60009081526006602052604090206001015490565b604051908152602001610330565b6103a861046a366004613371565b610e49565b6103a861047d3660046131a0565b610e71565b61044e6201637781565b6103a861049a366004613371565b610ed7565b61044e600080516020613f7e83398151915281565b6103a86104c23660046131d0565b610f7e565b6103a86104d5366004613335565b610fa5565b6103a86104e836600461339d565b610fc0565b6103a86104fb36600461340e565b611048565b61037d61050e3660046131a0565b61105c565b6103a8610521366004613429565b611067565b61035d61107b565b6103a861053c366004613469565b611109565b61044e61054f36600461340e565b61112b565b6103a86105623660046134fa565b611173565b6008546001600160a01b031661037d565b6103246105863660046131a0565b600e6020526000908152604090205460ff1681565b61037d6111c7565b6103a86105b136600461354c565b6111e0565b6103246105c4366004613371565b61120c565b61035d611237565b6105d9611246565b6040805165ffffffffffff938416815292909116602083015201610330565b61044e600081565b6103a861060e3660046135a2565b61129a565b61044e60095481565b6103a861062a3660046135f4565b6112a5565b61035d61063d3660046131a0565b6112bc565b6103246106503660046131a0565b6000908152600d602052604090205460ff1690565b61035d6106733660046131a0565b61142f565b61033e6114bf565b6103a861151e565b61069061155e565b604080516001600160a01b03909316835265ffffffffffff909116602083015201610330565b6103a861157f565b61044e600080516020613fde83398151915281565b6103a86106e1366004613371565b6115d0565b6103a86115f8565b6103a86106fc3660046136cf565b61160b565b61035d61070f3660046131a0565b6116a6565b61032461072236600461373a565b6116bb565b600b546103249060ff1681565b60006001600160e01b03198216632483248360e11b14806107595750610759826116e9565b92915050565b60606000805461076e90613764565b80601f016020809104026020016040519081016040528092919081815260200182805461079a90613764565b80156107e75780601f106107bc576101008083540402835291602001916107e7565b820191906000526020600020905b8154815290600101906020018083116107ca57829003601f168201915b5050505050905090565b60006107fc8261170e565b506000828152600460205260409020546001600160a01b0316610759565b610825828233611747565b5050565b600061083481611754565b61083c61175e565b50565b6108716040518060a0016040528060608152602001606081526020016060815260200160608152602001606081525090565b61087a82610c89565b506000828152600c602052604090819020815160a081019092528054829082906108a390613764565b80601f01602080910402602001604051908101604052809291908181526020018280546108cf90613764565b801561091c5780601f106108f15761010080835404028352916020019161091c565b820191906000526020600020905b8154815290600101906020018083116108ff57829003601f168201915b5050505050815260200160018201805461093590613764565b80601f016020809104026020016040519081016040528092919081815260200182805461096190613764565b80156109ae5780601f10610983576101008083540402835291602001916109ae565b820191906000526020600020905b81548152906001019060200180831161099157829003601f168201915b505050505081526020016002820180546109c790613764565b80601f01602080910402602001604051908101604052809291908181526020018280546109f390613764565b8015610a405780601f10610a1557610100808354040283529160200191610a40565b820191906000526020600020905b815481529060010190602001808311610a2357829003601f168201915b50505050508152602001600382018054610a5990613764565b80601f0160208091040260200160405190810160405280929190818152602001828054610a8590613764565b8015610ad25780601f10610aa757610100808354040283529160200191610ad2565b820191906000526020600020905b815481529060010190602001808311610ab557829003601f168201915b5050505050815260200160048201805480602002602001604051908101604052809291908181526020016000905b82821015610c565783829060005260206000209060020201604051806040016040529081600082018054610b3390613764565b80601f0160208091040260200160405190810160405280929190818152602001828054610b5f90613764565b8015610bac5780601f10610b8157610100808354040283529160200191610bac565b820191906000526020600020905b815481529060010190602001808311610b8f57829003601f168201915b50505050508152602001600182018054610bc590613764565b80601f0160208091040260200160405190810160405280929190818152602001828054610bf190613764565b8015610c3e5780601f10610c1357610100808354040283529160200191610c3e565b820191906000526020600020905b815481529060010190602001808311610c2157829003601f168201915b50505050508152505081526020019060010190610b00565b505050915250909392505050565b600080610c708361176b565b6003811115610c8157610c816132f7565b141592915050565b6000610c948261176b565b90506000816003811115610caa57610caa6132f7565b03610cd0576040516353abab0f60e01b8152600481018390526024015b60405180910390fd5b919050565b600080516020613f7e833981519152610ced81611754565b610cf6826117d6565b6000828152600c6020526040812090610d0f828261306f565b610d1d60018301600061306f565b610d2b60028301600061306f565b610d3960038301600061306f565b610d476004830160006130a9565b50506000828152600d6020526040808220805460ff191690555183917f043282acb0221e6f17ff7e7e56f271a3c5550143829948a198bc6025fa1e7b8891a26040518281527ff8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7906020015b60405180910390a15050565b6001600160a01b038216610de857604051633250574960e11b815260006004820152602401610cc7565b6000610df5838333611837565b9050836001600160a01b0316816001600160a01b031614610e43576040516364283d7b60e01b81526001600160a01b0380861660048301526024820184905282166044820152606401610cc7565b50505050565b81610e6757604051631fe1e13d60e11b815260040160405180910390fd5b6108258282611930565b600080516020613f7e833981519152610e8981611754565b610e9282610c89565b506000828152600e6020526040808220805460ff191660011790555183917ffd3328f5b1f143592b09179d213175ff381c5530dde9fadb207ecf55aa984ace91a25050565b81158015610ef257506008546001600160a01b038281169116145b15610f7457600080610f0261155e565b90925090506001600160a01b038216151580610f24575065ffffffffffff8116155b80610f3757504265ffffffffffff821610155b15610f5f576040516319ca5ebb60e01b815265ffffffffffff82166004820152602401610cc7565b50506007805465ffffffffffff60a01b191690555b6108258282611955565b600080516020613fde833981519152610f9681611754565b610fa08383611988565b505050565b610fa0838383604051806020016040528060008152506112a5565b600080516020613f7e833981519152610fd881611754565b600b5460ff1615610ffc5760405163b087bbf360e01b815260040160405180910390fd5b600a6110098385836137e0565b507f6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad838360405161103b92919061389f565b60405180910390a1505050565b600061105381611754565b610825826119a9565b60006107598261170e565b600061107281611754565b61082582611a1c565b600a805461108890613764565b80601f01602080910402602001604051908101604052809291908181526020018280546110b490613764565b80156111015780601f106110d657610100808354040283529160200191611101565b820191906000526020600020905b8154815290600101906020018083116110e457829003601f168201915b505050505081565b600080516020613f7e83398151915261112181611754565b610fa08383611a85565b60006001600160a01b038216611157576040516322718ad960e21b815260006004820152602401610cc7565b506001600160a01b031660009081526003602052604090205490565b600080516020613fde83398151915261118b81611754565b60005b828110156111c0576111b8858585848181106111ac576111ac6138ce565b90506020020135611988565b60010161118e565b5050505050565b60006111db6008546001600160a01b031690565b905090565b600080516020613fde8339815191526111f881611754565b6112028484611988565b610e438383611a85565b60009182526006602090815260408084206001600160a01b0393909316845291905290205460ff1690565b60606001805461076e90613764565b600854600090600160d01b900465ffffffffffff16801515801561127257504265ffffffffffff821610155b61127e57600080611292565b600854600160a01b900465ffffffffffff16815b915091509091565b610825338383611ce6565b6112b0848484610dbe565b610e4384848484611d85565b606060006112c983610c89565b9050600060018260038111156112e1576112e16132f7565b1461130d5760028260038111156112fa576112fa6132f7565b14611306576005611310565b6003611310565b60025b60ff169050600061132085611ea7565b90506000826001600160401b0381111561133c5761133c6135de565b6040519080825280601f01601f191660200182016040528015611366576020820181803683370190505b509050600082518461137891906138fa565b905060005b818110156113bb57600360fc1b83828151811061139c5761139c6138ce565b60200101906001600160f81b031916908160001a90535060010161137d565b5060005b8351811015611423578381815181106113da576113da6138ce565b01602001516001600160f81b031916836113f4838561390d565b81518110611404576114046138ce565b60200101906001600160f81b031916908160001a9053506001016113bf565b50909695505050505050565b606061143a8261170e565b506000828152600d602052604090205460ff161580156114675750600a805461146290613764565b151590505b1561149e57600a611477836112bc565b604051602001611488929190613920565b6040516020818303038152906040529050919050565b6114af6114aa83611f39565b61228c565b60405160200161148891906139a7565b600854600090600160d01b900465ffffffffffff1680151580156114ea57504265ffffffffffff8216105b61150557600754600160d01b900465ffffffffffff16611518565b600854600160a01b900465ffffffffffff165b91505090565b600061152861155e565b509050336001600160a01b0382161461155657604051636116401160e11b8152336004820152602401610cc7565b61083c6123eb565b6007546001600160a01b03811691600160a01b90910465ffffffffffff1690565b600080516020613f7e83398151915261159781611754565b600b805460ff191660011790556040517feef043febddf4e1d1cf1f72ff1407b84e036e805aa0934418cb82095da8d716490600090a150565b816115ee57604051631fe1e13d60e11b815260040160405180910390fd5b6108258282612484565b600061160381611754565b61083c6124a9565b600080516020613f7e83398151915261162381611754565b838214611646576040516001621398b960e31b0319815260040160405180910390fd5b60005b8481101561169e57611696868683818110611666576116666138ce565b9050602002013585858481811061167f5761167f6138ce565b905060200281019061169191906139ec565b611a85565b600101611649565b505050505050565b60606116b182610c89565b5061075982611f39565b6001600160a01b03918216600090815260056020908152604080832093909416825291909152205460ff1690565b60006001600160e01b031982166318a4c3c360e11b14806107595750610759826124b4565b6000818152600260205260408120546001600160a01b03168061075957604051637e27328960e01b815260048101849052602401610cc7565b610fa083838360016124d9565b61083c81336125df565b611769600080612618565b565b60006001821015801561177f575060638211155b1561178c57506001919050565b6064821015801561179f57506103e78211155b156117ac57506002919050565b61271082101580156117c157506201869f8211155b156117ce57506003919050565b506000919050565b6117df81610c89565b50600b5460ff16156118045760405163b087bbf360e01b815260040160405180910390fd5b6000818152600e602052604090205460ff161561083c57604051630c823b0160e31b815260048101829052602401610cc7565b6000828152600260205260408120546001600160a01b0390811690831615611864576118648184866126d8565b6001600160a01b038116156118a2576118816000856000806124d9565b6001600160a01b038116600090815260036020526040902080546000190190555b6001600160a01b038516156118d1576001600160a01b0385166000908152600360205260409020805460010190555b60008481526002602052604080822080546001600160a01b0319166001600160a01b0389811691821790925591518793918516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4949350505050565b60008281526006602052604090206001015461194b81611754565b610e43838361273c565b6001600160a01b038116331461197e5760405163334bd91960e11b815260040160405180910390fd5b610fa082826127ab565b61199181610c89565b5061199c82826127e8565b5050600980546001019055565b60006119b36114bf565b6119bc4261284d565b6119c69190613a0c565b90506119d28282612884565b60405165ffffffffffff821681526001600160a01b038316907f3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed69060200160405180910390a25050565b6000611a2782612903565b611a304261284d565b611a3a9190613a0c565b9050611a468282612618565b6040805165ffffffffffff8085168252831660208201527ff1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b9101610db2565b611a8e826117d6565b6000828152600c60205260409020611aa68280613a32565b8291611ab39190836137e0565b50611ac16020830183613a32565b6001830191611ad19190836137e0565b50611adf6040830183613a32565b6002830191611aef9190836137e0565b50611afd6060830183613a32565b6003830191611b0d9190836137e0565b50611b1c6004820160006130a9565b60005b611b2c6080840184613a78565b9050811015611c7557816004016040518060400160405280858060800190611b549190613a78565b85818110611b6457611b646138ce565b9050602002810190611b769190613ac1565b611b809080613a32565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250505090825250602001611bc76080870187613a78565b85818110611bd757611bd76138ce565b9050602002810190611be99190613ac1565b611bf7906020810190613a32565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920182905250939094525050835460018101855593815260209020825192936002020191829150611c529082613ad7565b5060208201516001820190611c679082613ad7565b505050806001019050611b1f565b506000838152600d6020526040808220805460ff191660011790555184917f043282acb0221e6f17ff7e7e56f271a3c5550143829948a198bc6025fa1e7b8891a26040518381527ff8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce79060200161103b565b6001600160a01b038216611d1857604051630b61174360e31b81526001600160a01b0383166004820152602401610cc7565b6001600160a01b03838116600081815260056020908152604080832094871680845294825291829020805460ff191686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b6001600160a01b0383163b15610e4357604051630a85bd0160e11b81526001600160a01b0384169063150b7a0290611dc7903390889087908790600401613b96565b6020604051808303816000875af1925050508015611e02575060408051601f3d908101601f19168201909252611dff91810190613bd3565b60015b611e6b573d808015611e30576040519150601f19603f3d011682016040523d82523d6000602084013e611e35565b606091505b508051600003611e6357604051633250574960e11b81526001600160a01b0385166004820152602401610cc7565b805181602001fd5b6001600160e01b03198116630a85bd0160e11b146111c057604051633250574960e11b81526001600160a01b0385166004820152602401610cc7565b60606000611eb48361294b565b60010190506000816001600160401b03811115611ed357611ed36135de565b6040519080825280601f01601f191660200182016040528015611efd576020820181803683370190505b5090508181016020015b600019016f181899199a1a9b1b9c1cb0b131b232b360811b600a86061a8153600a8504945084611f0757509392505050565b60606000611f4683610c89565b90506000611f53846112bc565b6000858152600c602052604081208054929350918290611f7290613764565b9050600003611fa05782604051602001611f8c9190613bf0565b60405160208183030381529060405261202c565b81548290611fad90613764565b80601f0160208091040260200160405190810160405280929190818152602001828054611fd990613764565b80156120265780601f10611ffb57610100808354040283529160200191612026565b820191906000526020600020905b81548152906001019060200180831161200957829003601f168201915b50505050505b9050600082600101805461203f90613764565b905060000361206657604051806060016040528060258152602001613ffe602591396120f4565b82600101805461207590613764565b80601f01602080910402602001604051908101604052809291908181526020018280546120a190613764565b80156120ee5780601f106120c3576101008083540402835291602001916120ee565b820191906000526020600020905b8154815290600101906020018083116120d157829003601f168201915b50505050505b9050600061210183612a23565b61210a83612a23565b60405160200161211b929190613c2b565b60408051601f1981840301815260208301909152600082526002860180549193509061214690613764565b159050612204576121e285600201805461215f90613764565b80601f016020809104026020016040519081016040528092919081815260200182805461218b90613764565b80156121d85780601f106121ad576101008083540402835291602001916121d8565b820191906000526020600020905b8154815290600101906020018083116121bb57829003601f168201915b5050505050612a23565b6040516020016121f29190613c9b565b60405160208183030381529060405290505b84600301805461221390613764565b159050612250578061222d86600301805461215f90613764565b60405160200161223e929190613cd9565b60405160208183030381529060405290505b818161225d888a89612d0a565b60405160200161226f93929190613d33565b604051602081830303815290604052975050505050505050919050565b606081516000036122ab57505060408051602081019091526000815290565b6000604051806060016040528060408152602001613f9e60409139905060006003845160026122da919061390d565b6122e49190613d9e565b6122ef906004613dc0565b6001600160401b03811115612306576123066135de565b6040519080825280601f01601f191660200182016040528015612330576020820181803683370190505b50905060018201602082018586518701602081018051600082525b828410156123a6576003840193508351603f8160121c168701518653600186019550603f81600c1c168701518653600186019550603f8160061c168701518653600186019550603f811687015186535060018501945061234b565b90525050855160039006600181146123c557600281146123d8576123e0565b603d6001830353603d60028303536123e0565b603d60018303535b509195945050505050565b6000806123f661155e565b9150915061240b8165ffffffffffff16151590565b158061241f57504265ffffffffffff821610155b15612447576040516319ca5ebb60e01b815265ffffffffffff82166004820152602401610cc7565b612463600061245e6008546001600160a01b031690565b6127ab565b5061246f60008361273c565b5050600780546001600160d01b031916905550565b60008281526006602052604090206001015461249f81611754565b610e4383836127ab565b611769600080612884565b60006001600160e01b03198216637965db0b60e01b1480610759575061075982612e13565b80806124ed57506001600160a01b03821615155b156125af5760006124fd8461170e565b90506001600160a01b038316158015906125295750826001600160a01b0316816001600160a01b031614155b801561253c575061253a81846116bb565b155b156125655760405163a9fbf51f60e01b81526001600160a01b0384166004820152602401610cc7565b81156125ad5783856001600160a01b0316826001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45b505b5050600090815260046020526040902080546001600160a01b0319166001600160a01b0392909216919091179055565b6125e9828261120c565b6108255760405163e2517d3f60e01b81526001600160a01b038216600482015260248101839052604401610cc7565b600854600160d01b900465ffffffffffff16801561269b574265ffffffffffff8216101561267157600854600780546001600160d01b0316600160a01b90920465ffffffffffff16600160d01b0291909117905561269b565b6040517f2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec590600090a15b50600880546001600160a01b0316600160a01b65ffffffffffff948516026001600160d01b031617600160d01b9290931691909102919091179055565b6126e3838383612e63565b610fa0576001600160a01b03831661271157604051637e27328960e01b815260048101829052602401610cc7565b60405163177e802f60e01b81526001600160a01b038316600482015260248101829052604401610cc7565b60008261279a5760006127576008546001600160a01b031690565b6001600160a01b03161461277e57604051631fe1e13d60e11b815260040160405180910390fd5b600880546001600160a01b0319166001600160a01b0384161790555b6127a48383612ec9565b9392505050565b6000821580156127c857506008546001600160a01b038381169116145b156127de57600880546001600160a01b03191690555b6127a48383612f5d565b6001600160a01b03821661281257604051633250574960e11b815260006004820152602401610cc7565b600061282083836000611837565b90506001600160a01b03811615610fa0576040516339e3563760e11b815260006004820152602401610cc7565b600065ffffffffffff821115612880576040516306dfcc6560e41b81526030600482015260248101839052604401610cc7565b5090565b600061288e61155e565b6007805465ffffffffffff8616600160a01b026001600160d01b03199091166001600160a01b0388161717905591506128d090508165ffffffffffff16151590565b15610fa0576040517f8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a960510990600090a1505050565b60008061290e6114bf565b90508065ffffffffffff168365ffffffffffff1611612936576129318382613dd7565b6127a4565b6127a465ffffffffffff841662069780612fca565b60008072184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b831061298a5772184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b830492506040015b6d04ee2d6d415b85acef810000000083106129b6576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc1000083106129d457662386f26fc10000830492506010015b6305f5e10083106129ec576305f5e100830492506008015b6127108310612a0057612710830492506004015b60648310612a12576064830492506002015b600a83106107595760010192915050565b80516060908290600090612a38906006613dc0565b6001600160401b03811115612a4f57612a4f6135de565b6040519080825280601f01601f191660200182016040528015612a79576020820181803683370190505b50905060006f181899199a1a9b1b9c1cb0b131b232b360811b815b8451811015612cff576000858281518110612ab157612ab16138ce565b01602001516001600160f81b0319169050601160f91b811480612ae15750601760fa1b6001600160f81b03198216145b15612b5957601760fa1b8585612af681613df6565b965081518110612b0857612b086138ce565b60200101906001600160f81b031916908160001a905350808585612b2b81613df6565b965081518110612b3d57612b3d6138ce565b60200101906001600160f81b031916908160001a905350612cf6565b602060f882901c1015612cc057601760fa1b8585612b7681613df6565b965081518110612b8857612b886138ce565b60200101906001600160f81b031916908160001a905350607560f81b8585612baf81613df6565b965081518110612bc157612bc16138ce565b60200101906001600160f81b031916908160001a905350600360fc1b8585612be881613df6565b965081518110612bfa57612bfa6138ce565b60200101906001600160f81b031916908160001a905350600360fc1b8585612c2181613df6565b965081518110612c3357612c336138ce565b60200101906001600160f81b031916908160001a9053508260fc82901c60108110612c6057612c606138ce565b1a60f81b8585612c6f81613df6565b965081518110612c8157612c816138ce565b60200101906001600160f81b031916908160001a90535082600f60f883901c1660108110612cb157612cb16138ce565b1a60f81b8585612b2b81613df6565b808585612ccc81613df6565b965081518110612cde57612cde6138ce565b60200101906001600160f81b031916908160001a9053505b50600101612a94565b505081529392505050565b6060612d1583612fe0565b612d1f8551611ea7565b604051602001612d30929190613e0f565b60408051601f19818403018152919052600483015490915060005b81811015612de85782612d88856004018381548110612d6c57612d6c6138ce565b9060005260206000209060020201600001805461215f90613764565b612dbc866004018481548110612da057612da06138ce565b9060005260206000209060020201600101805461215f90613764565b604051602001612dce93929190613ed1565b60408051601f198184030181529190529250600101612d4b565b5081604051602001612dfa9190613f58565b6040516020818303038152906040529150509392505050565b60006001600160e01b031982166380ac58cd60e01b1480612e4457506001600160e01b03198216635b5e139f60e01b145b8061075957506301ffc9a760e01b6001600160e01b0319831614610759565b60006001600160a01b03831615801590612ec15750826001600160a01b0316846001600160a01b03161480612e9d5750612e9d84846116bb565b80612ec157506000828152600460205260409020546001600160a01b038481169116145b949350505050565b6000612ed5838361120c565b612f555760008381526006602090815260408083206001600160a01b03861684529091529020805460ff19166001179055612f0d3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a4506001610759565b506000610759565b6000612f69838361120c565b15612f555760008381526006602090815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a4506001610759565b6000818310612fd957816127a4565b5090919050565b60606001826003811115612ff657612ff66132f7565b1461305257600282600381111561300f5761300f6132f7565b1461303357604051806040016040528060018152602001602160f91b815250610759565b604051806040016040528060018152602001604160f81b815250610759565b50506040805180820190915260018152605360f81b602082015290565b50805461307b90613764565b6000825580601f1061308b575050565b601f01602090049060005260206000209081019061083c91906130ca565b508054600082556002029060005260206000209081019061083c91906130df565b5b8082111561288057600081556001016130cb565b808211156128805760006130f3828261306f565b61310160018301600061306f565b506002016130df565b6001600160e01b03198116811461083c57600080fd5b60006020828403121561313257600080fd5b81356127a48161310a565b60005b83811015613158578181015183820152602001613140565b50506000910152565b6000815180845261317981602086016020860161313d565b601f01601f19169290920160200192915050565b6020815260006127a46020830184613161565b6000602082840312156131b257600080fd5b5035919050565b80356001600160a01b0381168114610cd057600080fd5b600080604083850312156131e357600080fd5b6131ec836131b9565b946020939093013593505050565b60006020808352835160a08285015261321660c0850182613161565b905081850151601f196040818785030160408801526132358484613161565b935060408801519250818785030160608801526132528484613161565b9350606088015192508187850301608088015261326f8484613161565b6080890151888203840160a08a015280518083529195508601935085850190600581901b8601870160005b828110156132e85785888303018452865180518684526132bc87850182613161565b918b0151848303858d01529190506132d48183613161565b988b0198958b01959350505060010161329a565b509a9950505050505050505050565b634e487b7160e01b600052602160045260246000fd5b602081016004831061332f57634e487b7160e01b600052602160045260246000fd5b91905290565b60008060006060848603121561334a57600080fd5b613353846131b9565b9250613361602085016131b9565b9150604084013590509250925092565b6000806040838503121561338457600080fd5b82359150613394602084016131b9565b90509250929050565b600080602083850312156133b057600080fd5b82356001600160401b03808211156133c757600080fd5b818501915085601f8301126133db57600080fd5b8135818111156133ea57600080fd5b8660208285010111156133fc57600080fd5b60209290920196919550909350505050565b60006020828403121561342057600080fd5b6127a4826131b9565b60006020828403121561343b57600080fd5b813565ffffffffffff811681146127a457600080fd5b600060a0828403121561346357600080fd5b50919050565b6000806040838503121561347c57600080fd5b8235915060208301356001600160401b0381111561349957600080fd5b6134a585828601613451565b9150509250929050565b60008083601f8401126134c157600080fd5b5081356001600160401b038111156134d857600080fd5b6020830191508360208260051b85010111156134f357600080fd5b9250929050565b60008060006040848603121561350f57600080fd5b613518846131b9565b925060208401356001600160401b0381111561353357600080fd5b61353f868287016134af565b9497909650939450505050565b60008060006060848603121561356157600080fd5b61356a846131b9565b92506020840135915060408401356001600160401b0381111561358c57600080fd5b61359886828701613451565b9150509250925092565b600080604083850312156135b557600080fd5b6135be836131b9565b9150602083013580151581146135d357600080fd5b809150509250929050565b634e487b7160e01b600052604160045260246000fd5b6000806000806080858703121561360a57600080fd5b613613856131b9565b9350613621602086016131b9565b92506040850135915060608501356001600160401b038082111561364457600080fd5b818701915087601f83011261365857600080fd5b81358181111561366a5761366a6135de565b604051601f8201601f19908116603f01168101908382118183101715613692576136926135de565b816040528281528a60208487010111156136ab57600080fd5b82602086016020830137600060208483010152809550505050505092959194509250565b600080600080604085870312156136e557600080fd5b84356001600160401b03808211156136fc57600080fd5b613708888389016134af565b9096509450602087013591508082111561372157600080fd5b5061372e878288016134af565b95989497509550505050565b6000806040838503121561374d57600080fd5b613756836131b9565b9150613394602084016131b9565b600181811c9082168061377857607f821691505b60208210810361346357634e487b7160e01b600052602260045260246000fd5b601f821115610fa0576000816000526020600020601f850160051c810160208610156137c15750805b601f850160051c820191505b8181101561169e578281556001016137cd565b6001600160401b038311156137f7576137f76135de565b61380b836138058354613764565b83613798565b6000601f84116001811461383f57600085156138275750838201355b600019600387901b1c1916600186901b1783556111c0565b600083815260209020601f19861690835b828110156138705786850135825560209485019460019092019101613850565b508682101561388d5760001960f88860031b161c19848701351681555b505060018560011b0183555050505050565b60208152816020820152818360408301376000818301604090810191909152601f909201601f19160101919050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b81810381811115610759576107596138e4565b80820180821115610759576107596138e4565b600080845461392e81613764565b60018281168015613946576001811461395b5761398a565b60ff198416875282151583028701945061398a565b8860005260208060002060005b858110156139815781548a820152908401908201613968565b50505082870194505b50505050835161399e81836020880161313d565b01949350505050565b7f646174613a6170706c69636174696f6e2f6a736f6e3b6261736536342c0000008152600082516139df81601d85016020870161313d565b91909101601d0192915050565b60008235609e19833603018112613a0257600080fd5b9190910192915050565b65ffffffffffff818116838216019080821115613a2b57613a2b6138e4565b5092915050565b6000808335601e19843603018112613a4957600080fd5b8301803591506001600160401b03821115613a6357600080fd5b6020019150368190038213156134f357600080fd5b6000808335601e19843603018112613a8f57600080fd5b8301803591506001600160401b03821115613aa957600080fd5b6020019150600581901b36038213156134f357600080fd5b60008235603e19833603018112613a0257600080fd5b81516001600160401b03811115613af057613af06135de565b613b0481613afe8454613764565b84613798565b602080601f831160018114613b395760008415613b215750858301515b600019600386901b1c1916600185901b17855561169e565b600085815260208120601f198616915b82811015613b6857888601518255948401946001909101908401613b49565b5085821015613b865787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6001600160a01b0385811682528416602082015260408101839052608060608201819052600090613bc990830184613161565b9695505050505050565b600060208284031215613be557600080fd5b81516127a48161310a565b72021b63ab290283937bb34b232b91021b7b2329606d1b815260008251613c1e81601385016020870161313d565b9190910160130192915050565b683d913730b6b2911d1160b91b81528251600090613c5081600985016020880161313d565b701116113232b9b1b934b83a34b7b7111d1160791b6009918401918201528351613c8181601a84016020880161313d565b601160f91b601a9290910191820152601b01949350505050565b69161134b6b0b3b2911d1160b11b81528151600090613cc181600a85016020870161313d565b601160f91b600a939091019283015250600b01919050565b60008351613ceb81846020880161313d565b70161132bc3a32b93730b62fbab936111d1160791b9083019081528351613d1981601184016020880161313d565b601160f91b60119290910191820152601201949350505050565b60008451613d4581846020890161313d565b845190830190613d5981836020890161313d565b6d161130ba3a3934b13aba32b9911d60911b91019081528351613d8381600e84016020880161313d565b607d60f81b600e9290910191820152600f0195945050505050565b600082613dbb57634e487b7160e01b600052601260045260246000fd5b500490565b8082028115828204841417610759576107596138e4565b65ffffffffffff828116828216039080821115613a2b57613a2b6138e4565b600060018201613e0857613e086138e4565b5060010190565b7f5b7b2274726169745f74797065223a224772616465222c2276616c7565223a22815260008351613e4781602085016020880161313d565b808301905062089f4b60ea1b60208201527f7b2274726169745f74797065223a22446967697473222c22646973706c61795f60238201527f74797065223a226e756d626572222c2276616c7565223a00000000000000000060438201528351613eb781605a84016020880161313d565b607d60f81b605a9290910191820152605b01949350505050565b60008451613ee381846020890161313d565b6f163d913a3930b4ba2fba3cb832911d1160811b9083019081528451613f1081601084016020890161313d565b6a1116113b30b63ab2911d1160a91b601092909101918201528351613f3c81601b84016020880161313d565b61227d60f01b601b9290910191820152601d0195945050505050565b60008251613f6a81846020870161313d565b605d60f81b92019182525060010191905056fe6bd6b5318a46e5fff572d5e4258a20774aab40cc35ac7680654b9081fcc82f804142434445464748494a4b4c4d4e4f505152535455565758595a6162636465666768696a6b6c6d6e6f707172737475767778797a303132333435363738392b2f9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6412050726f766964657220436f6465206f662074686520436c75652050726f746f636f6c2ea2646970667358221220843783c636403fffbec8fe54ad5da90e6fb50b3a4d77f20073396352fa62d39364736f6c63430008180033",
}

// ClueCodeABI is the input ABI used to generate the binding from.
// Deprecated: Use ClueCodeMetaData.ABI instead.
var ClueCodeABI = ClueCodeMetaData.ABI

// ClueCodeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ClueCodeMetaData.Bin instead.
var ClueCodeBin = ClueCodeMetaData.Bin

// DeployClueCode deploys a new Ethereum contract, binding an instance of ClueCode to it.
func DeployClueCode(auth *bind.TransactOpts, backend bind.ContractBackend, admin common.Address, adminDelay *big.Int, minter common.Address, initialURI string) (common.Address, *types.Transaction, *ClueCode, error) {
	parsed, err := ClueCodeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ClueCodeBin), backend, admin, adminDelay, minter, initialURI)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ClueCode{ClueCodeCaller: ClueCodeCaller{contract: contract}, ClueCodeTransactor: ClueCodeTransactor{contract: contract}, ClueCodeFilterer: ClueCodeFilterer{contract: contract}}, nil
}

// ClueCode is an auto generated Go binding around an Ethereum contract.
type ClueCode struct {
	ClueCodeCaller     // Read-only binding to the contract
	ClueCodeTransactor // Write-only binding to the contract
	ClueCodeFilterer   // Log filterer for contract events
}

// ClueCodeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ClueCodeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClueCodeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ClueCodeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClueCodeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ClueCodeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ClueCodeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ClueCodeSession struct {
	Contract     *ClueCode         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ClueCodeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ClueCodeCallerSession struct {
	Contract *ClueCodeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ClueCodeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ClueCodeTransactorSession struct {
	Contract     *ClueCodeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ClueCodeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ClueCodeRaw struct {
	Contract *ClueCode // Generic contract binding to access the raw methods on
}

// ClueCodeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ClueCodeCallerRaw struct {
	Contract *ClueCodeCaller // Generic read-only contract binding to access the raw methods on
}

// ClueCodeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ClueCodeTransactorRaw struct {
	Contract *ClueCodeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewClueCode creates a new instance of ClueCode, bound to a specific deployed contract.
func NewClueCode(address common.Address, backend bind.ContractBackend) (*ClueCode, error) {
	contract, err := bindClueCode(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ClueCode{ClueCodeCaller: ClueCodeCaller{contract: contract}, ClueCodeTransactor: ClueCodeTransactor{contract: contract}, ClueCodeFilterer: ClueCodeFilterer{contract: contract}}, nil
}

// NewClueCodeCaller creates a new read-only instance of ClueCode, bound to a specific deployed contract.
func NewClueCodeCaller(address common.Address, caller bind.ContractCaller) (*ClueCodeCaller, error) {
	contract, err := bindClueCode(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ClueCodeCaller{contract: contract}, nil
}

// NewClueCodeTransactor creates a new write-only instance of ClueCode, bound to a specific deployed contract.
func NewClueCodeTransactor(address common.Address, transactor bind.ContractTransactor) (*ClueCodeTransactor, error) {
	contract, err := bindClueCode(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ClueCodeTransactor{contract: contract}, nil
}

// NewClueCodeFilterer creates a new log filterer instance of ClueCode, bound to a specific deployed contract.
func NewClueCodeFilterer(address common.Address, filterer bind.ContractFilterer) (*ClueCodeFilterer, error) {
	contract, err := bindClueCode(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ClueCodeFilterer{contract: contract}, nil
}

// bindClueCode binds a generic wrapper to an already deployed contract.
func bindClueCode(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ClueCodeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ClueCode *ClueCodeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ClueCode.Contract.ClueCodeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ClueCode *ClueCodeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClueCode.Contract.ClueCodeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ClueCode *ClueCodeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ClueCode.Contract.ClueCodeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ClueCode *ClueCodeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ClueCode.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ClueCode *ClueCodeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClueCode.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ClueCode *ClueCodeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ClueCode.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ClueCode *ClueCodeCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ClueCode *ClueCodeSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ClueCode.Contract.DEFAULTADMINROLE(&_ClueCode.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_ClueCode *ClueCodeCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _ClueCode.Contract.DEFAULTADMINROLE(&_ClueCode.CallOpts)
}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_ClueCode *ClueCodeCaller) MAXSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "MAX_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_ClueCode *ClueCodeSession) MAXSUPPLY() (*big.Int, error) {
	return _ClueCode.Contract.MAXSUPPLY(&_ClueCode.CallOpts)
}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_ClueCode *ClueCodeCallerSession) MAXSUPPLY() (*big.Int, error) {
	return _ClueCode.Contract.MAXSUPPLY(&_ClueCode.CallOpts)
}

// METADATAROLE is a free data retrieval call binding the contract method 0x38841782.
//
// Solidity: function METADATA_ROLE() view returns(bytes32)
func (_ClueCode *ClueCodeCaller) METADATAROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "METADATA_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// METADATAROLE is a free data retrieval call binding the contract method 0x38841782.
//
// Solidity: function METADATA_ROLE() view returns(bytes32)
func (_ClueCode *ClueCodeSession) METADATAROLE() ([32]byte, error) {
	return _ClueCode.Contract.METADATAROLE(&_ClueCode.CallOpts)
}

// METADATAROLE is a free data retrieval call binding the contract method 0x38841782.
//
// Solidity: function METADATA_ROLE() view returns(bytes32)
func (_ClueCode *ClueCodeCallerSession) METADATAROLE() ([32]byte, error) {
	return _ClueCode.Contract.METADATAROLE(&_ClueCode.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_ClueCode *ClueCodeCaller) MINTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "MINTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_ClueCode *ClueCodeSession) MINTERROLE() ([32]byte, error) {
	return _ClueCode.Contract.MINTERROLE(&_ClueCode.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_ClueCode *ClueCodeCallerSession) MINTERROLE() ([32]byte, error) {
	return _ClueCode.Contract.MINTERROLE(&_ClueCode.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ClueCode *ClueCodeCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ClueCode *ClueCodeSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ClueCode.Contract.BalanceOf(&_ClueCode.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ClueCode *ClueCodeCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ClueCode.Contract.BalanceOf(&_ClueCode.CallOpts, owner)
}

// BaseURI is a free data retrieval call binding the contract method 0x6c0360eb.
//
// Solidity: function baseURI() view returns(string)
func (_ClueCode *ClueCodeCaller) BaseURI(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "baseURI")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// BaseURI is a free data retrieval call binding the contract method 0x6c0360eb.
//
// Solidity: function baseURI() view returns(string)
func (_ClueCode *ClueCodeSession) BaseURI() (string, error) {
	return _ClueCode.Contract.BaseURI(&_ClueCode.CallOpts)
}

// BaseURI is a free data retrieval call binding the contract method 0x6c0360eb.
//
// Solidity: function baseURI() view returns(string)
func (_ClueCode *ClueCodeCallerSession) BaseURI() (string, error) {
	return _ClueCode.Contract.BaseURI(&_ClueCode.CallOpts)
}

// Canonical is a free data retrieval call binding the contract method 0xc15b3417.
//
// Solidity: function canonical(uint256 code) pure returns(string)
func (_ClueCode *ClueCodeCaller) Canonical(opts *bind.CallOpts, code *big.Int) (string, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "canonical", code)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Canonical is a free data retrieval call binding the contract method 0xc15b3417.
//
// Solidity: function canonical(uint256 code) pure returns(string)
func (_ClueCode *ClueCodeSession) Canonical(code *big.Int) (string, error) {
	return _ClueCode.Contract.Canonical(&_ClueCode.CallOpts, code)
}

// Canonical is a free data retrieval call binding the contract method 0xc15b3417.
//
// Solidity: function canonical(uint256 code) pure returns(string)
func (_ClueCode *ClueCodeCallerSession) Canonical(code *big.Int) (string, error) {
	return _ClueCode.Contract.Canonical(&_ClueCode.CallOpts, code)
}

// CodeFrozen is a free data retrieval call binding the contract method 0x8a49dc59.
//
// Solidity: function codeFrozen(uint256 ) view returns(bool)
func (_ClueCode *ClueCodeCaller) CodeFrozen(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "codeFrozen", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CodeFrozen is a free data retrieval call binding the contract method 0x8a49dc59.
//
// Solidity: function codeFrozen(uint256 ) view returns(bool)
func (_ClueCode *ClueCodeSession) CodeFrozen(arg0 *big.Int) (bool, error) {
	return _ClueCode.Contract.CodeFrozen(&_ClueCode.CallOpts, arg0)
}

// CodeFrozen is a free data retrieval call binding the contract method 0x8a49dc59.
//
// Solidity: function codeFrozen(uint256 ) view returns(bool)
func (_ClueCode *ClueCodeCallerSession) CodeFrozen(arg0 *big.Int) (bool, error) {
	return _ClueCode.Contract.CodeFrozen(&_ClueCode.CallOpts, arg0)
}

// DefaultAdmin is a free data retrieval call binding the contract method 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (_ClueCode *ClueCodeCaller) DefaultAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "defaultAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DefaultAdmin is a free data retrieval call binding the contract method 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (_ClueCode *ClueCodeSession) DefaultAdmin() (common.Address, error) {
	return _ClueCode.Contract.DefaultAdmin(&_ClueCode.CallOpts)
}

// DefaultAdmin is a free data retrieval call binding the contract method 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (_ClueCode *ClueCodeCallerSession) DefaultAdmin() (common.Address, error) {
	return _ClueCode.Contract.DefaultAdmin(&_ClueCode.CallOpts)
}

// DefaultAdminDelay is a free data retrieval call binding the contract method 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (_ClueCode *ClueCodeCaller) DefaultAdminDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "defaultAdminDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DefaultAdminDelay is a free data retrieval call binding the contract method 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (_ClueCode *ClueCodeSession) DefaultAdminDelay() (*big.Int, error) {
	return _ClueCode.Contract.DefaultAdminDelay(&_ClueCode.CallOpts)
}

// DefaultAdminDelay is a free data retrieval call binding the contract method 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (_ClueCode *ClueCodeCallerSession) DefaultAdminDelay() (*big.Int, error) {
	return _ClueCode.Contract.DefaultAdminDelay(&_ClueCode.CallOpts)
}

// DefaultAdminDelayIncreaseWait is a free data retrieval call binding the contract method 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (_ClueCode *ClueCodeCaller) DefaultAdminDelayIncreaseWait(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "defaultAdminDelayIncreaseWait")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DefaultAdminDelayIncreaseWait is a free data retrieval call binding the contract method 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (_ClueCode *ClueCodeSession) DefaultAdminDelayIncreaseWait() (*big.Int, error) {
	return _ClueCode.Contract.DefaultAdminDelayIncreaseWait(&_ClueCode.CallOpts)
}

// DefaultAdminDelayIncreaseWait is a free data retrieval call binding the contract method 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (_ClueCode *ClueCodeCallerSession) DefaultAdminDelayIncreaseWait() (*big.Int, error) {
	return _ClueCode.Contract.DefaultAdminDelayIncreaseWait(&_ClueCode.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ClueCode *ClueCodeCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ClueCode *ClueCodeSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ClueCode.Contract.GetApproved(&_ClueCode.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ClueCode *ClueCodeCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ClueCode.Contract.GetApproved(&_ClueCode.CallOpts, tokenId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ClueCode *ClueCodeCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ClueCode *ClueCodeSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ClueCode.Contract.GetRoleAdmin(&_ClueCode.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_ClueCode *ClueCodeCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _ClueCode.Contract.GetRoleAdmin(&_ClueCode.CallOpts, role)
}

// GradeOf is a free data retrieval call binding the contract method 0x1ac9a7d2.
//
// Solidity: function gradeOf(uint256 code) pure returns(uint8 g)
func (_ClueCode *ClueCodeCaller) GradeOf(opts *bind.CallOpts, code *big.Int) (uint8, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "gradeOf", code)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GradeOf is a free data retrieval call binding the contract method 0x1ac9a7d2.
//
// Solidity: function gradeOf(uint256 code) pure returns(uint8 g)
func (_ClueCode *ClueCodeSession) GradeOf(code *big.Int) (uint8, error) {
	return _ClueCode.Contract.GradeOf(&_ClueCode.CallOpts, code)
}

// GradeOf is a free data retrieval call binding the contract method 0x1ac9a7d2.
//
// Solidity: function gradeOf(uint256 code) pure returns(uint8 g)
func (_ClueCode *ClueCodeCallerSession) GradeOf(code *big.Int) (uint8, error) {
	return _ClueCode.Contract.GradeOf(&_ClueCode.CallOpts, code)
}

// HasMetadata is a free data retrieval call binding the contract method 0xc7c02bf5.
//
// Solidity: function hasMetadata(uint256 code) view returns(bool)
func (_ClueCode *ClueCodeCaller) HasMetadata(opts *bind.CallOpts, code *big.Int) (bool, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "hasMetadata", code)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasMetadata is a free data retrieval call binding the contract method 0xc7c02bf5.
//
// Solidity: function hasMetadata(uint256 code) view returns(bool)
func (_ClueCode *ClueCodeSession) HasMetadata(code *big.Int) (bool, error) {
	return _ClueCode.Contract.HasMetadata(&_ClueCode.CallOpts, code)
}

// HasMetadata is a free data retrieval call binding the contract method 0xc7c02bf5.
//
// Solidity: function hasMetadata(uint256 code) view returns(bool)
func (_ClueCode *ClueCodeCallerSession) HasMetadata(code *big.Int) (bool, error) {
	return _ClueCode.Contract.HasMetadata(&_ClueCode.CallOpts, code)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ClueCode *ClueCodeCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ClueCode *ClueCodeSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ClueCode.Contract.HasRole(&_ClueCode.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_ClueCode *ClueCodeCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _ClueCode.Contract.HasRole(&_ClueCode.CallOpts, role, account)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ClueCode *ClueCodeCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ClueCode *ClueCodeSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ClueCode.Contract.IsApprovedForAll(&_ClueCode.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ClueCode *ClueCodeCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ClueCode.Contract.IsApprovedForAll(&_ClueCode.CallOpts, owner, operator)
}

// IsValidCode is a free data retrieval call binding the contract method 0x1a3eb71e.
//
// Solidity: function isValidCode(uint256 code) pure returns(bool)
func (_ClueCode *ClueCodeCaller) IsValidCode(opts *bind.CallOpts, code *big.Int) (bool, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "isValidCode", code)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidCode is a free data retrieval call binding the contract method 0x1a3eb71e.
//
// Solidity: function isValidCode(uint256 code) pure returns(bool)
func (_ClueCode *ClueCodeSession) IsValidCode(code *big.Int) (bool, error) {
	return _ClueCode.Contract.IsValidCode(&_ClueCode.CallOpts, code)
}

// IsValidCode is a free data retrieval call binding the contract method 0x1a3eb71e.
//
// Solidity: function isValidCode(uint256 code) pure returns(bool)
func (_ClueCode *ClueCodeCallerSession) IsValidCode(code *big.Int) (bool, error) {
	return _ClueCode.Contract.IsValidCode(&_ClueCode.CallOpts, code)
}

// Metadata is a free data retrieval call binding the contract method 0xe3684e39.
//
// Solidity: function metadata(uint256 code) view returns(string)
func (_ClueCode *ClueCodeCaller) Metadata(opts *bind.CallOpts, code *big.Int) (string, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "metadata", code)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Metadata is a free data retrieval call binding the contract method 0xe3684e39.
//
// Solidity: function metadata(uint256 code) view returns(string)
func (_ClueCode *ClueCodeSession) Metadata(code *big.Int) (string, error) {
	return _ClueCode.Contract.Metadata(&_ClueCode.CallOpts, code)
}

// Metadata is a free data retrieval call binding the contract method 0xe3684e39.
//
// Solidity: function metadata(uint256 code) view returns(string)
func (_ClueCode *ClueCodeCallerSession) Metadata(code *big.Int) (string, error) {
	return _ClueCode.Contract.Metadata(&_ClueCode.CallOpts, code)
}

// MetadataFrozen is a free data retrieval call binding the contract method 0xfb3cc6c2.
//
// Solidity: function metadataFrozen() view returns(bool)
func (_ClueCode *ClueCodeCaller) MetadataFrozen(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "metadataFrozen")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// MetadataFrozen is a free data retrieval call binding the contract method 0xfb3cc6c2.
//
// Solidity: function metadataFrozen() view returns(bool)
func (_ClueCode *ClueCodeSession) MetadataFrozen() (bool, error) {
	return _ClueCode.Contract.MetadataFrozen(&_ClueCode.CallOpts)
}

// MetadataFrozen is a free data retrieval call binding the contract method 0xfb3cc6c2.
//
// Solidity: function metadataFrozen() view returns(bool)
func (_ClueCode *ClueCodeCallerSession) MetadataFrozen() (bool, error) {
	return _ClueCode.Contract.MetadataFrozen(&_ClueCode.CallOpts)
}

// MetadataOf is a free data retrieval call binding the contract method 0x0ef7cc8e.
//
// Solidity: function metadataOf(uint256 code) view returns((string,string,string,string,(string,string)[]))
func (_ClueCode *ClueCodeCaller) MetadataOf(opts *bind.CallOpts, code *big.Int) (ClueCodeCodeMetadata, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "metadataOf", code)

	if err != nil {
		return *new(ClueCodeCodeMetadata), err
	}

	out0 := *abi.ConvertType(out[0], new(ClueCodeCodeMetadata)).(*ClueCodeCodeMetadata)

	return out0, err

}

// MetadataOf is a free data retrieval call binding the contract method 0x0ef7cc8e.
//
// Solidity: function metadataOf(uint256 code) view returns((string,string,string,string,(string,string)[]))
func (_ClueCode *ClueCodeSession) MetadataOf(code *big.Int) (ClueCodeCodeMetadata, error) {
	return _ClueCode.Contract.MetadataOf(&_ClueCode.CallOpts, code)
}

// MetadataOf is a free data retrieval call binding the contract method 0x0ef7cc8e.
//
// Solidity: function metadataOf(uint256 code) view returns((string,string,string,string,(string,string)[]))
func (_ClueCode *ClueCodeCallerSession) MetadataOf(code *big.Int) (ClueCodeCodeMetadata, error) {
	return _ClueCode.Contract.MetadataOf(&_ClueCode.CallOpts, code)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ClueCode *ClueCodeCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ClueCode *ClueCodeSession) Name() (string, error) {
	return _ClueCode.Contract.Name(&_ClueCode.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ClueCode *ClueCodeCallerSession) Name() (string, error) {
	return _ClueCode.Contract.Name(&_ClueCode.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ClueCode *ClueCodeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ClueCode *ClueCodeSession) Owner() (common.Address, error) {
	return _ClueCode.Contract.Owner(&_ClueCode.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ClueCode *ClueCodeCallerSession) Owner() (common.Address, error) {
	return _ClueCode.Contract.Owner(&_ClueCode.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ClueCode *ClueCodeCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ClueCode *ClueCodeSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ClueCode.Contract.OwnerOf(&_ClueCode.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ClueCode *ClueCodeCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ClueCode.Contract.OwnerOf(&_ClueCode.CallOpts, tokenId)
}

// PendingDefaultAdmin is a free data retrieval call binding the contract method 0xcf6eefb7.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (_ClueCode *ClueCodeCaller) PendingDefaultAdmin(opts *bind.CallOpts) (struct {
	NewAdmin common.Address
	Schedule *big.Int
}, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "pendingDefaultAdmin")

	outstruct := new(struct {
		NewAdmin common.Address
		Schedule *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NewAdmin = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Schedule = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// PendingDefaultAdmin is a free data retrieval call binding the contract method 0xcf6eefb7.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (_ClueCode *ClueCodeSession) PendingDefaultAdmin() (struct {
	NewAdmin common.Address
	Schedule *big.Int
}, error) {
	return _ClueCode.Contract.PendingDefaultAdmin(&_ClueCode.CallOpts)
}

// PendingDefaultAdmin is a free data retrieval call binding the contract method 0xcf6eefb7.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (_ClueCode *ClueCodeCallerSession) PendingDefaultAdmin() (struct {
	NewAdmin common.Address
	Schedule *big.Int
}, error) {
	return _ClueCode.Contract.PendingDefaultAdmin(&_ClueCode.CallOpts)
}

// PendingDefaultAdminDelay is a free data retrieval call binding the contract method 0xa1eda53c.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (_ClueCode *ClueCodeCaller) PendingDefaultAdminDelay(opts *bind.CallOpts) (struct {
	NewDelay *big.Int
	Schedule *big.Int
}, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "pendingDefaultAdminDelay")

	outstruct := new(struct {
		NewDelay *big.Int
		Schedule *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NewDelay = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Schedule = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// PendingDefaultAdminDelay is a free data retrieval call binding the contract method 0xa1eda53c.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (_ClueCode *ClueCodeSession) PendingDefaultAdminDelay() (struct {
	NewDelay *big.Int
	Schedule *big.Int
}, error) {
	return _ClueCode.Contract.PendingDefaultAdminDelay(&_ClueCode.CallOpts)
}

// PendingDefaultAdminDelay is a free data retrieval call binding the contract method 0xa1eda53c.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (_ClueCode *ClueCodeCallerSession) PendingDefaultAdminDelay() (struct {
	NewDelay *big.Int
	Schedule *big.Int
}, error) {
	return _ClueCode.Contract.PendingDefaultAdminDelay(&_ClueCode.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 id) view returns(bool)
func (_ClueCode *ClueCodeCaller) SupportsInterface(opts *bind.CallOpts, id [4]byte) (bool, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "supportsInterface", id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 id) view returns(bool)
func (_ClueCode *ClueCodeSession) SupportsInterface(id [4]byte) (bool, error) {
	return _ClueCode.Contract.SupportsInterface(&_ClueCode.CallOpts, id)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 id) view returns(bool)
func (_ClueCode *ClueCodeCallerSession) SupportsInterface(id [4]byte) (bool, error) {
	return _ClueCode.Contract.SupportsInterface(&_ClueCode.CallOpts, id)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ClueCode *ClueCodeCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ClueCode *ClueCodeSession) Symbol() (string, error) {
	return _ClueCode.Contract.Symbol(&_ClueCode.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ClueCode *ClueCodeCallerSession) Symbol() (string, error) {
	return _ClueCode.Contract.Symbol(&_ClueCode.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 code) view returns(string)
func (_ClueCode *ClueCodeCaller) TokenURI(opts *bind.CallOpts, code *big.Int) (string, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "tokenURI", code)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 code) view returns(string)
func (_ClueCode *ClueCodeSession) TokenURI(code *big.Int) (string, error) {
	return _ClueCode.Contract.TokenURI(&_ClueCode.CallOpts, code)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 code) view returns(string)
func (_ClueCode *ClueCodeCallerSession) TokenURI(code *big.Int) (string, error) {
	return _ClueCode.Contract.TokenURI(&_ClueCode.CallOpts, code)
}

// TotalMinted is a free data retrieval call binding the contract method 0xa2309ff8.
//
// Solidity: function totalMinted() view returns(uint256)
func (_ClueCode *ClueCodeCaller) TotalMinted(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ClueCode.contract.Call(opts, &out, "totalMinted")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalMinted is a free data retrieval call binding the contract method 0xa2309ff8.
//
// Solidity: function totalMinted() view returns(uint256)
func (_ClueCode *ClueCodeSession) TotalMinted() (*big.Int, error) {
	return _ClueCode.Contract.TotalMinted(&_ClueCode.CallOpts)
}

// TotalMinted is a free data retrieval call binding the contract method 0xa2309ff8.
//
// Solidity: function totalMinted() view returns(uint256)
func (_ClueCode *ClueCodeCallerSession) TotalMinted() (*big.Int, error) {
	return _ClueCode.Contract.TotalMinted(&_ClueCode.CallOpts)
}

// AcceptDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xcefc1429.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (_ClueCode *ClueCodeTransactor) AcceptDefaultAdminTransfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "acceptDefaultAdminTransfer")
}

// AcceptDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xcefc1429.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (_ClueCode *ClueCodeSession) AcceptDefaultAdminTransfer() (*types.Transaction, error) {
	return _ClueCode.Contract.AcceptDefaultAdminTransfer(&_ClueCode.TransactOpts)
}

// AcceptDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xcefc1429.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (_ClueCode *ClueCodeTransactorSession) AcceptDefaultAdminTransfer() (*types.Transaction, error) {
	return _ClueCode.Contract.AcceptDefaultAdminTransfer(&_ClueCode.TransactOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ClueCode *ClueCodeTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ClueCode *ClueCodeSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ClueCode.Contract.Approve(&_ClueCode.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ClueCode *ClueCodeTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ClueCode.Contract.Approve(&_ClueCode.TransactOpts, to, tokenId)
}

// BeginDefaultAdminTransfer is a paid mutator transaction binding the contract method 0x634e93da.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (_ClueCode *ClueCodeTransactor) BeginDefaultAdminTransfer(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "beginDefaultAdminTransfer", newAdmin)
}

// BeginDefaultAdminTransfer is a paid mutator transaction binding the contract method 0x634e93da.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (_ClueCode *ClueCodeSession) BeginDefaultAdminTransfer(newAdmin common.Address) (*types.Transaction, error) {
	return _ClueCode.Contract.BeginDefaultAdminTransfer(&_ClueCode.TransactOpts, newAdmin)
}

// BeginDefaultAdminTransfer is a paid mutator transaction binding the contract method 0x634e93da.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (_ClueCode *ClueCodeTransactorSession) BeginDefaultAdminTransfer(newAdmin common.Address) (*types.Transaction, error) {
	return _ClueCode.Contract.BeginDefaultAdminTransfer(&_ClueCode.TransactOpts, newAdmin)
}

// CancelDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xd602b9fd.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (_ClueCode *ClueCodeTransactor) CancelDefaultAdminTransfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "cancelDefaultAdminTransfer")
}

// CancelDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xd602b9fd.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (_ClueCode *ClueCodeSession) CancelDefaultAdminTransfer() (*types.Transaction, error) {
	return _ClueCode.Contract.CancelDefaultAdminTransfer(&_ClueCode.TransactOpts)
}

// CancelDefaultAdminTransfer is a paid mutator transaction binding the contract method 0xd602b9fd.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (_ClueCode *ClueCodeTransactorSession) CancelDefaultAdminTransfer() (*types.Transaction, error) {
	return _ClueCode.Contract.CancelDefaultAdminTransfer(&_ClueCode.TransactOpts)
}

// ChangeDefaultAdminDelay is a paid mutator transaction binding the contract method 0x649a5ec7.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (_ClueCode *ClueCodeTransactor) ChangeDefaultAdminDelay(opts *bind.TransactOpts, newDelay *big.Int) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "changeDefaultAdminDelay", newDelay)
}

// ChangeDefaultAdminDelay is a paid mutator transaction binding the contract method 0x649a5ec7.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (_ClueCode *ClueCodeSession) ChangeDefaultAdminDelay(newDelay *big.Int) (*types.Transaction, error) {
	return _ClueCode.Contract.ChangeDefaultAdminDelay(&_ClueCode.TransactOpts, newDelay)
}

// ChangeDefaultAdminDelay is a paid mutator transaction binding the contract method 0x649a5ec7.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (_ClueCode *ClueCodeTransactorSession) ChangeDefaultAdminDelay(newDelay *big.Int) (*types.Transaction, error) {
	return _ClueCode.Contract.ChangeDefaultAdminDelay(&_ClueCode.TransactOpts, newDelay)
}

// ClearMetadata is a paid mutator transaction binding the contract method 0x1cfaf219.
//
// Solidity: function clearMetadata(uint256 code) returns()
func (_ClueCode *ClueCodeTransactor) ClearMetadata(opts *bind.TransactOpts, code *big.Int) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "clearMetadata", code)
}

// ClearMetadata is a paid mutator transaction binding the contract method 0x1cfaf219.
//
// Solidity: function clearMetadata(uint256 code) returns()
func (_ClueCode *ClueCodeSession) ClearMetadata(code *big.Int) (*types.Transaction, error) {
	return _ClueCode.Contract.ClearMetadata(&_ClueCode.TransactOpts, code)
}

// ClearMetadata is a paid mutator transaction binding the contract method 0x1cfaf219.
//
// Solidity: function clearMetadata(uint256 code) returns()
func (_ClueCode *ClueCodeTransactorSession) ClearMetadata(code *big.Int) (*types.Transaction, error) {
	return _ClueCode.Contract.ClearMetadata(&_ClueCode.TransactOpts, code)
}

// FreezeCode is a paid mutator transaction binding the contract method 0x2fd718ac.
//
// Solidity: function freezeCode(uint256 code) returns()
func (_ClueCode *ClueCodeTransactor) FreezeCode(opts *bind.TransactOpts, code *big.Int) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "freezeCode", code)
}

// FreezeCode is a paid mutator transaction binding the contract method 0x2fd718ac.
//
// Solidity: function freezeCode(uint256 code) returns()
func (_ClueCode *ClueCodeSession) FreezeCode(code *big.Int) (*types.Transaction, error) {
	return _ClueCode.Contract.FreezeCode(&_ClueCode.TransactOpts, code)
}

// FreezeCode is a paid mutator transaction binding the contract method 0x2fd718ac.
//
// Solidity: function freezeCode(uint256 code) returns()
func (_ClueCode *ClueCodeTransactorSession) FreezeCode(code *big.Int) (*types.Transaction, error) {
	return _ClueCode.Contract.FreezeCode(&_ClueCode.TransactOpts, code)
}

// FreezeMetadata is a paid mutator transaction binding the contract method 0xd111515d.
//
// Solidity: function freezeMetadata() returns()
func (_ClueCode *ClueCodeTransactor) FreezeMetadata(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "freezeMetadata")
}

// FreezeMetadata is a paid mutator transaction binding the contract method 0xd111515d.
//
// Solidity: function freezeMetadata() returns()
func (_ClueCode *ClueCodeSession) FreezeMetadata() (*types.Transaction, error) {
	return _ClueCode.Contract.FreezeMetadata(&_ClueCode.TransactOpts)
}

// FreezeMetadata is a paid mutator transaction binding the contract method 0xd111515d.
//
// Solidity: function freezeMetadata() returns()
func (_ClueCode *ClueCodeTransactorSession) FreezeMetadata() (*types.Transaction, error) {
	return _ClueCode.Contract.FreezeMetadata(&_ClueCode.TransactOpts)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ClueCode *ClueCodeTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ClueCode *ClueCodeSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ClueCode.Contract.GrantRole(&_ClueCode.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_ClueCode *ClueCodeTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ClueCode.Contract.GrantRole(&_ClueCode.TransactOpts, role, account)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 code) returns()
func (_ClueCode *ClueCodeTransactor) Mint(opts *bind.TransactOpts, to common.Address, code *big.Int) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "mint", to, code)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 code) returns()
func (_ClueCode *ClueCodeSession) Mint(to common.Address, code *big.Int) (*types.Transaction, error) {
	return _ClueCode.Contract.Mint(&_ClueCode.TransactOpts, to, code)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 code) returns()
func (_ClueCode *ClueCodeTransactorSession) Mint(to common.Address, code *big.Int) (*types.Transaction, error) {
	return _ClueCode.Contract.Mint(&_ClueCode.TransactOpts, to, code)
}

// MintBatch is a paid mutator transaction binding the contract method 0x75ceb341.
//
// Solidity: function mintBatch(address to, uint256[] codes) returns()
func (_ClueCode *ClueCodeTransactor) MintBatch(opts *bind.TransactOpts, to common.Address, codes []*big.Int) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "mintBatch", to, codes)
}

// MintBatch is a paid mutator transaction binding the contract method 0x75ceb341.
//
// Solidity: function mintBatch(address to, uint256[] codes) returns()
func (_ClueCode *ClueCodeSession) MintBatch(to common.Address, codes []*big.Int) (*types.Transaction, error) {
	return _ClueCode.Contract.MintBatch(&_ClueCode.TransactOpts, to, codes)
}

// MintBatch is a paid mutator transaction binding the contract method 0x75ceb341.
//
// Solidity: function mintBatch(address to, uint256[] codes) returns()
func (_ClueCode *ClueCodeTransactorSession) MintBatch(to common.Address, codes []*big.Int) (*types.Transaction, error) {
	return _ClueCode.Contract.MintBatch(&_ClueCode.TransactOpts, to, codes)
}

// MintWithMetadata is a paid mutator transaction binding the contract method 0x918ba32a.
//
// Solidity: function mintWithMetadata(address to, uint256 code, (string,string,string,string,(string,string)[]) md) returns()
func (_ClueCode *ClueCodeTransactor) MintWithMetadata(opts *bind.TransactOpts, to common.Address, code *big.Int, md ClueCodeCodeMetadata) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "mintWithMetadata", to, code, md)
}

// MintWithMetadata is a paid mutator transaction binding the contract method 0x918ba32a.
//
// Solidity: function mintWithMetadata(address to, uint256 code, (string,string,string,string,(string,string)[]) md) returns()
func (_ClueCode *ClueCodeSession) MintWithMetadata(to common.Address, code *big.Int, md ClueCodeCodeMetadata) (*types.Transaction, error) {
	return _ClueCode.Contract.MintWithMetadata(&_ClueCode.TransactOpts, to, code, md)
}

// MintWithMetadata is a paid mutator transaction binding the contract method 0x918ba32a.
//
// Solidity: function mintWithMetadata(address to, uint256 code, (string,string,string,string,(string,string)[]) md) returns()
func (_ClueCode *ClueCodeTransactorSession) MintWithMetadata(to common.Address, code *big.Int, md ClueCodeCodeMetadata) (*types.Transaction, error) {
	return _ClueCode.Contract.MintWithMetadata(&_ClueCode.TransactOpts, to, code, md)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_ClueCode *ClueCodeTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_ClueCode *ClueCodeSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ClueCode.Contract.RenounceRole(&_ClueCode.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_ClueCode *ClueCodeTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ClueCode.Contract.RenounceRole(&_ClueCode.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ClueCode *ClueCodeTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ClueCode *ClueCodeSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ClueCode.Contract.RevokeRole(&_ClueCode.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_ClueCode *ClueCodeTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _ClueCode.Contract.RevokeRole(&_ClueCode.TransactOpts, role, account)
}

// RollbackDefaultAdminDelay is a paid mutator transaction binding the contract method 0x0aa6220b.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (_ClueCode *ClueCodeTransactor) RollbackDefaultAdminDelay(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "rollbackDefaultAdminDelay")
}

// RollbackDefaultAdminDelay is a paid mutator transaction binding the contract method 0x0aa6220b.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (_ClueCode *ClueCodeSession) RollbackDefaultAdminDelay() (*types.Transaction, error) {
	return _ClueCode.Contract.RollbackDefaultAdminDelay(&_ClueCode.TransactOpts)
}

// RollbackDefaultAdminDelay is a paid mutator transaction binding the contract method 0x0aa6220b.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (_ClueCode *ClueCodeTransactorSession) RollbackDefaultAdminDelay() (*types.Transaction, error) {
	return _ClueCode.Contract.RollbackDefaultAdminDelay(&_ClueCode.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ClueCode *ClueCodeTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ClueCode *ClueCodeSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ClueCode.Contract.SafeTransferFrom(&_ClueCode.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ClueCode *ClueCodeTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ClueCode.Contract.SafeTransferFrom(&_ClueCode.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_ClueCode *ClueCodeTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_ClueCode *ClueCodeSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _ClueCode.Contract.SafeTransferFrom0(&_ClueCode.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_ClueCode *ClueCodeTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _ClueCode.Contract.SafeTransferFrom0(&_ClueCode.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ClueCode *ClueCodeTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ClueCode *ClueCodeSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _ClueCode.Contract.SetApprovalForAll(&_ClueCode.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ClueCode *ClueCodeTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _ClueCode.Contract.SetApprovalForAll(&_ClueCode.TransactOpts, operator, approved)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string newBaseURI) returns()
func (_ClueCode *ClueCodeTransactor) SetBaseURI(opts *bind.TransactOpts, newBaseURI string) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "setBaseURI", newBaseURI)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string newBaseURI) returns()
func (_ClueCode *ClueCodeSession) SetBaseURI(newBaseURI string) (*types.Transaction, error) {
	return _ClueCode.Contract.SetBaseURI(&_ClueCode.TransactOpts, newBaseURI)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string newBaseURI) returns()
func (_ClueCode *ClueCodeTransactorSession) SetBaseURI(newBaseURI string) (*types.Transaction, error) {
	return _ClueCode.Contract.SetBaseURI(&_ClueCode.TransactOpts, newBaseURI)
}

// SetMetadata is a paid mutator transaction binding the contract method 0x6de17ce4.
//
// Solidity: function setMetadata(uint256 code, (string,string,string,string,(string,string)[]) md) returns()
func (_ClueCode *ClueCodeTransactor) SetMetadata(opts *bind.TransactOpts, code *big.Int, md ClueCodeCodeMetadata) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "setMetadata", code, md)
}

// SetMetadata is a paid mutator transaction binding the contract method 0x6de17ce4.
//
// Solidity: function setMetadata(uint256 code, (string,string,string,string,(string,string)[]) md) returns()
func (_ClueCode *ClueCodeSession) SetMetadata(code *big.Int, md ClueCodeCodeMetadata) (*types.Transaction, error) {
	return _ClueCode.Contract.SetMetadata(&_ClueCode.TransactOpts, code, md)
}

// SetMetadata is a paid mutator transaction binding the contract method 0x6de17ce4.
//
// Solidity: function setMetadata(uint256 code, (string,string,string,string,(string,string)[]) md) returns()
func (_ClueCode *ClueCodeTransactorSession) SetMetadata(code *big.Int, md ClueCodeCodeMetadata) (*types.Transaction, error) {
	return _ClueCode.Contract.SetMetadata(&_ClueCode.TransactOpts, code, md)
}

// SetMetadataBatch is a paid mutator transaction binding the contract method 0xdca99a28.
//
// Solidity: function setMetadataBatch(uint256[] codes, (string,string,string,string,(string,string)[])[] mds) returns()
func (_ClueCode *ClueCodeTransactor) SetMetadataBatch(opts *bind.TransactOpts, codes []*big.Int, mds []ClueCodeCodeMetadata) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "setMetadataBatch", codes, mds)
}

// SetMetadataBatch is a paid mutator transaction binding the contract method 0xdca99a28.
//
// Solidity: function setMetadataBatch(uint256[] codes, (string,string,string,string,(string,string)[])[] mds) returns()
func (_ClueCode *ClueCodeSession) SetMetadataBatch(codes []*big.Int, mds []ClueCodeCodeMetadata) (*types.Transaction, error) {
	return _ClueCode.Contract.SetMetadataBatch(&_ClueCode.TransactOpts, codes, mds)
}

// SetMetadataBatch is a paid mutator transaction binding the contract method 0xdca99a28.
//
// Solidity: function setMetadataBatch(uint256[] codes, (string,string,string,string,(string,string)[])[] mds) returns()
func (_ClueCode *ClueCodeTransactorSession) SetMetadataBatch(codes []*big.Int, mds []ClueCodeCodeMetadata) (*types.Transaction, error) {
	return _ClueCode.Contract.SetMetadataBatch(&_ClueCode.TransactOpts, codes, mds)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ClueCode *ClueCodeTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ClueCode.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ClueCode *ClueCodeSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ClueCode.Contract.TransferFrom(&_ClueCode.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ClueCode *ClueCodeTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ClueCode.Contract.TransferFrom(&_ClueCode.TransactOpts, from, to, tokenId)
}

// ClueCodeApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ClueCode contract.
type ClueCodeApprovalIterator struct {
	Event *ClueCodeApproval // Event containing the contract specifics and raw log

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
func (it *ClueCodeApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClueCodeApproval)
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
		it.Event = new(ClueCodeApproval)
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
func (it *ClueCodeApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClueCodeApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClueCodeApproval represents a Approval event raised by the ClueCode contract.
type ClueCodeApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_ClueCode *ClueCodeFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*ClueCodeApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ClueCode.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ClueCodeApprovalIterator{contract: _ClueCode.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_ClueCode *ClueCodeFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ClueCodeApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ClueCode.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClueCodeApproval)
				if err := _ClueCode.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_ClueCode *ClueCodeFilterer) ParseApproval(log types.Log) (*ClueCodeApproval, error) {
	event := new(ClueCodeApproval)
	if err := _ClueCode.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClueCodeApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the ClueCode contract.
type ClueCodeApprovalForAllIterator struct {
	Event *ClueCodeApprovalForAll // Event containing the contract specifics and raw log

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
func (it *ClueCodeApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClueCodeApprovalForAll)
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
		it.Event = new(ClueCodeApprovalForAll)
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
func (it *ClueCodeApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClueCodeApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClueCodeApprovalForAll represents a ApprovalForAll event raised by the ClueCode contract.
type ClueCodeApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ClueCode *ClueCodeFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*ClueCodeApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ClueCode.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &ClueCodeApprovalForAllIterator{contract: _ClueCode.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ClueCode *ClueCodeFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *ClueCodeApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ClueCode.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClueCodeApprovalForAll)
				if err := _ClueCode.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ClueCode *ClueCodeFilterer) ParseApprovalForAll(log types.Log) (*ClueCodeApprovalForAll, error) {
	event := new(ClueCodeApprovalForAll)
	if err := _ClueCode.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClueCodeBaseURIUpdatedIterator is returned from FilterBaseURIUpdated and is used to iterate over the raw logs and unpacked data for BaseURIUpdated events raised by the ClueCode contract.
type ClueCodeBaseURIUpdatedIterator struct {
	Event *ClueCodeBaseURIUpdated // Event containing the contract specifics and raw log

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
func (it *ClueCodeBaseURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClueCodeBaseURIUpdated)
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
		it.Event = new(ClueCodeBaseURIUpdated)
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
func (it *ClueCodeBaseURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClueCodeBaseURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClueCodeBaseURIUpdated represents a BaseURIUpdated event raised by the ClueCode contract.
type ClueCodeBaseURIUpdated struct {
	NewBaseURI string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBaseURIUpdated is a free log retrieval operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string newBaseURI)
func (_ClueCode *ClueCodeFilterer) FilterBaseURIUpdated(opts *bind.FilterOpts) (*ClueCodeBaseURIUpdatedIterator, error) {

	logs, sub, err := _ClueCode.contract.FilterLogs(opts, "BaseURIUpdated")
	if err != nil {
		return nil, err
	}
	return &ClueCodeBaseURIUpdatedIterator{contract: _ClueCode.contract, event: "BaseURIUpdated", logs: logs, sub: sub}, nil
}

// WatchBaseURIUpdated is a free log subscription operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string newBaseURI)
func (_ClueCode *ClueCodeFilterer) WatchBaseURIUpdated(opts *bind.WatchOpts, sink chan<- *ClueCodeBaseURIUpdated) (event.Subscription, error) {

	logs, sub, err := _ClueCode.contract.WatchLogs(opts, "BaseURIUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClueCodeBaseURIUpdated)
				if err := _ClueCode.contract.UnpackLog(event, "BaseURIUpdated", log); err != nil {
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

// ParseBaseURIUpdated is a log parse operation binding the contract event 0x6741b2fc379fad678116fe3d4d4b9a1a184ab53ba36b86ad0fa66340b1ab41ad.
//
// Solidity: event BaseURIUpdated(string newBaseURI)
func (_ClueCode *ClueCodeFilterer) ParseBaseURIUpdated(log types.Log) (*ClueCodeBaseURIUpdated, error) {
	event := new(ClueCodeBaseURIUpdated)
	if err := _ClueCode.contract.UnpackLog(event, "BaseURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClueCodeCodeMetadataFrozenIterator is returned from FilterCodeMetadataFrozen and is used to iterate over the raw logs and unpacked data for CodeMetadataFrozen events raised by the ClueCode contract.
type ClueCodeCodeMetadataFrozenIterator struct {
	Event *ClueCodeCodeMetadataFrozen // Event containing the contract specifics and raw log

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
func (it *ClueCodeCodeMetadataFrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClueCodeCodeMetadataFrozen)
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
		it.Event = new(ClueCodeCodeMetadataFrozen)
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
func (it *ClueCodeCodeMetadataFrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClueCodeCodeMetadataFrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClueCodeCodeMetadataFrozen represents a CodeMetadataFrozen event raised by the ClueCode contract.
type ClueCodeCodeMetadataFrozen struct {
	Code *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterCodeMetadataFrozen is a free log retrieval operation binding the contract event 0xfd3328f5b1f143592b09179d213175ff381c5530dde9fadb207ecf55aa984ace.
//
// Solidity: event CodeMetadataFrozen(uint256 indexed code)
func (_ClueCode *ClueCodeFilterer) FilterCodeMetadataFrozen(opts *bind.FilterOpts, code []*big.Int) (*ClueCodeCodeMetadataFrozenIterator, error) {

	var codeRule []interface{}
	for _, codeItem := range code {
		codeRule = append(codeRule, codeItem)
	}

	logs, sub, err := _ClueCode.contract.FilterLogs(opts, "CodeMetadataFrozen", codeRule)
	if err != nil {
		return nil, err
	}
	return &ClueCodeCodeMetadataFrozenIterator{contract: _ClueCode.contract, event: "CodeMetadataFrozen", logs: logs, sub: sub}, nil
}

// WatchCodeMetadataFrozen is a free log subscription operation binding the contract event 0xfd3328f5b1f143592b09179d213175ff381c5530dde9fadb207ecf55aa984ace.
//
// Solidity: event CodeMetadataFrozen(uint256 indexed code)
func (_ClueCode *ClueCodeFilterer) WatchCodeMetadataFrozen(opts *bind.WatchOpts, sink chan<- *ClueCodeCodeMetadataFrozen, code []*big.Int) (event.Subscription, error) {

	var codeRule []interface{}
	for _, codeItem := range code {
		codeRule = append(codeRule, codeItem)
	}

	logs, sub, err := _ClueCode.contract.WatchLogs(opts, "CodeMetadataFrozen", codeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClueCodeCodeMetadataFrozen)
				if err := _ClueCode.contract.UnpackLog(event, "CodeMetadataFrozen", log); err != nil {
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

// ParseCodeMetadataFrozen is a log parse operation binding the contract event 0xfd3328f5b1f143592b09179d213175ff381c5530dde9fadb207ecf55aa984ace.
//
// Solidity: event CodeMetadataFrozen(uint256 indexed code)
func (_ClueCode *ClueCodeFilterer) ParseCodeMetadataFrozen(log types.Log) (*ClueCodeCodeMetadataFrozen, error) {
	event := new(ClueCodeCodeMetadataFrozen)
	if err := _ClueCode.contract.UnpackLog(event, "CodeMetadataFrozen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClueCodeDefaultAdminDelayChangeCanceledIterator is returned from FilterDefaultAdminDelayChangeCanceled and is used to iterate over the raw logs and unpacked data for DefaultAdminDelayChangeCanceled events raised by the ClueCode contract.
type ClueCodeDefaultAdminDelayChangeCanceledIterator struct {
	Event *ClueCodeDefaultAdminDelayChangeCanceled // Event containing the contract specifics and raw log

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
func (it *ClueCodeDefaultAdminDelayChangeCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClueCodeDefaultAdminDelayChangeCanceled)
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
		it.Event = new(ClueCodeDefaultAdminDelayChangeCanceled)
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
func (it *ClueCodeDefaultAdminDelayChangeCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClueCodeDefaultAdminDelayChangeCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClueCodeDefaultAdminDelayChangeCanceled represents a DefaultAdminDelayChangeCanceled event raised by the ClueCode contract.
type ClueCodeDefaultAdminDelayChangeCanceled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminDelayChangeCanceled is a free log retrieval operation binding the contract event 0x2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (_ClueCode *ClueCodeFilterer) FilterDefaultAdminDelayChangeCanceled(opts *bind.FilterOpts) (*ClueCodeDefaultAdminDelayChangeCanceledIterator, error) {

	logs, sub, err := _ClueCode.contract.FilterLogs(opts, "DefaultAdminDelayChangeCanceled")
	if err != nil {
		return nil, err
	}
	return &ClueCodeDefaultAdminDelayChangeCanceledIterator{contract: _ClueCode.contract, event: "DefaultAdminDelayChangeCanceled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminDelayChangeCanceled is a free log subscription operation binding the contract event 0x2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (_ClueCode *ClueCodeFilterer) WatchDefaultAdminDelayChangeCanceled(opts *bind.WatchOpts, sink chan<- *ClueCodeDefaultAdminDelayChangeCanceled) (event.Subscription, error) {

	logs, sub, err := _ClueCode.contract.WatchLogs(opts, "DefaultAdminDelayChangeCanceled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClueCodeDefaultAdminDelayChangeCanceled)
				if err := _ClueCode.contract.UnpackLog(event, "DefaultAdminDelayChangeCanceled", log); err != nil {
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

// ParseDefaultAdminDelayChangeCanceled is a log parse operation binding the contract event 0x2b1fa2edafe6f7b9e97c1a9e0c3660e645beb2dcaa2d45bdbf9beaf5472e1ec5.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (_ClueCode *ClueCodeFilterer) ParseDefaultAdminDelayChangeCanceled(log types.Log) (*ClueCodeDefaultAdminDelayChangeCanceled, error) {
	event := new(ClueCodeDefaultAdminDelayChangeCanceled)
	if err := _ClueCode.contract.UnpackLog(event, "DefaultAdminDelayChangeCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClueCodeDefaultAdminDelayChangeScheduledIterator is returned from FilterDefaultAdminDelayChangeScheduled and is used to iterate over the raw logs and unpacked data for DefaultAdminDelayChangeScheduled events raised by the ClueCode contract.
type ClueCodeDefaultAdminDelayChangeScheduledIterator struct {
	Event *ClueCodeDefaultAdminDelayChangeScheduled // Event containing the contract specifics and raw log

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
func (it *ClueCodeDefaultAdminDelayChangeScheduledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClueCodeDefaultAdminDelayChangeScheduled)
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
		it.Event = new(ClueCodeDefaultAdminDelayChangeScheduled)
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
func (it *ClueCodeDefaultAdminDelayChangeScheduledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClueCodeDefaultAdminDelayChangeScheduledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClueCodeDefaultAdminDelayChangeScheduled represents a DefaultAdminDelayChangeScheduled event raised by the ClueCode contract.
type ClueCodeDefaultAdminDelayChangeScheduled struct {
	NewDelay       *big.Int
	EffectSchedule *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminDelayChangeScheduled is a free log retrieval operation binding the contract event 0xf1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (_ClueCode *ClueCodeFilterer) FilterDefaultAdminDelayChangeScheduled(opts *bind.FilterOpts) (*ClueCodeDefaultAdminDelayChangeScheduledIterator, error) {

	logs, sub, err := _ClueCode.contract.FilterLogs(opts, "DefaultAdminDelayChangeScheduled")
	if err != nil {
		return nil, err
	}
	return &ClueCodeDefaultAdminDelayChangeScheduledIterator{contract: _ClueCode.contract, event: "DefaultAdminDelayChangeScheduled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminDelayChangeScheduled is a free log subscription operation binding the contract event 0xf1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (_ClueCode *ClueCodeFilterer) WatchDefaultAdminDelayChangeScheduled(opts *bind.WatchOpts, sink chan<- *ClueCodeDefaultAdminDelayChangeScheduled) (event.Subscription, error) {

	logs, sub, err := _ClueCode.contract.WatchLogs(opts, "DefaultAdminDelayChangeScheduled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClueCodeDefaultAdminDelayChangeScheduled)
				if err := _ClueCode.contract.UnpackLog(event, "DefaultAdminDelayChangeScheduled", log); err != nil {
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

// ParseDefaultAdminDelayChangeScheduled is a log parse operation binding the contract event 0xf1038c18cf84a56e432fdbfaf746924b7ea511dfe03a6506a0ceba4888788d9b.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (_ClueCode *ClueCodeFilterer) ParseDefaultAdminDelayChangeScheduled(log types.Log) (*ClueCodeDefaultAdminDelayChangeScheduled, error) {
	event := new(ClueCodeDefaultAdminDelayChangeScheduled)
	if err := _ClueCode.contract.UnpackLog(event, "DefaultAdminDelayChangeScheduled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClueCodeDefaultAdminTransferCanceledIterator is returned from FilterDefaultAdminTransferCanceled and is used to iterate over the raw logs and unpacked data for DefaultAdminTransferCanceled events raised by the ClueCode contract.
type ClueCodeDefaultAdminTransferCanceledIterator struct {
	Event *ClueCodeDefaultAdminTransferCanceled // Event containing the contract specifics and raw log

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
func (it *ClueCodeDefaultAdminTransferCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClueCodeDefaultAdminTransferCanceled)
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
		it.Event = new(ClueCodeDefaultAdminTransferCanceled)
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
func (it *ClueCodeDefaultAdminTransferCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClueCodeDefaultAdminTransferCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClueCodeDefaultAdminTransferCanceled represents a DefaultAdminTransferCanceled event raised by the ClueCode contract.
type ClueCodeDefaultAdminTransferCanceled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminTransferCanceled is a free log retrieval operation binding the contract event 0x8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109.
//
// Solidity: event DefaultAdminTransferCanceled()
func (_ClueCode *ClueCodeFilterer) FilterDefaultAdminTransferCanceled(opts *bind.FilterOpts) (*ClueCodeDefaultAdminTransferCanceledIterator, error) {

	logs, sub, err := _ClueCode.contract.FilterLogs(opts, "DefaultAdminTransferCanceled")
	if err != nil {
		return nil, err
	}
	return &ClueCodeDefaultAdminTransferCanceledIterator{contract: _ClueCode.contract, event: "DefaultAdminTransferCanceled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminTransferCanceled is a free log subscription operation binding the contract event 0x8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109.
//
// Solidity: event DefaultAdminTransferCanceled()
func (_ClueCode *ClueCodeFilterer) WatchDefaultAdminTransferCanceled(opts *bind.WatchOpts, sink chan<- *ClueCodeDefaultAdminTransferCanceled) (event.Subscription, error) {

	logs, sub, err := _ClueCode.contract.WatchLogs(opts, "DefaultAdminTransferCanceled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClueCodeDefaultAdminTransferCanceled)
				if err := _ClueCode.contract.UnpackLog(event, "DefaultAdminTransferCanceled", log); err != nil {
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

// ParseDefaultAdminTransferCanceled is a log parse operation binding the contract event 0x8886ebfc4259abdbc16601dd8fb5678e54878f47b3c34836cfc51154a9605109.
//
// Solidity: event DefaultAdminTransferCanceled()
func (_ClueCode *ClueCodeFilterer) ParseDefaultAdminTransferCanceled(log types.Log) (*ClueCodeDefaultAdminTransferCanceled, error) {
	event := new(ClueCodeDefaultAdminTransferCanceled)
	if err := _ClueCode.contract.UnpackLog(event, "DefaultAdminTransferCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClueCodeDefaultAdminTransferScheduledIterator is returned from FilterDefaultAdminTransferScheduled and is used to iterate over the raw logs and unpacked data for DefaultAdminTransferScheduled events raised by the ClueCode contract.
type ClueCodeDefaultAdminTransferScheduledIterator struct {
	Event *ClueCodeDefaultAdminTransferScheduled // Event containing the contract specifics and raw log

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
func (it *ClueCodeDefaultAdminTransferScheduledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClueCodeDefaultAdminTransferScheduled)
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
		it.Event = new(ClueCodeDefaultAdminTransferScheduled)
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
func (it *ClueCodeDefaultAdminTransferScheduledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClueCodeDefaultAdminTransferScheduledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClueCodeDefaultAdminTransferScheduled represents a DefaultAdminTransferScheduled event raised by the ClueCode contract.
type ClueCodeDefaultAdminTransferScheduled struct {
	NewAdmin       common.Address
	AcceptSchedule *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDefaultAdminTransferScheduled is a free log retrieval operation binding the contract event 0x3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed6.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (_ClueCode *ClueCodeFilterer) FilterDefaultAdminTransferScheduled(opts *bind.FilterOpts, newAdmin []common.Address) (*ClueCodeDefaultAdminTransferScheduledIterator, error) {

	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _ClueCode.contract.FilterLogs(opts, "DefaultAdminTransferScheduled", newAdminRule)
	if err != nil {
		return nil, err
	}
	return &ClueCodeDefaultAdminTransferScheduledIterator{contract: _ClueCode.contract, event: "DefaultAdminTransferScheduled", logs: logs, sub: sub}, nil
}

// WatchDefaultAdminTransferScheduled is a free log subscription operation binding the contract event 0x3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed6.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (_ClueCode *ClueCodeFilterer) WatchDefaultAdminTransferScheduled(opts *bind.WatchOpts, sink chan<- *ClueCodeDefaultAdminTransferScheduled, newAdmin []common.Address) (event.Subscription, error) {

	var newAdminRule []interface{}
	for _, newAdminItem := range newAdmin {
		newAdminRule = append(newAdminRule, newAdminItem)
	}

	logs, sub, err := _ClueCode.contract.WatchLogs(opts, "DefaultAdminTransferScheduled", newAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClueCodeDefaultAdminTransferScheduled)
				if err := _ClueCode.contract.UnpackLog(event, "DefaultAdminTransferScheduled", log); err != nil {
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

// ParseDefaultAdminTransferScheduled is a log parse operation binding the contract event 0x3377dc44241e779dd06afab5b788a35ca5f3b778836e2990bdb26a2a4b2e5ed6.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (_ClueCode *ClueCodeFilterer) ParseDefaultAdminTransferScheduled(log types.Log) (*ClueCodeDefaultAdminTransferScheduled, error) {
	event := new(ClueCodeDefaultAdminTransferScheduled)
	if err := _ClueCode.contract.UnpackLog(event, "DefaultAdminTransferScheduled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClueCodeMetadataFrozenIterator is returned from FilterMetadataFrozen and is used to iterate over the raw logs and unpacked data for MetadataFrozen events raised by the ClueCode contract.
type ClueCodeMetadataFrozenIterator struct {
	Event *ClueCodeMetadataFrozen // Event containing the contract specifics and raw log

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
func (it *ClueCodeMetadataFrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClueCodeMetadataFrozen)
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
		it.Event = new(ClueCodeMetadataFrozen)
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
func (it *ClueCodeMetadataFrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClueCodeMetadataFrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClueCodeMetadataFrozen represents a MetadataFrozen event raised by the ClueCode contract.
type ClueCodeMetadataFrozen struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterMetadataFrozen is a free log retrieval operation binding the contract event 0xeef043febddf4e1d1cf1f72ff1407b84e036e805aa0934418cb82095da8d7164.
//
// Solidity: event MetadataFrozen()
func (_ClueCode *ClueCodeFilterer) FilterMetadataFrozen(opts *bind.FilterOpts) (*ClueCodeMetadataFrozenIterator, error) {

	logs, sub, err := _ClueCode.contract.FilterLogs(opts, "MetadataFrozen")
	if err != nil {
		return nil, err
	}
	return &ClueCodeMetadataFrozenIterator{contract: _ClueCode.contract, event: "MetadataFrozen", logs: logs, sub: sub}, nil
}

// WatchMetadataFrozen is a free log subscription operation binding the contract event 0xeef043febddf4e1d1cf1f72ff1407b84e036e805aa0934418cb82095da8d7164.
//
// Solidity: event MetadataFrozen()
func (_ClueCode *ClueCodeFilterer) WatchMetadataFrozen(opts *bind.WatchOpts, sink chan<- *ClueCodeMetadataFrozen) (event.Subscription, error) {

	logs, sub, err := _ClueCode.contract.WatchLogs(opts, "MetadataFrozen")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClueCodeMetadataFrozen)
				if err := _ClueCode.contract.UnpackLog(event, "MetadataFrozen", log); err != nil {
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

// ParseMetadataFrozen is a log parse operation binding the contract event 0xeef043febddf4e1d1cf1f72ff1407b84e036e805aa0934418cb82095da8d7164.
//
// Solidity: event MetadataFrozen()
func (_ClueCode *ClueCodeFilterer) ParseMetadataFrozen(log types.Log) (*ClueCodeMetadataFrozen, error) {
	event := new(ClueCodeMetadataFrozen)
	if err := _ClueCode.contract.UnpackLog(event, "MetadataFrozen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClueCodeMetadataSetIterator is returned from FilterMetadataSet and is used to iterate over the raw logs and unpacked data for MetadataSet events raised by the ClueCode contract.
type ClueCodeMetadataSetIterator struct {
	Event *ClueCodeMetadataSet // Event containing the contract specifics and raw log

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
func (it *ClueCodeMetadataSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClueCodeMetadataSet)
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
		it.Event = new(ClueCodeMetadataSet)
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
func (it *ClueCodeMetadataSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClueCodeMetadataSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClueCodeMetadataSet represents a MetadataSet event raised by the ClueCode contract.
type ClueCodeMetadataSet struct {
	Code *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterMetadataSet is a free log retrieval operation binding the contract event 0x043282acb0221e6f17ff7e7e56f271a3c5550143829948a198bc6025fa1e7b88.
//
// Solidity: event MetadataSet(uint256 indexed code)
func (_ClueCode *ClueCodeFilterer) FilterMetadataSet(opts *bind.FilterOpts, code []*big.Int) (*ClueCodeMetadataSetIterator, error) {

	var codeRule []interface{}
	for _, codeItem := range code {
		codeRule = append(codeRule, codeItem)
	}

	logs, sub, err := _ClueCode.contract.FilterLogs(opts, "MetadataSet", codeRule)
	if err != nil {
		return nil, err
	}
	return &ClueCodeMetadataSetIterator{contract: _ClueCode.contract, event: "MetadataSet", logs: logs, sub: sub}, nil
}

// WatchMetadataSet is a free log subscription operation binding the contract event 0x043282acb0221e6f17ff7e7e56f271a3c5550143829948a198bc6025fa1e7b88.
//
// Solidity: event MetadataSet(uint256 indexed code)
func (_ClueCode *ClueCodeFilterer) WatchMetadataSet(opts *bind.WatchOpts, sink chan<- *ClueCodeMetadataSet, code []*big.Int) (event.Subscription, error) {

	var codeRule []interface{}
	for _, codeItem := range code {
		codeRule = append(codeRule, codeItem)
	}

	logs, sub, err := _ClueCode.contract.WatchLogs(opts, "MetadataSet", codeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClueCodeMetadataSet)
				if err := _ClueCode.contract.UnpackLog(event, "MetadataSet", log); err != nil {
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

// ParseMetadataSet is a log parse operation binding the contract event 0x043282acb0221e6f17ff7e7e56f271a3c5550143829948a198bc6025fa1e7b88.
//
// Solidity: event MetadataSet(uint256 indexed code)
func (_ClueCode *ClueCodeFilterer) ParseMetadataSet(log types.Log) (*ClueCodeMetadataSet, error) {
	event := new(ClueCodeMetadataSet)
	if err := _ClueCode.contract.UnpackLog(event, "MetadataSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClueCodeMetadataUpdateIterator is returned from FilterMetadataUpdate and is used to iterate over the raw logs and unpacked data for MetadataUpdate events raised by the ClueCode contract.
type ClueCodeMetadataUpdateIterator struct {
	Event *ClueCodeMetadataUpdate // Event containing the contract specifics and raw log

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
func (it *ClueCodeMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClueCodeMetadataUpdate)
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
		it.Event = new(ClueCodeMetadataUpdate)
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
func (it *ClueCodeMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClueCodeMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClueCodeMetadataUpdate represents a MetadataUpdate event raised by the ClueCode contract.
type ClueCodeMetadataUpdate struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMetadataUpdate is a free log retrieval operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_ClueCode *ClueCodeFilterer) FilterMetadataUpdate(opts *bind.FilterOpts) (*ClueCodeMetadataUpdateIterator, error) {

	logs, sub, err := _ClueCode.contract.FilterLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &ClueCodeMetadataUpdateIterator{contract: _ClueCode.contract, event: "MetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchMetadataUpdate is a free log subscription operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_ClueCode *ClueCodeFilterer) WatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *ClueCodeMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _ClueCode.contract.WatchLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClueCodeMetadataUpdate)
				if err := _ClueCode.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
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

// ParseMetadataUpdate is a log parse operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_ClueCode *ClueCodeFilterer) ParseMetadataUpdate(log types.Log) (*ClueCodeMetadataUpdate, error) {
	event := new(ClueCodeMetadataUpdate)
	if err := _ClueCode.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClueCodeRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the ClueCode contract.
type ClueCodeRoleAdminChangedIterator struct {
	Event *ClueCodeRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *ClueCodeRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClueCodeRoleAdminChanged)
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
		it.Event = new(ClueCodeRoleAdminChanged)
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
func (it *ClueCodeRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClueCodeRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClueCodeRoleAdminChanged represents a RoleAdminChanged event raised by the ClueCode contract.
type ClueCodeRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ClueCode *ClueCodeFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*ClueCodeRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _ClueCode.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &ClueCodeRoleAdminChangedIterator{contract: _ClueCode.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ClueCode *ClueCodeFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *ClueCodeRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _ClueCode.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClueCodeRoleAdminChanged)
				if err := _ClueCode.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_ClueCode *ClueCodeFilterer) ParseRoleAdminChanged(log types.Log) (*ClueCodeRoleAdminChanged, error) {
	event := new(ClueCodeRoleAdminChanged)
	if err := _ClueCode.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClueCodeRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the ClueCode contract.
type ClueCodeRoleGrantedIterator struct {
	Event *ClueCodeRoleGranted // Event containing the contract specifics and raw log

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
func (it *ClueCodeRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClueCodeRoleGranted)
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
		it.Event = new(ClueCodeRoleGranted)
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
func (it *ClueCodeRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClueCodeRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClueCodeRoleGranted represents a RoleGranted event raised by the ClueCode contract.
type ClueCodeRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ClueCode *ClueCodeFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ClueCodeRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ClueCode.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ClueCodeRoleGrantedIterator{contract: _ClueCode.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ClueCode *ClueCodeFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *ClueCodeRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ClueCode.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClueCodeRoleGranted)
				if err := _ClueCode.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_ClueCode *ClueCodeFilterer) ParseRoleGranted(log types.Log) (*ClueCodeRoleGranted, error) {
	event := new(ClueCodeRoleGranted)
	if err := _ClueCode.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClueCodeRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the ClueCode contract.
type ClueCodeRoleRevokedIterator struct {
	Event *ClueCodeRoleRevoked // Event containing the contract specifics and raw log

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
func (it *ClueCodeRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClueCodeRoleRevoked)
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
		it.Event = new(ClueCodeRoleRevoked)
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
func (it *ClueCodeRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClueCodeRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClueCodeRoleRevoked represents a RoleRevoked event raised by the ClueCode contract.
type ClueCodeRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ClueCode *ClueCodeFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*ClueCodeRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ClueCode.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &ClueCodeRoleRevokedIterator{contract: _ClueCode.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ClueCode *ClueCodeFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *ClueCodeRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ClueCode.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClueCodeRoleRevoked)
				if err := _ClueCode.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_ClueCode *ClueCodeFilterer) ParseRoleRevoked(log types.Log) (*ClueCodeRoleRevoked, error) {
	event := new(ClueCodeRoleRevoked)
	if err := _ClueCode.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ClueCodeTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ClueCode contract.
type ClueCodeTransferIterator struct {
	Event *ClueCodeTransfer // Event containing the contract specifics and raw log

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
func (it *ClueCodeTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ClueCodeTransfer)
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
		it.Event = new(ClueCodeTransfer)
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
func (it *ClueCodeTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ClueCodeTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ClueCodeTransfer represents a Transfer event raised by the ClueCode contract.
type ClueCodeTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_ClueCode *ClueCodeFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*ClueCodeTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ClueCode.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ClueCodeTransferIterator{contract: _ClueCode.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_ClueCode *ClueCodeFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ClueCodeTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ClueCode.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ClueCodeTransfer)
				if err := _ClueCode.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_ClueCode *ClueCodeFilterer) ParseTransfer(log types.Log) (*ClueCodeTransfer, error) {
	event := new(ClueCodeTransfer)
	if err := _ClueCode.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
