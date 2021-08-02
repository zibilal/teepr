package teepr

import (
	"encoding/hex"
	"github.com/google/uuid"
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

type TestType struct {
	Id        AppId
	Name      string
	Email     string
	CreatedAt time.Time
}

func TestCutomTypeAppId(t *testing.T) {
	tmp, err := hex.DecodeString("1DD1B664F14E11EBACE1ACDE48001122")
	if err != nil {
		t.Fatal(err)
	}
	tmpId, err := uuid.FromBytes(tmp)
	t.Log("Testing Custom Type AppId")
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
}
