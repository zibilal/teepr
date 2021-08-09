package teepr

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"
	"time"
)

const (
	DefaultDateLayout   = "2006-01-02 15:04:05"
)

func Teepr(input interface{}, output interface{}, customValues ...func(interface{}) (interface{}, error)) (err error) {

	// it is ok to be panicked
	defer func() {
		if r := recover(); r != nil {
			log.Println("[Teepr]unable to continue iterate value ", fmt.Sprintf("%v", r))
		}
	}()

	if input == nil || output == nil {
		return nil
	}

	ival := reflect.Indirect(reflect.ValueOf(input))
	if !ival.IsValid() {
		return nil
	}

	ityp := ival.Type()
	oval := reflect.Indirect(reflect.ValueOf(output))
	otyp := oval.Type()

	switch ival.Kind() {
	case reflect.Map:

		if oval.Kind() != reflect.Map && oval.Kind() != reflect.Struct && oval.Kind() != reflect.Ptr {
			return fmt.Errorf("[Teepr]expecting output type of map or struct")
		}

		for _, k := range ival.MapKeys() {
			mival := ival.MapIndex(k)

			if oval.Kind() == reflect.Ptr {
				if oval.IsNil() {
					tmpOval := reflect.New(oval.Type().Elem())
					oval.Set(tmpOval)
					oval = oval.Elem()
					otyp = otyp.Elem()
				}
			}

			if oval.Kind() == reflect.Struct {
				foval := oval.FieldByName(k.String())
				if !foval.IsValid() {
					for s := 0; s < oval.NumField(); s++ {
						oftype := otyp.Field(s)
						oftag := string(oftype.Tag)
						ofsplit := strings.Split(oftag, " ")

						for _, split := range ofsplit {
							osplit := strings.Split(split, ":")
							if len(osplit) == 2 {
								oosplit := strings.Replace(strings.Split(osplit[1], ",")[0], "\"", "", -1)
								if k.String() == oosplit {
									foval = oval.Field(s)
									break
								}
							}
						}

					}
				}
				if !foval.IsValid() {
					continue
				}

				if istr, ok := mival.Interface().(string); ok && foval.Kind() == reflect.String {
					foval.Set(reflect.ValueOf(istr))
				} else if istr, ok := mival.Interface().(string); ok && foval.Type().String() == "time.Time" {
					dlayout := strings.Split(DefaultDateLayout, ",")
					if len(dlayout) > 0 {
						for _, l := range dlayout {
							t , err := time.Parse(l, istr)
							if err != nil {
								break
							} else {
								foval.Set(reflect.ValueOf(t))
							}
						}
					}
				} else if iint, ok := mival.Interface().(int); ok && foval.Kind() == reflect.Int {
					foval.Set(reflect.ValueOf(iint))
				} else if iint8, ok := mival.Interface().(int8); ok && foval.Kind() == reflect.Int8 {
					foval.Set(reflect.ValueOf(iint8))
				} else if iint16, ok := mival.Interface().(int16); ok && foval.Kind() == reflect.Int16 {
					foval.Set(reflect.ValueOf(iint16))
				} else if iint32, ok := mival.Interface().(int32); ok && foval.Kind() == reflect.Int32 {
					foval.Set(reflect.ValueOf(iint32))
				} else if iint64, ok := mival.Interface().(int64); ok && foval.Kind() == reflect.Int64 {
					foval.Set(reflect.ValueOf(iint64))
				} else if ifloat32, ok := mival.Interface().(float32); ok && foval.Kind() == reflect.Float32 {
					foval.Set(reflect.ValueOf(ifloat32))
				} else if ifloat64, ok := mival.Interface().(float64); ok && foval.Kind() == reflect.Float64 {
					foval.Set(reflect.ValueOf(ifloat64))
				} else if itimestamp, ok := mival.Interface().(time.Time); ok && foval.Type().String() == "time.Time" {
					foval.Set(reflect.ValueOf(itimestamp))
				} else if iffloat64, ok := mival.Interface().(float64);ok {
					switch foval.Kind() {
					case reflect.Int64:
						foval.Set(reflect.ValueOf(int64(iffloat64)))
					case reflect.Int32:
						foval.Set(reflect.ValueOf(int32(iffloat64)))
					case reflect.Int16:
						foval.Set(reflect.ValueOf(int16(iffloat64)))
					case reflect.Int8:
						foval.Set(reflect.ValueOf(int8(iffloat64)))
					case reflect.Int:
						foval.Set(reflect.ValueOf(int(iffloat64)))
					case reflect.Uint64:
						foval.Set(reflect.ValueOf(uint64(iffloat64)))
					case reflect.Uint32:
						foval.Set(reflect.ValueOf(uint32(iffloat64)))
					case reflect.Uint16:
						foval.Set(reflect.ValueOf(uint16(iffloat64)))
					case reflect.Uint8:
						foval.Set(reflect.ValueOf(uint8(iffloat64)))
					case reflect.Uint:
						foval.Set(reflect.ValueOf(uint(iffloat64)))
					case reflect.Float32:
						foval.Set(reflect.ValueOf(float32(iffloat64)))
					case reflect.Interface:
						foval.Set(reflect.ValueOf(mival.Interface()))
					}
				} else if foval.Type().String() == mival.Type().String() {
					foval.Set(mival)
				} else if mival.Kind() == reflect.Interface {
					var isHandled bool
					for _, c := range customValues {
						result, resultError := c(mival.Interface())
						if resultError == nil && reflect.ValueOf(result).Type().String() == foval.Type().String() {
							foval.Set(reflect.ValueOf(result))
							isHandled = true
						}
					}

					if !isHandled {
						elemival := reflect.Indirect(mival.Elem())

						if elemival.Kind() == reflect.Slice && foval.Kind() == reflect.Slice {
							mSlice := reflect.MakeSlice(foval.Type(), 0, elemival.Len())
							for idx := 0; idx < elemival.Len(); idx++ {
								theOutput := reflect.New(foval.Type().Elem())

								tmpelemival := elemival.Index(idx)
								if elemival.Index(idx).Kind() == reflect.Interface {
									tmpelemival = elemival.Index(idx).Elem()
								}
								err = Teepr(tmpelemival.Interface(), theOutput.Interface(), customValues...)
								mSlice = reflect.Append(mSlice, theOutput.Elem())
							}
							foval.Set(mSlice)
						} else {
							if mival.Interface() != nil {
								pval := reflect.Indirect(mival.Elem())
								err = Teepr(pval.Interface(), foval.Addr().Interface(), customValues...)
								if err != nil {
									log.Println("[Teepr]", err.Error())
									return
								}
							}
						}
					}
				} else if foval.Type().String() == mival.Type().String() {
					foval.Set(mival)
				}
			} else { // assumes output of type Map
				if ityp.Elem().String() == otyp.Elem().String() {
					oval.SetMapIndex(k, mival)
				} else {
					switch otyp.Elem().String() {
					case "int":
						if itmp, ok := mival.Interface().(int); ok {
							oval.SetMapIndex(k, reflect.ValueOf(itmp))
						} else if mival.Kind() == reflect.String {
							var itmp64 int64
							itmp64, err = strconv.ParseInt(mival.Interface().(string), 10, 64)
							if err != nil {
								log.Println("[Teepr]", err.Error())
								return
							}
							oval.SetMapIndex(k, reflect.ValueOf(int(itmp64)))
						}
					case "int8":
						if itmp, ok := mival.Interface().(int8); ok {
							oval.SetMapIndex(k, reflect.ValueOf(itmp))
						} else if mival.Kind() == reflect.String {
							var itmp64 int64
							itmp64, err = strconv.ParseInt(mival.Interface().(string), 10, 64)
							if err != nil {
								log.Println("[Teepr]", err.Error())
								return
							}
							oval.SetMapIndex(k, reflect.ValueOf(int8(itmp64)))
						}
					case "int16":
						if itmp, ok := mival.Interface().(int16); ok {
							oval.SetMapIndex(k, reflect.ValueOf(itmp))
						} else if mival.Kind() == reflect.String {
							var itmp64 int64
							itmp64, err = strconv.ParseInt(mival.Interface().(string), 10, 64)
							if err != nil {
								log.Println("[Teepr]", err.Error())
								return
							}
							oval.SetMapIndex(k, reflect.ValueOf(int16(itmp64)))
						}
					case "int32":
						if itmp, ok := mival.Interface().(int32); ok {
							oval.SetMapIndex(k, reflect.ValueOf(itmp))
						} else if mival.Kind() == reflect.String {
							var itmp64 int64
							itmp64, err = strconv.ParseInt(mival.Interface().(string), 10, 64)
							if err != nil {
								log.Println("[Teepr]", err.Error())
								return
							}
							oval.SetMapIndex(k, reflect.ValueOf(int32(itmp64)))
						}
					case "int64":
						if itmp, ok := mival.Interface().(int64); ok {
							oval.SetMapIndex(k, reflect.ValueOf(itmp))
						} else if mival.Kind() == reflect.String {
							var itmp64 int64
							itmp64, err = strconv.ParseInt(mival.Interface().(string), 10, 64)
							if err != nil {
								log.Println("[Teepr]", err.Error())
								return
							}
							oval.SetMapIndex(k, reflect.ValueOf(itmp64))
						}
					case "uint":
						if itmp, ok := mival.Interface().(uint); ok {
							oval.SetMapIndex(k, reflect.ValueOf(itmp))
						} else if mival.Kind() == reflect.String {
							var itmp64 uint64
							itmp64, err = strconv.ParseUint(mival.Interface().(string), 10, 64)
							if err != nil {
								log.Println("[Teepr]", err.Error())
								return
							}
							oval.SetMapIndex(k, reflect.ValueOf(uint(itmp64)))
						}
					case "uint8":
						if itmp, ok := mival.Interface().(uint8); ok {
							oval.SetMapIndex(k, reflect.ValueOf(itmp))
						} else if mival.Kind() == reflect.String {
							var itmp64 uint64
							itmp64, err = strconv.ParseUint(mival.Interface().(string), 10, 64)
							if err != nil {
								log.Println("[Teepr]", err.Error())
								return
							}
							oval.SetMapIndex(k, reflect.ValueOf(uint8(itmp64)))
						}
					case "uint16":
						if itmp, ok := mival.Interface().(uint16); ok {
							oval.SetMapIndex(k, reflect.ValueOf(itmp))
						} else if mival.Kind() == reflect.String {
							var itmp64 uint64
							itmp64, err = strconv.ParseUint(mival.Interface().(string), 10, 64)
							if err != nil {
								log.Println("[Teepr]", err.Error())
								return
							}
							oval.SetMapIndex(k, reflect.ValueOf(uint16(itmp64)))
						}
					case "uint32":
						if itmp, ok := mival.Interface().(uint32); ok {
							oval.SetMapIndex(k, reflect.ValueOf(itmp))
						} else if mival.Kind() == reflect.String {
							var itmp64 uint64
							itmp64, err = strconv.ParseUint(mival.Interface().(string), 10, 64)
							if err != nil {
								log.Println("[Teepr]", err.Error())
								return
							}
							oval.SetMapIndex(k, reflect.ValueOf(uint32(itmp64)))
						}
					case "uint64":
						if itmp, ok := mival.Interface().(uint64); ok {
							oval.SetMapIndex(k, reflect.ValueOf(itmp))
						} else if mival.Kind() == reflect.String {
							var itmp64 uint64
							itmp64, err = strconv.ParseUint(mival.Interface().(string), 10, 64)
							if err != nil {
								log.Println("[Teepr]", err.Error())
								return
							}
							oval.SetMapIndex(k, reflect.ValueOf(itmp64))
						}
					case "float32":
						if itmp, ok := mival.Interface().(float32); ok {
							oval.SetMapIndex(k, reflect.ValueOf(itmp))
						} else if mival.Kind() == reflect.String {
							var ftmp64 float64
							ftmp64, err = strconv.ParseFloat(mival.Interface().(string), 64)
							if err != nil {
								log.Println("[Teepr]", err.Error())
								return
							}
							oval.SetMapIndex(k, reflect.ValueOf(float32(ftmp64)))
						}
					case "float64":
						if itmp, ok := mival.Interface().(float32); ok {
							oval.SetMapIndex(k, reflect.ValueOf(itmp))
						} else if mival.Kind() == reflect.String {
							var ftmp64 float64
							ftmp64, err = strconv.ParseFloat(mival.Interface().(string), 64)
							if err != nil {
								log.Println("[Teepr]", err.Error())
								return
							}
							oval.SetMapIndex(k, reflect.ValueOf(ftmp64))
						}
					case "interface {}":
						oval.SetMapIndex(k, mival)
					default:
						if otyp.Elem().Kind() == reflect.Struct {
							vvtyp := reflect.New(otyp.Elem())
							eerr := Teepr(mival.Interface(), vvtyp.Interface(), customValues...)

							if eerr != nil {
								log.Println("[Teepr]", err.Error())
								return eerr
							}

							oval.SetMapIndex(k, vvtyp.Elem())
						} else {
							panic("unsupported type pairs")
						}
					}
				}
			}
		}


		return
	case reflect.Struct:

		if oval.Kind() != reflect.Struct {
			return fmt.Errorf("expecting output type of struct")
		} else {

			for i := 0; i < ival.NumField(); i++ {

				fin := ival.Field(i)
				ftin := ityp.Field(i)

				var fout reflect.Value
				var ftout reflect.StructField

				if fout = oval.FieldByName(ftin.Name); !fout.IsValid() {
					for j := 0; j < oval.NumField(); j++ {
						ftout = otyp.Field(j)
						if itag, otag := ftin.Tag, ftout.Tag; itag != "" && otag != "" {
							var scanner1 scanner.Scanner
							scanner1.Init(strings.NewReader(string(itag)))

							for tok := scanner1.Scan(); tok != scanner.EOF; tok = scanner1.Scan() {
								switch tok {
								case scanner.String:
									text := scanner1.TokenText()
									if strings.Contains(string(otag), text) {
										fout = oval.Field(j)
										break
									}
								}
							}
						}
					}
				}

				if !fout.IsValid() {
					continue
				}

				if fout.Kind() == reflect.Interface {
					fout.Set(fin)
				} else if fin.Type().String() == "time.Time" {
					if fout.Type().String() == "time.Time" {
						fout.Set(fin)
					} else if fout.Kind() == reflect.String {
						dTime := fin.Interface().(time.Time)
						str := dTime.Format(DefaultDateLayout)
						fout.Set(reflect.ValueOf(str))
					} else if fout.Type().String() == "mysql.NullTime" {
						dTime := fin.Interface().(time.Time)
						dNullTime := mysql.NullTime{Time:dTime}
						fout.Set(reflect.ValueOf(dNullTime))
					}
				} else if fin.Type().String() == "sql.NullString" {
					if fout.Kind() == reflect.String {
						data := fin.Interface().(sql.NullString)
						fout.Set(reflect.ValueOf(data.String))
					}
				} else if fin.Type().String() == "sql.NullInt64" {
					data := fin.Interface().(sql.NullInt64)
					switch fout.Interface().(type) {
					case int64:
						fout.Set(reflect.ValueOf(data.Int64))
					case int32:
						fout.Set(reflect.ValueOf(int32(data.Int64)))
					case int16:
						fout.Set(reflect.ValueOf(int16(data.Int64)))
					case int8:
						fout.Set(reflect.ValueOf(int8(data.Int64)))
					case int:
						fout.Set(reflect.ValueOf(int(data.Int64)))
					case uint64:
						fout.Set(reflect.ValueOf(uint64(data.Int64)) )
					case uint32:
						fout.Set(reflect.ValueOf(uint32(data.Int64)))
					case uint16:
						fout.Set(reflect.ValueOf(uint16(data.Int64)))
					case uint8:
						fout.Set(reflect.ValueOf(uint8(data.Int64)))
					case uint:
						fout.Set(reflect.ValueOf(uint(data.Int64)))
					}
				} else if fin.Type().String() == "sql.NullFloat64" {
					data := fin.Interface().(sql.NullFloat64)
					switch fout.Interface().(type) {
					case float64:
						fout.Set(reflect.ValueOf(data.Float64))
					case float32:
						fout.Set(reflect.ValueOf(float32(data.Float64)))
					}
				} else if fin.Type().String() == "mysql.NullTime" {
					if fout.Type().String() == "time.Time" {
						data := fin.Interface().(mysql.NullTime)
						fout.Set(reflect.ValueOf(data.Time))
					}
				} else if fin.Kind() == reflect.Map {
					err = Teepr(fin.Interface(), fout.Interface(), customValues...)
					if err != nil {
						log.Println("[Teepr]", err.Error())
						return err
					}
				} else {

					if fout.IsValid() && fin.IsValid() {
						var atype reflect.Type
						var abool bool
						if fout.Kind() == reflect.Ptr {
							atype = fout.Type().Elem()
							abool = true
						} else {
							atype = fout.Type()
							abool = false
						}
						iout := reflect.New(atype)

						err = Teepr(fin.Interface(), iout.Interface(), customValues...)
						if err != nil {
							log.Println("[Teepr]", err.Error())
							return
						}
						if abool {
							fout.Set(iout)
						} else {
							fout.Set(iout.Elem())
						}
					}
				}

			}
		}

		return nil
	case reflect.Slice:

		if oval.Kind() == reflect.Interface {
			oval.Set(ival)
		} else if oval.Kind() == reflect.Slice {
			outSlice := reflect.MakeSlice(reflect.SliceOf(otyp.Elem()), 0, ival.Len())
			for i := 0; i < ival.Len(); i++ {

				oItem := reflect.New(otyp.Elem())
				iItem := ival.Index(i)
				err = Teepr(iItem.Interface(), oItem.Interface(), customValues...)
				if err != nil {
					log.Println("[Teepr]", err.Error())
					return
				}
				outSlice = reflect.Append(outSlice, oItem.Elem())
			}
			oval.Set(outSlice)
		}
	case reflect.Array:
		if oval.Kind() == reflect.Interface {
			oval.Set(ival)
		}else if oval.Kind() == reflect.Array {
			oval.Set(ival)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:

		if oval.Kind() == ival.Kind() {
			oval.Set(ival)
		} else if ival.Kind() == reflect.String && oval.Type().String() == "sql.NullString" {
			outString := sql.NullString{String: ival.Interface().(string)}
			oval.Set(reflect.ValueOf(outString))
		} else if ival.Kind() == reflect.Int64 && oval.Type().String() == "sql.NullInt64" {
			outInt64 := sql.NullInt64{Int64:ival.Interface().(int64)}
			oval.Set(reflect.ValueOf(outInt64))
		} else if ival.Kind() == reflect.Int32 && oval.Type().String() == "sql.NullInt64" {
			tmp := ival.Interface().(int32)
			outInt64 := sql.NullInt64{Int64:int64(tmp)}
			oval.Set(reflect.ValueOf(outInt64))
		} else if ival.Kind() == reflect.Int16 && oval.Type().String() == "sql.NullInt64" {
			tmp := ival.Interface().(int16)
			outInt64 := sql.NullInt64{Int64:int64(tmp)}
			oval.Set(reflect.ValueOf(outInt64))
		} else if ival.Kind() == reflect.Int8 && oval.Type().String() == "sql.NullInt64" {
			tmp := ival.Interface().(int8)
			outInt64 := sql.NullInt64{Int64:int64(tmp)}
			oval.Set(reflect.ValueOf(outInt64))
		} else if ival.Kind() == reflect.Int && oval.Type().String() == "sql.NullInt64" {
			tmp := ival.Interface().(int)
			outInt64 := sql.NullInt64{Int64: int64(tmp)}
			oval.Set(reflect.ValueOf(outInt64))
		} else if ival.Kind() == reflect.Uint64 && oval.Type().String() == "sql.NullInt64" {
			tmp := ival.Interface().(uint64)
			outInt64 := sql.NullInt64{Int64: int64(tmp)}
			oval.Set(reflect.ValueOf(outInt64))
		} else if ival.Kind() == reflect.Uint32 && oval.Type().String() == "sql.NullInt64" {
			tmp := ival.Interface().(uint32)
			outInt64 := sql.NullInt64{Int64: int64(tmp)}
			oval.Set(reflect.ValueOf(outInt64))
		} else if ival.Kind() == reflect.Uint16 && oval.Type().String() == "sql.NullInt64" {
			tmp := ival.Interface().(uint16)
			outInt64 := sql.NullInt64{Int64: int64(tmp)}
			oval.Set(reflect.ValueOf(outInt64))
		} else if ival.Kind() == reflect.Uint8 && oval.Type().String() == "sql.NullInt64" {
			tmp := ival.Interface().(uint8)
			outInt64 := sql.NullInt64{Int64: int64(tmp)}
			oval.Set(reflect.ValueOf(outInt64))
		} else if ival.Kind() == reflect.Uint && oval.Type().String() == "sql.NullInt64" {
			tmp := ival.Interface().(uint)
			outInt64 := sql.NullInt64{Int64: int64(tmp)}
			oval.Set(reflect.ValueOf(outInt64))
		} else if ival.Kind() == reflect.Float64 && oval.Type().String() == "sql.NullFloat64" {
			tmp := ival.Interface().(float64)
			outFloat64 := sql.NullFloat64{Float64:tmp}
			oval.Set(reflect.ValueOf(outFloat64))
		} else if ival.Kind() == reflect.Float32 && oval.Type().String() == "sql.NullFloat32" {
			tmp := ival.Interface().(float32)
			outFloat64 := sql.NullFloat64{Float64:float64(tmp)}
			oval.Set(reflect.ValueOf(outFloat64))
		} else if ival.Kind() == reflect.Float64 {
			switch oval.Kind() {
			case reflect.Int:
				if dval, ok := ival.Interface().(float64); ok {
					ial := int(dval)
					oval.Set(reflect.ValueOf(ial))
				}
			case reflect.Int8:
				if dval, ok := ival.Interface().(float64); ok {
					ial := int8(dval)
					oval.Set(reflect.ValueOf(ial))
				}
			case reflect.Int16:
				if dval, ok := ival.Interface().(float64); ok {
					ial := int16(dval)
					oval.Set(reflect.ValueOf(ial))
				}
			case reflect.Int32:
				if dval, ok := ival.Interface().(float64); ok {
					ial := int32(dval)
					oval.Set(reflect.ValueOf(ial))
				}
			case reflect.Int64:
				if dval, ok := ival.Interface().(float64); ok {
					ial := int64(dval)
					oval.Set(reflect.ValueOf(ial))
				}
			case reflect.Uint:
				if dval, ok := ival.Interface().(float64); ok {
					ial := uint(dval)
					oval.Set(reflect.ValueOf(ial))
				}
			case reflect.Uint8:
				if dval, ok := ival.Interface().(float64); ok {
					ial := uint8(dval)
					oval.Set(reflect.ValueOf(ial))
				}
			case reflect.Uint16:
				if dval, ok := ival.Interface().(float64); ok {
					ial := uint16(dval)
					oval.Set(reflect.ValueOf(ial))
				}
			case reflect.Uint32:
				if dval, ok := ival.Interface().(float64); ok {
					ial := uint32(dval)
					oval.Set(reflect.ValueOf(ial))
				}
			case reflect.Uint64:
				if dval, ok := ival.Interface().(uint64); ok {
					ial := uint64(dval)
					oval.Set(reflect.ValueOf(ial))
				}
			case reflect.Float32:
				if dval, ok := ival.Interface().(float64); ok {
					ial := float32(dval)
					oval.Set(reflect.ValueOf(ial))
				}
			case reflect.Float64:
				if dval, ok := ival.Interface().(float64); ok {
					oval.Set(reflect.ValueOf(dval))
				}
			case reflect.Bool:
				if dval, ok := ival.Interface().(bool); ok {
					oval.Set(reflect.ValueOf(dval))
				}
			}
		} else {
			theVal := append(make([]reflect.Value, 0), ival)
			insideOutput := reflect.ValueOf(output)
			insideOutput.MethodByName("Parse").Call(theVal)
		}

		return nil
	case reflect.Interface:
		pival := ival.Elem()
		err = Teepr(pival.Interface(), output, customValues...)
		if err != nil {
			log.Println("[Teepr]Error: ", err)
			return
		}

	default:
		err = fmt.Errorf("unsupported type %T", input)
		panic(err)
	}

	return
}

// IsEmpty is an helper function to decide whether a value is empty or not
// This function is mean to be used to decide whether a struct variable is empty or not
func IsEmpty(t interface{}) bool {
	return t == nil || reflect.DeepEqual(t, reflect.Zero(reflect.TypeOf(t)).Interface())
}

type Parser interface {
	Parse(input interface{}) error
}
