package findinbytes

import (
	"bytes"
	"fmt"
	"log"
)

// save all []bytes via get json
type JsonObject struct {
	Data [][]byte
}

// return pointer of *jsonObject
func NewByteJsonObj() *JsonObject {
	return &JsonObject{
		Data: [][]byte(nil),
	}
}

// append []bytes  in data
func (j *JsonObject) AppendByte(inJsonBytes ...[]byte) *JsonObject {

	for i := 0; i < len(inJsonBytes); i++ {
		//log.Printf("inJsonBytes: %s\n", inJsonBytes)
		if inJsonBytes[i] == nil || len(inJsonBytes[i]) < 7 {
			log.Panicln("cant'n send nill slice or empty json to append")
		}

		log.Printf("inJsonBytes [i]: %s\n", inJsonBytes[i])
		ok := j.IsValidJsonBytes(inJsonBytes[i])
		if !ok {
			log.Panicf("input json bytes is invalid format at index: %d, input:%s\n", i, inJsonBytes[i])
		}
		j.Data = append(j.Data, inJsonBytes[i])
	}
	log.Println("len of inJsonBytes: ", len(inJsonBytes))

	//log.Printf("j.Data: %s\n", j.Data)
	//log.Printf("j.Size: %d\n", j.Size)
	return j
}

// check json bytes with 1st and last byte must be in '{' and '}'
// minimum bytes json len : 7 e.g : {"a":4}
func (j *JsonObject) IsValidJsonBytes(inJsonBytes []byte) bool {
	if inJsonBytes == nil || len(inJsonBytes) < 7 {
		log.Println("cant'n send nill or empty json bytes slice to append")
		return false
	}
	var braceL byte = '{'
	var braceR byte = '}'
	n := len(inJsonBytes)
	return inJsonBytes[0] == braceL && inJsonBytes[n-1] == braceR
}

// get all data
func (j *JsonObject) GetAllData() [][]byte {
	return j.Data
}

// get  records as given indices from [][]data
func (j *JsonObject) GetDataByIndicesIn2D(indices []int) (data [][]byte, err error) {
	n := len(indices)
	size := j.GetDataSize()

	if n > size {
		return nil, fmt.Errorf("indices size is greater than size of data")
	}
	if indices == nil || n == 0 {
		return nil, fmt.Errorf("given indices are nil")
	}
	for _, val := range indices {
		data = append(data, j.Data[val])
	}
	return data, nil
}

// get  records as given indices from [][]data in 1d slice , to seperate the value use the Char ';' samicolumn
func (j *JsonObject) GetDataByIndicesIn1D(indices []int) (data []byte, err error) {
	n := len(indices)
	size := j.GetDataSize()

	if n > size {
		return nil, fmt.Errorf("indices size is greater than size of data")
	}
	if indices == nil || n == 0 {
		return nil, fmt.Errorf("given indices are nil")
	}
	for i, val := range indices {
		data = append(data, j.Data[val]...)
		if i != len(indices)-1 {
			data = append(data, ';')
		}
	}
	return data, nil
}

// get [][]data size, if not empty otherwise 0
func (j *JsonObject) GetDataSize() int {
	return len(j.Data)
}

// ?
// delete exact bytes in [][]data
func (j *JsonObject) DeleteByteById(in []byte) (bool, error) {

	valid := j.IsValidJsonBytes(in)
	n := len(in)

	size := j.GetDataSize()
	if in == nil || n == 0 || !valid {
		return false, fmt.Errorf("given bytes are nil or invalid format")

	}
	if size == 0 {
		return false, fmt.Errorf("don't try in empty data")
	}
	inByte, err := RemoveBraceTopBottom(in)
	if err != nil {
		return false, fmt.Errorf("given bytes are nil or invalid format")
	}

	flag := false

	if size == 1 {
		ok := bytes.Contains(j.Data[0], inByte)
		if ok {
			j.Data = nil
			return true, nil
		}
	}

	newData := [][]byte{}
	for i := 0; i < size; i++ {
		ok := bytes.Contains(j.Data[i], inByte)
		log.Printf("} ok:%v \n", ok)
		if ok {
			flag = true
			newData = append(newData, j.Data[:i]...)
			newData = append(newData, j.Data[i+1:]...)
			log.Printf("} ok:%v,  newData:%s\n", ok, newData)
		}
	}
	if flag {
		j.Data = newData
		return size != j.GetDataSize(), nil
	}
	log.Printf("} j.Data:%s\n", newData)
	return false, fmt.Errorf("could not deleted")
}

// find all field and value  exists in all data. return data indices
// e.g data := [{"name":"name_field","age":15},{"name":"name_field","age":34}]
// find := {"name":"name_field","age":15}
// find 1st fields and value exists in data[0], data[1]
// return data's indices [0,1]
func (j *JsonObject) GetIndexesExistsAllFields(find []byte) (idx []int, err error) {
	log.Println("} find :", string(find))

	ok := j.IsValidJsonBytes(find)
	if !ok {
		return nil, fmt.Errorf("input json bytes is invalid format")
	}

	if j.Data == nil {
		return nil, fmt.Errorf("data is empty")
	}

	var existsIndexs []int
	splitDatas, err := SplitBytesByComma(find)

	log.Println("} len splitDatas:", len(splitDatas))
	log.Println("} splitDatas:", splitDatas)
	if err != nil {
		return nil, err
	}
	flag := false
	size := j.GetDataSize()
	// iterate the j.data[] that find bytes  is exists
	for i := 0; i < size; i++ {

		//check splitDatas's every json bytes field exists in j.Data[idx]
		for s := 0; s < len(splitDatas); s++ {
			log.Println("} s :", splitDatas[s], "row: ", j.Data[i])
			log.Println("} s :", string(splitDatas[s]), "row: ", string(j.Data[i]))

			//check: sometime found row as  [] or string ""
			if splitDatas[s] != nil && string(splitDatas[s]) != "" {
				ok := bytes.Contains(j.Data[i], splitDatas[s])
				if ok {
					existsIndexs = append(existsIndexs, i)
					flag = true
					break
				}
			}
		}
	}
	if flag {
		return existsIndexs, nil
	}

	return nil, fmt.Errorf("colud not found records")
}

// find exact field and value  exists in all data. return data indices.
// e.g data := [{"name":"name_field","age":15},{"name":"name_field","age":34}]
// find := {"name":"name_field"}
// find exact field and value exists in data[0]
// return data's index [0]
func (j *JsonObject) FindGetExactExistsIndexs(find []byte) (idx int, err error) {
	ok := j.IsValidJsonBytes(find)
	if !ok {
		return -1, fmt.Errorf("input json bytes is invalid format")
	}

	log.Println("} FindGetExactExistsIndexs, find", find)
	log.Println("} j.size, ", j.GetDataSize())
	// iterate the j.data[] that find bytes  is exists
	n := j.GetDataSize()
	log.Println("} n, ", n)
	for i := 0; i < n; i++ {
		log.Println("} in  loop ")
		//check: sometime found row as  [] or string ""
		if find != nil && string(find) != "" {
			log.Println("} in  loop , in if, ")
			f := find[1 : len(find)-1]
			log.Println("f :", f, "row: ", j.Data[i])
			ok := bytes.Contains(j.Data[i], f)
			if ok {
				return i, nil
			}
		}
	}
	log.Println("} out loop, ")
	return -1, fmt.Errorf("record not found records")
}

// remove top,bottom braces in bytes
func RemoveBraceTopBottom(data_row []byte) (raw_bytes []byte, err error) {
	n := len(data_row)

	// data_row[0]  have  '{'
	// data_row[n-1] have '}'
	//minimum need of len more than 2 bute need 7
	//e.g. : {"a":2}, its lenght is 7, so, it is required.
	if n < 7 || data_row == nil {
		return nil, fmt.Errorf("bytes lenght should be greater than 7")
	}
	var braceL byte = '{'
	var braceR byte = '}'
	if data_row[0] == braceL && data_row[n-1] == braceR {
		raw_bytes = data_row[1 : n-1]
		return
	}
	return nil, fmt.Errorf("braces could not removed")
}

// every new bytes row check : if row[idx] == "" , skip the row,
// find next if rows[idx] != "" { do....}.
func SplitBytesByComma(data_row []byte) (splitData [][]byte, err error) {
	row, er := RemoveBraceTopBottom(data_row)
	if er != nil {
		return nil, er
	}
	var bytbyt [][]byte
	// if len(data_row) < 7 || data_row == nil {
	// 	return nil, fmt.Errorf("cant'n send nill or empty json bytes slice to split")
	// }
	// var braceL byte = '{'
	// var braceR byte = '}'
	new_data_rows := bytes.Split(row, []byte(","))

	//openBrace := new_data_rows[0][0] //to get '{'  byte from data[0] in in openBrace
	//rowsLen := len(new_data_rows)
	//lastRowLen := len(new_data_rows[rowsLen-1])

	if len(new_data_rows) == 1 {
		// bitSlice := data_row
		bytbyt = append(bytbyt, row)
		return bytbyt, nil
	}
	// if openBrace == '{' && len(new_data_rows) == 1 {
	// 	new_data_rows[0] = new_data_rows[0][1 : lastRowLen-1]
	// } else
	// if openBrace == braceL {
	// 	new_data_rows[0] = new_data_rows[0][1:]
	// 	new_data_rows[rowsLen-1] = new_data_rows[rowsLen-1][:lastRowLen-1]
	// }
	return new_data_rows, nil
}
