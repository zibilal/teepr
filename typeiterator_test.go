package teepr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestIsEmpty(t *testing.T) {
	t.Log("Test IsEmpty function")
	{
		v1 := 12
		v2 := ""
		var err error
		var v3 interface{}

		if !IsEmpty(v1) {
			t.Logf("%s expected v1 not empty", success)
		} else {
			t.Fatalf("%s expected v1 not empty, got empty", failed)
		}

		if IsEmpty(v2) {
			t.Logf("%s expected v2 is empty", success)
		} else {
			t.Fatalf("%s expected v2 is empty, got %s", failed, v2)
		}

		if IsEmpty(err) {
			t.Logf("%s expected err is empty", success)
		} else {
			t.Fatalf("%s expected err is empty, got %+v", failed, err)
		}

		if IsEmpty(v3) {
			t.Logf("%s expected v3 is empty", success)
		} else {
			t.Fatalf("%s expected v3 is empty, got %v", failed, v3)
		}
	}
}

func TestSimpleTest(t *testing.T) {
	t.Log("Simple struct to struct testing")
	{
		input := struct {
			ItemName string `itag:"the_item"`
		}{
			"First Item",
		}

		output := struct {
			Name string `transform:"item_name" json:"name" transform:"the_item"`
		}{}

		err := TypeIterator(input, &output)

		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)

			if output.Name == "First Item" {
				t.Logf("%s expected output = First Item", success)
			} else {
				t.Fatalf("%s expected output = First Item, got %s", failed, output.Name)
			}
		}
	}
}

func TestTypeWithTime(t *testing.T) {
	t.Log("Testing map[strign]interface{} with string time")
	{
		output := struct {
			Name      string    `json:"name"`
			DateBirth time.Time `json:"date_birth"`
		}{}
		input := make(map[string]interface{})
		input["name"] = "Test Second"
		input["date_birth"] = "1977-12-11 12:21:50"

		err := TypeIterator(input, &output)
		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		}

		if IsEmpty(output) {
			t.Fatalf("%s expected output not empty", failed)
		}

		b, _ := json.MarshalIndent(output, "", "\t")
		t.Logf("%s Result: %s", success, string(b))
	}
}

func TestTypeIterator(t *testing.T) {
	t.Log("Testing type iterator")
	{
		ival1 := 20
		oval1 := 0

		err := TypeIterator(ival1, &oval1)

		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		if ival1 != oval1 {
			t.Fatalf("%s expected oval == ival = %v, got %v", failed, ival1, oval1)
		} else {
			t.Logf("%s expected oval == ival = %v, got %v", success, ival1, oval1)
		}

		ival2 := 20.50
		oval2 := 0.0

		err = TypeIterator(ival2, &oval2)

		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		if ival2 != oval2 {
			t.Fatalf("%s expected oval == ival = %v, got %v", failed, ival2, oval2)
		} else {
			t.Logf("%s expected oval == ival = %v, got %v", success, ival2, oval2)
		}

		ival3 := "dua puluh lima"
		oval3 := ""

		err = TypeIterator(ival3, &oval3)

		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		t.Logf("ival %v : oval %v", ival3, oval3)
	}
}

func TestTypeIteratorWithStructInputOutput(t *testing.T) {

	t.Log("testing input and output a simple struct")
	{
		ival := struct {
			Name  string  `json:"name"`
			Point float64 `json:"point"`
		}{
			"Test1", 12.2,
		}

		oval := struct {
			FullName string  `val:"name"`
			Point    float64 `val:"point"`
		}{}

		err := TypeIterator(ival, &oval)

		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		if oval.FullName == ival.Name {
			t.Logf("%s expected oval.FullName = ival.Name", success)
		} else {
			t.Fatalf("%s expected oval.FullName = ival.Name, got 1: %s, 2: %s", failed, oval.FullName, ival.Name)
		}

		if oval.Point == ival.Point {
			t.Logf("%s expected oval.Point = ival.Point", success)
		} else {
			t.Fatalf("%s expected oval.Point = ival.Point, got 1: %f, 2: %f", failed, oval.Point, ival.Point)
		}
	}

	t.Log("testing input and output nested struct")
	{
		ival := struct {
			Name    string
			Point   float64
			Profile interface{}
		}{
			Name:  "Example 1",
			Point: 56.89,
			Profile: struct {
				DeviceId  string
				ApiSecret string
				Document  string
			}{
				"1231341234", "api-bbfdsf32324df343", "http://example.com/doc.pdf",
			},
		}

		oval := struct {
			Name    string
			Point   float64
			Profile interface{}
		}{}

		err := TypeIterator(ival, &oval)
		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		if IsEmpty(oval) {
			t.Fatalf("%s expected oval not empty", failed)
		} else {
			b, err := json.MarshalIndent(oval, "", "\t")
			if err != nil {
				t.Fatalf("%s expected error nil, got %s", failed, err.Error())
			} else {
				t.Logf("%s expected rror nil", success)
			}

			t.Logf("Output %s", string(b))
		}
	}

	t.Log("input with nested struct")
	{
		userEx := UserExample{
			FirstName: "firstex",
			LastName:  "lastex",
			Email:     "firstduo@example.com",
			Title:     "title",
			Authentication: Authentication{
				Username:  "xiexample",
				APISecret: "apisecret",
				APIToken:  "apitoken",
				ServiceDetail: ServiceDetail{
					Service: "aservicenm",
					Cost:    15000.00,
				},
			},
		}

		profEx := ProfileExample{}

		err := TypeIterator(userEx, &profEx)
		if err != nil {
			t.Fatalf("%s expected err not nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error not nil", success)
		}

		if IsEmpty(profEx) {
			t.Fatalf("%s expected profEx not empty", failed)
		} else {
			b, err := json.MarshalIndent(profEx, "", "\t")
			if err != nil {
				t.Fatalf("%s expected error nil, got %s", failed, err.Error())
			} else {
				t.Logf("%s expected error nil", success)
			}

			t.Logf("%s", string(b))
		}
	}

}

func TestTypeIteratorWithMapInput(t *testing.T) {
	t.Log("Testing TypeIterator with map input, output bytes.Buffer")
	{
		map1 := make(map[string]ServiceDetail)
		map1["1"] = ServiceDetail{}
		map1["2"] = ServiceDetail{}
		map1["3"] = ServiceDetail{
			"Premium MS Spring", 30000,
		}

		map2 := bytes.NewBufferString("")
		err := TypeIterator(map1, map2)

		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		if IsEmpty(map2) {
			t.Fatalf("%s expected map2 not empty", failed)
		} else {
			if err != nil {
				t.Fatalf("%s expected error empty, got %s", success, err.Error())
			}
			t.Logf("%s expected map2 not empty, result: %s", success, map2.String())
		}
	}

	t.Log("Testing TypeIterator with map input, output map of struct")
	{
		map1 := make(map[string]ServiceDetail)
		map1["1"] = ServiceDetail{
			"Premium MS Order", 15000,
		}
		map1["2"] = ServiceDetail{
			"Premium MS Deposit", 25000,
		}
		map1["3"] = ServiceDetail{
			"Premium MS Spring", 30000,
		}

		map2 := make(map[string]ServiceCost)

		err := TypeIterator(map1, &map2)

		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		//if IsEmpty(map2) || len(map2) == 0 {
		if IsEmpty(map2) {
			t.Fatalf("%s expected map2 not empty", failed)
		} else {
			b, err := json.MarshalIndent(map2, "", "\t")
			if err != nil {
				t.Fatalf("%s expected error empty, got %s", success, err.Error())
			}
			t.Logf("%s expected map2 not empty, result: %s", success, string(b))
		}
	}
}

func TestTypeIteratorStructWithSliceToStructWithSlice(t *testing.T) {
	t.Log("Struct with interface")
	{
		order := OrderEx{
			Id:      "o123",
			Created: time.Now(),
			Updated: time.Now(),
			Status:  "OrderCreated",
			Items: []OrderItem{
				{
					Id:       "itm123",
					ItemName: "XL 2 Giga",
					Price:    150000,
				}, {
					Id:       "itm124",
					ItemName: "XL 5 Giga",
					Price:    300000,
				},
			},
		}

		orderOutput := struct {
			Id      string
			Created time.Time
			Updated time.Time
			Status  string
			Items   interface{}
		}{}

		err := TypeIterator(order, &orderOutput)

		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		if IsEmpty(orderOutput) {
			t.Fatalf("%s expected orderOutput not empty", failed)
		} else {
			b, err := json.MarshalIndent(orderOutput, "", "\t")
			if err != nil {
				t.Fatalf("%s expected error nil, got %s", failed, err.Error())
			}
			t.Logf("%s expected orderOutput not emtpy, got %s", success, string(b))
		}
	}

	t.Log("Struct with slice of struct")
	{
		order := OrderEx{
			Id:      "o123",
			Created: time.Now(),
			Updated: time.Now(),
			Status:  "OrderCreated",
			Items: []OrderItem{
				{
					Id:       "itm123",
					ItemName: "XL 2 Giga",
					Price:    150000,
				}, {
					Id:       "itm124",
					ItemName: "XL 5 Giga",
					Price:    300000,
				},
			},
		}

		orderOutput := struct {
			Id      string
			Created time.Time
			Updated time.Time
			Status  string
			Items   []Item
		}{}

		err := TypeIterator(order, &orderOutput)

		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		if IsEmpty(orderOutput) {
			t.Fatalf("%s expected orderOutput not empty", failed)
		} else {
			b, err := json.MarshalIndent(orderOutput, "", "\t")
			if err != nil {
				t.Fatalf("%s expected error nil, got %s", failed, err.Error())
			}
			t.Logf("%s expected orderOutput not emtpy, got %s", success, string(b))
		}
	}
}

func TestMapTypedonTypeIterator(t *testing.T) {
	t.Log("Testing TypeIterator on Map input")
	{
		dataString := ` {"_id" : {"$oid" : "5b2ca6d5fc0eab2244d31567"},"aggregate_id" : "000000010","created_at" : {"$date" : 1529652949638},"updated_at" : {"$date" : 1529652949638},"events" : [{"event_id" : "b363a221-3893-44e6-b08e-35c3ff3bfe6d","reference" : "201806220235490045346127","event_type" : "OrderCreated","aggregate_id" : "000000010","created_at" : {"$date" : 1529652949638},"updated_at" : {"$date" : 1529652949638},"version" : 1,"payload" : {"created_at" : {"$date" : 1529652949638},"expired_date" : {"$date" : 1529652949638},"customer_id" : {"$numberLong" : "123455676"},"campaign_id" : "","id" : "000000010","reference" : "201806220235490045346127","shipment_type" : 0,"shipping_cost" : 0.0,"status" : "Order Created","subtotal_price" : 0.0,"user_id" : {"$numberLong" : "12445"},"vendor_id" : {"$numberLong" : "15544"},"order_items" : [{"attributes" : "{\"Customer\":{\"cellphone_number\":\"0818780077\"},\"PurchaseReferral\":{\"action\":\"pulsa\"}}","description" : "XL 5 giga Ultimate Internate","item_image" : "","commission" : 3500.0,"name" : "","product_code" : "C443243234","price" : 25000.0,"quantity" : 2,"reseller_price" : 26000.0,"item_id" : {"$numberLong" : "15540"}}, {"attributes" : "{\"Customer\":{\"cellphone_number\":\"0812780077\"},\"PurchaseReferral\":{\"action\":\"pulsa\"}}","description" : "Telkomsel flash 5 giga Ultimate Internate","item_image" : "","commission" : 3500.0,"name" : "","product_code" : "C443243234","price" : 25000.0,"quantity" : 2,"reseller_price" : 26000.0,"item_id" : {"$numberLong" : "15540"}}],"agent_id" : {"$numberLong" : "0"},"device_id" : "5566478997710","total_commission" : 0.0,"shipping_trx_id" : "","channel" : "","total_price" : 0.0,"payment_type" : "","merchant_trx_id" : "","cart_id" : ""}}]}`
		dataMap := make(map[string]interface{})

		err := json.Unmarshal([]byte(dataString), &dataMap)
		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		dataAggregate := OrderAggregate{}
		handleDateTag := func(input interface{}) (interface{}, error) {
			if minput, ok := input.(map[string]interface{}); ok {
				d, found := minput["$date"]
				if found {
					if dfloat64, ok := d.(float64); ok {
						dint64 := int64(dfloat64) / 1000

						return time.Unix(dint64, 0), nil
					}
				}
			}

			return nil, fmt.Errorf("unable to handle data: %v", input)
		}

		handleIdTag := func(input interface{}) (interface{}, error) {
			if minput, ok := input.(map[string]interface{}); ok {
				s, found := minput["$oid"]
				if found {
					if str, ok := s.(string); ok {
						return str, nil
					}
				}
			}

			return nil, fmt.Errorf("unable to handle data: %v", input)
		}

		handleLongTag := func(input interface{}) (interface{}, error) {
			if minput, ok := input.(map[string]interface{}); ok {
				d, found := minput["$numberLong"]
				if found {

					if dstring, ok := d.(string); ok {
						return strconv.ParseUint(dstring, 10, 64)
					}
				}
			}

			return nil, fmt.Errorf("unable to handle data: %v", input)
		}
		err = TypeIterator(dataMap, &dataAggregate, handleDateTag, handleIdTag, handleLongTag)
		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
			b, err := json.MarshalIndent(dataAggregate, "", "\t")
			if err != nil {
				t.Fatalf("%s expected error nil, got %s", failed, err.Error())
			} else {
				t.Logf("RESULT: %s", string(b))
			}
		}
	}
}

func TestIterateStructWithInterfaceTypedField(t *testing.T) {
	now := time.Now()
	checkout := &Checkout{
		CreatedAt:  now,
		UpdatedAt:  now,
		Status:     "Process Checkout",
		CheckoutBy: "user1",
		Attributes: map[string]string{
			"phone_number": "081123480",
			"vendor":       "Numerindo",
		},
	}

	checkoutEvent := CheckoutEvent{
		AggregateId: "1313211234",
		CreatedAt:   now,
		UpdatedAt:   now,
		Checkouts:   []Encoder{checkout},
	}

	t.Log("Testing Checkout event iterate type")
	{

		ival := reflect.Indirect(reflect.ValueOf(checkoutEvent))
		ityp := ival.Type()

		for i := 0; i < ival.NumField(); i++ {
			fval := ival.Field(i)
			ftyp := ityp.Field(i)

			if fval.Kind() == reflect.Interface {
				pval := reflect.Indirect(fval.Elem())
				ptyp := pval.Type()
				t.Log("Name:", ftyp.Name, "Type:", ftyp.Type, "Value", pval, "Elem", pval.Type())

				if pval.Kind() == reflect.Map {
					ppval := reflect.Indirect(pval.Elem())
					t.Log("\tM:Name:", ftyp.Name, "Type:", ftyp.Type, "Value", ppval, "Elem", ppval.Type())
				} else if pval.Kind() == reflect.Slice {
					spval := pval.Index(i)
					t.Log("\tSl:Name:", ftyp.Name, "Type:", ftyp.Type, "Value", spval, "Elem", spval.Type())
				} else if pval.Kind() == reflect.Struct {
					for j := 0; j < pval.NumField(); j++ {
						fpval := pval.Field(j)
						fptyp := ptyp.Field(j)
						t.Log("\tSt:Name:", fptyp.Name, "Type:", fptyp.Type, "Value", fpval)
					}
				}
			} else if fval.Kind() == reflect.Slice {
				t.Log("Slice: Name:", ftyp.Name, "Type:", ftyp.Type, "Value", fval, "Elem", fval.Type())
			} else {
				t.Log("Default: Name:", ftyp.Name, "Type:", ftyp.Type, "Value", fval)
			}

		}
	}

	t.Log("Testing With TypeIterator")
	{

		checkoutOutput := struct {
			AggregateId string             `bson:"aggregate_id" json:"aggregate_id"`
			CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
			UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
			Checkouts   []CheckoutResponse `bson:"events" json:"events"`
		}{}

		if err := TypeIterator(checkoutEvent, &checkoutOutput); err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			if IsEmpty(checkoutOutput) {
				t.Fatalf("%s expected checkoutOutput is not empty", failed)
			} else {
				if b, err := json.MarshalIndent(checkoutOutput, "", "\t"); err != nil {
					t.Fatalf("%s expected error nil, got %s", failed, err)
				} else {
					t.Logf("%s expected checkoutOutput is not empty, RESULT: %s", success, string(b))
				}
			}
		}
	}

	t.Log("Testing from an map input")
	{
		mval := make(map[string]interface{})
		b, err := json.Marshal(checkoutEvent)
		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}
		err = json.Unmarshal(b, &mval)
		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		t.Log("Mval", mval)

		checkoutOutput := struct {
			AggregateId string             `bson:"aggregate_id" json:"aggregate_id"`
			CreatedAt   string             `bson:"created_at" json:"created_at"`
			UpdatedAt   string             `bson:"updated_at" json:"updated_at"`
			Checkouts   []CheckoutResponse `bson:"events" json:"events"`
		}{}

		if err := TypeIterator(mval, &checkoutOutput); err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			if IsEmpty(checkoutOutput) {
				t.Fatalf("%s expected checkoutOutput is not empty", failed)
			} else {
				if b, err := json.MarshalIndent(checkoutOutput, "", "\t"); err != nil {
					t.Fatalf("%s expected error nil, got %s", failed, err)
				} else {
					t.Logf("%s expected checkoutOutput is not empty, RESULT: %s", success, string(b))
				}
			}
		}
	}
}

type CheckoutEvent struct {
	AggregateId string    `bson:"aggregate_id" json:"aggregate_id"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	Checkouts   []Encoder `bson:"events" json:"events"`
}

type Checkout struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Status     string
	CheckoutBy string
	Attributes interface{}
}

type CheckoutResponse struct {
	CreatedAt  string
	UpdatedAt  string
	Status     string
	CheckoutBy string
	Attributes interface{}
}

func (e *Checkout) Unmarshal(b []byte) error {
	return nil
}

func (e *Checkout) Marshal(data Encoder) ([]byte, error) {
	return nil, nil
}

type Encoder interface {
	Unmarshal(b []byte) error
	Marshal(data Encoder) ([]byte, error)
}

type OrderAggregate struct {
	AggregateId string       `bson:"aggregate_id" json:"aggregate_id"`
	CreatedAt   time.Time    `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time    `bson:"updated_at" json:"updated_at"`
	Events      []OrderEvent `bson:"events" json:"events"`
}

type OrderEvent struct {
	ID          string    `bson:"event_id" json:"event_id"`
	Reference   string    `bson:"reference" json:"reference"`
	EventType   string    `bson:"event_type" json:"event_type"`
	AggregateID string    `bson:"aggregate_id" json:"aggregate_id"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	Version     int       `bson:"version" json:"version"`
	Payload     Order     `bson:"payload" json:"payload"`
}

type Order struct {
	CreatedAt       time.Time     `json:"created_at" bson:"created_at"`
	ExpiryDate      time.Time     `json:"expiry_date" bson:"expired_date"`
	CustomerID      uint64        `json:"customer_id" bson:"customer_id"`
	CampaignID      string        `json:"campaign_id" bson:"campaign_id"`
	ID              string        `json:"id" bson:"id"`
	Reference       string        `json:"reference" bson:"reference"`
	ShipmentType    int           `json:"shipment_type" bson:"shipment_type"`
	ShippingCost    float32       `json:"shipping_cost" bson:"shipping_cost"`
	Status          string        `json:"status" bson:"status"`
	SubtotalPrice   float32       `json:"subtotal_price" bson:"subtotal_price"`
	UserID          uint64        `json:"user_id" bson:"user_id"`
	VendorID        uint64        `json:"vendor_id" bson:"vendor_id"`
	OrderItems      []OrderItemOp `json:"order_items" bson:"order_items"`
	AgentID         uint64        `json:"agent_id" bson:"agent_id"`
	DeviceID        string        `json:"device_id" bson:"device_id"`
	TotalCommission float32       `json:"total_commission" bson:"total_commission"`
	ShippingTrxID   string        `json:"shipping_trx_id" bson:"shipping_trx_id"`
	Channel         string        `json:"channel" bson:"channel"`
	TotalPrice      float32       `json:"total_price" bson:"total_price"`
	PaymentType     string        `json:"payment_type" bson:"payment_type"`
	MerchantTrxID   string        `json:"merchant_trx_id" bson:"merchant_trx_id"`
	CartID          string        `json:"cart_id" bson:"cart_id"`
}

type OrderItemOp struct {
	Attributes    string  `json:"attributes" bson:"attributes" transform:"attributes"`
	Description   string  `json:"description" bson:"description"`
	ItemImage     string  `json:"item_image" bson:"item_image" transform:"item_image"`
	Commission    float32 `json:"commission" bson:"commission" transform:"commission"`
	Name          string  `json:"name" bson:"name" transform:"item_name"`
	ProductCode   string  `json:"product_code" bson:"product_code" transform:"item_reference_id"`
	Price         float32 `json:"price" bson:"price" transform:"price"`
	Quantity      int     `json:"quantity" bson:"quantity" transform:"quantity"`
	ResellerPrice float64 `json:"reseller_price" bson:"reseller_price"`
	ItemID        uint64  `json:"item_id" bson:"item_id" transform:"item_id"`
}

type OrderItemMulti struct {
	Attributes    string  `json:"attributes" transform:"attr"`
	Description   string  `json:"description" bson:"description"`
	ItemImage     string  `json:"item_image" bson:"item_image" transform:"item_image"`
	Commission    float64 `json:"commission" bson:"commission" transform:"commission"`
	Name          string  `json:"name" bson:"name" transform:"item_name"`
	ProductCode   string  `json:"product_code" bson:"product_code" transform:"item_reference_id"`
	Price         float64 `json:"price" bson:"price" transform:"price"`
	Quantity      int     `json:"quantity" bson:"quantity" transform:"quantity"`
	ResellerPrice float64 `json:"reseller_price" bson:"reseller_price"`
	ItemID        uint64  `json:"item_id" bson:"item_id" transform:"item_id"`
}

type ItemMulti struct {
	Attr          string  `json:"attributes" transform:"attr"`
	Description   string  `json:"description" bson:"description"`
	ItemImage     string  `json:"item_image" bson:"item_image" transform:"item_image"`
	Commission    float64 `json:"commission" bson:"commission" transform:"commission"`
	Name          string  `json:"name" bson:"name" transform:"item_name"`
	ProductCode   string  `json:"product_code" bson:"product_code" transform:"item_reference_id"`
	Price         float64 `json:"price" bson:"price" transform:"price"`
	Quantity      int     `json:"quantity" bson:"quantity" transform:"quantity"`
	ResellerPrice float64 `json:"reseller_price" bson:"reseller_price"`
	ItemID        uint64  `json:"item_id" bson:"item_id" transform:"item_id"`
}

type OrderEx struct {
	Id      string
	Updated time.Time
	Created time.Time
	Status  string
	Items   []OrderItem
}

type OrderItem struct {
	Id       string
	ItemName string
	Price    float64
}

type Item struct {
	Id       string
	ItemName string
	Price    float64
}

type User struct {
	Id           string         `sqltype:"id"`
	Username     string         `sqltype:"username"`
	FirstName    string         `sqltype:"first_name"`
	LastName     string         `sqltype:"last_name"`
	Email        string         `sqltype:"email"`
	SecretDetail SecretDetailEx `sqltype:"secret_detail"`
}

type SecretDetailEx struct {
	Id        string `json:"id"`
	APISecret string `json:"api_secret"`
	APIToken  string `json:"api_token"`
}

type OrderExample struct {
	Id      int64
	Updated time.Time
	Created time.Time
	Status  string
}

type ProfileExample struct {
	FirstName      string
	LastName       string
	Email          string
	Title          string
	Authentication SecretDetail
}

type SecretDetail struct {
	APISecret string `json:"api_secret"`
	APIToken  string `json:"api_token"`
}

type UserExample struct {
	FirstName      string
	LastName       string
	Email          string
	Title          string
	Authentication Authentication
	CreatedAt      time.Time
}

type Authentication struct {
	Username      string
	APISecret     string
	APIToken      string
	ServiceDetail ServiceDetail
}

type ServiceDetail struct {
	Service string  `ksmg:"service_name"`
	Cost    float64 `kmsg:"service_cost"`
}

type ServiceCost struct {
	Service string  `json:"service_name"`
	Cost    float64 `json:"service_cost"`
}
