package utilities

import (
	"context"
	"errors"
	"strconv"

	"google.golang.org/grpc/metadata"
)

func GetStringValue(ctx context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("No Metadata")
	}

	values := md[key]
	if len(values) == 0 {
		return "", errors.New("Key Not Found")
	}

	return values[0], nil
}

func GetInt64Value(ctx context.Context, key string) (int64, error) {
	strVal, err := GetStringValue(ctx, key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(strVal, 10, 64)
}

func GetInt32Value(ctx context.Context, key string) (int32, error) {
	strVal, err := GetStringValue(ctx, key)
	if err != nil {
		return 0, err
	}

	val, err := strconv.ParseInt(strVal, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(val), err
}

func GetIntValue(ctx context.Context, key string) (int, error) {
	strVal, err := GetStringValue(ctx, key)
	if err != nil {
		return 0, err
	}

	val, err := strconv.ParseInt(strVal, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(val), err
}
