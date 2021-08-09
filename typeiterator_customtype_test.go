package teepr

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type MyInt int

func (i MyInt) String() string {
	return strconv.Itoa(int(i))
}

type TestingEx struct {
	Id      string
	Name    string
	ANumber MyInt
}

func TestCustomTypeInt(t *testing.T) {
	t.Log("Testing Custom Type Int")
	{
		output := struct {
			Id      string
			Name    string
			ANumber MyInt
		}{}

		input := TestingEx{
			"00088112266777",
			"Testing",
			MyInt(1),
		}

		err := Teepr(input, &output)
		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		}
		t.Logf("%s Result: %v", success, output)
	}
}

type AppId uuid.UUID

func (id *AppId) Parse(input interface{}) error {
	tmp, ok := input.(string)
	if !ok {
		return errors.New(fmt.Sprintf("Unable to parse input of type %v", reflect.TypeOf(input)))
	}
	v, err := hex.DecodeString(tmp)
	if err != nil {
		return err
	}
	result, err := uuid.FromBytes(v)
	if err != nil {
		return err
	}
	*id = AppId(result)
	return nil
}

func (id *AppId) String() string {
	fmt.Println("String()")
	tmp := uuid.UUID(*id)
	b, _ := tmp.MarshalBinary()
	return hex.EncodeToString(b)
}

type TestType struct {
	Id        AppId
	Name      string
	Email     string
	CreatedAt time.Time
}

func TestParserType(t *testing.T) {
	strId := "1DD1B664F14E11EBACE1ACDE48001122"
	tmp, err := hex.DecodeString(strId)
	tmpId, err := uuid.FromBytes(tmp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("------------------------------------------------------")
	t.Log("Testing input of type Parser and output of type Parser")
	t.Log("------------------------------------------------------")
	{
		theId := AppId(tmpId)
		input := struct {
			Id    Parser
			Name  string
			Email string
		}{
			&theId, "A Name", "aname@example.com",
		}
		output := struct {
			Id    Parser
			Name  string
			Email string
		}{}
		err := Teepr(input, &output)
		if err != nil {
			t.Fatal(err)
		}

		if IsEmpty(output.Id) {
			t.Fatalf("%s Expected Id is not empty", failed)
		}
		t.Logf("%s Result:%v", success, output)

		appId := AppId{}

		err = Teepr(output.Id, &appId)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s Expected error nil", success)
		t.Logf("%s Result AppId %v", success, uuid.UUID(appId).String())
	}

	t.Log("-----------------------------------------------------")
	t.Log("Testing input of type interface and output of type Parser")
	t.Log("-----------------------------------------------------")
	{
		input := struct {
			Id    interface{}
			Name  string
			Email string
		}{
			strId, "A Name", "aname@example.com",
		}

		output := struct {
			Id    Parser
			Name  string
			Email string
		}{}
		err := Teepr(input, &output)
		if err != nil {
			t.Fatalf("%s Expected Error nil, got %v", failed, err)
		}
		if IsEmpty(output.Id) {
			t.Logf("%s Expected output.Id is empty", success)
		}
		if IsEmpty(output) {
			t.Logf("%s Expected output is empty", success)
		}
	}

	t.Log("-----------------------------------------------------")
	t.Log("Testing input of type interface and output of type string")
	t.Log("-----------------------------------------------------")
	{
		input := struct {
			Id    interface{}
			Name  string
			Email string
		}{
			strId, "A Name", "aname@example.com",
		}

		output := struct {
			Id    string
			Name  string
			Email string
		}{}
		err := Teepr(input, &output)
		if err != nil {
			t.Fatalf("%s Expected Error nil, got %v", failed, err)
		}
		if IsEmpty(output.Id) {
			t.Fatalf("%s Expected output.Id is not empty", failed)
		}
		t.Logf("%s Result output: %v", success, output)
	}

	t.Log("-----------------------------------------------------")
	t.Log("Testing input of type string and output of type AppId")
	t.Log("-----------------------------------------------------")
	{
		input := struct {
			Id    string
			Name  string
			Email string
		}{
			strId, "A Name", "aname@example.com",
		}

		output := struct {
			Id    AppId
			Name  string
			Email string
		}{}
		err := Teepr(input, &output)
		if err != nil {
			t.Fatalf("%s Expected Error nil, got %v", failed, err)
		}
		if IsEmpty(output.Id) {
			t.Fatalf("%s Expected output.Id is not empty", failed)
		}
		t.Logf("%s Output Id %s", success, output.Id.String())
		t.Logf("%s Result output: %v", success, output)
	}
}

func TestCustomTypeAppId(t *testing.T) {
	strId := "1DD1B664F14E11EBACE1ACDE48001122"
	tmp, err := hex.DecodeString(strId)
	if err != nil {
		t.Fatal(err)
	}
	tmpId, err := uuid.FromBytes(tmp)
	t.Log("------------------------------------------------------------")
	t.Log("Testing Custom Type AppId, input and output on the same type")
	t.Log("------------------------------------------------------------")
	{
		input := struct {
			Id    AppId
			Name  string
			Email string
		}{
			AppId(tmpId),
			"A Name",
			"aname@example.com",
		}

		output := TestType{}

		err := Teepr(input, &output)
		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		}
		t.Logf("%s Result: %v", success, output)
	}
	t.Log("-----------------------------------------------------------------------------")
	t.Log("Testing Custom Type AppId, input of type AppId and output of type interface{}")
	t.Log("-----------------------------------------------------------------------------")
	{
		input := struct {
			Id    AppId
			Name  string
			Email string
		}{
			AppId(tmpId), "A Name", "aname@example.com",
		}

		output := struct {
			Id    interface{}
			Name  string
			Email string
		}{}
		err := Teepr(input, &output)
		if err != nil {
			t.Fatalf("%s Expected error nil, got %s", failed, err.Error())
		}
		if IsEmpty(output.Id) {
			t.Fatalf("%s Expected Id not empty", failed)
		}
		t.Logf("%s Result %v", success, output)
	}

	t.Log("-----------------------------------------------------------------------------")
	t.Log("Testing Custom Type AppId, input of type interface{} and output of type AppId")
	t.Log("-----------------------------------------------------------------------------")
	{
		input := struct {
			Id    interface{}
			Name  string
			Email string
		}{
			strId, "A Name", "aname@example.com",
		}

		output := struct {
			Id    AppId
			Name  string
			Email string
		}{}
		err := Teepr(input, &output)
		if err != nil {
			t.Fatalf("%s Expected error nil, got %s", failed, err.Error())
		}
		if IsEmpty(output.Id) {
			t.Fatalf("%s Expected Id not empty", failed)
		}
		t.Logf("%s result %v", success, output)
	}
}
