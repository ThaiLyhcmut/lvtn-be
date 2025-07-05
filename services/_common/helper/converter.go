package helper

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/structpb"
)

func StructToDoc(s *structpb.Struct) bson.M {
	if s == nil {
		return bson.M{}
	}
	
	doc := bson.M{}
	for k, v := range s.Fields {
		doc[k] = StructValueToInterface(v)
	}
	return doc
}

func StructValueToInterface(v *structpb.Value) interface{} {
	if v == nil {
		return nil
	}
	
	switch v.Kind.(type) {
	case *structpb.Value_NullValue:
		return nil
	case *structpb.Value_NumberValue:
		return v.GetNumberValue()
	case *structpb.Value_StringValue:
		return v.GetStringValue()
	case *structpb.Value_BoolValue:
		return v.GetBoolValue()
	case *structpb.Value_StructValue:
		return StructToDoc(v.GetStructValue())
	case *structpb.Value_ListValue:
		list := v.GetListValue()
		result := make([]interface{}, len(list.Values))
		for i, item := range list.Values {
			result[i] = StructValueToInterface(item)
		}
		return result
	default:
		return nil
	}
}

func DocToStruct(doc bson.M) (*structpb.Struct, error) {
	// Convert BSON document to a clean map first
	cleanMap := make(map[string]interface{})
	for k, v := range doc {
		cleanMap[k] = convertBSONValue(v)
	}
	
	// Then convert to structpb
	return structpb.NewStruct(cleanMap)
}

// convertBSONValue converts BSON types to clean Go types
func convertBSONValue(v interface{}) interface{} {
	switch val := v.(type) {
	case primitive.ObjectID:
		return val.Hex()
	case time.Time:
		return val.Format(time.RFC3339)
	case primitive.DateTime:
		return time.Unix(int64(val)/1000, int64(val)%1000*1000000).Format(time.RFC3339)
	case bson.M:
		cleanMap := make(map[string]interface{})
		for k, v := range val {
			cleanMap[k] = convertBSONValue(v)
		}
		return cleanMap
	case bson.A:
		cleanArray := make([]interface{}, len(val))
		for i, item := range val {
			cleanArray[i] = convertBSONValue(item)
		}
		return cleanArray
	case []interface{}:
		cleanArray := make([]interface{}, len(val))
		for i, item := range val {
			cleanArray[i] = convertBSONValue(item)
		}
		return cleanArray
	case int64:
		// Convert int64 to float64 for JSON compatibility
		return float64(val)
	case int32:
		return float64(val)
	case int:
		return float64(val)
	default:
		return v
	}
}

func InterfaceToStructValue(v interface{}) (*structpb.Value, error) {
	// First convert BSON types to clean types
	cleanValue := convertBSONValue(v)
	
	// Then use structpb's built-in conversion
	return structpb.NewValue(cleanValue)
}