findinbytes is a package, using it the CRUD operation handle easy with any json bytes:


### example
This is an example of how to use and give json data using it.
* fn:
  ```sh
    var data = findinbytes.NewByteJsonObj() ;

    func (j *data) AppendByte(inJsonBytes ...[]byte) *JsonObject ;

    func (j *data) IsValidJsonBytes(inJsonBytes []byte) bool;

    func (j *JsonObject) GetAllData() [][]byte ;

    func (j *JsonObject) GetDataByIndicesIn2D(indices []int) (data [][]byte, err error) ;

    func (j *JsonObject) GetDataByIndicesIn1D(indices []int) (data []byte, err error) ;

    func (j *JsonObject) DeleteByteById(in []byte) (bool, error) ;
  
  ```



