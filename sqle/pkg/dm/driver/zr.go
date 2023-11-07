/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"dm/util"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	Dm_build_675 = "7.6.0.0"

	Dm_build_676 = "7.0.0.9"

	Dm_build_677 = "8.0.0.73"

	Dm_build_678 = "7.1.2.128"

	Dm_build_679 = "7.1.5.144"

	Dm_build_680 = "7.1.6.123"

	Dm_build_681 = 2

	Dm_build_682 = 32768 - 128

	Dm_build_683 = 0x20000000

	Dm_build_684 int16 = 1

	Dm_build_685 int16 = 2

	Dm_build_686 int16 = 3

	Dm_build_687 int16 = 4

	Dm_build_688 int16 = 5

	Dm_build_689 int16 = 6

	Dm_build_690 int16 = 7

	Dm_build_691 int16 = 8

	Dm_build_692 int16 = 9

	Dm_build_693 int16 = 13

	Dm_build_694 int16 = 14

	Dm_build_695 int16 = 15

	Dm_build_696 int16 = 17

	Dm_build_697 int16 = 21

	Dm_build_698 int16 = 24

	Dm_build_699 int16 = 27

	Dm_build_700 int16 = 29

	Dm_build_701 int16 = 30

	Dm_build_702 int16 = 31

	Dm_build_703 int16 = 32

	Dm_build_704 int16 = 44

	Dm_build_705 int16 = 52

	Dm_build_706 int16 = 60

	Dm_build_707 int16 = 71

	Dm_build_708 int16 = 90

	Dm_build_709 int16 = 91

	Dm_build_710 int16 = 200

	Dm_build_711 = 64

	Dm_build_712 = 20

	Dm_build_713 = 0

	Dm_build_714 = 4

	Dm_build_715 = 6

	Dm_build_716 = 10

	Dm_build_717 = 14

	Dm_build_718 = 18

	Dm_build_719 = 19

	Dm_build_720 = 128

	Dm_build_721 = 256

	Dm_build_722 = 0xffff

	Dm_build_723 int32 = 2

	Dm_build_724 int32 = 5

	Dm_build_725 = -1

	Dm_build_726 uint16 = 0xFFFE

	Dm_build_727 uint16 = uint16(Dm_build_726 - 3)

	Dm_build_728 uint16 = Dm_build_726

	Dm_build_729 int32 = 0xFFFF

	Dm_build_730 int32 = 0x80

	Dm_build_731 byte = 0x60

	Dm_build_732 uint16 = uint16(Dm_build_728)

	Dm_build_733 uint16 = uint16(Dm_build_729)

	Dm_build_734 int16 = 0x00

	Dm_build_735 int16 = 0x03

	Dm_build_736 int32 = 0x80

	Dm_build_737 byte = 0

	Dm_build_738 byte = 1

	Dm_build_739 byte = 2

	Dm_build_740 byte = 3

	Dm_build_741 byte = 4

	Dm_build_742 byte = Dm_build_737

	Dm_build_743 int = 10

	Dm_build_744 int32 = 32

	Dm_build_745 int32 = 65536

	Dm_build_746 byte = 0

	Dm_build_747 byte = 1

	Dm_build_748 int32 = 0x00000000

	Dm_build_749 int32 = 0x00000020

	Dm_build_750 int32 = 0x00000040

	Dm_build_751 int32 = 0x00000FFF

	Dm_build_752 int32 = 0

	Dm_build_753 int32 = 1

	Dm_build_754 int32 = 2

	Dm_build_755 int32 = 3

	Dm_build_756 = 8192

	Dm_build_757 = 1

	Dm_build_758 = 2

	Dm_build_759 = 0

	Dm_build_760 = 0

	Dm_build_761 = 1

	Dm_build_762 = -1

	Dm_build_763 int16 = 0

	Dm_build_764 int16 = 1

	Dm_build_765 int16 = 2

	Dm_build_766 int16 = 3

	Dm_build_767 int16 = 4

	Dm_build_768 int16 = 127

	Dm_build_769 int16 = Dm_build_768 + 20

	Dm_build_770 int16 = Dm_build_768 + 21

	Dm_build_771 int16 = Dm_build_768 + 22

	Dm_build_772 int16 = Dm_build_768 + 24

	Dm_build_773 int16 = Dm_build_768 + 25

	Dm_build_774 int16 = Dm_build_768 + 26

	Dm_build_775 int16 = Dm_build_768 + 30

	Dm_build_776 int16 = Dm_build_768 + 31

	Dm_build_777 int16 = Dm_build_768 + 32

	Dm_build_778 int16 = Dm_build_768 + 33

	Dm_build_779 int16 = Dm_build_768 + 35

	Dm_build_780 int16 = Dm_build_768 + 38

	Dm_build_781 int16 = Dm_build_768 + 39

	Dm_build_782 int16 = Dm_build_768 + 51

	Dm_build_783 int16 = Dm_build_768 + 71

	Dm_build_784 int16 = Dm_build_768 + 124

	Dm_build_785 int16 = Dm_build_768 + 125

	Dm_build_786 int16 = Dm_build_768 + 126

	Dm_build_787 int16 = Dm_build_768 + 127

	Dm_build_788 int16 = Dm_build_768 + 128

	Dm_build_789 int16 = Dm_build_768 + 129

	Dm_build_790 byte = 0

	Dm_build_791 byte = 2

	Dm_build_792 = 2048

	Dm_build_793 = -1

	Dm_build_794 = 0

	Dm_build_795 = 16000

	Dm_build_796 = 32000

	Dm_build_797 = 0x00000000

	Dm_build_798 = 0x00000020

	Dm_build_799 = 0x00000040

	Dm_build_800 = 0x00000FFF

	Dm_build_801 = 4
)

var Dm_build_802 = [8][256]uint32{

	{0x00000000, 0x77073096, 0xee0e612c, 0x990951ba, 0x076dc419, 0x706af48f, 0xe963a535,
		0x9e6495a3, 0x0edb8832, 0x79dcb8a4, 0xe0d5e91e, 0x97d2d988, 0x09b64c2b,
		0x7eb17cbd, 0xe7b82d07, 0x90bf1d91, 0x1db71064, 0x6ab020f2, 0xf3b97148,
		0x84be41de, 0x1adad47d, 0x6ddde4eb, 0xf4d4b551, 0x83d385c7, 0x136c9856,
		0x646ba8c0, 0xfd62f97a, 0x8a65c9ec, 0x14015c4f, 0x63066cd9, 0xfa0f3d63,
		0x8d080df5, 0x3b6e20c8, 0x4c69105e, 0xd56041e4, 0xa2677172, 0x3c03e4d1,
		0x4b04d447, 0xd20d85fd, 0xa50ab56b, 0x35b5a8fa, 0x42b2986c, 0xdbbbc9d6,
		0xacbcf940, 0x32d86ce3, 0x45df5c75, 0xdcd60dcf, 0xabd13d59, 0x26d930ac,
		0x51de003a, 0xc8d75180, 0xbfd06116, 0x21b4f4b5, 0x56b3c423, 0xcfba9599,
		0xb8bda50f, 0x2802b89e, 0x5f058808, 0xc60cd9b2, 0xb10be924, 0x2f6f7c87,
		0x58684c11, 0xc1611dab, 0xb6662d3d, 0x76dc4190, 0x01db7106, 0x98d220bc,
		0xefd5102a, 0x71b18589, 0x06b6b51f, 0x9fbfe4a5, 0xe8b8d433, 0x7807c9a2,
		0x0f00f934, 0x9609a88e, 0xe10e9818, 0x7f6a0dbb, 0x086d3d2d, 0x91646c97,
		0xe6635c01, 0x6b6b51f4, 0x1c6c6162, 0x856530d8, 0xf262004e, 0x6c0695ed,
		0x1b01a57b, 0x8208f4c1, 0xf50fc457, 0x65b0d9c6, 0x12b7e950, 0x8bbeb8ea,
		0xfcb9887c, 0x62dd1ddf, 0x15da2d49, 0x8cd37cf3, 0xfbd44c65, 0x4db26158,
		0x3ab551ce, 0xa3bc0074, 0xd4bb30e2, 0x4adfa541, 0x3dd895d7, 0xa4d1c46d,
		0xd3d6f4fb, 0x4369e96a, 0x346ed9fc, 0xad678846, 0xda60b8d0, 0x44042d73,
		0x33031de5, 0xaa0a4c5f, 0xdd0d7cc9, 0x5005713c, 0x270241aa, 0xbe0b1010,
		0xc90c2086, 0x5768b525, 0x206f85b3, 0xb966d409, 0xce61e49f, 0x5edef90e,
		0x29d9c998, 0xb0d09822, 0xc7d7a8b4, 0x59b33d17, 0x2eb40d81, 0xb7bd5c3b,
		0xc0ba6cad, 0xedb88320, 0x9abfb3b6, 0x03b6e20c, 0x74b1d29a, 0xead54739,
		0x9dd277af, 0x04db2615, 0x73dc1683, 0xe3630b12, 0x94643b84, 0x0d6d6a3e,
		0x7a6a5aa8, 0xe40ecf0b, 0x9309ff9d, 0x0a00ae27, 0x7d079eb1, 0xf00f9344,
		0x8708a3d2, 0x1e01f268, 0x6906c2fe, 0xf762575d, 0x806567cb, 0x196c3671,
		0x6e6b06e7, 0xfed41b76, 0x89d32be0, 0x10da7a5a, 0x67dd4acc, 0xf9b9df6f,
		0x8ebeeff9, 0x17b7be43, 0x60b08ed5, 0xd6d6a3e8, 0xa1d1937e, 0x38d8c2c4,
		0x4fdff252, 0xd1bb67f1, 0xa6bc5767, 0x3fb506dd, 0x48b2364b, 0xd80d2bda,
		0xaf0a1b4c, 0x36034af6, 0x41047a60, 0xdf60efc3, 0xa867df55, 0x316e8eef,
		0x4669be79, 0xcb61b38c, 0xbc66831a, 0x256fd2a0, 0x5268e236, 0xcc0c7795,
		0xbb0b4703, 0x220216b9, 0x5505262f, 0xc5ba3bbe, 0xb2bd0b28, 0x2bb45a92,
		0x5cb36a04, 0xc2d7ffa7, 0xb5d0cf31, 0x2cd99e8b, 0x5bdeae1d, 0x9b64c2b0,
		0xec63f226, 0x756aa39c, 0x026d930a, 0x9c0906a9, 0xeb0e363f, 0x72076785,
		0x05005713, 0x95bf4a82, 0xe2b87a14, 0x7bb12bae, 0x0cb61b38, 0x92d28e9b,
		0xe5d5be0d, 0x7cdcefb7, 0x0bdbdf21, 0x86d3d2d4, 0xf1d4e242, 0x68ddb3f8,
		0x1fda836e, 0x81be16cd, 0xf6b9265b, 0x6fb077e1, 0x18b74777, 0x88085ae6,
		0xff0f6a70, 0x66063bca, 0x11010b5c, 0x8f659eff, 0xf862ae69, 0x616bffd3,
		0x166ccf45, 0xa00ae278, 0xd70dd2ee, 0x4e048354, 0x3903b3c2, 0xa7672661,
		0xd06016f7, 0x4969474d, 0x3e6e77db, 0xaed16a4a, 0xd9d65adc, 0x40df0b66,
		0x37d83bf0, 0xa9bcae53, 0xdebb9ec5, 0x47b2cf7f, 0x30b5ffe9, 0xbdbdf21c,
		0xcabac28a, 0x53b39330, 0x24b4a3a6, 0xbad03605, 0xcdd70693, 0x54de5729,
		0x23d967bf, 0xb3667a2e, 0xc4614ab8, 0x5d681b02, 0x2a6f2b94, 0xb40bbe37,
		0xc30c8ea1, 0x5a05df1b, 0x2d02ef8d},

	{0x00000000, 0x191b3141, 0x32366282, 0x2b2d53c3, 0x646cc504, 0x7d77f445, 0x565aa786,
		0x4f4196c7, 0xc8d98a08, 0xd1c2bb49, 0xfaefe88a, 0xe3f4d9cb, 0xacb54f0c,
		0xb5ae7e4d, 0x9e832d8e, 0x87981ccf, 0x4ac21251, 0x53d92310, 0x78f470d3,
		0x61ef4192, 0x2eaed755, 0x37b5e614, 0x1c98b5d7, 0x05838496, 0x821b9859,
		0x9b00a918, 0xb02dfadb, 0xa936cb9a, 0xe6775d5d, 0xff6c6c1c, 0xd4413fdf,
		0xcd5a0e9e, 0x958424a2, 0x8c9f15e3, 0xa7b24620, 0xbea97761, 0xf1e8e1a6,
		0xe8f3d0e7, 0xc3de8324, 0xdac5b265, 0x5d5daeaa, 0x44469feb, 0x6f6bcc28,
		0x7670fd69, 0x39316bae, 0x202a5aef, 0x0b07092c, 0x121c386d, 0xdf4636f3,
		0xc65d07b2, 0xed705471, 0xf46b6530, 0xbb2af3f7, 0xa231c2b6, 0x891c9175,
		0x9007a034, 0x179fbcfb, 0x0e848dba, 0x25a9de79, 0x3cb2ef38, 0x73f379ff,
		0x6ae848be, 0x41c51b7d, 0x58de2a3c, 0xf0794f05, 0xe9627e44, 0xc24f2d87,
		0xdb541cc6, 0x94158a01, 0x8d0ebb40, 0xa623e883, 0xbf38d9c2, 0x38a0c50d,
		0x21bbf44c, 0x0a96a78f, 0x138d96ce, 0x5ccc0009, 0x45d73148, 0x6efa628b,
		0x77e153ca, 0xbabb5d54, 0xa3a06c15, 0x888d3fd6, 0x91960e97, 0xded79850,
		0xc7cca911, 0xece1fad2, 0xf5facb93, 0x7262d75c, 0x6b79e61d, 0x4054b5de,
		0x594f849f, 0x160e1258, 0x0f152319, 0x243870da, 0x3d23419b, 0x65fd6ba7,
		0x7ce65ae6, 0x57cb0925, 0x4ed03864, 0x0191aea3, 0x188a9fe2, 0x33a7cc21,
		0x2abcfd60, 0xad24e1af, 0xb43fd0ee, 0x9f12832d, 0x8609b26c, 0xc94824ab,
		0xd05315ea, 0xfb7e4629, 0xe2657768, 0x2f3f79f6, 0x362448b7, 0x1d091b74,
		0x04122a35, 0x4b53bcf2, 0x52488db3, 0x7965de70, 0x607eef31, 0xe7e6f3fe,
		0xfefdc2bf, 0xd5d0917c, 0xcccba03d, 0x838a36fa, 0x9a9107bb, 0xb1bc5478,
		0xa8a76539, 0x3b83984b, 0x2298a90a, 0x09b5fac9, 0x10aecb88, 0x5fef5d4f,
		0x46f46c0e, 0x6dd93fcd, 0x74c20e8c, 0xf35a1243, 0xea412302, 0xc16c70c1,
		0xd8774180, 0x9736d747, 0x8e2de606, 0xa500b5c5, 0xbc1b8484, 0x71418a1a,
		0x685abb5b, 0x4377e898, 0x5a6cd9d9, 0x152d4f1e, 0x0c367e5f, 0x271b2d9c,
		0x3e001cdd, 0xb9980012, 0xa0833153, 0x8bae6290, 0x92b553d1, 0xddf4c516,
		0xc4eff457, 0xefc2a794, 0xf6d996d5, 0xae07bce9, 0xb71c8da8, 0x9c31de6b,
		0x852aef2a, 0xca6b79ed, 0xd37048ac, 0xf85d1b6f, 0xe1462a2e, 0x66de36e1,
		0x7fc507a0, 0x54e85463, 0x4df36522, 0x02b2f3e5, 0x1ba9c2a4, 0x30849167,
		0x299fa026, 0xe4c5aeb8, 0xfdde9ff9, 0xd6f3cc3a, 0xcfe8fd7b, 0x80a96bbc,
		0x99b25afd, 0xb29f093e, 0xab84387f, 0x2c1c24b0, 0x350715f1, 0x1e2a4632,
		0x07317773, 0x4870e1b4, 0x516bd0f5, 0x7a468336, 0x635db277, 0xcbfad74e,
		0xd2e1e60f, 0xf9ccb5cc, 0xe0d7848d, 0xaf96124a, 0xb68d230b, 0x9da070c8,
		0x84bb4189, 0x03235d46, 0x1a386c07, 0x31153fc4, 0x280e0e85, 0x674f9842,
		0x7e54a903, 0x5579fac0, 0x4c62cb81, 0x8138c51f, 0x9823f45e, 0xb30ea79d,
		0xaa1596dc, 0xe554001b, 0xfc4f315a, 0xd7626299, 0xce7953d8, 0x49e14f17,
		0x50fa7e56, 0x7bd72d95, 0x62cc1cd4, 0x2d8d8a13, 0x3496bb52, 0x1fbbe891,
		0x06a0d9d0, 0x5e7ef3ec, 0x4765c2ad, 0x6c48916e, 0x7553a02f, 0x3a1236e8,
		0x230907a9, 0x0824546a, 0x113f652b, 0x96a779e4, 0x8fbc48a5, 0xa4911b66,
		0xbd8a2a27, 0xf2cbbce0, 0xebd08da1, 0xc0fdde62, 0xd9e6ef23, 0x14bce1bd,
		0x0da7d0fc, 0x268a833f, 0x3f91b27e, 0x70d024b9, 0x69cb15f8, 0x42e6463b,
		0x5bfd777a, 0xdc656bb5, 0xc57e5af4, 0xee530937, 0xf7483876, 0xb809aeb1,
		0xa1129ff0, 0x8a3fcc33, 0x9324fd72},

	{0x00000000, 0x01c26a37, 0x0384d46e, 0x0246be59, 0x0709a8dc, 0x06cbc2eb, 0x048d7cb2,
		0x054f1685, 0x0e1351b8, 0x0fd13b8f, 0x0d9785d6, 0x0c55efe1, 0x091af964,
		0x08d89353, 0x0a9e2d0a, 0x0b5c473d, 0x1c26a370, 0x1de4c947, 0x1fa2771e,
		0x1e601d29, 0x1b2f0bac, 0x1aed619b, 0x18abdfc2, 0x1969b5f5, 0x1235f2c8,
		0x13f798ff, 0x11b126a6, 0x10734c91, 0x153c5a14, 0x14fe3023, 0x16b88e7a,
		0x177ae44d, 0x384d46e0, 0x398f2cd7, 0x3bc9928e, 0x3a0bf8b9, 0x3f44ee3c,
		0x3e86840b, 0x3cc03a52, 0x3d025065, 0x365e1758, 0x379c7d6f, 0x35dac336,
		0x3418a901, 0x3157bf84, 0x3095d5b3, 0x32d36bea, 0x331101dd, 0x246be590,
		0x25a98fa7, 0x27ef31fe, 0x262d5bc9, 0x23624d4c, 0x22a0277b, 0x20e69922,
		0x2124f315, 0x2a78b428, 0x2bbade1f, 0x29fc6046, 0x283e0a71, 0x2d711cf4,
		0x2cb376c3, 0x2ef5c89a, 0x2f37a2ad, 0x709a8dc0, 0x7158e7f7, 0x731e59ae,
		0x72dc3399, 0x7793251c, 0x76514f2b, 0x7417f172, 0x75d59b45, 0x7e89dc78,
		0x7f4bb64f, 0x7d0d0816, 0x7ccf6221, 0x798074a4, 0x78421e93, 0x7a04a0ca,
		0x7bc6cafd, 0x6cbc2eb0, 0x6d7e4487, 0x6f38fade, 0x6efa90e9, 0x6bb5866c,
		0x6a77ec5b, 0x68315202, 0x69f33835, 0x62af7f08, 0x636d153f, 0x612bab66,
		0x60e9c151, 0x65a6d7d4, 0x6464bde3, 0x662203ba, 0x67e0698d, 0x48d7cb20,
		0x4915a117, 0x4b531f4e, 0x4a917579, 0x4fde63fc, 0x4e1c09cb, 0x4c5ab792,
		0x4d98dda5, 0x46c49a98, 0x4706f0af, 0x45404ef6, 0x448224c1, 0x41cd3244,
		0x400f5873, 0x4249e62a, 0x438b8c1d, 0x54f16850, 0x55330267, 0x5775bc3e,
		0x56b7d609, 0x53f8c08c, 0x523aaabb, 0x507c14e2, 0x51be7ed5, 0x5ae239e8,
		0x5b2053df, 0x5966ed86, 0x58a487b1, 0x5deb9134, 0x5c29fb03, 0x5e6f455a,
		0x5fad2f6d, 0xe1351b80, 0xe0f771b7, 0xe2b1cfee, 0xe373a5d9, 0xe63cb35c,
		0xe7fed96b, 0xe5b86732, 0xe47a0d05, 0xef264a38, 0xeee4200f, 0xeca29e56,
		0xed60f461, 0xe82fe2e4, 0xe9ed88d3, 0xebab368a, 0xea695cbd, 0xfd13b8f0,
		0xfcd1d2c7, 0xfe976c9e, 0xff5506a9, 0xfa1a102c, 0xfbd87a1b, 0xf99ec442,
		0xf85cae75, 0xf300e948, 0xf2c2837f, 0xf0843d26, 0xf1465711, 0xf4094194,
		0xf5cb2ba3, 0xf78d95fa, 0xf64fffcd, 0xd9785d60, 0xd8ba3757, 0xdafc890e,
		0xdb3ee339, 0xde71f5bc, 0xdfb39f8b, 0xddf521d2, 0xdc374be5, 0xd76b0cd8,
		0xd6a966ef, 0xd4efd8b6, 0xd52db281, 0xd062a404, 0xd1a0ce33, 0xd3e6706a,
		0xd2241a5d, 0xc55efe10, 0xc49c9427, 0xc6da2a7e, 0xc7184049, 0xc25756cc,
		0xc3953cfb, 0xc1d382a2, 0xc011e895, 0xcb4dafa8, 0xca8fc59f, 0xc8c97bc6,
		0xc90b11f1, 0xcc440774, 0xcd866d43, 0xcfc0d31a, 0xce02b92d, 0x91af9640,
		0x906dfc77, 0x922b422e, 0x93e92819, 0x96a63e9c, 0x976454ab, 0x9522eaf2,
		0x94e080c5, 0x9fbcc7f8, 0x9e7eadcf, 0x9c381396, 0x9dfa79a1, 0x98b56f24,
		0x99770513, 0x9b31bb4a, 0x9af3d17d, 0x8d893530, 0x8c4b5f07, 0x8e0de15e,
		0x8fcf8b69, 0x8a809dec, 0x8b42f7db, 0x89044982, 0x88c623b5, 0x839a6488,
		0x82580ebf, 0x801eb0e6, 0x81dcdad1, 0x8493cc54, 0x8551a663, 0x8717183a,
		0x86d5720d, 0xa9e2d0a0, 0xa820ba97, 0xaa6604ce, 0xaba46ef9, 0xaeeb787c,
		0xaf29124b, 0xad6fac12, 0xacadc625, 0xa7f18118, 0xa633eb2f, 0xa4755576,
		0xa5b73f41, 0xa0f829c4, 0xa13a43f3, 0xa37cfdaa, 0xa2be979d, 0xb5c473d0,
		0xb40619e7, 0xb640a7be, 0xb782cd89, 0xb2cddb0c, 0xb30fb13b, 0xb1490f62,
		0xb08b6555, 0xbbd72268, 0xba15485f, 0xb853f606, 0xb9919c31, 0xbcde8ab4,
		0xbd1ce083, 0xbf5a5eda, 0xbe9834ed},

	{0x00000000, 0xb8bc6765, 0xaa09c88b, 0x12b5afee, 0x8f629757, 0x37def032, 0x256b5fdc,
		0x9dd738b9, 0xc5b428ef, 0x7d084f8a, 0x6fbde064, 0xd7018701, 0x4ad6bfb8,
		0xf26ad8dd, 0xe0df7733, 0x58631056, 0x5019579f, 0xe8a530fa, 0xfa109f14,
		0x42acf871, 0xdf7bc0c8, 0x67c7a7ad, 0x75720843, 0xcdce6f26, 0x95ad7f70,
		0x2d111815, 0x3fa4b7fb, 0x8718d09e, 0x1acfe827, 0xa2738f42, 0xb0c620ac,
		0x087a47c9, 0xa032af3e, 0x188ec85b, 0x0a3b67b5, 0xb28700d0, 0x2f503869,
		0x97ec5f0c, 0x8559f0e2, 0x3de59787, 0x658687d1, 0xdd3ae0b4, 0xcf8f4f5a,
		0x7733283f, 0xeae41086, 0x525877e3, 0x40edd80d, 0xf851bf68, 0xf02bf8a1,
		0x48979fc4, 0x5a22302a, 0xe29e574f, 0x7f496ff6, 0xc7f50893, 0xd540a77d,
		0x6dfcc018, 0x359fd04e, 0x8d23b72b, 0x9f9618c5, 0x272a7fa0, 0xbafd4719,
		0x0241207c, 0x10f48f92, 0xa848e8f7, 0x9b14583d, 0x23a83f58, 0x311d90b6,
		0x89a1f7d3, 0x1476cf6a, 0xaccaa80f, 0xbe7f07e1, 0x06c36084, 0x5ea070d2,
		0xe61c17b7, 0xf4a9b859, 0x4c15df3c, 0xd1c2e785, 0x697e80e0, 0x7bcb2f0e,
		0xc377486b, 0xcb0d0fa2, 0x73b168c7, 0x6104c729, 0xd9b8a04c, 0x446f98f5,
		0xfcd3ff90, 0xee66507e, 0x56da371b, 0x0eb9274d, 0xb6054028, 0xa4b0efc6,
		0x1c0c88a3, 0x81dbb01a, 0x3967d77f, 0x2bd27891, 0x936e1ff4, 0x3b26f703,
		0x839a9066, 0x912f3f88, 0x299358ed, 0xb4446054, 0x0cf80731, 0x1e4da8df,
		0xa6f1cfba, 0xfe92dfec, 0x462eb889, 0x549b1767, 0xec277002, 0x71f048bb,
		0xc94c2fde, 0xdbf98030, 0x6345e755, 0x6b3fa09c, 0xd383c7f9, 0xc1366817,
		0x798a0f72, 0xe45d37cb, 0x5ce150ae, 0x4e54ff40, 0xf6e89825, 0xae8b8873,
		0x1637ef16, 0x048240f8, 0xbc3e279d, 0x21e91f24, 0x99557841, 0x8be0d7af,
		0x335cb0ca, 0xed59b63b, 0x55e5d15e, 0x47507eb0, 0xffec19d5, 0x623b216c,
		0xda874609, 0xc832e9e7, 0x708e8e82, 0x28ed9ed4, 0x9051f9b1, 0x82e4565f,
		0x3a58313a, 0xa78f0983, 0x1f336ee6, 0x0d86c108, 0xb53aa66d, 0xbd40e1a4,
		0x05fc86c1, 0x1749292f, 0xaff54e4a, 0x322276f3, 0x8a9e1196, 0x982bbe78,
		0x2097d91d, 0x78f4c94b, 0xc048ae2e, 0xd2fd01c0, 0x6a4166a5, 0xf7965e1c,
		0x4f2a3979, 0x5d9f9697, 0xe523f1f2, 0x4d6b1905, 0xf5d77e60, 0xe762d18e,
		0x5fdeb6eb, 0xc2098e52, 0x7ab5e937, 0x680046d9, 0xd0bc21bc, 0x88df31ea,
		0x3063568f, 0x22d6f961, 0x9a6a9e04, 0x07bda6bd, 0xbf01c1d8, 0xadb46e36,
		0x15080953, 0x1d724e9a, 0xa5ce29ff, 0xb77b8611, 0x0fc7e174, 0x9210d9cd,
		0x2aacbea8, 0x38191146, 0x80a57623, 0xd8c66675, 0x607a0110, 0x72cfaefe,
		0xca73c99b, 0x57a4f122, 0xef189647, 0xfdad39a9, 0x45115ecc, 0x764dee06,
		0xcef18963, 0xdc44268d, 0x64f841e8, 0xf92f7951, 0x41931e34, 0x5326b1da,
		0xeb9ad6bf, 0xb3f9c6e9, 0x0b45a18c, 0x19f00e62, 0xa14c6907, 0x3c9b51be,
		0x842736db, 0x96929935, 0x2e2efe50, 0x2654b999, 0x9ee8defc, 0x8c5d7112,
		0x34e11677, 0xa9362ece, 0x118a49ab, 0x033fe645, 0xbb838120, 0xe3e09176,
		0x5b5cf613, 0x49e959fd, 0xf1553e98, 0x6c820621, 0xd43e6144, 0xc68bceaa,
		0x7e37a9cf, 0xd67f4138, 0x6ec3265d, 0x7c7689b3, 0xc4caeed6, 0x591dd66f,
		0xe1a1b10a, 0xf3141ee4, 0x4ba87981, 0x13cb69d7, 0xab770eb2, 0xb9c2a15c,
		0x017ec639, 0x9ca9fe80, 0x241599e5, 0x36a0360b, 0x8e1c516e, 0x866616a7,
		0x3eda71c2, 0x2c6fde2c, 0x94d3b949, 0x090481f0, 0xb1b8e695, 0xa30d497b,
		0x1bb12e1e, 0x43d23e48, 0xfb6e592d, 0xe9dbf6c3, 0x516791a6, 0xccb0a91f,
		0x740cce7a, 0x66b96194, 0xde0506f1},

	{0x00000000, 0x3d6029b0, 0x7ac05360, 0x47a07ad0, 0xf580a6c0, 0xc8e08f70, 0x8f40f5a0,
		0xb220dc10, 0x30704bc1, 0x0d106271, 0x4ab018a1, 0x77d03111, 0xc5f0ed01,
		0xf890c4b1, 0xbf30be61, 0x825097d1, 0x60e09782, 0x5d80be32, 0x1a20c4e2,
		0x2740ed52, 0x95603142, 0xa80018f2, 0xefa06222, 0xd2c04b92, 0x5090dc43,
		0x6df0f5f3, 0x2a508f23, 0x1730a693, 0xa5107a83, 0x98705333, 0xdfd029e3,
		0xe2b00053, 0xc1c12f04, 0xfca106b4, 0xbb017c64, 0x866155d4, 0x344189c4,
		0x0921a074, 0x4e81daa4, 0x73e1f314, 0xf1b164c5, 0xccd14d75, 0x8b7137a5,
		0xb6111e15, 0x0431c205, 0x3951ebb5, 0x7ef19165, 0x4391b8d5, 0xa121b886,
		0x9c419136, 0xdbe1ebe6, 0xe681c256, 0x54a11e46, 0x69c137f6, 0x2e614d26,
		0x13016496, 0x9151f347, 0xac31daf7, 0xeb91a027, 0xd6f18997, 0x64d15587,
		0x59b17c37, 0x1e1106e7, 0x23712f57, 0x58f35849, 0x659371f9, 0x22330b29,
		0x1f532299, 0xad73fe89, 0x9013d739, 0xd7b3ade9, 0xead38459, 0x68831388,
		0x55e33a38, 0x124340e8, 0x2f236958, 0x9d03b548, 0xa0639cf8, 0xe7c3e628,
		0xdaa3cf98, 0x3813cfcb, 0x0573e67b, 0x42d39cab, 0x7fb3b51b, 0xcd93690b,
		0xf0f340bb, 0xb7533a6b, 0x8a3313db, 0x0863840a, 0x3503adba, 0x72a3d76a,
		0x4fc3feda, 0xfde322ca, 0xc0830b7a, 0x872371aa, 0xba43581a, 0x9932774d,
		0xa4525efd, 0xe3f2242d, 0xde920d9d, 0x6cb2d18d, 0x51d2f83d, 0x167282ed,
		0x2b12ab5d, 0xa9423c8c, 0x9422153c, 0xd3826fec, 0xeee2465c, 0x5cc29a4c,
		0x61a2b3fc, 0x2602c92c, 0x1b62e09c, 0xf9d2e0cf, 0xc4b2c97f, 0x8312b3af,
		0xbe729a1f, 0x0c52460f, 0x31326fbf, 0x7692156f, 0x4bf23cdf, 0xc9a2ab0e,
		0xf4c282be, 0xb362f86e, 0x8e02d1de, 0x3c220dce, 0x0142247e, 0x46e25eae,
		0x7b82771e, 0xb1e6b092, 0x8c869922, 0xcb26e3f2, 0xf646ca42, 0x44661652,
		0x79063fe2, 0x3ea64532, 0x03c66c82, 0x8196fb53, 0xbcf6d2e3, 0xfb56a833,
		0xc6368183, 0x74165d93, 0x49767423, 0x0ed60ef3, 0x33b62743, 0xd1062710,
		0xec660ea0, 0xabc67470, 0x96a65dc0, 0x248681d0, 0x19e6a860, 0x5e46d2b0,
		0x6326fb00, 0xe1766cd1, 0xdc164561, 0x9bb63fb1, 0xa6d61601, 0x14f6ca11,
		0x2996e3a1, 0x6e369971, 0x5356b0c1, 0x70279f96, 0x4d47b626, 0x0ae7ccf6,
		0x3787e546, 0x85a73956, 0xb8c710e6, 0xff676a36, 0xc2074386, 0x4057d457,
		0x7d37fde7, 0x3a978737, 0x07f7ae87, 0xb5d77297, 0x88b75b27, 0xcf1721f7,
		0xf2770847, 0x10c70814, 0x2da721a4, 0x6a075b74, 0x576772c4, 0xe547aed4,
		0xd8278764, 0x9f87fdb4, 0xa2e7d404, 0x20b743d5, 0x1dd76a65, 0x5a7710b5,
		0x67173905, 0xd537e515, 0xe857cca5, 0xaff7b675, 0x92979fc5, 0xe915e8db,
		0xd475c16b, 0x93d5bbbb, 0xaeb5920b, 0x1c954e1b, 0x21f567ab, 0x66551d7b,
		0x5b3534cb, 0xd965a31a, 0xe4058aaa, 0xa3a5f07a, 0x9ec5d9ca, 0x2ce505da,
		0x11852c6a, 0x562556ba, 0x6b457f0a, 0x89f57f59, 0xb49556e9, 0xf3352c39,
		0xce550589, 0x7c75d999, 0x4115f029, 0x06b58af9, 0x3bd5a349, 0xb9853498,
		0x84e51d28, 0xc34567f8, 0xfe254e48, 0x4c059258, 0x7165bbe8, 0x36c5c138,
		0x0ba5e888, 0x28d4c7df, 0x15b4ee6f, 0x521494bf, 0x6f74bd0f, 0xdd54611f,
		0xe03448af, 0xa794327f, 0x9af41bcf, 0x18a48c1e, 0x25c4a5ae, 0x6264df7e,
		0x5f04f6ce, 0xed242ade, 0xd044036e, 0x97e479be, 0xaa84500e, 0x4834505d,
		0x755479ed, 0x32f4033d, 0x0f942a8d, 0xbdb4f69d, 0x80d4df2d, 0xc774a5fd,
		0xfa148c4d, 0x78441b9c, 0x4524322c, 0x028448fc, 0x3fe4614c, 0x8dc4bd5c,
		0xb0a494ec, 0xf704ee3c, 0xca64c78c},

	{0x00000000, 0xcb5cd3a5, 0x4dc8a10b, 0x869472ae, 0x9b914216, 0x50cd91b3, 0xd659e31d,
		0x1d0530b8, 0xec53826d, 0x270f51c8, 0xa19b2366, 0x6ac7f0c3, 0x77c2c07b,
		0xbc9e13de, 0x3a0a6170, 0xf156b2d5, 0x03d6029b, 0xc88ad13e, 0x4e1ea390,
		0x85427035, 0x9847408d, 0x531b9328, 0xd58fe186, 0x1ed33223, 0xef8580f6,
		0x24d95353, 0xa24d21fd, 0x6911f258, 0x7414c2e0, 0xbf481145, 0x39dc63eb,
		0xf280b04e, 0x07ac0536, 0xccf0d693, 0x4a64a43d, 0x81387798, 0x9c3d4720,
		0x57619485, 0xd1f5e62b, 0x1aa9358e, 0xebff875b, 0x20a354fe, 0xa6372650,
		0x6d6bf5f5, 0x706ec54d, 0xbb3216e8, 0x3da66446, 0xf6fab7e3, 0x047a07ad,
		0xcf26d408, 0x49b2a6a6, 0x82ee7503, 0x9feb45bb, 0x54b7961e, 0xd223e4b0,
		0x197f3715, 0xe82985c0, 0x23755665, 0xa5e124cb, 0x6ebdf76e, 0x73b8c7d6,
		0xb8e41473, 0x3e7066dd, 0xf52cb578, 0x0f580a6c, 0xc404d9c9, 0x4290ab67,
		0x89cc78c2, 0x94c9487a, 0x5f959bdf, 0xd901e971, 0x125d3ad4, 0xe30b8801,
		0x28575ba4, 0xaec3290a, 0x659ffaaf, 0x789aca17, 0xb3c619b2, 0x35526b1c,
		0xfe0eb8b9, 0x0c8e08f7, 0xc7d2db52, 0x4146a9fc, 0x8a1a7a59, 0x971f4ae1,
		0x5c439944, 0xdad7ebea, 0x118b384f, 0xe0dd8a9a, 0x2b81593f, 0xad152b91,
		0x6649f834, 0x7b4cc88c, 0xb0101b29, 0x36846987, 0xfdd8ba22, 0x08f40f5a,
		0xc3a8dcff, 0x453cae51, 0x8e607df4, 0x93654d4c, 0x58399ee9, 0xdeadec47,
		0x15f13fe2, 0xe4a78d37, 0x2ffb5e92, 0xa96f2c3c, 0x6233ff99, 0x7f36cf21,
		0xb46a1c84, 0x32fe6e2a, 0xf9a2bd8f, 0x0b220dc1, 0xc07ede64, 0x46eaacca,
		0x8db67f6f, 0x90b34fd7, 0x5bef9c72, 0xdd7beedc, 0x16273d79, 0xe7718fac,
		0x2c2d5c09, 0xaab92ea7, 0x61e5fd02, 0x7ce0cdba, 0xb7bc1e1f, 0x31286cb1,
		0xfa74bf14, 0x1eb014d8, 0xd5ecc77d, 0x5378b5d3, 0x98246676, 0x852156ce,
		0x4e7d856b, 0xc8e9f7c5, 0x03b52460, 0xf2e396b5, 0x39bf4510, 0xbf2b37be,
		0x7477e41b, 0x6972d4a3, 0xa22e0706, 0x24ba75a8, 0xefe6a60d, 0x1d661643,
		0xd63ac5e6, 0x50aeb748, 0x9bf264ed, 0x86f75455, 0x4dab87f0, 0xcb3ff55e,
		0x006326fb, 0xf135942e, 0x3a69478b, 0xbcfd3525, 0x77a1e680, 0x6aa4d638,
		0xa1f8059d, 0x276c7733, 0xec30a496, 0x191c11ee, 0xd240c24b, 0x54d4b0e5,
		0x9f886340, 0x828d53f8, 0x49d1805d, 0xcf45f2f3, 0x04192156, 0xf54f9383,
		0x3e134026, 0xb8873288, 0x73dbe12d, 0x6eded195, 0xa5820230, 0x2316709e,
		0xe84aa33b, 0x1aca1375, 0xd196c0d0, 0x5702b27e, 0x9c5e61db, 0x815b5163,
		0x4a0782c6, 0xcc93f068, 0x07cf23cd, 0xf6999118, 0x3dc542bd, 0xbb513013,
		0x700de3b6, 0x6d08d30e, 0xa65400ab, 0x20c07205, 0xeb9ca1a0, 0x11e81eb4,
		0xdab4cd11, 0x5c20bfbf, 0x977c6c1a, 0x8a795ca2, 0x41258f07, 0xc7b1fda9,
		0x0ced2e0c, 0xfdbb9cd9, 0x36e74f7c, 0xb0733dd2, 0x7b2fee77, 0x662adecf,
		0xad760d6a, 0x2be27fc4, 0xe0beac61, 0x123e1c2f, 0xd962cf8a, 0x5ff6bd24,
		0x94aa6e81, 0x89af5e39, 0x42f38d9c, 0xc467ff32, 0x0f3b2c97, 0xfe6d9e42,
		0x35314de7, 0xb3a53f49, 0x78f9ecec, 0x65fcdc54, 0xaea00ff1, 0x28347d5f,
		0xe368aefa, 0x16441b82, 0xdd18c827, 0x5b8cba89, 0x90d0692c, 0x8dd55994,
		0x46898a31, 0xc01df89f, 0x0b412b3a, 0xfa1799ef, 0x314b4a4a, 0xb7df38e4,
		0x7c83eb41, 0x6186dbf9, 0xaada085c, 0x2c4e7af2, 0xe712a957, 0x15921919,
		0xdececabc, 0x585ab812, 0x93066bb7, 0x8e035b0f, 0x455f88aa, 0xc3cbfa04,
		0x089729a1, 0xf9c19b74, 0x329d48d1, 0xb4093a7f, 0x7f55e9da, 0x6250d962,
		0xa90c0ac7, 0x2f987869, 0xe4c4abcc},

	{0x00000000, 0xa6770bb4, 0x979f1129, 0x31e81a9d, 0xf44f2413, 0x52382fa7, 0x63d0353a,
		0xc5a73e8e, 0x33ef4e67, 0x959845d3, 0xa4705f4e, 0x020754fa, 0xc7a06a74,
		0x61d761c0, 0x503f7b5d, 0xf64870e9, 0x67de9cce, 0xc1a9977a, 0xf0418de7,
		0x56368653, 0x9391b8dd, 0x35e6b369, 0x040ea9f4, 0xa279a240, 0x5431d2a9,
		0xf246d91d, 0xc3aec380, 0x65d9c834, 0xa07ef6ba, 0x0609fd0e, 0x37e1e793,
		0x9196ec27, 0xcfbd399c, 0x69ca3228, 0x582228b5, 0xfe552301, 0x3bf21d8f,
		0x9d85163b, 0xac6d0ca6, 0x0a1a0712, 0xfc5277fb, 0x5a257c4f, 0x6bcd66d2,
		0xcdba6d66, 0x081d53e8, 0xae6a585c, 0x9f8242c1, 0x39f54975, 0xa863a552,
		0x0e14aee6, 0x3ffcb47b, 0x998bbfcf, 0x5c2c8141, 0xfa5b8af5, 0xcbb39068,
		0x6dc49bdc, 0x9b8ceb35, 0x3dfbe081, 0x0c13fa1c, 0xaa64f1a8, 0x6fc3cf26,
		0xc9b4c492, 0xf85cde0f, 0x5e2bd5bb, 0x440b7579, 0xe27c7ecd, 0xd3946450,
		0x75e36fe4, 0xb044516a, 0x16335ade, 0x27db4043, 0x81ac4bf7, 0x77e43b1e,
		0xd19330aa, 0xe07b2a37, 0x460c2183, 0x83ab1f0d, 0x25dc14b9, 0x14340e24,
		0xb2430590, 0x23d5e9b7, 0x85a2e203, 0xb44af89e, 0x123df32a, 0xd79acda4,
		0x71edc610, 0x4005dc8d, 0xe672d739, 0x103aa7d0, 0xb64dac64, 0x87a5b6f9,
		0x21d2bd4d, 0xe47583c3, 0x42028877, 0x73ea92ea, 0xd59d995e, 0x8bb64ce5,
		0x2dc14751, 0x1c295dcc, 0xba5e5678, 0x7ff968f6, 0xd98e6342, 0xe86679df,
		0x4e11726b, 0xb8590282, 0x1e2e0936, 0x2fc613ab, 0x89b1181f, 0x4c162691,
		0xea612d25, 0xdb8937b8, 0x7dfe3c0c, 0xec68d02b, 0x4a1fdb9f, 0x7bf7c102,
		0xdd80cab6, 0x1827f438, 0xbe50ff8c, 0x8fb8e511, 0x29cfeea5, 0xdf879e4c,
		0x79f095f8, 0x48188f65, 0xee6f84d1, 0x2bc8ba5f, 0x8dbfb1eb, 0xbc57ab76,
		0x1a20a0c2, 0x8816eaf2, 0x2e61e146, 0x1f89fbdb, 0xb9fef06f, 0x7c59cee1,
		0xda2ec555, 0xebc6dfc8, 0x4db1d47c, 0xbbf9a495, 0x1d8eaf21, 0x2c66b5bc,
		0x8a11be08, 0x4fb68086, 0xe9c18b32, 0xd82991af, 0x7e5e9a1b, 0xefc8763c,
		0x49bf7d88, 0x78576715, 0xde206ca1, 0x1b87522f, 0xbdf0599b, 0x8c184306,
		0x2a6f48b2, 0xdc27385b, 0x7a5033ef, 0x4bb82972, 0xedcf22c6, 0x28681c48,
		0x8e1f17fc, 0xbff70d61, 0x198006d5, 0x47abd36e, 0xe1dcd8da, 0xd034c247,
		0x7643c9f3, 0xb3e4f77d, 0x1593fcc9, 0x247be654, 0x820cede0, 0x74449d09,
		0xd23396bd, 0xe3db8c20, 0x45ac8794, 0x800bb91a, 0x267cb2ae, 0x1794a833,
		0xb1e3a387, 0x20754fa0, 0x86024414, 0xb7ea5e89, 0x119d553d, 0xd43a6bb3,
		0x724d6007, 0x43a57a9a, 0xe5d2712e, 0x139a01c7, 0xb5ed0a73, 0x840510ee,
		0x22721b5a, 0xe7d525d4, 0x41a22e60, 0x704a34fd, 0xd63d3f49, 0xcc1d9f8b,
		0x6a6a943f, 0x5b828ea2, 0xfdf58516, 0x3852bb98, 0x9e25b02c, 0xafcdaab1,
		0x09baa105, 0xfff2d1ec, 0x5985da58, 0x686dc0c5, 0xce1acb71, 0x0bbdf5ff,
		0xadcafe4b, 0x9c22e4d6, 0x3a55ef62, 0xabc30345, 0x0db408f1, 0x3c5c126c,
		0x9a2b19d8, 0x5f8c2756, 0xf9fb2ce2, 0xc813367f, 0x6e643dcb, 0x982c4d22,
		0x3e5b4696, 0x0fb35c0b, 0xa9c457bf, 0x6c636931, 0xca146285, 0xfbfc7818,
		0x5d8b73ac, 0x03a0a617, 0xa5d7ada3, 0x943fb73e, 0x3248bc8a, 0xf7ef8204,
		0x519889b0, 0x6070932d, 0xc6079899, 0x304fe870, 0x9638e3c4, 0xa7d0f959,
		0x01a7f2ed, 0xc400cc63, 0x6277c7d7, 0x539fdd4a, 0xf5e8d6fe, 0x647e3ad9,
		0xc209316d, 0xf3e12bf0, 0x55962044, 0x90311eca, 0x3646157e, 0x07ae0fe3,
		0xa1d90457, 0x579174be, 0xf1e67f0a, 0xc00e6597, 0x66796e23, 0xa3de50ad,
		0x05a95b19, 0x34414184, 0x92364a30},

	{0x00000000, 0xccaa009e, 0x4225077d, 0x8e8f07e3, 0x844a0efa, 0x48e00e64, 0xc66f0987,
		0x0ac50919, 0xd3e51bb5, 0x1f4f1b2b, 0x91c01cc8, 0x5d6a1c56, 0x57af154f,
		0x9b0515d1, 0x158a1232, 0xd92012ac, 0x7cbb312b, 0xb01131b5, 0x3e9e3656,
		0xf23436c8, 0xf8f13fd1, 0x345b3f4f, 0xbad438ac, 0x767e3832, 0xaf5e2a9e,
		0x63f42a00, 0xed7b2de3, 0x21d12d7d, 0x2b142464, 0xe7be24fa, 0x69312319,
		0xa59b2387, 0xf9766256, 0x35dc62c8, 0xbb53652b, 0x77f965b5, 0x7d3c6cac,
		0xb1966c32, 0x3f196bd1, 0xf3b36b4f, 0x2a9379e3, 0xe639797d, 0x68b67e9e,
		0xa41c7e00, 0xaed97719, 0x62737787, 0xecfc7064, 0x205670fa, 0x85cd537d,
		0x496753e3, 0xc7e85400, 0x0b42549e, 0x01875d87, 0xcd2d5d19, 0x43a25afa,
		0x8f085a64, 0x562848c8, 0x9a824856, 0x140d4fb5, 0xd8a74f2b, 0xd2624632,
		0x1ec846ac, 0x9047414f, 0x5ced41d1, 0x299dc2ed, 0xe537c273, 0x6bb8c590,
		0xa712c50e, 0xadd7cc17, 0x617dcc89, 0xeff2cb6a, 0x2358cbf4, 0xfa78d958,
		0x36d2d9c6, 0xb85dde25, 0x74f7debb, 0x7e32d7a2, 0xb298d73c, 0x3c17d0df,
		0xf0bdd041, 0x5526f3c6, 0x998cf358, 0x1703f4bb, 0xdba9f425, 0xd16cfd3c,
		0x1dc6fda2, 0x9349fa41, 0x5fe3fadf, 0x86c3e873, 0x4a69e8ed, 0xc4e6ef0e,
		0x084cef90, 0x0289e689, 0xce23e617, 0x40ace1f4, 0x8c06e16a, 0xd0eba0bb,
		0x1c41a025, 0x92cea7c6, 0x5e64a758, 0x54a1ae41, 0x980baedf, 0x1684a93c,
		0xda2ea9a2, 0x030ebb0e, 0xcfa4bb90, 0x412bbc73, 0x8d81bced, 0x8744b5f4,
		0x4beeb56a, 0xc561b289, 0x09cbb217, 0xac509190, 0x60fa910e, 0xee7596ed,
		0x22df9673, 0x281a9f6a, 0xe4b09ff4, 0x6a3f9817, 0xa6959889, 0x7fb58a25,
		0xb31f8abb, 0x3d908d58, 0xf13a8dc6, 0xfbff84df, 0x37558441, 0xb9da83a2,
		0x7570833c, 0x533b85da, 0x9f918544, 0x111e82a7, 0xddb48239, 0xd7718b20,
		0x1bdb8bbe, 0x95548c5d, 0x59fe8cc3, 0x80de9e6f, 0x4c749ef1, 0xc2fb9912,
		0x0e51998c, 0x04949095, 0xc83e900b, 0x46b197e8, 0x8a1b9776, 0x2f80b4f1,
		0xe32ab46f, 0x6da5b38c, 0xa10fb312, 0xabcaba0b, 0x6760ba95, 0xe9efbd76,
		0x2545bde8, 0xfc65af44, 0x30cfafda, 0xbe40a839, 0x72eaa8a7, 0x782fa1be,
		0xb485a120, 0x3a0aa6c3, 0xf6a0a65d, 0xaa4de78c, 0x66e7e712, 0xe868e0f1,
		0x24c2e06f, 0x2e07e976, 0xe2ade9e8, 0x6c22ee0b, 0xa088ee95, 0x79a8fc39,
		0xb502fca7, 0x3b8dfb44, 0xf727fbda, 0xfde2f2c3, 0x3148f25d, 0xbfc7f5be,
		0x736df520, 0xd6f6d6a7, 0x1a5cd639, 0x94d3d1da, 0x5879d144, 0x52bcd85d,
		0x9e16d8c3, 0x1099df20, 0xdc33dfbe, 0x0513cd12, 0xc9b9cd8c, 0x4736ca6f,
		0x8b9ccaf1, 0x8159c3e8, 0x4df3c376, 0xc37cc495, 0x0fd6c40b, 0x7aa64737,
		0xb60c47a9, 0x3883404a, 0xf42940d4, 0xfeec49cd, 0x32464953, 0xbcc94eb0,
		0x70634e2e, 0xa9435c82, 0x65e95c1c, 0xeb665bff, 0x27cc5b61, 0x2d095278,
		0xe1a352e6, 0x6f2c5505, 0xa386559b, 0x061d761c, 0xcab77682, 0x44387161,
		0x889271ff, 0x825778e6, 0x4efd7878, 0xc0727f9b, 0x0cd87f05, 0xd5f86da9,
		0x19526d37, 0x97dd6ad4, 0x5b776a4a, 0x51b26353, 0x9d1863cd, 0x1397642e,
		0xdf3d64b0, 0x83d02561, 0x4f7a25ff, 0xc1f5221c, 0x0d5f2282, 0x079a2b9b,
		0xcb302b05, 0x45bf2ce6, 0x89152c78, 0x50353ed4, 0x9c9f3e4a, 0x121039a9,
		0xdeba3937, 0xd47f302e, 0x18d530b0, 0x965a3753, 0x5af037cd, 0xff6b144a,
		0x33c114d4, 0xbd4e1337, 0x71e413a9, 0x7b211ab0, 0xb78b1a2e, 0x39041dcd,
		0xf5ae1d53, 0x2c8e0fff, 0xe0240f61, 0x6eab0882, 0xa201081c, 0xa8c40105,
		0x646e019b, 0xeae10678, 0x264b06e6}}

type dm_build_803 interface {
	dm_build_804()
	dm_build_805() error
	dm_build_806()
	dm_build_807(imsg dm_build_803) error
	dm_build_808() error
	dm_build_809() (interface{}, error)
	dm_build_810()
	dm_build_811(imsg dm_build_803) (interface{}, error)
	dm_build_812()
	dm_build_813() error
	dm_build_814() byte
	dm_build_815(buffer *Dm_build_78, startOff int32, endOff int32) uint32
	dm_build_816() int32
	dm_build_817(length int32)
	dm_build_818() int16
}

type dm_build_819 struct {
	dm_build_820 *dm_build_410

	dm_build_821 int16

	dm_build_822 int32

	dm_build_823 *DmStatement
}

func (dm_build_825 *dm_build_819) dm_build_824(dm_build_826 *dm_build_410, dm_build_827 int16) *dm_build_819 {
	dm_build_825.dm_build_820 = dm_build_826
	dm_build_825.dm_build_821 = dm_build_827
	return dm_build_825
}

func (dm_build_829 *dm_build_819) dm_build_828(dm_build_830 *dm_build_410, dm_build_831 int16, dm_build_832 *DmStatement) *dm_build_819 {
	dm_build_829.dm_build_824(dm_build_830, dm_build_831).dm_build_823 = dm_build_832
	return dm_build_829
}

func dm_build_833(dm_build_834 *dm_build_410, dm_build_835 int16) *dm_build_819 {
	return new(dm_build_819).dm_build_824(dm_build_834, dm_build_835)
}

func dm_build_836(dm_build_837 *dm_build_410, dm_build_838 int16, dm_build_839 *DmStatement) *dm_build_819 {
	return new(dm_build_819).dm_build_828(dm_build_837, dm_build_838, dm_build_839)
}

func (dm_build_841 *dm_build_819) dm_build_804() {
	dm_build_841.dm_build_820.dm_build_413.Dm_build_92(0)
	dm_build_841.dm_build_820.dm_build_413.Dm_build_103(Dm_build_711, true, true)
}

func (dm_build_843 *dm_build_819) dm_build_805() error {
	return nil
}

func (dm_build_845 *dm_build_819) dm_build_806() {
	if dm_build_845.dm_build_823 == nil {
		dm_build_845.dm_build_820.dm_build_413.Dm_build_269(Dm_build_713, 0)
	} else {
		dm_build_845.dm_build_820.dm_build_413.Dm_build_269(Dm_build_713, dm_build_845.dm_build_823.id)
	}

	dm_build_845.dm_build_820.dm_build_413.Dm_build_265(Dm_build_714, dm_build_845.dm_build_821)
	dm_build_845.dm_build_820.dm_build_413.Dm_build_269(Dm_build_715, int32(dm_build_845.dm_build_820.dm_build_413.Dm_build_90()-Dm_build_711))
}

func (dm_build_847 *dm_build_819) dm_build_808() error {
	dm_build_847.dm_build_820.dm_build_413.Dm_build_95(0)
	dm_build_847.dm_build_820.dm_build_413.Dm_build_103(Dm_build_711, false, true)
	return dm_build_847.dm_build_852()
}

func (dm_build_849 *dm_build_819) dm_build_809() (interface{}, error) {
	return nil, nil
}

func (dm_build_851 *dm_build_819) dm_build_810() {
}

func (dm_build_853 *dm_build_819) dm_build_852() error {
	dm_build_853.dm_build_822 = dm_build_853.dm_build_820.dm_build_413.Dm_build_347(Dm_build_716)
	if dm_build_853.dm_build_822 < 0 && dm_build_853.dm_build_822 != EC_RN_EXCEED_ROWSET_SIZE.ErrCode {
		return (&DmError{dm_build_853.dm_build_822, dm_build_853.dm_build_854(), nil, ""}).throw()
	} else if dm_build_853.dm_build_822 > 0 {

	} else if dm_build_853.dm_build_821 == Dm_build_710 || dm_build_853.dm_build_821 == Dm_build_684 {
		dm_build_853.dm_build_854()
	}

	return nil
}

func (dm_build_855 *dm_build_819) dm_build_854() string {

	dm_build_856 := dm_build_855.dm_build_820.dm_build_414.getServerEncoding()

	if dm_build_856 != "" && dm_build_856 == ENCODING_EUCKR && Locale != LANGUAGE_EN {
		dm_build_856 = ENCODING_GB18030
	}

	dm_build_855.dm_build_820.dm_build_413.Dm_build_103(int(dm_build_855.dm_build_820.dm_build_413.Dm_build_203()), false, true)

	dm_build_855.dm_build_820.dm_build_413.Dm_build_103(int(dm_build_855.dm_build_820.dm_build_413.Dm_build_203()), false, true)

	dm_build_855.dm_build_820.dm_build_413.Dm_build_103(int(dm_build_855.dm_build_820.dm_build_413.Dm_build_203()), false, true)

	return dm_build_855.dm_build_820.dm_build_413.Dm_build_245(dm_build_856, dm_build_855.dm_build_820.dm_build_414)
}

func (dm_build_858 *dm_build_819) dm_build_807(dm_build_859 dm_build_803) (dm_build_860 error) {
	dm_build_859.dm_build_804()
	if dm_build_860 = dm_build_859.dm_build_805(); dm_build_860 != nil {
		return dm_build_860
	}
	dm_build_859.dm_build_806()
	return nil
}

func (dm_build_862 *dm_build_819) dm_build_811(dm_build_863 dm_build_803) (dm_build_864 interface{}, dm_build_865 error) {
	dm_build_865 = dm_build_863.dm_build_808()
	if dm_build_865 != nil {
		return nil, dm_build_865
	}
	dm_build_864, dm_build_865 = dm_build_863.dm_build_809()
	if dm_build_865 != nil {
		return nil, dm_build_865
	}
	dm_build_863.dm_build_810()
	return dm_build_864, nil
}

func (dm_build_867 *dm_build_819) dm_build_812() {
	if dm_build_867.dm_build_820.dm_build_419 {

		var orgLen = dm_build_867.dm_build_816()

		dm_build_867.dm_build_817(orgLen + Dm_build_801)
		var crc = dm_build_867.dm_build_815(dm_build_867.dm_build_820.dm_build_413, 0, Dm_build_711+orgLen)
		dm_build_867.dm_build_820.dm_build_413.Dm_build_141(crc)
	} else {
		dm_build_867.dm_build_820.dm_build_413.Dm_build_261(Dm_build_719, dm_build_867.dm_build_814())
	}
}

func (dm_build_869 *dm_build_819) dm_build_813() error {
	if dm_build_869.dm_build_820.dm_build_419 {

		var bodyLen = dm_build_869.dm_build_816() - Dm_build_801
		var msgLen = Dm_build_711 + bodyLen
		var recv = dm_build_869.dm_build_820.dm_build_413.Dm_build_365(int(msgLen))
		var calc = dm_build_869.dm_build_815(dm_build_869.dm_build_820.dm_build_413, 0, msgLen)
		if recv != calc {
			return ECGO_MSG_CHECK_ERROR.throw()
		}

		dm_build_869.dm_build_817(bodyLen)
		dm_build_869.dm_build_820.dm_build_413.Dm_build_92(int(msgLen))
		return nil
	} else {
		var recv = dm_build_869.dm_build_820.dm_build_413.Dm_build_341(Dm_build_719)
		var calc = dm_build_869.dm_build_814()
		if recv != calc {
			return ECGO_MSG_CHECK_ERROR.throw()
		}
		return nil
	}
}

func (dm_build_871 *dm_build_819) dm_build_814() byte {
	dm_build_872 := dm_build_871.dm_build_820.dm_build_413.Dm_build_341(0)

	for i := 1; i < Dm_build_719; i++ {
		dm_build_872 ^= dm_build_871.dm_build_820.dm_build_413.Dm_build_341(i)
	}

	return dm_build_872
}

func (dm_build_874 *dm_build_819) dm_build_815(dm_build_875 *Dm_build_78, dm_build_876 int32, dm_build_877 int32) uint32 {

	var dm_build_878 uint32 = 0xFFFFFFFF
	var dm_build_879 = dm_build_876
	var dm_build_880 = dm_build_877 - dm_build_876
	var dm_build_881, dm_build_882 uint32

	for dm_build_880 >= 8 {
		dm_build_881 = dm_build_875.Dm_build_365(int(dm_build_879)) ^ dm_build_878
		dm_build_879 += ULINT_SIZE

		dm_build_882 = dm_build_875.Dm_build_365(int(dm_build_879))
		dm_build_879 += ULINT_SIZE

		dm_build_878 = Dm_build_802[7][dm_build_881&0xFF] ^ Dm_build_802[6][(dm_build_881>>8)&0xFF] ^
			Dm_build_802[5][(dm_build_881>>16)&0xFF] ^ Dm_build_802[4][(dm_build_881>>24)&0xFF] ^
			Dm_build_802[3][dm_build_882&0xFF] ^ Dm_build_802[2][(dm_build_882>>8)&0xFF] ^
			Dm_build_802[1][(dm_build_882>>16)&0xFF] ^ Dm_build_802[0][(dm_build_882>>24)&0xFF]
		dm_build_880 -= 8
	}

	for dm_build_880 > 0 {
		dm_build_878 = ((dm_build_878 >> 8) & 0x00FFFFFF) ^ Dm_build_802[0][(dm_build_878&0xFF)^uint32(dm_build_875.Dm_build_359(int(dm_build_879)))]
		dm_build_879++
		dm_build_880--
	}
	return ^dm_build_878
}

func (dm_build_884 *dm_build_819) dm_build_816() int32 {
	return dm_build_884.dm_build_820.dm_build_413.Dm_build_347(Dm_build_715)
}

func (dm_build_886 *dm_build_819) dm_build_817(dm_build_887 int32) {
	dm_build_886.dm_build_820.dm_build_413.Dm_build_269(Dm_build_715, dm_build_887)
}

func (dm_build_889 *dm_build_819) dm_build_818() int16 {
	return dm_build_889.dm_build_821
}

type dm_build_890 struct {
	dm_build_819
}

func dm_build_891(dm_build_892 *dm_build_410) *dm_build_890 {
	dm_build_893 := new(dm_build_890)
	dm_build_893.dm_build_824(dm_build_892, Dm_build_691)
	return dm_build_893
}

type dm_build_894 struct {
	dm_build_819
	dm_build_895 string
}

func dm_build_896(dm_build_897 *dm_build_410, dm_build_898 *DmStatement, dm_build_899 string) *dm_build_894 {
	dm_build_900 := new(dm_build_894)
	dm_build_900.dm_build_828(dm_build_897, Dm_build_699, dm_build_898)
	dm_build_900.dm_build_895 = dm_build_899
	dm_build_900.dm_build_823.cursorName = dm_build_899
	return dm_build_900
}

func (dm_build_902 *dm_build_894) dm_build_805() error {
	dm_build_902.dm_build_820.dm_build_413.Dm_build_191(dm_build_902.dm_build_895, dm_build_902.dm_build_820.dm_build_414.getServerEncoding(), dm_build_902.dm_build_820.dm_build_414)
	dm_build_902.dm_build_820.dm_build_413.Dm_build_129(1)
	return nil
}

type Dm_build_903 struct {
	dm_build_919
	dm_build_904 []OptParameter
}

func dm_build_905(dm_build_906 *dm_build_410, dm_build_907 *DmStatement, dm_build_908 []OptParameter) *Dm_build_903 {
	dm_build_909 := new(Dm_build_903)
	dm_build_909.dm_build_828(dm_build_906, Dm_build_709, dm_build_907)
	dm_build_909.dm_build_904 = dm_build_908
	return dm_build_909
}

func (dm_build_911 *Dm_build_903) dm_build_805() error {
	dm_build_912 := len(dm_build_911.dm_build_904)

	dm_build_911.dm_build_933(int32(dm_build_912), 1)

	dm_build_911.dm_build_820.dm_build_413.Dm_build_191(dm_build_911.dm_build_823.nativeSql, dm_build_911.dm_build_823.dmConn.getServerEncoding(), dm_build_911.dm_build_823.dmConn)

	for _, param := range dm_build_911.dm_build_904 {
		dm_build_911.dm_build_820.dm_build_413.Dm_build_121(param.ioType)
		dm_build_911.dm_build_820.dm_build_413.Dm_build_129(int32(param.tp))
		dm_build_911.dm_build_820.dm_build_413.Dm_build_129(int32(param.prec))
		dm_build_911.dm_build_820.dm_build_413.Dm_build_129(int32(param.scale))
	}

	for _, param := range dm_build_911.dm_build_904 {
		if param.bytes == nil {
			dm_build_911.dm_build_820.dm_build_413.Dm_build_137(Dm_build_728)
		} else {
			dm_build_911.dm_build_820.dm_build_413.Dm_build_167(param.bytes[:len(param.bytes)])
		}
	}
	return nil
}

func (dm_build_914 *Dm_build_903) dm_build_809() (interface{}, error) {
	return dm_build_914.dm_build_919.dm_build_809()
}

const (
	Dm_build_915 int = 0x01

	Dm_build_916 int = 0x02

	Dm_build_917 int = 0x04

	Dm_build_918 int = 0x08
)

type dm_build_919 struct {
	dm_build_819
	dm_build_920 [][]interface{}
	dm_build_921 []parameter
	dm_build_922 bool
}

func dm_build_923(dm_build_924 *dm_build_410, dm_build_925 int16, dm_build_926 *DmStatement) *dm_build_919 {
	dm_build_927 := new(dm_build_919)
	dm_build_927.dm_build_828(dm_build_924, dm_build_925, dm_build_926)
	dm_build_927.dm_build_922 = true
	return dm_build_927
}

func dm_build_928(dm_build_929 *dm_build_410, dm_build_930 *DmStatement, dm_build_931 [][]interface{}) *dm_build_919 {
	dm_build_932 := new(dm_build_919)

	if dm_build_929.dm_build_414.Execute2 {
		dm_build_932.dm_build_828(dm_build_929, Dm_build_693, dm_build_930)
	} else {
		dm_build_932.dm_build_828(dm_build_929, Dm_build_689, dm_build_930)
	}

	dm_build_932.dm_build_921 = dm_build_930.params
	dm_build_932.dm_build_920 = dm_build_931
	dm_build_932.dm_build_922 = true
	return dm_build_932
}

func (dm_build_934 *dm_build_919) dm_build_933(dm_build_935 int32, dm_build_936 int64) {

	dm_build_937 := Dm_build_712

	if dm_build_934.dm_build_820.dm_build_414.autoCommit {
		dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_261(dm_build_937, 1)
	} else {
		dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_261(dm_build_937, 0)
	}

	dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_289(dm_build_937, uint16(dm_build_935))

	dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_261(dm_build_937, 1)

	dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_273(dm_build_937, dm_build_936)

	dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_273(dm_build_937, dm_build_934.dm_build_823.cursorUpdateRow)

	if dm_build_934.dm_build_823.maxRows <= 0 || dm_build_934.dm_build_823.dmConn.dmConnector.enRsCache {
		dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_273(dm_build_937, INT64_MAX)
	} else {
		dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_273(dm_build_937, dm_build_934.dm_build_823.maxRows)
	}

	dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_261(dm_build_937, 1)

	if dm_build_934.dm_build_820.dm_build_414.dmConnector.continueBatchOnError {
		dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_261(dm_build_937, 1)
	} else {
		dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_261(dm_build_937, 0)
	}

	dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_261(dm_build_937, 0)

	dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_261(dm_build_937, 0)

	if dm_build_934.dm_build_823.queryTimeout == 0 {
		dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_269(dm_build_937, -1)
	} else {
		dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_269(dm_build_937, dm_build_934.dm_build_823.queryTimeout)
	}

	dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_269(dm_build_937, dm_build_934.dm_build_820.dm_build_414.dmConnector.batchAllowMaxErrors)

	if dm_build_934.dm_build_823.innerExec {
		dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_261(dm_build_937, 1)
	} else {
		dm_build_937 += dm_build_934.dm_build_820.dm_build_413.Dm_build_261(dm_build_937, 0)
	}
}

func (dm_build_939 *dm_build_919) dm_build_805() error {
	var dm_build_940 int32
	var dm_build_941 int64

	if dm_build_939.dm_build_921 != nil {
		dm_build_940 = int32(len(dm_build_939.dm_build_921))
	} else {
		dm_build_940 = 0
	}

	if dm_build_939.dm_build_920 != nil {
		dm_build_941 = int64(len(dm_build_939.dm_build_920))
	} else {
		dm_build_941 = 0
	}

	dm_build_939.dm_build_933(dm_build_940, dm_build_941)

	if dm_build_940 > 0 {
		err := dm_build_939.dm_build_942(dm_build_939.dm_build_921)
		if err != nil {
			return err
		}
		if dm_build_939.dm_build_920 != nil && len(dm_build_939.dm_build_920) > 0 {
			for _, paramObject := range dm_build_939.dm_build_920 {
				if err := dm_build_939.dm_build_945(paramObject); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (dm_build_943 *dm_build_919) dm_build_942(dm_build_944 []parameter) error {
	for _, param := range dm_build_944 {

		if param.colType == CURSOR && param.ioType == IO_TYPE_OUT {
			dm_build_943.dm_build_820.dm_build_413.Dm_build_121(IO_TYPE_INOUT)
		} else {
			dm_build_943.dm_build_820.dm_build_413.Dm_build_121(param.ioType)
		}

		dm_build_943.dm_build_820.dm_build_413.Dm_build_129(param.colType)

		lprec := param.prec
		lscale := param.scale
		typeDesc := param.typeDescriptor
		switch param.colType {
		case ARRAY, SARRAY:
			tmp, err := getPackArraySize(typeDesc)
			if err != nil {
				return err
			}
			lprec = int32(tmp)
		case PLTYPE_RECORD:
			tmp, err := getPackRecordSize(typeDesc)
			if err != nil {
				return err
			}
			lprec = int32(tmp)
		case CLASS:
			tmp, err := getPackClassSize(typeDesc)
			if err != nil {
				return err
			}
			lprec = int32(tmp)
		case BLOB:
			if isComplexType(int(param.colType), int(param.scale)) {
				lprec = int32(typeDesc.getObjId())
				if lprec == 4 {
					lprec = int32(typeDesc.getOuterId())
				}
			}
		}

		dm_build_943.dm_build_820.dm_build_413.Dm_build_129(lprec)

		dm_build_943.dm_build_820.dm_build_413.Dm_build_129(lscale)

		switch param.colType {
		case ARRAY, SARRAY:
			err := packArray(typeDesc, dm_build_943.dm_build_820.dm_build_413)
			if err != nil {
				return err
			}

		case PLTYPE_RECORD:
			err := packRecord(typeDesc, dm_build_943.dm_build_820.dm_build_413)
			if err != nil {
				return err
			}

		case CLASS:
			err := packClass(typeDesc, dm_build_943.dm_build_820.dm_build_413)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func (dm_build_946 *dm_build_919) dm_build_945(dm_build_947 []interface{}) error {
	for i := 0; i < len(dm_build_946.dm_build_921); i++ {

		if dm_build_946.dm_build_921[i].colType == CURSOR {
			dm_build_946.dm_build_820.dm_build_413.Dm_build_125(ULINT_SIZE)
			dm_build_946.dm_build_820.dm_build_413.Dm_build_129(dm_build_946.dm_build_921[i].cursorStmt.id)
			continue
		}

		if dm_build_946.dm_build_921[i].ioType == IO_TYPE_OUT {
			continue
		}

		switch dm_build_947[i].(type) {
		case []byte:
			if dataBytes, ok := dm_build_947[i].([]byte); ok {
				if len(dataBytes) > Dm_build_722 {
					return ECGO_DATA_TOO_LONG.throw()
				}
				dm_build_946.dm_build_820.dm_build_413.Dm_build_167(dataBytes)
			}
		case int:
			if dm_build_947[i] == ParamDataEnum_Null {
				dm_build_946.dm_build_820.dm_build_413.Dm_build_137(Dm_build_728)
			} else if dm_build_947[i] == ParamDataEnum_OFF_ROW {
				dm_build_946.dm_build_820.dm_build_413.Dm_build_125(0)
			}
		case lobCtl:
			dm_build_946.dm_build_820.dm_build_413.Dm_build_137(uint16(Dm_build_727))
			dm_build_946.dm_build_820.dm_build_413.Dm_build_157(dm_build_947[i].(lobCtl).value)
		default:
			return fmt.Errorf("Bind param data failed by invalid param data type: ")
		}
	}

	return nil
}

func (dm_build_949 *dm_build_919) dm_build_809() (interface{}, error) {
	dm_build_950 := execRetInfo{}
	dm_build_951 := dm_build_949.dm_build_823.dmConn

	dm_build_952 := Dm_build_712

	dm_build_950.retSqlType = dm_build_949.dm_build_820.dm_build_413.Dm_build_344(dm_build_952)
	dm_build_952 += USINT_SIZE

	dm_build_953 := dm_build_949.dm_build_820.dm_build_413.Dm_build_344(dm_build_952)
	dm_build_952 += USINT_SIZE

	dm_build_950.updateCount = dm_build_949.dm_build_820.dm_build_413.Dm_build_350(dm_build_952)
	dm_build_952 += DDWORD_SIZE

	dm_build_954 := dm_build_949.dm_build_820.dm_build_413.Dm_build_362(dm_build_952)
	dm_build_952 += USINT_SIZE

	dm_build_950.rsUpdatable = dm_build_949.dm_build_820.dm_build_413.Dm_build_341(dm_build_952) != 0
	dm_build_952 += BYTE_SIZE

	dm_build_955 := dm_build_949.dm_build_820.dm_build_413.Dm_build_344(dm_build_952)
	dm_build_952 += ULINT_SIZE

	dm_build_950.printLen = dm_build_949.dm_build_820.dm_build_413.Dm_build_347(dm_build_952)
	dm_build_952 += ULINT_SIZE

	var dm_build_956 int16 = -1
	if dm_build_950.retSqlType == Dm_build_778 || dm_build_950.retSqlType == Dm_build_779 {
		dm_build_950.rowid = 0

		dm_build_950.rsBdta = dm_build_949.dm_build_820.dm_build_413.Dm_build_341(dm_build_952) == Dm_build_791
		dm_build_952 += BYTE_SIZE

		dm_build_956 = dm_build_949.dm_build_820.dm_build_413.Dm_build_344(dm_build_952)
		dm_build_952 += USINT_SIZE
		dm_build_952 += 5
	} else {
		dm_build_950.rowid = dm_build_949.dm_build_820.dm_build_413.Dm_build_350(dm_build_952)
		dm_build_952 += DDWORD_SIZE
	}

	dm_build_950.execId = dm_build_949.dm_build_820.dm_build_413.Dm_build_347(dm_build_952)
	dm_build_952 += ULINT_SIZE

	dm_build_950.rsCacheOffset = dm_build_949.dm_build_820.dm_build_413.Dm_build_347(dm_build_952)
	dm_build_952 += ULINT_SIZE

	dm_build_957 := dm_build_949.dm_build_820.dm_build_413.Dm_build_341(dm_build_952)
	dm_build_952 += BYTE_SIZE
	dm_build_958 := (dm_build_957 & 0x01) == 0x01
	dm_build_959 := (dm_build_957 & 0x02) == 0x02

	dm_build_951.TrxStatus = dm_build_949.dm_build_820.dm_build_413.Dm_build_347(dm_build_952)
	dm_build_951.setTrxFinish(dm_build_951.TrxStatus)
	dm_build_952 += ULINT_SIZE

	if dm_build_950.printLen > 0 {
		bytes := dm_build_949.dm_build_820.dm_build_413.Dm_build_224(int(dm_build_950.printLen))
		dm_build_950.printMsg = Dm_build_1298.Dm_build_1455(bytes, 0, len(bytes), dm_build_951.getServerEncoding(), dm_build_951)
	}

	if dm_build_954 > 0 {
		dm_build_950.outParamDatas = dm_build_949.dm_build_960(int(dm_build_954))
	}

	switch dm_build_950.retSqlType {
	case Dm_build_780:
		dm_build_951.dmConnector.localTimezone = dm_build_949.dm_build_820.dm_build_413.Dm_build_200()
	case Dm_build_778:
		dm_build_950.hasResultSet = true
		if dm_build_953 > 0 {
			dm_build_949.dm_build_823.columns = dm_build_949.dm_build_969(int(dm_build_953), dm_build_950.rsBdta)
		}
		dm_build_949.dm_build_979(&dm_build_950, len(dm_build_949.dm_build_823.columns), int(dm_build_955), int(dm_build_956))
	case Dm_build_779:
		if dm_build_953 > 0 || dm_build_955 > 0 {
			dm_build_950.hasResultSet = true
		}
		if dm_build_953 > 0 {
			dm_build_949.dm_build_823.columns = dm_build_949.dm_build_969(int(dm_build_953), dm_build_950.rsBdta)
		}
		dm_build_949.dm_build_979(&dm_build_950, len(dm_build_949.dm_build_823.columns), int(dm_build_955), int(dm_build_956))
	case Dm_build_781:
		dm_build_951.IsoLevel = int32(dm_build_949.dm_build_820.dm_build_413.Dm_build_200())
		dm_build_951.ReadOnly = dm_build_949.dm_build_820.dm_build_413.Dm_build_197() == 1
	case Dm_build_774:
		dm_build_951.Schema = dm_build_949.dm_build_820.dm_build_413.Dm_build_245(dm_build_951.getServerEncoding(), dm_build_951)
	case Dm_build_771:
		dm_build_950.explain = dm_build_949.dm_build_820.dm_build_413.Dm_build_245(dm_build_951.getServerEncoding(), dm_build_951)

	case Dm_build_775, Dm_build_777, Dm_build_776:
		if dm_build_958 {

			counts := dm_build_949.dm_build_820.dm_build_413.Dm_build_203()
			rowCounts := make([]int64, counts)
			for i := 0; i < int(counts); i++ {
				rowCounts[i] = dm_build_949.dm_build_820.dm_build_413.Dm_build_206()
			}
			dm_build_950.updateCounts = rowCounts
		}

		if dm_build_959 {
			rows := dm_build_949.dm_build_820.dm_build_413.Dm_build_203()

			var lastInsertId int64
			for i := 0; i < int(rows); i++ {
				lastInsertId = dm_build_949.dm_build_820.dm_build_413.Dm_build_206()
			}
			dm_build_950.lastInsertId = lastInsertId

		} else if dm_build_950.updateCount == 1 {
			dm_build_950.lastInsertId = dm_build_950.rowid
		}

		if dm_build_949.dm_build_822 == EC_BP_WITH_ERROR.ErrCode {
			dm_build_949.dm_build_985(dm_build_950.updateCounts)
		}
	case Dm_build_784:
		len := dm_build_949.dm_build_820.dm_build_413.Dm_build_215()
		dm_build_951.FormatDate = dm_build_949.dm_build_820.dm_build_413.Dm_build_240(int(len), dm_build_951.getServerEncoding(), dm_build_951)
	case Dm_build_786:

		len := dm_build_949.dm_build_820.dm_build_413.Dm_build_215()
		dm_build_951.FormatTimestamp = dm_build_949.dm_build_820.dm_build_413.Dm_build_240(int(len), dm_build_951.getServerEncoding(), dm_build_951)
	case Dm_build_787:

		len := dm_build_949.dm_build_820.dm_build_413.Dm_build_215()
		dm_build_951.FormatTimestampTZ = dm_build_949.dm_build_820.dm_build_413.Dm_build_240(int(len), dm_build_951.getServerEncoding(), dm_build_951)
	case Dm_build_785:
		len := dm_build_949.dm_build_820.dm_build_413.Dm_build_215()
		dm_build_951.FormatTime = dm_build_949.dm_build_820.dm_build_413.Dm_build_240(int(len), dm_build_951.getServerEncoding(), dm_build_951)
	case Dm_build_788:
		len := dm_build_949.dm_build_820.dm_build_413.Dm_build_215()
		dm_build_951.FormatTimeTZ = dm_build_949.dm_build_820.dm_build_413.Dm_build_240(int(len), dm_build_951.getServerEncoding(), dm_build_951)
	case Dm_build_789:
		dm_build_951.OracleDateLanguage = dm_build_949.dm_build_820.dm_build_413.Dm_build_215()
	}

	return &dm_build_950, nil
}

func (dm_build_961 *dm_build_919) dm_build_960(dm_build_962 int) [][]byte {
	dm_build_963 := make([]int, dm_build_962)

	dm_build_964 := 0
	for i := 0; i < len(dm_build_961.dm_build_921); i++ {
		if dm_build_961.dm_build_921[i].ioType == IO_TYPE_INOUT || dm_build_961.dm_build_921[i].ioType == IO_TYPE_OUT {
			dm_build_963[dm_build_964] = i
			dm_build_964++
		}
	}

	dm_build_965 := make([][]byte, len(dm_build_961.dm_build_921))
	var dm_build_966 int32
	var dm_build_967 bool
	var dm_build_968 []byte = nil
	for i := 0; i < dm_build_962; i++ {
		dm_build_967 = false
		dm_build_966 = int32(dm_build_961.dm_build_820.dm_build_413.Dm_build_218())

		if dm_build_966 == int32(Dm_build_728) {
			dm_build_966 = 0
			dm_build_967 = true
		} else if dm_build_966 == int32(Dm_build_729) {
			dm_build_966 = dm_build_961.dm_build_820.dm_build_413.Dm_build_203()
		}

		if dm_build_967 {
			dm_build_965[dm_build_963[i]] = nil
		} else {
			dm_build_968 = dm_build_961.dm_build_820.dm_build_413.Dm_build_224(int(dm_build_966))
			dm_build_965[dm_build_963[i]] = dm_build_968
		}
	}

	return dm_build_965
}

func (dm_build_970 *dm_build_919) dm_build_969(dm_build_971 int, dm_build_972 bool) []column {
	dm_build_973 := dm_build_970.dm_build_820.dm_build_414.getServerEncoding()
	var dm_build_974, dm_build_975, dm_build_976, dm_build_977 int16
	dm_build_978 := make([]column, dm_build_971)
	for i := 0; i < dm_build_971; i++ {
		dm_build_978[i].InitColumn()

		dm_build_978[i].colType = dm_build_970.dm_build_820.dm_build_413.Dm_build_203()

		dm_build_978[i].prec = dm_build_970.dm_build_820.dm_build_413.Dm_build_203()

		dm_build_978[i].scale = dm_build_970.dm_build_820.dm_build_413.Dm_build_203()

		dm_build_978[i].nullable = dm_build_970.dm_build_820.dm_build_413.Dm_build_203() != 0

		itemFlag := dm_build_970.dm_build_820.dm_build_413.Dm_build_200()
		dm_build_978[i].lob = int(itemFlag)&Dm_build_916 != 0
		dm_build_978[i].identity = int(itemFlag)&Dm_build_915 != 0
		dm_build_978[i].readonly = int(itemFlag)&Dm_build_917 != 0

		dm_build_970.dm_build_820.dm_build_413.Dm_build_103(4, false, true)

		dm_build_970.dm_build_820.dm_build_413.Dm_build_103(2, false, true)

		dm_build_974 = dm_build_970.dm_build_820.dm_build_413.Dm_build_200()

		dm_build_975 = dm_build_970.dm_build_820.dm_build_413.Dm_build_200()

		dm_build_976 = dm_build_970.dm_build_820.dm_build_413.Dm_build_200()

		dm_build_977 = dm_build_970.dm_build_820.dm_build_413.Dm_build_200()
		dm_build_978[i].name = dm_build_970.dm_build_820.dm_build_413.Dm_build_240(int(dm_build_974), dm_build_973, dm_build_970.dm_build_820.dm_build_414)
		dm_build_978[i].typeName = dm_build_970.dm_build_820.dm_build_413.Dm_build_240(int(dm_build_975), dm_build_973, dm_build_970.dm_build_820.dm_build_414)
		dm_build_978[i].tableName = dm_build_970.dm_build_820.dm_build_413.Dm_build_240(int(dm_build_976), dm_build_973, dm_build_970.dm_build_820.dm_build_414)
		dm_build_978[i].schemaName = dm_build_970.dm_build_820.dm_build_413.Dm_build_240(int(dm_build_977), dm_build_973, dm_build_970.dm_build_820.dm_build_414)

		if dm_build_970.dm_build_823.readBaseColName {
			dm_build_978[i].baseName = dm_build_970.dm_build_820.dm_build_413.Dm_build_253(dm_build_973, dm_build_970.dm_build_820.dm_build_414)
		}

		if dm_build_978[i].lob {
			dm_build_978[i].lobTabId = dm_build_970.dm_build_820.dm_build_413.Dm_build_203()
			dm_build_978[i].lobColId = dm_build_970.dm_build_820.dm_build_413.Dm_build_200()
		}

	}

	for i := 0; i < dm_build_971; i++ {

		if isComplexType(int(dm_build_978[i].colType), int(dm_build_978[i].scale)) {
			strDesc := newTypeDescriptor(dm_build_970.dm_build_820.dm_build_414)
			strDesc.unpack(dm_build_970.dm_build_820.dm_build_413)
			dm_build_978[i].typeDescriptor = strDesc
		}
	}

	return dm_build_978
}

func (dm_build_980 *dm_build_919) dm_build_979(dm_build_981 *execRetInfo, dm_build_982 int, dm_build_983 int, dm_build_984 int) {
	if dm_build_983 > 0 {
		startOffset := dm_build_980.dm_build_820.dm_build_413.Dm_build_98()
		if dm_build_981.rsBdta {
			dm_build_981.rsDatas = dm_build_980.dm_build_998(dm_build_980.dm_build_823.columns, dm_build_984)
		} else {
			datas := make([][][]byte, dm_build_983)

			for i := 0; i < dm_build_983; i++ {

				datas[i] = make([][]byte, dm_build_982+1)

				dm_build_980.dm_build_820.dm_build_413.Dm_build_103(2, false, true)

				datas[i][0] = dm_build_980.dm_build_820.dm_build_413.Dm_build_224(LINT64_SIZE)

				dm_build_980.dm_build_820.dm_build_413.Dm_build_103(2*dm_build_982, false, true)

				for j := 1; j < dm_build_982+1; j++ {

					colLen := dm_build_980.dm_build_820.dm_build_413.Dm_build_218()
					if colLen == Dm_build_732 {
						datas[i][j] = nil
					} else if colLen != Dm_build_733 {
						datas[i][j] = dm_build_980.dm_build_820.dm_build_413.Dm_build_224(int(colLen))
					} else {
						datas[i][j] = dm_build_980.dm_build_820.dm_build_413.Dm_build_228()
					}
				}
			}

			dm_build_981.rsDatas = datas
		}
		dm_build_981.rsSizeof = dm_build_980.dm_build_820.dm_build_413.Dm_build_98() - startOffset
	}

	if dm_build_981.rsCacheOffset > 0 {
		tbCount := dm_build_980.dm_build_820.dm_build_413.Dm_build_200()

		ids := make([]int32, tbCount)
		tss := make([]int64, tbCount)

		for i := 0; i < int(tbCount); i++ {
			ids[i] = dm_build_980.dm_build_820.dm_build_413.Dm_build_203()
			tss[i] = dm_build_980.dm_build_820.dm_build_413.Dm_build_206()
		}

		dm_build_981.tbIds = ids[:]
		dm_build_981.tbTss = tss[:]
	}
}

func (dm_build_986 *dm_build_919) dm_build_985(dm_build_987 []int64) error {

	dm_build_986.dm_build_820.dm_build_413.Dm_build_103(4, false, true)

	dm_build_988 := dm_build_986.dm_build_820.dm_build_413.Dm_build_203()

	dm_build_989 := make([]string, 0, 8)
	for i := 0; i < int(dm_build_988); i++ {
		irow := dm_build_986.dm_build_820.dm_build_413.Dm_build_203()

		dm_build_987[irow] = -3

		code := dm_build_986.dm_build_820.dm_build_413.Dm_build_203()

		errStr := dm_build_986.dm_build_820.dm_build_413.Dm_build_253(dm_build_986.dm_build_820.dm_build_414.getServerEncoding(), dm_build_986.dm_build_820.dm_build_414)

		dm_build_989 = append(dm_build_989, "row["+strconv.Itoa(int(irow))+"]:"+strconv.Itoa(int(code))+", "+errStr)
	}

	if len(dm_build_989) > 0 {
		builder := &strings.Builder{}
		for _, str := range dm_build_989 {
			builder.WriteString(util.LINE_SEPARATOR)
			builder.WriteString(str)
		}
		EC_BP_WITH_ERROR.ErrText += builder.String()
		return EC_BP_WITH_ERROR.throw()
	}
	return nil
}

const (
	Dm_build_990 = 0

	Dm_build_991 = Dm_build_990 + ULINT_SIZE

	Dm_build_992 = Dm_build_991 + USINT_SIZE

	Dm_build_993 = Dm_build_992 + ULINT_SIZE

	Dm_build_994 = Dm_build_993 + ULINT_SIZE

	Dm_build_995 = Dm_build_994 + BYTE_SIZE

	Dm_build_996 = -2

	Dm_build_997 = -3
)

func (dm_build_999 *dm_build_919) dm_build_998(dm_build_1000 []column, dm_build_1001 int) [][][]byte {

	dm_build_1002 := dm_build_999.dm_build_820.dm_build_413.Dm_build_221()

	dm_build_1003 := dm_build_999.dm_build_820.dm_build_413.Dm_build_218()

	var dm_build_1004 bool

	if dm_build_1001 >= 0 && int(dm_build_1003) == len(dm_build_1000)+1 {
		dm_build_1004 = true
	} else {
		dm_build_1004 = false
	}

	dm_build_999.dm_build_820.dm_build_413.Dm_build_103(ULINT_SIZE, false, true)

	dm_build_999.dm_build_820.dm_build_413.Dm_build_103(ULINT_SIZE, false, true)

	dm_build_999.dm_build_820.dm_build_413.Dm_build_103(BYTE_SIZE, false, true)

	dm_build_1005 := make([]uint16, dm_build_1003)
	for icol := 0; icol < int(dm_build_1003); icol++ {
		dm_build_1005[icol] = dm_build_999.dm_build_820.dm_build_413.Dm_build_218()
	}

	dm_build_1006 := make([]uint32, dm_build_1003)
	dm_build_1007 := make([][][]byte, dm_build_1002)

	for i := uint32(0); i < dm_build_1002; i++ {
		dm_build_1007[i] = make([][]byte, len(dm_build_1000)+1)
	}

	for icol := 0; icol < int(dm_build_1003); icol++ {
		dm_build_1006[icol] = dm_build_999.dm_build_820.dm_build_413.Dm_build_221()
	}

	for icol := 0; icol < int(dm_build_1003); icol++ {

		dataCol := icol + 1
		if dm_build_1004 && icol == dm_build_1001 {
			dataCol = 0
		} else if dm_build_1004 && icol > dm_build_1001 {
			dataCol = icol
		}

		allNotNull := dm_build_999.dm_build_820.dm_build_413.Dm_build_203() == 1
		var isNull []bool = nil
		if !allNotNull {
			isNull = make([]bool, dm_build_1002)
			for irow := uint32(0); irow < dm_build_1002; irow++ {
				isNull[irow] = dm_build_999.dm_build_820.dm_build_413.Dm_build_197() == 0
			}
		}

		for irow := uint32(0); irow < dm_build_1002; irow++ {
			if allNotNull || !isNull[irow] {
				dm_build_1007[irow][dataCol] = dm_build_999.dm_build_1008(int(dm_build_1005[icol]))
			}
		}
	}

	if !dm_build_1004 && dm_build_1001 >= 0 {
		for irow := uint32(0); irow < dm_build_1002; irow++ {
			dm_build_1007[irow][0] = dm_build_1007[irow][dm_build_1001+1]
		}
	}

	return dm_build_1007
}

func (dm_build_1009 *dm_build_919) dm_build_1008(dm_build_1010 int) []byte {

	dm_build_1011 := dm_build_1009.dm_build_1014(dm_build_1010)

	dm_build_1012 := int32(0)
	if dm_build_1011 == Dm_build_996 {
		dm_build_1012 = dm_build_1009.dm_build_820.dm_build_413.Dm_build_203()
		dm_build_1011 = int(dm_build_1009.dm_build_820.dm_build_413.Dm_build_203())
	} else if dm_build_1011 == Dm_build_997 {
		dm_build_1011 = int(dm_build_1009.dm_build_820.dm_build_413.Dm_build_203())
	}

	dm_build_1013 := dm_build_1009.dm_build_820.dm_build_413.Dm_build_224(dm_build_1011 + int(dm_build_1012))
	if dm_build_1012 == 0 {
		return dm_build_1013
	}

	for i := dm_build_1011; i < len(dm_build_1013); i++ {
		dm_build_1013[i] = ' '
	}
	return dm_build_1013
}

func (dm_build_1015 *dm_build_919) dm_build_1014(dm_build_1016 int) int {

	dm_build_1017 := 0
	switch dm_build_1016 {
	case INT:
	case BIT:
	case TINYINT:
	case SMALLINT:
	case BOOLEAN:
	case NULL:
		dm_build_1017 = 4

	case BIGINT:

		dm_build_1017 = 8

	case CHAR:
	case VARCHAR2:
	case VARCHAR:
	case BINARY:
	case VARBINARY:
	case BLOB:
	case CLOB:
		dm_build_1017 = Dm_build_996

	case DECIMAL:
		dm_build_1017 = Dm_build_997

	case REAL:
		dm_build_1017 = 4

	case DOUBLE:
		dm_build_1017 = 8

	case DATE:
	case TIME:
	case DATETIME:
	case TIME_TZ:
	case DATETIME_TZ:
		dm_build_1017 = 12

	case INTERVAL_YM:
		dm_build_1017 = 12

	case INTERVAL_DT:
		dm_build_1017 = 24

	default:
		dm_build_1017 = 0
	}
	return dm_build_1017
}

const (
	Dm_build_1018 = Dm_build_712

	Dm_build_1019 = Dm_build_1018 + DDWORD_SIZE

	Dm_build_1020 = Dm_build_1019 + LINT64_SIZE

	Dm_build_1021 = Dm_build_1020 + USINT_SIZE

	Dm_build_1022 = Dm_build_712

	Dm_build_1023 = Dm_build_1022 + DDWORD_SIZE
)

type dm_build_1024 struct {
	dm_build_919
	dm_build_1025 *innerRows
	dm_build_1026 int64
	dm_build_1027 int64
}

func dm_build_1028(dm_build_1029 *dm_build_410, dm_build_1030 *innerRows, dm_build_1031 int64, dm_build_1032 int64) *dm_build_1024 {
	dm_build_1033 := new(dm_build_1024)
	dm_build_1033.dm_build_828(dm_build_1029, Dm_build_690, dm_build_1030.dmStmt)
	dm_build_1033.dm_build_1025 = dm_build_1030
	dm_build_1033.dm_build_1026 = dm_build_1031
	dm_build_1033.dm_build_1027 = dm_build_1032
	return dm_build_1033
}

func (dm_build_1035 *dm_build_1024) dm_build_805() error {

	dm_build_1035.dm_build_820.dm_build_413.Dm_build_273(Dm_build_1018, dm_build_1035.dm_build_1026)

	dm_build_1035.dm_build_820.dm_build_413.Dm_build_273(Dm_build_1019, dm_build_1035.dm_build_1027)

	dm_build_1035.dm_build_820.dm_build_413.Dm_build_265(Dm_build_1020, dm_build_1035.dm_build_1025.id)

	dm_build_1036 := dm_build_1035.dm_build_1025.dmStmt.dmConn.dmConnector.bufPrefetch
	var dm_build_1037 int32
	if dm_build_1035.dm_build_1025.sizeOfRow != 0 && dm_build_1035.dm_build_1025.fetchSize != 0 {
		if dm_build_1035.dm_build_1025.sizeOfRow*dm_build_1035.dm_build_1025.fetchSize > int(INT32_MAX) {
			dm_build_1037 = INT32_MAX
		} else {
			dm_build_1037 = int32(dm_build_1035.dm_build_1025.sizeOfRow * dm_build_1035.dm_build_1025.fetchSize)
		}

		if dm_build_1037 < Dm_build_744 {
			dm_build_1036 = int(Dm_build_744)
		} else if dm_build_1037 > Dm_build_745 {
			dm_build_1036 = int(Dm_build_745)
		} else {
			dm_build_1036 = int(dm_build_1037)
		}

		dm_build_1035.dm_build_820.dm_build_413.Dm_build_269(Dm_build_1021, int32(dm_build_1036))
	}

	return nil
}

func (dm_build_1039 *dm_build_1024) dm_build_809() (interface{}, error) {
	dm_build_1040 := execRetInfo{}
	dm_build_1040.rsBdta = dm_build_1039.dm_build_1025.isBdta

	dm_build_1040.updateCount = dm_build_1039.dm_build_820.dm_build_413.Dm_build_350(Dm_build_1022)
	dm_build_1041 := dm_build_1039.dm_build_820.dm_build_413.Dm_build_347(Dm_build_1023)

	dm_build_1039.dm_build_979(&dm_build_1040, len(dm_build_1039.dm_build_1025.columns), int(dm_build_1041), -1)

	return &dm_build_1040, nil
}

type dm_build_1042 struct {
	dm_build_819
	dm_build_1043 *lob
	dm_build_1044 int
	dm_build_1045 int
}

func dm_build_1046(dm_build_1047 *dm_build_410, dm_build_1048 *lob, dm_build_1049 int, dm_build_1050 int) *dm_build_1042 {
	dm_build_1051 := new(dm_build_1042)
	dm_build_1051.dm_build_824(dm_build_1047, Dm_build_703)
	dm_build_1051.dm_build_1043 = dm_build_1048
	dm_build_1051.dm_build_1044 = dm_build_1049
	dm_build_1051.dm_build_1045 = dm_build_1050
	return dm_build_1051
}

func (dm_build_1053 *dm_build_1042) dm_build_805() error {

	dm_build_1053.dm_build_820.dm_build_413.Dm_build_121(byte(dm_build_1053.dm_build_1043.lobFlag))

	dm_build_1053.dm_build_820.dm_build_413.Dm_build_129(dm_build_1053.dm_build_1043.tabId)

	dm_build_1053.dm_build_820.dm_build_413.Dm_build_125(dm_build_1053.dm_build_1043.colId)

	dm_build_1053.dm_build_820.dm_build_413.Dm_build_145(uint64(dm_build_1053.dm_build_1043.blobId))

	dm_build_1053.dm_build_820.dm_build_413.Dm_build_125(dm_build_1053.dm_build_1043.groupId)

	dm_build_1053.dm_build_820.dm_build_413.Dm_build_125(dm_build_1053.dm_build_1043.fileId)

	dm_build_1053.dm_build_820.dm_build_413.Dm_build_129(dm_build_1053.dm_build_1043.pageNo)

	dm_build_1053.dm_build_820.dm_build_413.Dm_build_125(dm_build_1053.dm_build_1043.curFileId)

	dm_build_1053.dm_build_820.dm_build_413.Dm_build_129(dm_build_1053.dm_build_1043.curPageNo)

	dm_build_1053.dm_build_820.dm_build_413.Dm_build_129(dm_build_1053.dm_build_1043.totalOffset)

	dm_build_1053.dm_build_820.dm_build_413.Dm_build_129(int32(dm_build_1053.dm_build_1044))

	dm_build_1053.dm_build_820.dm_build_413.Dm_build_129(int32(dm_build_1053.dm_build_1045))

	if dm_build_1053.dm_build_820.dm_build_414.NewLobFlag {
		dm_build_1053.dm_build_820.dm_build_413.Dm_build_145(uint64(dm_build_1053.dm_build_1043.rowId))
		dm_build_1053.dm_build_820.dm_build_413.Dm_build_125(dm_build_1053.dm_build_1043.exGroupId)
		dm_build_1053.dm_build_820.dm_build_413.Dm_build_125(dm_build_1053.dm_build_1043.exFileId)
		dm_build_1053.dm_build_820.dm_build_413.Dm_build_129(dm_build_1053.dm_build_1043.exPageNo)
	}

	return nil
}

func (dm_build_1055 *dm_build_1042) dm_build_809() (interface{}, error) {

	dm_build_1055.dm_build_1043.readOver = dm_build_1055.dm_build_820.dm_build_413.Dm_build_197() == 1
	var dm_build_1056 = dm_build_1055.dm_build_820.dm_build_413.Dm_build_221()
	if dm_build_1056 <= 0 {
		return []byte(nil), nil
	}
	dm_build_1055.dm_build_1043.curFileId = dm_build_1055.dm_build_820.dm_build_413.Dm_build_200()
	dm_build_1055.dm_build_1043.curPageNo = dm_build_1055.dm_build_820.dm_build_413.Dm_build_203()
	dm_build_1055.dm_build_1043.totalOffset = dm_build_1055.dm_build_820.dm_build_413.Dm_build_203()

	return dm_build_1055.dm_build_820.dm_build_413.Dm_build_234(int(dm_build_1056)), nil
}

type dm_build_1057 struct {
	dm_build_819
	dm_build_1058 *lob
}

func dm_build_1059(dm_build_1060 *dm_build_410, dm_build_1061 *lob) *dm_build_1057 {
	dm_build_1062 := new(dm_build_1057)
	dm_build_1062.dm_build_824(dm_build_1060, Dm_build_700)
	dm_build_1062.dm_build_1058 = dm_build_1061
	return dm_build_1062
}

func (dm_build_1064 *dm_build_1057) dm_build_805() error {

	dm_build_1064.dm_build_820.dm_build_413.Dm_build_121(byte(dm_build_1064.dm_build_1058.lobFlag))

	dm_build_1064.dm_build_820.dm_build_413.Dm_build_145(uint64(dm_build_1064.dm_build_1058.blobId))

	dm_build_1064.dm_build_820.dm_build_413.Dm_build_125(dm_build_1064.dm_build_1058.groupId)

	dm_build_1064.dm_build_820.dm_build_413.Dm_build_125(dm_build_1064.dm_build_1058.fileId)

	dm_build_1064.dm_build_820.dm_build_413.Dm_build_129(dm_build_1064.dm_build_1058.pageNo)

	if dm_build_1064.dm_build_820.dm_build_414.NewLobFlag {
		dm_build_1064.dm_build_820.dm_build_413.Dm_build_129(dm_build_1064.dm_build_1058.tabId)
		dm_build_1064.dm_build_820.dm_build_413.Dm_build_125(dm_build_1064.dm_build_1058.colId)
		dm_build_1064.dm_build_820.dm_build_413.Dm_build_145(uint64(dm_build_1064.dm_build_1058.rowId))

		dm_build_1064.dm_build_820.dm_build_413.Dm_build_125(dm_build_1064.dm_build_1058.exGroupId)
		dm_build_1064.dm_build_820.dm_build_413.Dm_build_125(dm_build_1064.dm_build_1058.exFileId)
		dm_build_1064.dm_build_820.dm_build_413.Dm_build_129(dm_build_1064.dm_build_1058.exPageNo)
	}

	return nil
}

func (dm_build_1066 *dm_build_1057) dm_build_809() (interface{}, error) {

	if dm_build_1066.dm_build_820.dm_build_413.Dm_build_100(false) == 8 {
		return dm_build_1066.dm_build_820.dm_build_413.Dm_build_206(), nil
	} else {
		return int64(dm_build_1066.dm_build_820.dm_build_413.Dm_build_221()), nil
	}
}

type dm_build_1067 struct {
	dm_build_819
	dm_build_1068 *lob
	dm_build_1069 int
}

func dm_build_1070(dm_build_1071 *dm_build_410, dm_build_1072 *lob, dm_build_1073 int) *dm_build_1067 {
	dm_build_1074 := new(dm_build_1067)
	dm_build_1074.dm_build_824(dm_build_1071, Dm_build_702)
	dm_build_1074.dm_build_1068 = dm_build_1072
	dm_build_1074.dm_build_1069 = dm_build_1073
	return dm_build_1074
}

func (dm_build_1076 *dm_build_1067) dm_build_805() error {

	dm_build_1076.dm_build_820.dm_build_413.Dm_build_121(byte(dm_build_1076.dm_build_1068.lobFlag))

	dm_build_1076.dm_build_820.dm_build_413.Dm_build_145(uint64(dm_build_1076.dm_build_1068.blobId))

	dm_build_1076.dm_build_820.dm_build_413.Dm_build_125(dm_build_1076.dm_build_1068.groupId)

	dm_build_1076.dm_build_820.dm_build_413.Dm_build_125(dm_build_1076.dm_build_1068.fileId)

	dm_build_1076.dm_build_820.dm_build_413.Dm_build_129(dm_build_1076.dm_build_1068.pageNo)

	dm_build_1076.dm_build_820.dm_build_413.Dm_build_129(dm_build_1076.dm_build_1068.tabId)
	dm_build_1076.dm_build_820.dm_build_413.Dm_build_125(dm_build_1076.dm_build_1068.colId)
	dm_build_1076.dm_build_820.dm_build_413.Dm_build_145(uint64(dm_build_1076.dm_build_1068.rowId))
	dm_build_1076.dm_build_820.dm_build_413.Dm_build_157(Dm_build_1298.Dm_build_1500(uint32(dm_build_1076.dm_build_1069)))

	if dm_build_1076.dm_build_820.dm_build_414.NewLobFlag {
		dm_build_1076.dm_build_820.dm_build_413.Dm_build_125(dm_build_1076.dm_build_1068.exGroupId)
		dm_build_1076.dm_build_820.dm_build_413.Dm_build_125(dm_build_1076.dm_build_1068.exFileId)
		dm_build_1076.dm_build_820.dm_build_413.Dm_build_129(dm_build_1076.dm_build_1068.exPageNo)
	}
	return nil
}

func (dm_build_1078 *dm_build_1067) dm_build_809() (interface{}, error) {

	dm_build_1079 := dm_build_1078.dm_build_820.dm_build_413.Dm_build_221()
	dm_build_1078.dm_build_1068.blobId = dm_build_1078.dm_build_820.dm_build_413.Dm_build_206()
	dm_build_1078.dm_build_1068.resetCurrentInfo()
	return int64(dm_build_1079), nil
}

const (
	Dm_build_1080 = Dm_build_712

	Dm_build_1081 = Dm_build_1080 + ULINT_SIZE

	Dm_build_1082 = Dm_build_1081 + ULINT_SIZE

	Dm_build_1083 = Dm_build_1082 + ULINT_SIZE

	Dm_build_1084 = Dm_build_1083 + BYTE_SIZE

	Dm_build_1085 = Dm_build_1084 + USINT_SIZE

	Dm_build_1086 = Dm_build_1085 + ULINT_SIZE

	Dm_build_1087 = Dm_build_1086 + BYTE_SIZE

	Dm_build_1088 = Dm_build_1087 + BYTE_SIZE

	Dm_build_1089 = Dm_build_1088 + BYTE_SIZE

	Dm_build_1090 = Dm_build_712

	Dm_build_1091 = Dm_build_1090 + ULINT_SIZE

	Dm_build_1092 = Dm_build_1091 + ULINT_SIZE

	Dm_build_1093 = Dm_build_1092 + BYTE_SIZE

	Dm_build_1094 = Dm_build_1093 + ULINT_SIZE

	Dm_build_1095 = Dm_build_1094 + BYTE_SIZE

	Dm_build_1096 = Dm_build_1095 + BYTE_SIZE

	Dm_build_1097 = Dm_build_1096 + USINT_SIZE

	Dm_build_1098 = Dm_build_1097 + USINT_SIZE

	Dm_build_1099 = Dm_build_1098 + BYTE_SIZE

	Dm_build_1100 = Dm_build_1099 + USINT_SIZE

	Dm_build_1101 = Dm_build_1100 + BYTE_SIZE

	Dm_build_1102 = Dm_build_1101 + BYTE_SIZE

	Dm_build_1103 = Dm_build_1102 + ULINT_SIZE

	Dm_build_1104 = Dm_build_1103 + USINT_SIZE
)

type dm_build_1105 struct {
	dm_build_819

	dm_build_1106 *DmConnection

	dm_build_1107 bool
}

func dm_build_1108(dm_build_1109 *dm_build_410) *dm_build_1105 {
	dm_build_1110 := new(dm_build_1105)
	dm_build_1110.dm_build_824(dm_build_1109, Dm_build_684)
	dm_build_1110.dm_build_1106 = dm_build_1109.dm_build_414
	return dm_build_1110
}

func (dm_build_1112 *dm_build_1105) dm_build_805() error {

	if dm_build_1112.dm_build_1106.dmConnector.newClientType {
		dm_build_1112.dm_build_820.dm_build_413.Dm_build_269(Dm_build_1080, Dm_build_724)
	} else {
		dm_build_1112.dm_build_820.dm_build_413.Dm_build_269(Dm_build_1080, Dm_build_723)
	}

	dm_build_1112.dm_build_820.dm_build_413.Dm_build_269(Dm_build_1081, g2dbIsoLevel(dm_build_1112.dm_build_1106.IsoLevel))
	dm_build_1112.dm_build_820.dm_build_413.Dm_build_269(Dm_build_1082, int32(Locale))
	dm_build_1112.dm_build_820.dm_build_413.Dm_build_265(Dm_build_1084, dm_build_1112.dm_build_1106.dmConnector.localTimezone)

	if dm_build_1112.dm_build_1106.ReadOnly {
		dm_build_1112.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1083, Dm_build_747)
	} else {
		dm_build_1112.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1083, Dm_build_746)
	}

	dm_build_1112.dm_build_820.dm_build_413.Dm_build_269(Dm_build_1085, int32(dm_build_1112.dm_build_1106.dmConnector.sessionTimeout))

	if dm_build_1112.dm_build_1106.dmConnector.mppLocal {
		dm_build_1112.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1086, 1)
	} else {
		dm_build_1112.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1086, 0)
	}

	if dm_build_1112.dm_build_1106.dmConnector.rwSeparate {
		dm_build_1112.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1087, 1)
	} else {
		dm_build_1112.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1087, 0)
	}

	if dm_build_1112.dm_build_1106.NewLobFlag {
		dm_build_1112.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1088, 1)
	} else {
		dm_build_1112.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1088, 0)
	}

	dm_build_1112.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1089, dm_build_1112.dm_build_1106.dmConnector.osAuthType)

	dm_build_1113 := dm_build_1112.dm_build_1106.getServerEncoding()

	if dm_build_1112.dm_build_820.dm_build_420 != "" {

	}

	dm_build_1114 := Dm_build_1298.Dm_build_1511(dm_build_1112.dm_build_1106.dmConnector.user, dm_build_1113, dm_build_1112.dm_build_820.dm_build_414)
	dm_build_1115 := Dm_build_1298.Dm_build_1511(dm_build_1112.dm_build_1106.dmConnector.password, dm_build_1113, dm_build_1112.dm_build_820.dm_build_414)
	if len(dm_build_1114) > Dm_build_720 {
		return ECGO_USERNAME_TOO_LONG.throw()
	}
	if len(dm_build_1115) > Dm_build_720 {
		return ECGO_PASSWORD_TOO_LONG.throw()
	}

	if dm_build_1112.dm_build_820.dm_build_416 && dm_build_1112.dm_build_1106.dmConnector.loginCertificate != "" {

	} else if dm_build_1112.dm_build_820.dm_build_416 {
		dm_build_1114 = dm_build_1112.dm_build_820.dm_build_415.Encrypt(dm_build_1114, false)
		dm_build_1115 = dm_build_1112.dm_build_820.dm_build_415.Encrypt(dm_build_1115, false)
	}

	dm_build_1112.dm_build_820.dm_build_413.Dm_build_161(dm_build_1114)
	dm_build_1112.dm_build_820.dm_build_413.Dm_build_161(dm_build_1115)

	dm_build_1112.dm_build_820.dm_build_413.Dm_build_173(dm_build_1112.dm_build_1106.dmConnector.appName, dm_build_1113, dm_build_1112.dm_build_820.dm_build_414)
	dm_build_1112.dm_build_820.dm_build_413.Dm_build_173(dm_build_1112.dm_build_1106.dmConnector.osName, dm_build_1113, dm_build_1112.dm_build_820.dm_build_414)

	if hostName, err := os.Hostname(); err != nil {
		dm_build_1112.dm_build_820.dm_build_413.Dm_build_173(hostName, dm_build_1113, dm_build_1112.dm_build_820.dm_build_414)
	} else {
		dm_build_1112.dm_build_820.dm_build_413.Dm_build_173("", dm_build_1113, dm_build_1112.dm_build_820.dm_build_414)
	}

	if dm_build_1112.dm_build_1106.dmConnector.rwStandby {
		dm_build_1112.dm_build_820.dm_build_413.Dm_build_121(1)
	} else {
		dm_build_1112.dm_build_820.dm_build_413.Dm_build_121(0)
	}

	return nil
}

func (dm_build_1117 *dm_build_1105) dm_build_809() (interface{}, error) {

	dm_build_1117.dm_build_1106.MaxRowSize = dm_build_1117.dm_build_820.dm_build_413.Dm_build_347(Dm_build_1090)
	dm_build_1117.dm_build_1106.DDLAutoCommit = dm_build_1117.dm_build_820.dm_build_413.Dm_build_341(Dm_build_1092) == 1
	dm_build_1117.dm_build_1106.IsoLevel = dm_build_1117.dm_build_820.dm_build_413.Dm_build_347(Dm_build_1093)
	dm_build_1117.dm_build_1106.dmConnector.caseSensitive = dm_build_1117.dm_build_820.dm_build_413.Dm_build_341(Dm_build_1094) == 1
	dm_build_1117.dm_build_1106.BackslashEscape = dm_build_1117.dm_build_820.dm_build_413.Dm_build_341(Dm_build_1095) == 1
	dm_build_1117.dm_build_1106.SvrStat = int32(dm_build_1117.dm_build_820.dm_build_413.Dm_build_344(Dm_build_1097))
	dm_build_1117.dm_build_1106.SvrMode = int32(dm_build_1117.dm_build_820.dm_build_413.Dm_build_344(Dm_build_1096))
	dm_build_1117.dm_build_1106.ConstParaOpt = dm_build_1117.dm_build_820.dm_build_413.Dm_build_341(Dm_build_1098) == 1
	dm_build_1117.dm_build_1106.DbTimezone = dm_build_1117.dm_build_820.dm_build_413.Dm_build_344(Dm_build_1099)
	dm_build_1117.dm_build_1106.NewLobFlag = dm_build_1117.dm_build_820.dm_build_413.Dm_build_341(Dm_build_1101) == 1

	if dm_build_1117.dm_build_1106.dmConnector.bufPrefetch == 0 {
		dm_build_1117.dm_build_1106.dmConnector.bufPrefetch = int(dm_build_1117.dm_build_820.dm_build_413.Dm_build_347(Dm_build_1102))
	}

	dm_build_1117.dm_build_1106.LifeTimeRemainder = dm_build_1117.dm_build_820.dm_build_413.Dm_build_344(Dm_build_1103)
	dm_build_1117.dm_build_1106.dscControl = dm_build_1117.dm_build_820.dm_build_413.Dm_build_341(Dm_build_1104) == 1

	dm_build_1118 := dm_build_1117.dm_build_1106.getServerEncoding()

	dm_build_1117.dm_build_1106.InstanceName = dm_build_1117.dm_build_820.dm_build_413.Dm_build_245(dm_build_1118, dm_build_1117.dm_build_820.dm_build_414)

	var dm_build_1119 = dm_build_1117.dm_build_820.dm_build_413.Dm_build_203()
	if dm_build_1119 == 0 && dm_build_1117.dm_build_1106.MsgVersion > 0 {
		dm_build_1117.dm_build_1106.Schema = strings.ToUpper(dm_build_1117.dm_build_1106.dmConnector.user)
	} else {
		dm_build_1117.dm_build_1106.Schema = dm_build_1117.dm_build_820.dm_build_413.Dm_build_240(int(dm_build_1119), dm_build_1118, dm_build_1117.dm_build_820.dm_build_414)
	}

	dm_build_1117.dm_build_1106.LastLoginIP = dm_build_1117.dm_build_820.dm_build_413.Dm_build_245(dm_build_1118, dm_build_1117.dm_build_820.dm_build_414)
	dm_build_1117.dm_build_1106.LastLoginTime = dm_build_1117.dm_build_820.dm_build_413.Dm_build_245(dm_build_1118, dm_build_1117.dm_build_820.dm_build_414)
	dm_build_1117.dm_build_1106.FailedAttempts = dm_build_1117.dm_build_820.dm_build_413.Dm_build_203()
	dm_build_1117.dm_build_1106.LoginWarningID = dm_build_1117.dm_build_820.dm_build_413.Dm_build_203()
	dm_build_1117.dm_build_1106.GraceTimeRemainder = dm_build_1117.dm_build_820.dm_build_413.Dm_build_203()
	dm_build_1117.dm_build_1106.Guid = dm_build_1117.dm_build_820.dm_build_413.Dm_build_245(dm_build_1118, dm_build_1117.dm_build_820.dm_build_414)
	dm_build_1117.dm_build_1106.DbName = dm_build_1117.dm_build_820.dm_build_413.Dm_build_245(dm_build_1118, dm_build_1117.dm_build_820.dm_build_414)

	if dm_build_1117.dm_build_820.dm_build_413.Dm_build_341(Dm_build_1100) == 1 {
		dm_build_1117.dm_build_1106.StandbyHost = dm_build_1117.dm_build_820.dm_build_413.Dm_build_245(dm_build_1118, dm_build_1117.dm_build_820.dm_build_414)
		dm_build_1117.dm_build_1106.StandbyPort = dm_build_1117.dm_build_820.dm_build_413.Dm_build_203()
		dm_build_1117.dm_build_1106.StandbyCount = int32(dm_build_1117.dm_build_820.dm_build_413.Dm_build_218())
	}

	if dm_build_1117.dm_build_820.dm_build_413.Dm_build_100(false) > 0 {
		dm_build_1117.dm_build_1106.SessionID = dm_build_1117.dm_build_820.dm_build_413.Dm_build_206()
	}

	if dm_build_1117.dm_build_820.dm_build_413.Dm_build_100(false) > 0 {
		if dm_build_1117.dm_build_820.dm_build_413.Dm_build_197() == 1 {

			dm_build_1117.dm_build_1106.FormatDate = "DD-MON-YY"

			dm_build_1117.dm_build_1106.FormatTime = "HH12.MI.SS.FF6 AM"

			dm_build_1117.dm_build_1106.FormatTimestamp = "DD-MON-YY HH12.MI.SS.FF6 AM"

			dm_build_1117.dm_build_1106.FormatTimestampTZ = "DD-MON-YY HH12.MI.SS.FF6 AM +TZH:TZM"

			dm_build_1117.dm_build_1106.FormatTimeTZ = "HH12.MI.SS.FF6 AM +TZH:TZM"
		}
	}

	if dm_build_1117.dm_build_820.dm_build_413.Dm_build_100(false) > 0 {

		format := dm_build_1117.dm_build_820.dm_build_413.Dm_build_249(dm_build_1118, dm_build_1117.dm_build_820.dm_build_414)
		if format != "" {
			dm_build_1117.dm_build_1106.FormatDate = format
		}
		format = dm_build_1117.dm_build_820.dm_build_413.Dm_build_249(dm_build_1118, dm_build_1117.dm_build_820.dm_build_414)
		if format != "" {
			dm_build_1117.dm_build_1106.FormatTime = format
		}
		format = dm_build_1117.dm_build_820.dm_build_413.Dm_build_249(dm_build_1118, dm_build_1117.dm_build_820.dm_build_414)
		if format != "" {
			dm_build_1117.dm_build_1106.FormatTimestamp = format
		}
		format = dm_build_1117.dm_build_820.dm_build_413.Dm_build_249(dm_build_1118, dm_build_1117.dm_build_820.dm_build_414)
		if format != "" {
			dm_build_1117.dm_build_1106.FormatTimestampTZ = format
		}
		format = dm_build_1117.dm_build_820.dm_build_413.Dm_build_249(dm_build_1118, dm_build_1117.dm_build_820.dm_build_414)
		if format != "" {
			dm_build_1117.dm_build_1106.FormatTimeTZ = format
		}
	}

	return nil, nil
}

const (
	Dm_build_1120 = Dm_build_712
)

type dm_build_1121 struct {
	dm_build_919
	dm_build_1122 int16
}

func dm_build_1123(dm_build_1124 *dm_build_410, dm_build_1125 *DmStatement, dm_build_1126 int16) *dm_build_1121 {
	dm_build_1127 := new(dm_build_1121)
	dm_build_1127.dm_build_828(dm_build_1124, Dm_build_704, dm_build_1125)
	dm_build_1127.dm_build_1122 = dm_build_1126
	return dm_build_1127
}

func (dm_build_1129 *dm_build_1121) dm_build_805() error {
	dm_build_1129.dm_build_820.dm_build_413.Dm_build_265(Dm_build_1120, dm_build_1129.dm_build_1122)
	return nil
}

func (dm_build_1131 *dm_build_1121) dm_build_809() (interface{}, error) {
	return dm_build_1131.dm_build_919.dm_build_809()
}

const (
	Dm_build_1132 = Dm_build_712
	Dm_build_1133 = Dm_build_1132 + USINT_SIZE
)

type dm_build_1134 struct {
	dm_build_919
	dm_build_1135 []parameter
}

func dm_build_1136(dm_build_1137 *dm_build_410, dm_build_1138 *DmStatement, dm_build_1139 []parameter) *dm_build_1134 {
	dm_build_1140 := new(dm_build_1134)
	dm_build_1140.dm_build_828(dm_build_1137, Dm_build_708, dm_build_1138)
	dm_build_1140.dm_build_1135 = dm_build_1139
	return dm_build_1140
}

func (dm_build_1142 *dm_build_1134) dm_build_805() error {

	if dm_build_1142.dm_build_1135 == nil {
		dm_build_1142.dm_build_820.dm_build_413.Dm_build_265(Dm_build_1132, 0)
	} else {
		dm_build_1142.dm_build_820.dm_build_413.Dm_build_265(Dm_build_1132, int16(len(dm_build_1142.dm_build_1135)))
	}

	dm_build_1142.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1133, 0)

	return dm_build_1142.dm_build_942(dm_build_1142.dm_build_1135)
}

type dm_build_1143 struct {
	dm_build_919
	dm_build_1144 bool
	dm_build_1145 int16
}

func dm_build_1146(dm_build_1147 *dm_build_410, dm_build_1148 *DmStatement, dm_build_1149 bool, dm_build_1150 int16) *dm_build_1143 {
	dm_build_1151 := new(dm_build_1143)
	dm_build_1151.dm_build_828(dm_build_1147, Dm_build_688, dm_build_1148)
	dm_build_1151.dm_build_1144 = dm_build_1149
	dm_build_1151.dm_build_1145 = dm_build_1150
	return dm_build_1151
}

func (dm_build_1153 *dm_build_1143) dm_build_805() error {

	dm_build_1154 := Dm_build_712

	if dm_build_1153.dm_build_820.dm_build_414.autoCommit {
		dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_261(dm_build_1154, 1)
	} else {
		dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_261(dm_build_1154, 0)
	}

	if dm_build_1153.dm_build_1144 {
		dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_261(dm_build_1154, 1)
	} else {
		dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_261(dm_build_1154, 0)
	}

	dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_261(dm_build_1154, 0)

	dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_261(dm_build_1154, 1)

	if dm_build_1153.dm_build_820.dm_build_414.CompatibleOracle() {
		dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_261(dm_build_1154, 0)
	} else {
		dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_261(dm_build_1154, 1)
	}

	dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_265(dm_build_1154, dm_build_1153.dm_build_1145)

	if dm_build_1153.dm_build_823.maxRows <= 0 || dm_build_1153.dm_build_820.dm_build_414.dmConnector.enRsCache {
		dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_273(dm_build_1154, INT64_MAX)
	} else {
		dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_273(dm_build_1154, dm_build_1153.dm_build_823.maxRows)
	}

	if dm_build_1153.dm_build_820.dm_build_414.dmConnector.isBdtaRS {
		dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_261(dm_build_1154, Dm_build_791)
	} else {
		dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_261(dm_build_1154, Dm_build_790)
	}

	dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_265(dm_build_1154, 0)

	dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_261(dm_build_1154, 1)

	dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_261(dm_build_1154, 0)

	dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_261(dm_build_1154, 0)

	dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_269(dm_build_1154, dm_build_1153.dm_build_823.queryTimeout)

	if dm_build_1153.dm_build_823.innerExec {
		dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_261(dm_build_1154, 1)
	} else {
		dm_build_1154 += dm_build_1153.dm_build_820.dm_build_413.Dm_build_261(dm_build_1154, 0)
	}

	dm_build_1153.dm_build_820.dm_build_413.Dm_build_191(dm_build_1153.dm_build_823.nativeSql, dm_build_1153.dm_build_820.dm_build_414.getServerEncoding(), dm_build_1153.dm_build_820.dm_build_414)

	return nil
}

func (dm_build_1156 *dm_build_1143) dm_build_809() (interface{}, error) {

	if dm_build_1156.dm_build_1144 {
		return dm_build_1156.dm_build_919.dm_build_809()
	}

	dm_build_1157 := NewExceInfo()
	dm_build_1158 := Dm_build_712

	dm_build_1157.retSqlType = dm_build_1156.dm_build_820.dm_build_413.Dm_build_344(dm_build_1158)
	dm_build_1158 += USINT_SIZE

	dm_build_1159 := dm_build_1156.dm_build_820.dm_build_413.Dm_build_362(dm_build_1158)
	dm_build_1158 += USINT_SIZE

	dm_build_1160 := dm_build_1156.dm_build_820.dm_build_413.Dm_build_344(dm_build_1158)
	dm_build_1158 += USINT_SIZE

	dm_build_1156.dm_build_820.dm_build_413.Dm_build_350(dm_build_1158)
	dm_build_1158 += DDWORD_SIZE

	dm_build_1156.dm_build_820.dm_build_414.TrxStatus = dm_build_1156.dm_build_820.dm_build_413.Dm_build_347(dm_build_1158)
	dm_build_1158 += ULINT_SIZE

	if dm_build_1159 > 0 {
		dm_build_1156.dm_build_823.params = dm_build_1156.dm_build_1161(int(dm_build_1159))
		dm_build_1156.dm_build_823.paramCount = int32(dm_build_1159)
	} else {
		dm_build_1156.dm_build_823.params = make([]parameter, 0)
		dm_build_1156.dm_build_823.paramCount = 0
	}

	if dm_build_1160 > 0 {
		dm_build_1156.dm_build_823.columns = dm_build_1156.dm_build_969(int(dm_build_1160), dm_build_1157.rsBdta)
	} else {
		dm_build_1156.dm_build_823.columns = make([]column, 0)
	}

	return dm_build_1157, nil
}

func (dm_build_1162 *dm_build_1143) dm_build_1161(dm_build_1163 int) []parameter {

	var dm_build_1164, dm_build_1165, dm_build_1166, dm_build_1167 int16

	dm_build_1168 := make([]parameter, dm_build_1163)
	for i := 0; i < dm_build_1163; i++ {

		dm_build_1168[i].InitParameter()

		dm_build_1168[i].colType = dm_build_1162.dm_build_820.dm_build_413.Dm_build_203()

		dm_build_1168[i].prec = dm_build_1162.dm_build_820.dm_build_413.Dm_build_203()

		dm_build_1168[i].scale = dm_build_1162.dm_build_820.dm_build_413.Dm_build_203()

		dm_build_1168[i].nullable = dm_build_1162.dm_build_820.dm_build_413.Dm_build_203() != 0

		itemFlag := dm_build_1162.dm_build_820.dm_build_413.Dm_build_200()

		if int(itemFlag)&Dm_build_918 != 0 {
			dm_build_1168[i].typeFlag = TYPE_FLAG_RECOMMEND
		} else {
			dm_build_1168[i].typeFlag = TYPE_FLAG_EXACT
		}

		dm_build_1168[i].lob = int(itemFlag)&Dm_build_916 != 0

		dm_build_1162.dm_build_820.dm_build_413.Dm_build_203()

		dm_build_1168[i].ioType = byte(dm_build_1162.dm_build_820.dm_build_413.Dm_build_200())

		dm_build_1164 = dm_build_1162.dm_build_820.dm_build_413.Dm_build_200()

		dm_build_1165 = dm_build_1162.dm_build_820.dm_build_413.Dm_build_200()

		dm_build_1166 = dm_build_1162.dm_build_820.dm_build_413.Dm_build_200()

		dm_build_1167 = dm_build_1162.dm_build_820.dm_build_413.Dm_build_200()
		dm_build_1168[i].name = dm_build_1162.dm_build_820.dm_build_413.Dm_build_240(int(dm_build_1164), dm_build_1162.dm_build_820.dm_build_414.getServerEncoding(), dm_build_1162.dm_build_820.dm_build_414)
		dm_build_1168[i].typeName = dm_build_1162.dm_build_820.dm_build_413.Dm_build_240(int(dm_build_1165), dm_build_1162.dm_build_820.dm_build_414.getServerEncoding(), dm_build_1162.dm_build_820.dm_build_414)
		dm_build_1168[i].tableName = dm_build_1162.dm_build_820.dm_build_413.Dm_build_240(int(dm_build_1166), dm_build_1162.dm_build_820.dm_build_414.getServerEncoding(), dm_build_1162.dm_build_820.dm_build_414)
		dm_build_1168[i].schemaName = dm_build_1162.dm_build_820.dm_build_413.Dm_build_240(int(dm_build_1167), dm_build_1162.dm_build_820.dm_build_414.getServerEncoding(), dm_build_1162.dm_build_820.dm_build_414)

		if dm_build_1168[i].lob {
			dm_build_1168[i].lobTabId = dm_build_1162.dm_build_820.dm_build_413.Dm_build_203()
			dm_build_1168[i].lobColId = dm_build_1162.dm_build_820.dm_build_413.Dm_build_200()
		}
	}

	for i := 0; i < dm_build_1163; i++ {

		if isComplexType(int(dm_build_1168[i].colType), int(dm_build_1168[i].scale)) {

			strDesc := newTypeDescriptor(dm_build_1162.dm_build_820.dm_build_414)
			strDesc.unpack(dm_build_1162.dm_build_820.dm_build_413)
			dm_build_1168[i].typeDescriptor = strDesc
		}
	}

	return dm_build_1168
}

const (
	Dm_build_1169 = Dm_build_712
)

type dm_build_1170 struct {
	dm_build_819
	dm_build_1171 int16
	dm_build_1172 *Dm_build_0
	dm_build_1173 int32
}

func dm_build_1174(dm_build_1175 *dm_build_410, dm_build_1176 *DmStatement, dm_build_1177 int16, dm_build_1178 *Dm_build_0, dm_build_1179 int32) *dm_build_1170 {
	dm_build_1180 := new(dm_build_1170)
	dm_build_1180.dm_build_828(dm_build_1175, Dm_build_694, dm_build_1176)
	dm_build_1180.dm_build_1171 = dm_build_1177
	dm_build_1180.dm_build_1172 = dm_build_1178
	dm_build_1180.dm_build_1173 = dm_build_1179
	return dm_build_1180
}

func (dm_build_1182 *dm_build_1170) dm_build_805() error {
	dm_build_1182.dm_build_820.dm_build_413.Dm_build_265(Dm_build_1169, dm_build_1182.dm_build_1171)

	dm_build_1182.dm_build_820.dm_build_413.Dm_build_129(dm_build_1182.dm_build_1173)

	if dm_build_1182.dm_build_820.dm_build_414.NewLobFlag {
		dm_build_1182.dm_build_820.dm_build_413.Dm_build_129(-1)
	}
	dm_build_1182.dm_build_1172.Dm_build_7(dm_build_1182.dm_build_820.dm_build_413, int(dm_build_1182.dm_build_1173))
	return nil
}

type dm_build_1183 struct {
	dm_build_819
}

func dm_build_1184(dm_build_1185 *dm_build_410) *dm_build_1183 {
	dm_build_1186 := new(dm_build_1183)
	dm_build_1186.dm_build_824(dm_build_1185, Dm_build_692)
	return dm_build_1186
}

type dm_build_1187 struct {
	dm_build_819
	dm_build_1188 int32
}

func dm_build_1189(dm_build_1190 *dm_build_410, dm_build_1191 int32) *dm_build_1187 {
	dm_build_1192 := new(dm_build_1187)
	dm_build_1192.dm_build_824(dm_build_1190, Dm_build_705)
	dm_build_1192.dm_build_1188 = dm_build_1191
	return dm_build_1192
}

func (dm_build_1194 *dm_build_1187) dm_build_805() error {

	dm_build_1195 := Dm_build_712
	dm_build_1195 += dm_build_1194.dm_build_820.dm_build_413.Dm_build_269(dm_build_1195, g2dbIsoLevel(dm_build_1194.dm_build_1188))
	return nil
}

type dm_build_1196 struct {
	dm_build_819
	dm_build_1197 *lob
	dm_build_1198 byte
	dm_build_1199 int
	dm_build_1200 []byte
	dm_build_1201 int
	dm_build_1202 int
}

func dm_build_1203(dm_build_1204 *dm_build_410, dm_build_1205 *lob, dm_build_1206 byte, dm_build_1207 int, dm_build_1208 []byte,
	dm_build_1209 int, dm_build_1210 int) *dm_build_1196 {
	dm_build_1211 := new(dm_build_1196)
	dm_build_1211.dm_build_824(dm_build_1204, Dm_build_701)
	dm_build_1211.dm_build_1197 = dm_build_1205
	dm_build_1211.dm_build_1198 = dm_build_1206
	dm_build_1211.dm_build_1199 = dm_build_1207
	dm_build_1211.dm_build_1200 = dm_build_1208
	dm_build_1211.dm_build_1201 = dm_build_1209
	dm_build_1211.dm_build_1202 = dm_build_1210
	return dm_build_1211
}

func (dm_build_1213 *dm_build_1196) dm_build_805() error {

	dm_build_1213.dm_build_820.dm_build_413.Dm_build_121(byte(dm_build_1213.dm_build_1197.lobFlag))
	dm_build_1213.dm_build_820.dm_build_413.Dm_build_121(dm_build_1213.dm_build_1198)
	dm_build_1213.dm_build_820.dm_build_413.Dm_build_145(uint64(dm_build_1213.dm_build_1197.blobId))
	dm_build_1213.dm_build_820.dm_build_413.Dm_build_125(dm_build_1213.dm_build_1197.groupId)
	dm_build_1213.dm_build_820.dm_build_413.Dm_build_125(dm_build_1213.dm_build_1197.fileId)
	dm_build_1213.dm_build_820.dm_build_413.Dm_build_129(dm_build_1213.dm_build_1197.pageNo)
	dm_build_1213.dm_build_820.dm_build_413.Dm_build_125(dm_build_1213.dm_build_1197.curFileId)
	dm_build_1213.dm_build_820.dm_build_413.Dm_build_129(dm_build_1213.dm_build_1197.curPageNo)
	dm_build_1213.dm_build_820.dm_build_413.Dm_build_129(dm_build_1213.dm_build_1197.totalOffset)
	dm_build_1213.dm_build_820.dm_build_413.Dm_build_129(dm_build_1213.dm_build_1197.tabId)
	dm_build_1213.dm_build_820.dm_build_413.Dm_build_125(dm_build_1213.dm_build_1197.colId)
	dm_build_1213.dm_build_820.dm_build_413.Dm_build_145(uint64(dm_build_1213.dm_build_1197.rowId))

	dm_build_1213.dm_build_820.dm_build_413.Dm_build_129(int32(dm_build_1213.dm_build_1199))
	dm_build_1213.dm_build_820.dm_build_413.Dm_build_129(int32(dm_build_1213.dm_build_1202))
	dm_build_1213.dm_build_820.dm_build_413.Dm_build_157(dm_build_1213.dm_build_1200)

	if dm_build_1213.dm_build_820.dm_build_414.NewLobFlag {
		dm_build_1213.dm_build_820.dm_build_413.Dm_build_125(dm_build_1213.dm_build_1197.exGroupId)
		dm_build_1213.dm_build_820.dm_build_413.Dm_build_125(dm_build_1213.dm_build_1197.exFileId)
		dm_build_1213.dm_build_820.dm_build_413.Dm_build_129(dm_build_1213.dm_build_1197.exPageNo)
	}
	return nil
}

func (dm_build_1215 *dm_build_1196) dm_build_809() (interface{}, error) {

	var dm_build_1216 = dm_build_1215.dm_build_820.dm_build_413.Dm_build_203()
	dm_build_1215.dm_build_1197.blobId = dm_build_1215.dm_build_820.dm_build_413.Dm_build_206()
	dm_build_1215.dm_build_1197.fileId = dm_build_1215.dm_build_820.dm_build_413.Dm_build_200()
	dm_build_1215.dm_build_1197.pageNo = dm_build_1215.dm_build_820.dm_build_413.Dm_build_203()
	dm_build_1215.dm_build_1197.curFileId = dm_build_1215.dm_build_820.dm_build_413.Dm_build_200()
	dm_build_1215.dm_build_1197.curPageNo = dm_build_1215.dm_build_820.dm_build_413.Dm_build_203()
	dm_build_1215.dm_build_1197.totalOffset = dm_build_1215.dm_build_820.dm_build_413.Dm_build_203()
	return dm_build_1216, nil
}

const (
	Dm_build_1217 = Dm_build_712

	Dm_build_1218 = Dm_build_1217 + ULINT_SIZE

	Dm_build_1219 = Dm_build_1218 + ULINT_SIZE

	Dm_build_1220 = Dm_build_1219 + BYTE_SIZE

	Dm_build_1221 = Dm_build_1220 + BYTE_SIZE

	Dm_build_1222 = Dm_build_1221 + BYTE_SIZE

	Dm_build_1223 = Dm_build_1222 + BYTE_SIZE

	Dm_build_1224 = Dm_build_1223 + BYTE_SIZE

	Dm_build_1225 = Dm_build_1224 + BYTE_SIZE

	Dm_build_1226 = Dm_build_1225 + BYTE_SIZE

	Dm_build_1227 = Dm_build_712

	Dm_build_1228 = Dm_build_1227 + ULINT_SIZE

	Dm_build_1229 = Dm_build_1228 + ULINT_SIZE

	Dm_build_1230 = Dm_build_1229 + ULINT_SIZE

	Dm_build_1231 = Dm_build_1230 + ULINT_SIZE

	Dm_build_1232 = Dm_build_1231 + ULINT_SIZE

	Dm_build_1233 = Dm_build_1232 + BYTE_SIZE

	Dm_build_1234 = Dm_build_1233 + BYTE_SIZE

	Dm_build_1235 = Dm_build_1234 + BYTE_SIZE

	Dm_build_1236 = Dm_build_1235 + BYTE_SIZE

	Dm_build_1237 = Dm_build_1236 + BYTE_SIZE

	Dm_build_1238 = Dm_build_1237 + USINT_SIZE

	Dm_build_1239 = Dm_build_1238 + BYTE_SIZE
)

type dm_build_1240 struct {
	dm_build_819
	dm_build_1241 *DmConnection
	dm_build_1242 int
	Dm_build_1243 int32
	Dm_build_1244 []byte
	dm_build_1245 byte
}

func dm_build_1246(dm_build_1247 *dm_build_410) *dm_build_1240 {
	dm_build_1248 := new(dm_build_1240)
	dm_build_1248.dm_build_824(dm_build_1247, Dm_build_710)
	dm_build_1248.dm_build_1241 = dm_build_1247.dm_build_414
	return dm_build_1248
}

func dm_build_1249(dm_build_1250 string, dm_build_1251 string) int {
	dm_build_1252 := strings.Split(dm_build_1250, ".")
	dm_build_1253 := strings.Split(dm_build_1251, ".")

	for i, serStr := range dm_build_1252 {
		ser, _ := strconv.ParseInt(serStr, 10, 32)
		global, _ := strconv.ParseInt(dm_build_1253[i], 10, 32)
		if ser < global {
			return -1
		} else if ser == global {
			continue
		} else {
			return 1
		}
	}

	return 0
}

func (dm_build_1255 *dm_build_1240) dm_build_805() error {

	dm_build_1255.dm_build_820.dm_build_413.Dm_build_269(Dm_build_1217, int32(0))
	dm_build_1255.dm_build_820.dm_build_413.Dm_build_269(Dm_build_1218, int32(dm_build_1255.dm_build_1241.dmConnector.compress))

	if dm_build_1255.dm_build_1241.dmConnector.loginEncrypt {
		dm_build_1255.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1220, 2)
		dm_build_1255.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1219, 1)
	} else {
		dm_build_1255.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1220, 0)
		dm_build_1255.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1219, 0)
	}

	if dm_build_1255.dm_build_1241.dmConnector.isBdtaRS {
		dm_build_1255.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1221, Dm_build_791)
	} else {
		dm_build_1255.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1221, Dm_build_790)
	}

	dm_build_1255.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1222, byte(dm_build_1255.dm_build_1241.dmConnector.compressID))

	if dm_build_1255.dm_build_1241.dmConnector.loginCertificate != "" {
		dm_build_1255.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1223, 1)
	} else {
		dm_build_1255.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1223, 0)
	}

	dm_build_1255.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1224, 0)
	dm_build_1255.dm_build_820.dm_build_413.Dm_build_261(Dm_build_1225, 1)
	dm_build_1255.dm_build_820.dm_build_413.Dm_build_289(Dm_build_1226, uint16(dm_build_1255.dm_build_1241.MsgVersion))

	dm_build_1256 := dm_build_1255.dm_build_1241.getServerEncoding()
	dm_build_1255.dm_build_820.dm_build_413.Dm_build_173(Dm_build_675, dm_build_1256, dm_build_1255.dm_build_820.dm_build_414)

	var dm_build_1257 byte
	if dm_build_1255.dm_build_1241.dmConnector.uKeyName != "" {
		dm_build_1257 = 1
	} else {
		dm_build_1257 = 0
	}

	dm_build_1255.dm_build_820.dm_build_413.Dm_build_121(0)

	if dm_build_1257 == 1 {

	}

	if dm_build_1255.dm_build_1241.dmConnector.loginEncrypt {
		clientPubKey, err := dm_build_1255.dm_build_820.dm_build_655()
		if err != nil {
			return err
		}
		dm_build_1255.dm_build_820.dm_build_413.Dm_build_161(clientPubKey)
	}

	return nil
}

func (dm_build_1259 *dm_build_1240) dm_build_808() error {
	dm_build_1259.dm_build_820.dm_build_413.Dm_build_95(0)
	dm_build_1259.dm_build_820.dm_build_413.Dm_build_103(Dm_build_711, false, true)
	return nil
}

func (dm_build_1261 *dm_build_1240) dm_build_809() (interface{}, error) {

	dm_build_1261.dm_build_1241.sslEncrypt = int(dm_build_1261.dm_build_820.dm_build_413.Dm_build_347(Dm_build_1227))
	dm_build_1261.dm_build_1241.GlobalServerSeries = int(dm_build_1261.dm_build_820.dm_build_413.Dm_build_347(Dm_build_1228))

	switch dm_build_1261.dm_build_820.dm_build_413.Dm_build_347(Dm_build_1229) {
	case 1:
		dm_build_1261.dm_build_1241.serverEncoding = ENCODING_UTF8
	case 2:
		dm_build_1261.dm_build_1241.serverEncoding = ENCODING_EUCKR
	default:
		dm_build_1261.dm_build_1241.serverEncoding = ENCODING_GB18030
	}

	dm_build_1261.dm_build_1241.dmConnector.compress = int(dm_build_1261.dm_build_820.dm_build_413.Dm_build_347(Dm_build_1230))
	dm_build_1262 := dm_build_1261.dm_build_820.dm_build_413.Dm_build_341(Dm_build_1232)
	dm_build_1263 := dm_build_1261.dm_build_820.dm_build_413.Dm_build_341(Dm_build_1233)
	dm_build_1261.dm_build_1241.dmConnector.isBdtaRS = dm_build_1261.dm_build_820.dm_build_413.Dm_build_341(Dm_build_1234) > 0
	dm_build_1261.dm_build_1241.dmConnector.compressID = int8(dm_build_1261.dm_build_820.dm_build_413.Dm_build_341(Dm_build_1235))

	dm_build_1261.dm_build_820.dm_build_419 = dm_build_1261.dm_build_820.dm_build_413.Dm_build_341(Dm_build_1237) == 1
	dm_build_1261.dm_build_1241.dmConnector.newClientType = dm_build_1261.dm_build_820.dm_build_413.Dm_build_341(Dm_build_1238) == 1
	dm_build_1261.dm_build_1241.MsgVersion = int32(dm_build_1261.dm_build_820.dm_build_413.Dm_build_362(Dm_build_1239))

	dm_build_1264 := dm_build_1261.dm_build_852()
	if dm_build_1264 != nil {
		return nil, dm_build_1264
	}

	dm_build_1265 := dm_build_1261.dm_build_820.dm_build_413.Dm_build_245(dm_build_1261.dm_build_1241.getServerEncoding(), dm_build_1261.dm_build_820.dm_build_414)
	if dm_build_1249(dm_build_1265, Dm_build_676) < 0 {
		return nil, ECGO_ERROR_SERVER_VERSION.throw()
	}

	dm_build_1261.dm_build_1241.ServerVersion = dm_build_1265
	dm_build_1261.dm_build_1241.Malini2 = dm_build_1249(dm_build_1265, Dm_build_677) > 0
	dm_build_1261.dm_build_1241.Execute2 = dm_build_1249(dm_build_1265, Dm_build_678) > 0
	dm_build_1261.dm_build_1241.LobEmptyCompOrcl = dm_build_1249(dm_build_1265, Dm_build_679) > 0

	if dm_build_1261.dm_build_820.dm_build_414.dmConnector.uKeyName != "" {
		dm_build_1261.dm_build_1245 = 1
	} else {
		dm_build_1261.dm_build_1245 = 0
	}

	if dm_build_1261.dm_build_1245 == 1 {
		dm_build_1261.dm_build_820.dm_build_420 = dm_build_1261.dm_build_820.dm_build_413.Dm_build_240(16, dm_build_1261.dm_build_1241.getServerEncoding(), dm_build_1261.dm_build_820.dm_build_414)
	}

	dm_build_1261.dm_build_1242 = -1
	dm_build_1266 := false
	dm_build_1267 := false
	dm_build_1261.Dm_build_1243 = -1
	if dm_build_1263 > 0 {
		dm_build_1261.dm_build_1242 = int(dm_build_1261.dm_build_820.dm_build_413.Dm_build_203())
	}

	if dm_build_1262 > 0 {

		if dm_build_1261.dm_build_1242 == -1 {
			dm_build_1266 = true
		} else {
			dm_build_1267 = true
		}

		dm_build_1261.Dm_build_1244 = dm_build_1261.dm_build_820.dm_build_413.Dm_build_228()
	}

	if dm_build_1263 == 2 {
		dm_build_1261.Dm_build_1243 = dm_build_1261.dm_build_820.dm_build_413.Dm_build_203()
	}
	dm_build_1261.dm_build_820.dm_build_416 = dm_build_1266
	dm_build_1261.dm_build_820.dm_build_417 = dm_build_1267

	return nil, nil
}

type dm_build_1268 struct {
	dm_build_819
}

func dm_build_1269(dm_build_1270 *dm_build_410, dm_build_1271 *DmStatement) *dm_build_1268 {
	dm_build_1272 := new(dm_build_1268)
	dm_build_1272.dm_build_828(dm_build_1270, Dm_build_686, dm_build_1271)
	return dm_build_1272
}

func (dm_build_1274 *dm_build_1268) dm_build_805() error {

	dm_build_1274.dm_build_820.dm_build_413.Dm_build_261(Dm_build_712, 1)
	return nil
}

func (dm_build_1276 *dm_build_1268) dm_build_809() (interface{}, error) {

	dm_build_1276.dm_build_823.id = dm_build_1276.dm_build_820.dm_build_413.Dm_build_347(Dm_build_713)

	dm_build_1276.dm_build_823.readBaseColName = dm_build_1276.dm_build_820.dm_build_413.Dm_build_341(Dm_build_712) == 1
	return nil, nil
}

type dm_build_1277 struct {
	dm_build_819
	dm_build_1278 int32
}

func dm_build_1279(dm_build_1280 *dm_build_410, dm_build_1281 int32) *dm_build_1277 {
	dm_build_1282 := new(dm_build_1277)
	dm_build_1282.dm_build_824(dm_build_1280, Dm_build_687)
	dm_build_1282.dm_build_1278 = dm_build_1281
	return dm_build_1282
}

func (dm_build_1284 *dm_build_1277) dm_build_806() {
	dm_build_1284.dm_build_819.dm_build_806()
	dm_build_1284.dm_build_820.dm_build_413.Dm_build_269(Dm_build_713, dm_build_1284.dm_build_1278)
}

type dm_build_1285 struct {
	dm_build_819
	dm_build_1286 []uint32
}

func dm_build_1287(dm_build_1288 *dm_build_410, dm_build_1289 []uint32) *dm_build_1285 {
	dm_build_1290 := new(dm_build_1285)
	dm_build_1290.dm_build_824(dm_build_1288, Dm_build_707)
	dm_build_1290.dm_build_1286 = dm_build_1289
	return dm_build_1290
}

func (dm_build_1292 *dm_build_1285) dm_build_805() error {

	dm_build_1292.dm_build_820.dm_build_413.Dm_build_289(Dm_build_712, uint16(len(dm_build_1292.dm_build_1286)))

	for _, tableID := range dm_build_1292.dm_build_1286 {
		dm_build_1292.dm_build_820.dm_build_413.Dm_build_141(uint32(tableID))
	}

	return nil
}

func (dm_build_1294 *dm_build_1285) dm_build_809() (interface{}, error) {
	dm_build_1295 := dm_build_1294.dm_build_820.dm_build_413.Dm_build_362(Dm_build_712)
	if dm_build_1295 <= 0 {
		return nil, nil
	}

	dm_build_1296 := make([]int64, dm_build_1295)
	for i := 0; i < int(dm_build_1295); i++ {
		dm_build_1296[i] = dm_build_1294.dm_build_820.dm_build_413.Dm_build_206()
	}

	return dm_build_1296, nil
}
