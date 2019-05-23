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

// Defines default configs used for initializing nodes in dev settings.

package dev

import (
	"github.com/kardiachain/go-kardia/node"
	"fmt"
)

const (
	// GenesisAccount used for matchEth tx
	MockSmartContractCallSenderAccount = "0x7cefC13B6E2aedEeDFB7Cb6c32457240746BAEe5"
)

// Nodes are used for testing authorized node in private case
// From 0-9: authorized which are listed in kvm/smc/Permission.sol
// While 10 is not listed mean it is unauthorized.
var Nodes = []map[string]interface{}{
	{
		"key": "8843ebcb1021b00ae9a644db6617f9c6d870e5fd53624cefe374c1d2d710fd06",
		"votingPower": 100,
		"listenAddr": "[::]:3000",
	},
	{
		"key": "77cfc693f7861a6e1ea817c593c04fbc9b63d4d3146c5753c008cfc67cffca79",
		"votingPower": 100,
		"listenAddr": "[::]:3001",
	},
	{
		"key": "98de1df1e242afb02bd5dc01fbcacddcc9a4d41df95a66f629139560ca6e4dbb",
		"votingPower": 100,
		"listenAddr": "[::]:3002",
	},
	{
		"key": "32f5c0aef7f9172044a472478421c63fd8492640ff2d0eaab9562389db3a8efe",
		"votingPower": 100,
		"listenAddr": "[::]:3003",
	},
	{
		"key": "68b53a92d846baafdc782cb9cad65d77020c8d747eca7b621370b52b18c91f9a",
		"votingPower": 100,
		"listenAddr": "[::]:3004",
	},
	{
		"key": "049de018e08c3bcd59c1a21f0cf7de8f17fe51f8ce7d9c2120d17b1f0251b265",
		"votingPower": 100,
		"listenAddr": "[::]:3005",
	},
	{
		"key": "9fdd56a3c2a536dc8f981d935f0f3f2ea04e125547fdfffa37e157ce86ff1007",
		"votingPower": 100,
		"listenAddr": "[::]:3006",
	},
	{
		"key": "ae1a52546294bed6e734185775dbc84009de00bdf51b709471e2415c31ceeed7",
		"votingPower": 100,
		"listenAddr": "[::]:3007",
	},
	{
		"key": "b34bd81838a4a335fb3403d0bf616eca1eb9a4b4716c7dda7c617503cfeaab67",
		"votingPower": 100,
		"listenAddr": "[::]:3008",
	},
	{
		"key": "0cf7ae0332a891044659ace49a0732fa07c2872b4aef479945501f385a23e689",
		"votingPower": 100,
		"listenAddr": "[::]:3009",
	},
	{
		"key": "2003be66077b0873f5bedb32a596530eb8a0c908c32dda7771f169ee137c1f82",
		"votingPower": 100,
		"listenAddr": "[::]:3010",
	},
	{
		"key": "9dce5ec0b40e363e898f296c01345c12a0961f1cccad098964c73ed85fef5850",
		"votingPower": 100,
		"listenAddr": "[::]:3011",
	},
	{
		"key": "f0b2f6f24b70481a51712639badf0e5587545080dc53e0664770adb9881823fb",
		"votingPower": 100,
		"listenAddr": "[::]:3012",
	},
	{
		"key": "83731e17afb0da61c0026eaf780364eee367c50a5225ece89a63ad94a4a1f088",
		"votingPower": 100,
		"listenAddr": "[::]:3013",
	},
	{
		"key": "fc09d3f004b1ee430fee60568aa29748e277e76f1f372eea9d2b9ff1e27bdfdb",
		"votingPower": 100,
		"listenAddr": "[::]:3014",
	},
	{
		"key": "5605dd5f4db003c396956b4b80c093c472ccef4021181aa910125d7c57324152",
		"votingPower": 100,
		"listenAddr": "[::]:3015",
	},
	// the key below is used for test un-authorized node (private case)
	{
		"key": "0cf7ae0332a891044659ace49a0732fa07c2872b4aef479945501f385a23e690",
		"votingPower": 0,
		"listenAddr": "[::]:3016",
	},
	// Additional node key
	{"key": "1bb159419a15971f4b426cec5d00593c5048f1aadf43b761e17e1b48bc14e293", "votingPower": 100, "listenAddr": "[::]:3017"},
	{"key": "5f696991ed981a4b8b90fb827d1450dbd51571fdba61c076ca46fa652b684733", "votingPower": 100, "listenAddr": "[::]:3018"},
	{"key": "02d11b38e8d10bea10782ce18769a15fd386dde6b65f1a2be94ea1ea38c0bb5a", "votingPower": 100, "listenAddr": "[::]:3019"},
	{"key": "cb62868a15ab31a8e3c08d566fd9a57c2077dfc043ad4ddc0f6302fa6c9a7fb8", "votingPower": 100, "listenAddr": "[::]:3020"},
	{"key": "b40f7b4bfb2eefa25ffb89854ec74749944fec199e11520284896368bda5c7a7", "votingPower": 100, "listenAddr": "[::]:3021"},
	{"key": "080afbebf8774e5e901a3770848ab047e2112146365a7a48ab5b74e5892d6192", "votingPower": 100, "listenAddr": "[::]:3022"},
	{"key": "8e6b4527cc262cb01fbaedbbcf2fd811f1913b78de1ab59a59ba3e0998cf6038", "votingPower": 100, "listenAddr": "[::]:3023"},
	{"key": "bf6b2210ef5e98ff9bb02784cf7c60b231590c79db15183badcb6c4366e7fb61", "votingPower": 100, "listenAddr": "[::]:3024"},
	{"key": "7ac2d222eeba6b658878f9795a0df55dcd083b68e3e24317eb90ec9eaf0e64cc", "votingPower": 100, "listenAddr": "[::]:3025"},
	{"key": "72aa3a3c11f60b5dae1037c37b932f3773d75ce2c0ed0233f7c433c6687173e7", "votingPower": 100, "listenAddr": "[::]:3026"},
	{"key": "9148fec5b6a9e3f7c3b444dc828070cf038d14ccdd16ee43f3c19e52ac2874f8", "votingPower": 100, "listenAddr": "[::]:3027"},
	{"key": "24abc6535788d5b0dd0f44ed0f55e582c80f0dc732f7552c0194f45c6689e79a", "votingPower": 100, "listenAddr": "[::]:3028"},
	{"key": "b9ed85c2baac73bbb31ebd50d81288d0f66530cf703f3cb48ed122193a7f29f2", "votingPower": 100, "listenAddr": "[::]:3029"},
	{"key": "dd4efc5e39a8e5e5827f811016f96d0ba8b1838b8c5927f845aae30f71968083", "votingPower": 100, "listenAddr": "[::]:3030"},
	{"key": "6fe9c191f69bd7dbb5f83dbc5730bdc6d7d737b72f62397f2fe8f9283f0ecfef", "votingPower": 100, "listenAddr": "[::]:3031"},
	{"key": "3a88f0ef12b99885ba07f97860562cc6f52fc15f9709cc2bf85b1eb5621702e0", "votingPower": 100, "listenAddr": "[::]:3032"},
	{"key": "1a60a0f456e1cfacefcc79264bd8605738b4ca34031eabb7913b6ed467080cb7", "votingPower": 100, "listenAddr": "[::]:3033"},
	{"key": "0ff45f46e6c176bede483650e967b0e08a586a7b648ac597ead3ee22d2a1e315", "votingPower": 100, "listenAddr": "[::]:3034"},
	{"key": "3e743aac55bca357ff35325782c43bd3a502ae7c7412b99967ccc693e7c18679", "votingPower": 100, "listenAddr": "[::]:3035"},
	{"key": "feff7673f83d35f4d86ad73d2eb559745ef9b02585912f4867700587601018a7", "votingPower": 100, "listenAddr": "[::]:3036"},
	{"key": "28cda0a8b54f48f117433b99b774755502cd9048a987460a11a482c20c90168d", "votingPower": 100, "listenAddr": "[::]:3037"},
	{"key": "46988e5c869a3ed0638cf2e5ba5763627032123e36eb6805bbe63ac6872bcb8a", "votingPower": 100, "listenAddr": "[::]:3038"},
	{"key": "2777c45f4194653271fc8e2ff863cf48730872d36fd33bee271abc5c082dbd72", "votingPower": 100, "listenAddr": "[::]:3039"},
	{"key": "3be42586d1e99ce2368bfd8364423c80383384f4befe19c2c432959ef81d5b9a", "votingPower": 100, "listenAddr": "[::]:3040"},
	{"key": "b2c78ad8dbd09e7f3e02deb8c0881615d99185c0a414a843ec31ef7269f76d87", "votingPower": 100, "listenAddr": "[::]:3041"},
	{"key": "edb63f5935bea5191d3be528562c99f95712a4ce227d186e6c62055a86fd69f4", "votingPower": 100, "listenAddr": "[::]:3042"},
	{"key": "86c050d0af5e565419fa06516540d5132dcf9b3ad470ea7cb5fc440ac7c9aebb", "votingPower": 100, "listenAddr": "[::]:3043"},
	{"key": "12100610a1f48d40a3921f00fca05a6035e91e750e2b6a121c6b895d87a76303", "votingPower": 100, "listenAddr": "[::]:3044"},
	{"key": "b4771275a5d61859b7c5584d7f3a4a4c9d00aeab9734fe42de3e127a7372f665", "votingPower": 100, "listenAddr": "[::]:3045"},
	{"key": "67acf96df6d6560fffe9e093d8c1f1f2cbe18e812a4e4152933c665f80e2c876", "votingPower": 100, "listenAddr": "[::]:3046"},
	{"key": "a4327d0b840cf742183d842ac0eea1225cad39bb6e2f4af3738f4a0b495cac5f", "votingPower": 100, "listenAddr": "[::]:3047"},
	{"key": "bca1a95745a52f24f8100a770ebbb4ddc39490af0a8944605f73e198179edfb9", "votingPower": 100, "listenAddr": "[::]:3048"},
	{"key": "ed6c540fa4aa42ad9c34f113eea400ee01d6344dda4bee31fa6177a7577ffbc5", "votingPower": 100, "listenAddr": "[::]:3049"},
	{"key": "c18c313e48a6d11658b5f8eb963d15d83ffa74b32c02e446fffd5e91eb2d4d99", "votingPower": 100, "listenAddr": "[::]:3050"},
	{"key": "411e5bae7d61e28cbd3ccdf6ca3ddcbee150fe650915f47288e1caf16fd68a19", "votingPower": 100, "listenAddr": "[::]:3051"},
}

// GetNodeMetadataByIndex return NodeMetadata from nodes
func GetNodeMetadataByIndex(idx int) (*node.NodeMetadata, error) {
	if idx < 0 || idx >= len(Nodes) {
		return nil, fmt.Errorf("node index must be within 0 to %v", len(Nodes) - 1)
	}
	key := Nodes[idx]["key"].(string)
	votingPower := int64(Nodes[idx]["votingPower"].(int))
	listenAddr := Nodes[idx]["listenAddr"].(string)

	n, err := node.NewNodeMetadata(&key, nil, votingPower, listenAddr)
	if err != nil {
		return nil, err
	}
	return n, nil
}
