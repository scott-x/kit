# kit

```bash
go get "github.com/scott-x/kit"
```

### API

**restrict**

- `func NewDBField(arg string, value interface{}) *DBField`
- `func NewDBFieldWithSingleRestrict(c *gin.Context, param string, convert2int bool, key RestrictKey) *DBField `
- `func (field *DBField) AddRestrict(key RestrictKey, value interface{})`


**RestrictKey options**

- `EQ`    
- `LT`                    
- `GT`                     
- `LIKE`                     
- `L_Like`                   
- `R_Like`                
- `LOCATE`

**response**

- `func New(err error, sucess_msg string, data interface{}) *Response`
- `func (r *Response) Do(c *gin.Context)`

**msql**

- `func HandleStmtExec(result sql.Result, err error) error`