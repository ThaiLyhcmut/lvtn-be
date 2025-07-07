# MongoDB Data Generator

Tool để tạo dữ liệu mẫu cho hệ thống quản lý luận văn với MongoDB.

## Cài đặt

```bash
cd mongodb-data-generator
npm install
```

## Sử dụng

### 1. Tạo dữ liệu

```bash
# Tạo dữ liệu với số lượng mặc định
npm run generate

# Tùy chỉnh số lượng bản ghi
node scripts/generate-data.js --users 200 --theses 100 --submissions 300

# Xem tất cả options
node scripts/generate-data.js --help
```

Options:
- `-u, --users <number>`: Số lượng users (mặc định: 100)
- `-t, --theses <number>`: Số lượng luận văn (mặc định: 50)
- `-s, --submissions <number>`: Số lượng bài nộp (mặc định: 150)
- `-r, --reviews <number>`: Số lượng đánh giá (mặc định: 150)
- `-d, --defenses <number>`: Số lượng lịch bảo vệ (mặc định: 30)
- `-a, --archived <number>`: Số lượng luận văn lưu trữ (mặc định: 20)

### 2. Import vào MongoDB

```bash
# Import với cấu hình mặc định (đã cấu hình sẵn connection string)
npm run import

# Tùy chỉnh connection và database
node scripts/import-data.js --uri mongodb://localhost:27017 --database my_thesis_db

# Drop collections cũ trước khi import
node scripts/import-data.js --drop

# Import collections cụ thể
node scripts/import-data.js --collections users,theses,submissions

# Không tạo indexes
node scripts/import-data.js --no-indexes
```

### 3. Xóa dữ liệu đã tạo

```bash
npm run clean
```

## Cấu trúc dữ liệu

Dữ liệu được tạo trong thư mục `data/` với các file:
- `roles.json`: Vai trò người dùng
- `departments.json`: Khoa/phòng ban
- `users.json`: Người dùng (admin, giảng viên, sinh viên)
- `thesis_statuses.json`: Trạng thái luận văn
- `theses.json`: Luận văn
- `supervisor_assignments.json`: Phân công giảng viên
- `submissions.json`: Bài nộp
- `reviews.json`: Đánh giá
- `defense_schedules.json`: Lịch bảo vệ
- `defense_scores.json`: Điểm bảo vệ
- `event_logs.json`: Log hoạt động
- `archived_theses.json`: Luận văn lưu trữ
- `archived_submissions.json`: Bài nộp lưu trữ
- `archived_reviews.json`: Đánh giá lưu trữ

## Indexes

Indexes được tự động tạo cho các collections theo cấu hình trong `config/indexes.js`.