package imbencode

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"strconv"
	"strings"
)

// Encode takes a numeric code (in the valid format) and returns a string in ATDF format, ready to be printed
// with the USPS IMB font
/**
 * Code extracted from TCPDF and modified to operate standalone. See http://www.tcpdf.org
 * Original conversion to PHP from Ian Simpson
 * Converted to golang by Scott Hillock
 *
 * Some methods in PHP are much shorter in golang, but are preserved for sake of comparison
 *
 * To use import "github.com/scotthillock/uspsimbencoder"
 * then encoded := imbencode.Encode(string)
 */
func Encode(input string) string {

	/**
	 * IMB - Intelligent Mail Barcode - Onecode - USPS-B-3200
	 * Intelligent Mail barcode is a 65-bar code for use on mail in the United States.
	 * The fields are described as follows:<ul><li>The Barcode Identifier shall be assigned by USPS to encode the presort identification that is currently printed in human readable form on the optional endorsement line (OEL) as well as for future USPS use. This shall be two digits, with the second digit in the range of 0–4. The allowable encoding ranges shall be 00–04, 10–14, 20–24, 30–34, 40–44, 50–54, 60–64, 70–74, 80–84, and 90–94.</li><li>The Service Type Identifier shall be assigned by USPS for any combination of services requested on the mailpiece. The allowable encoding range shall be 000http://it2.php.net/manual/en/function.dechex.php–999. Each 3-digit value shall correspond to a particular mail class with a particular combination of service(s). Each service program, such as OneCode Confirm and OneCode ACS, shall provide the list of Service Type Identifier values.</li><li>The Mailer or Customer Identifier shall be assigned by USPS as a unique, 6 or 9 digit number that identifies a business entity. The allowable encoding range for the 6 digit Mailer ID shall be 000000- 899999, while the allowable encoding range for the 9 digit Mailer ID shall be 900000000-999999999.</li><li>The Serial or Sequence Number shall be assigned by the mailer for uniquely identifying and tracking mailpieces. The allowable encoding range shall be 000000000–999999999 when used with a 6 digit Mailer ID and 000000-999999 when used with a 9 digit Mailer ID. e. The Delivery Point ZIP Code shall be assigned by the mailer for routing the mailpiece. This shall replace POSTNET for routing the mailpiece to its final delivery point. The length may be 0, 5, 9, or 11 digits. The allowable encoding ranges shall be no ZIP Code, 00000–99999,  000000000–999999999, and 00000000000–99999999999.</li></ul>
	 */

	ascChr := []int{4, 0, 2, 6, 3, 5, 1, 9, 8, 7, 1, 2, 0, 6, 4, 8, 2, 9, 5, 3, 0, 1, 3, 7, 4, 6, 8, 9, 2, 0, 5, 1, 9, 4, 3, 8, 6, 7, 1, 2, 4, 3, 9, 5, 7, 8, 3, 0, 2, 1, 4, 0, 9, 1, 7, 0, 2, 4, 6, 3, 7, 1, 9, 5, 8}
	dscChr := []int{7, 1, 9, 5, 8, 0, 2, 4, 6, 3, 5, 8, 9, 7, 3, 0, 6, 1, 7, 4, 6, 8, 9, 2, 5, 1, 7, 5, 4, 3, 8, 7, 6, 0, 2, 5, 4, 9, 3, 0, 1, 6, 8, 2, 0, 4, 5, 9, 6, 7, 5, 2, 6, 3, 8, 5, 1, 9, 8, 7, 4, 0, 2, 6, 3}
	ascPos := []int{3, 0, 8, 11, 1, 12, 8, 11, 10, 6, 4, 12, 2, 7, 9, 6, 7, 9, 2, 8, 4, 0, 12, 7, 10, 9, 0, 7, 10, 5, 7, 9, 6, 8, 2, 12, 1, 4, 2, 0, 1, 5, 4, 6, 12, 1, 0, 9, 4, 7, 5, 10, 2, 6, 9, 11, 2, 12, 6, 7, 5, 11, 0, 3, 2}
	dscPos := []int{2, 10, 12, 5, 9, 1, 5, 4, 3, 9, 11, 5, 10, 1, 6, 3, 4, 1, 10, 0, 2, 11, 8, 6, 1, 12, 3, 8, 6, 4, 4, 11, 0, 6, 1, 9, 11, 5, 3, 7, 3, 10, 7, 11, 8, 2, 10, 3, 5, 8, 0, 3, 12, 11, 8, 4, 5, 1, 3, 0, 7, 12, 9, 8, 10}

	trackingNumber := input[0:20]
	codeLength := len(input)
	routingCode := input[20:codeLength]
	rcBI := new(big.Int)
	fmt.Sscan(routingCode, rcBI)

	// Conversion of Routing Code
	binaryCode := new(big.Int)
	switch rcl := len(routingCode); rcl {
	case 0:
		{
			binaryCode.SetString("0", 10)
			break
		}
	case 5:
		{
			binaryCode.Add(rcBI, big.NewInt(1))
			break
		}
	case 9:
		{
			binaryCode.Add(rcBI, big.NewInt(100001))
			break
		}
	case 11:
		{
			binaryCode.Add(rcBI, big.NewInt(1000100001))
			break
		}
	default:
		{
			break
		}
	}

	binaryCode.Mul(binaryCode, big.NewInt(10))
	temp, _ := strconv.ParseInt(trackingNumber[0:1], 10, 64)
	binaryCode.Add(binaryCode, big.NewInt(temp))
	binaryCode.Mul(binaryCode, big.NewInt(5))
	temp, _ = strconv.ParseInt(trackingNumber[1:2], 10, 64)
	binaryCode.Add(binaryCode, big.NewInt(temp))
	bcAsString := binaryCode.String() + trackingNumber[2:20]
	binaryCode.SetString(bcAsString, 10)
	// convert to hexadecimal
	bcAsHex := decToHex(binaryCode)
	// pad to get 13 bytes
	bcAsHex = fmt.Sprintf("%026v", bcAsHex)

	// convert string to array of bytes
	var binaryCodeArr [13]string
	for i := 0; i < 13; i++ {
		start := i * 2
		end := (i + 1) * 2
		binaryCodeArr[i] = bcAsHex[start:end]
	}

	// calculate frame check sequence
	fcs := imbCrc11fcs(binaryCodeArr)

	// exclude first 2 bits from first byte
	firstByte := fmt.Sprintf("%2s", decToHex(big.NewInt(int64((hexToDec(binaryCodeArr[0])<<2)>>2))))
	length := len(bcAsHex)
	binaryCode102bit := firstByte + bcAsHex[2:length]

	// convert binary data to codewords
	var codeWords [10]int64
	data, _ := new(big.Int).SetString(binaryCode102bit, 16)
	mod, _ := new(big.Int).SetString(binaryCode102bit, 16)

	codeWords[0], _ = strconv.ParseInt((mod.Mod(data, big.NewInt(636))).String(), 10, 10)
	codeWords[0] = codeWords[0] * 2

	data = data.Div(data, big.NewInt(636))

	for i := 1; i < 9; i++ {
		cwResult := new(big.Int)
		cwResult = mod.Mod(data, big.NewInt(1365))
		cwrI, _ := strconv.ParseInt(cwResult.String(), 10, 64)
		codeWords[i] = cwrI
		data = new(big.Int).Div(data, big.NewInt(1365))
	}

	codeWords[9], _ = strconv.ParseInt(data.String(), 10, 64)

	if (fcs >> 10) == 1 {
		codeWords[9] = codeWords[9] + 659
	}

	// generate lookup tables
	table2of13 := imbTables(2, 78)
	table5of13 := imbTables(5, 1287)

	// convert codewords to characters
	var characters []int
	bitMask := 512
	chrCode := 0

	for _, val := range codeWords {
		if val <= 1286 {
			chrCode = table5of13[val]
		} else {
			index := val - 1287
			chrCode = table2of13[index]
		}
		if (fcs & bitMask) > 0 {
			chrCode = ((^chrCode) & 8191)
		}
		characters = append(characters, chrCode)
		bitMask /= 2

	}

	characters = arrayReverse(characters)

	// build bars
	var out []string
	var c string
	for i := 0; i <= 64; i++ {
		ap2 := math.Exp2(float64(ascPos[i]))
		dp2 := math.Exp2(float64(dscPos[i]))

		asc := ((characters[ascChr[i]] & int(ap2)) > 0)
		dsc := ((characters[dscChr[i]] & int(dp2)) > 0)

		if asc && dsc {
			// full bar (F)
			c = "F"
		} else if asc {
			// ascender (A)
			c = "A"
		} else if dsc {
			// descender (D)
			c = "D"
		} else {
			// tracker (T)
			c = "T"
		}
		out = append(out, c)

	}
	return strings.Join(out, "")
}

//Convert big Int to hexadecimal representation
func decToHex(i *big.Int) string {
	hex := fmt.Sprintf("%X", i)
	return hex
}

//Convert hexadecimal to integer
func hexToDec(s string) int {
	result, _ := strconv.ParseInt(s, 16, 64)
	return int(result)
}

//Intelligent Mail Barcode calculation of Frame Check Sequence
func imbCrc11fcs(codeArr [13]string) int {
	genpoly := 0x0F35 // generator polynomial
	fcs := 0x07FF     // Frame Check Sequence

	// do most significant byte skipping the 2 most significant bits
	data := hexToDec(codeArr[0]) << 5

	for bit := 2; bit < 8; bit++ {
		if (fcs^data)&0x400 != 0 {
			fcs = (fcs << 1) ^ genpoly
		} else {
			fcs = (fcs << 1)
		}
		fcs &= 0x7FF
		data <<= 1
	}

	// do rest of bytes
	for byte := 1; byte < 13; byte++ {
		data = hexToDec(codeArr[byte]) << 3
		for bit := 0; bit < 8; bit++ {
			if (fcs^data)&0x400 != 0 {
				fcs = (fcs << 1) ^ genpoly
			} else {
				fcs = (fcs << 1)
			}
			fcs &= 0x7FF
			data <<= 1
		}
	}

	return fcs
}

//Reverse unsigned short value
func imbReverseUs(num int) int {
	rev := 0
	for i := 0; i < 16; i++ {
		rev <<= 1
		rev |= (num & 1)
		num >>= 1
	}
	return rev
}

//generate Nof13 tables used for Intelligent Mail Barcode
func imbTables(n, size int) []int {

	table := make([]int, size)

	lli := 0        // LUT lower index
	lui := size - 1 // LUT upper index
	for count := 0; count < 8192; count++ {
		uc := uint(count)
		bitCount := bits.OnesCount(uc)
		// if we don't have the right number of bits on, go on to the next value
		if bitCount == n {
			reverse := imbReverseUs(count) >> 3
			// if the reverse is less than count, we have already visited this pair before
			if reverse >= count {
				// If count is symmetric, place it at the first free slot from the end of the list.
				// Otherwise, place it at the first free slot from the beginning of the list AND place $reverse ath the next free slot from the beginning of the list
				if reverse == count {
					table = append(table)
					table[lui] = count
					lui--
				} else {
					table[lli] = count
					lli++
					table[lli] = reverse
					lli++
				}
			}
		}
	}
	return table
}

func arrayReverse(arr []int) []int {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
