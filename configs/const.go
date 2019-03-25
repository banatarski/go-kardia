/*
 *  Copyright 2018 KardiaChain
 *  This file is part of the go-kardia library.
 *
 *  The go-kardia library is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Lesser General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  The go-kardia library is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 *  GNU Lesser General Public License for more details.
 *
 *  You should have received a copy of the GNU Lesser General Public License
 *  along with the go-kardia library. If not, see <http://www.gnu.org/licenses/>.
 */

package configs

import "errors"

// All const related to cross-chain demos including coin exchange and candidate exchange
// this will be dynamic and removed when run on production
const (
	// constants related to currency exchange
	KardiaNewExchangeSmcIndex          = 3
	CompleteFunction                   = "completeRequest"
	ExternalDepositFunction            = "deposit"
	ETH2NEO                            = "ETH-NEO"
	NEO2ETH                            = "NEO-ETH"
	ETH                                = "ETH"
	NEO                                = "NEO"
	AddOrderFunction                   = "addOrder"
	RateETH                            = 100000000
	RateNEO                            = 6482133
	ExchangeDataCompleteRequestIDIndex = 0
	ExchangeDataCompletePairIndex      = 1
	NumOfCompleteRequestDataField      = 2

	// Constantg related to exchange v2 which support original tx id
	ExchangeV2SourcePairIndex        = 0
	ExchangeV2DestPairIndex          = 1
	ExchangeV2SourceAddressIndex     = 2
	ExchangeV2DestAddressIndex       = 3
	ExchangeV2OriginalTxIdIndex      = 4
	ExchangeV2AmountIndex            = 5
	ExchangeV2TimestampIndex         = 6
	ExchangeV2NumOfExchangeDataField = 7
	ExchangeV2ReleaseFieldsSeparator = "|"
	ExchangeV2ReleaseToTypeIndex     = 0
	ExchangeV2ReleaseAddressesIndex  = 1
	ExchangeV2ReleaseAmountsIndex    = 2
	ExchangeV2ReleaseTxIdsIndex      = 3
	ExchangeV2ReleaseValuesSepatator = ";"

	// Constants related to candidate exchange, Kardia part
	KardiaCandidateExchangeSmcIndex    = 6
	KardiaForwardRequestFunction       = "forwardRequest"
	KardiaForwardResponseFunction      = "forwardResponse"
	KardiaForwardResponseFields        = 4
	KardiaForwardResponseEmailIndex    = 0
	KardiaForwardResponseResponseIndex = 1
	KardiaForwardResponseFromOrgIndex  = 2
	KardiaForwardResponseToOrgIndex    = 3
	KardiaForwardRequestFields         = 3
	KardiaForwardRequestEmailIndex     = 0
	KardiaForwardRequestFromOrgIndex   = 1
	KardiaForwardRequestToOrgIndex     = 2

	// Constants related to candidate exchange, private chain part
	PrivateChainCandidateDBSmcIndex                     = 5
	PrivateChainCandidateRequestCompletedFields         = 4
	PrivateChainCandidateRequestCompletedFromOrgIDIndex = 0
	PrivateChainCandidateRequestCompletedToOrgIDIndex   = 1
	PrivateChainCandidateRequestCompletedEmailIndex     = 2
	PrivateChainCandidateRequestCompletedContentIndex   = 3
	PrivateChainRequestInfoFunction                     = "requestCandidateInfo"
	PrivateChainCompleteRequestFunction                 = "completeRequest"
	PrivateChainCandidateRequestFields                  = 3
	PrivateChainCandidateRequestEmailIndex              = 0
	PrivateChainCandidateRequestFromOrgIndex            = 1
	PrivateChainCandidateRequestToOrgIndex              = 2
)

var (
	ErrTypeConversionFailed             = errors.New("fail type conversion")
	ErrInsufficientCandidateRequestData = errors.New("insufficient candidate request data")
	ErrNoMatchedRequest                 = errors.New("request has no matched opponent")
	ErrNotImplemented                   = errors.New("this function is not implemented yet")
	ErrInsufficientExchangeData = errors.New("insufficient exchange external data")
	ErrUnsupportedMethod        = errors.New("method is not supported by dual logic")
	ErrCreateKardiaTx           = errors.New("fail to create Kardia's Tx from DualEvent")
	ErrAddKardiaTx              = errors.New("fail to add Tx to Kardia's TxPool")
	ErrFailedGetState           = errors.New("fail to get Kardia state")
	ErrFailedGetEventData       = errors.New("fail to get event external data")
)
