package crypto

import (
	"errors"
	"github.com/hyperledger/fabric/core/crypto/primitives"
	"github.com/hyperledger/fabric/core/crypto/utils"
	obc "github.com/hyperledger/fabric/protos"
)

type validatorConfidentialityProcessorV1_1 struct {
	validator *validatorImpl
}

func (cp *validatorConfidentialityProcessorV1_1) getVersion() string {
	return "1.1"
}

func (cp *validatorConfidentialityProcessorV1_1) preValidation(ctx TransactionContext) (*obc.Transaction, error) {
	tx := ctx.GetTransaction()

	// TODO: check the flags
	if tx.Nonce == nil || len(tx.Nonce) == 0 {
		return nil, errors.New("Failed decrypting payload. Invalid nonce.")
	}

	// Clone tx
	tx, err := cp.validator.cloneTransaction(tx)
	if err != nil {
		cp.validator.error("Failed deep cloning [%s].", err.Error())
		return nil, err
	}

	// Derive root key
	// validator.enrollChainKey is an AES key represented as byte array
	enrollChainKey := cp.validator.enrollSymChainKey

	key := primitives.HMAC(enrollChainKey, tx.Nonce)

	// Decrypt ChaincodeID
	chaincodeIDKey := primitives.HMACTruncated(key, []byte{2}, primitives.AESKeyLength)
	chaincodeID, err := primitives.CBCPKCS7Decrypt(chaincodeIDKey, utils.Clone(tx.ChaincodeID))
	if err != nil {
		cp.validator.error("Failed decrypting chaincode [%s].", err.Error())
		return nil, err
	}
	tx.ChaincodeID = chaincodeID

	return tx, nil
}

func (cp *validatorConfidentialityProcessorV1_1) preExecution(ctx TransactionContext) (*obc.Transaction, error) {
	tx := ctx.GetTransaction()

	// TODO: check the flags
	if tx.Nonce == nil || len(tx.Nonce) == 0 {
		return nil, errors.New("Failed decrypting payload. Invalid nonce.")
	}

	// Clone tx
	tx, err := cp.validator.cloneTransaction(tx)
	if err != nil {
		cp.validator.error("Failed deep cloning [%s].", err.Error())
		return nil, err
	}

	// Derive root key
	// validator.enrollChainKey is an AES key represented as byte array
	enrollChainKey := cp.validator.enrollSymChainKey

	key := primitives.HMAC(enrollChainKey, tx.Nonce)

	//	validator.log.Info("Deriving from  ", utils.EncodeBase64(validator.peer.node.enrollChainKey))
	//	validator.log.Info("Nonce  ", utils.EncodeBase64(tx.Nonce))
	//	validator.log.Info("Derived key  ", utils.EncodeBase64(key))
	//	validator.log.Info("Encrypted Payload  ", utils.EncodeBase64(tx.EncryptedPayload))
	//	validator.log.Info("Encrypted ChaincodeID  ", utils.EncodeBase64(tx.EncryptedChaincodeID))

	// Decrypt Payload
	payloadKey := primitives.HMACTruncated(key, []byte{1}, primitives.AESKeyLength)
	payload, err := primitives.CBCPKCS7Decrypt(payloadKey, utils.Clone(tx.Payload))
	if err != nil {
		cp.validator.error("Failed decrypting payload [%s].", err.Error())
		return nil, err
	}
	tx.Payload = payload

	// ChaincodeID has been already decrypted by preValidation
	// Decrypt ChaincodeID
	//chaincodeIDKey := primitives.HMACTruncated(key, []byte{2}, primitives.AESKeyLength)
	//chaincodeID, err := primitives.CBCPKCS7Decrypt(chaincodeIDKey, utils.Clone(tx.ChaincodeID))
	//if err != nil {
	//	cp.validator.error("Failed decrypting chaincode [%s].", err.Error())
	//	return nil, err
	//}
	//tx.ChaincodeID = chaincodeID

	// Decrypt metadata
	if len(tx.Metadata) != 0 {
		metadataKey := primitives.HMACTruncated(key, []byte{3}, primitives.AESKeyLength)
		metadata, err := primitives.CBCPKCS7Decrypt(metadataKey, utils.Clone(tx.Metadata))
		if err != nil {
			cp.validator.error("Failed decrypting metadata [%s].", err.Error())
			return nil, err
		}
		tx.Metadata = metadata
	}

	return tx, nil
}
