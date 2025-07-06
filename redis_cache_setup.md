# Redis Cache Setup

## Cấu trúc đã implement

### 1. Redis Adapter (`services/adapter/redis.go`)
- **HashKey**: Tạo key bằng SHA256 hash từ service + method + request data
- **Get/Set**: Lưu/lấy data dạng bytes
- **GetJSON/SetJSON**: Lưu/lấy data dạng JSON
- **TTL support**: Có thể set TTL custom cho từng key

### 2. Cache Interceptor (`services/interceptor/cache.go`)
- **gRPC Interceptor**: Tự động cache response của các read operations
- **Skip write operations**: Login, Create, Update, Delete không được cache
- **Cache key**: Hash từ method name + request message

### 3. Service Integration
Đã tích hợp vào:
- **Common Service** (port 50051)
- **Auth Service** (port 50052)

## Sử dụng

### Environment Variables
```bash
export REDIS_URI="redis://localhost:6379"
export MONGO_URI="mongodb://thaily:Th@i2004@localhost:27017"
export MONGO_DB="mongorest"
```

### Command Line Options
```bash
# Run với cache enabled (default)
make service-common

# Run với cache disabled
go run services/_common/main.go --enable-cache=false

# Custom cache TTL (default 5 phút)
go run services/_common/main.go --cache-ttl=10m

# Custom Redis URI
go run services/_common/main.go --redis-uri="redis://localhost:6380"
```

## Test Cache

### 1. Start Redis
```bash
docker run -d -p 6379:6379 redis:latest
```

### 2. Run service với cache
```bash
make service-common
```

### 3. Test query 2 lần
```bash
# Lần 1: Cache miss (chậm)
time grpcurl -plaintext -d '{"entity_type": "users", "page": 1, "page_size": 10}' localhost:50051 common.CommonService/Query

# Lần 2: Cache hit (nhanh)
time grpcurl -plaintext -d '{"entity_type": "users", "page": 1, "page_size": 10}' localhost:50051 common.CommonService/Query
```

### 4. Monitor cache
```bash
redis-cli monitor
```

## Cache Strategy

1. **Read operations**: Tự động cache với TTL 5 phút
2. **Write operations**: Không cache, có thể invalidate related keys
3. **Key generation**: SHA256 hash để tránh collision
4. **Graceful degradation**: Service vẫn chạy nếu Redis fail

## Performance Impact

- **Cache hit**: < 1ms response time
- **Cache miss**: Normal query time + ~1ms cache write
- **Memory usage**: Tùy thuộc vào data size và TTL