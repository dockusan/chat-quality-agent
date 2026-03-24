# REST API

CQA cung cấp REST API đầy đủ. Tất cả endpoint yêu cầu JWT token trong header `Authorization: Bearer <token>`.

## Base URL

```
https://cqa.yourdomain.com/api/v1
```

## Setup (lần đầu)

### Kiểm tra trạng thái setup
```
GET /setup/status
Response: { "needs_setup": true }
```

### Tạo tài khoản admin đầu tiên
```
POST /setup
Body: { "email": "...", "password": "...", "name": "..." }
Response: { "access_token": "..." }
```
Chỉ hoạt động khi chưa có user nào trong hệ thống.

## Authentication

### Đăng nhập
```
POST /auth/login
Body: { "email": "...", "password": "..." }
Response: { "access_token": "..." }
```

### Refresh token
```
POST /auth/refresh
Cookie: refresh_token (HttpOnly)
```

### Đăng xuất
```
POST /auth/logout
```

## Tenant endpoints

Tất cả endpoint dưới đây nằm trong scope tenant: `/api/v1/tenants/:tenantId/...`

### Kênh chat
| Method | Path | Mô tả |
|--------|------|-------|
| GET | `/channels` | Danh sách kênh |
| POST | `/channels` | Thêm kênh mới |
| GET | `/channels/:id` | Chi tiết kênh |
| PUT | `/channels/:id` | Cập nhật kênh |
| DELETE | `/channels/:id` | Xóa kênh |
| POST | `/channels/:id/sync` | Đồng bộ tin nhắn |
| POST | `/channels/:id/test` | Test kết nối |

### Cuộc hội thoại
| Method | Path | Mô tả |
|--------|------|-------|
| GET | `/conversations` | Danh sách cuộc hội thoại |
| GET | `/conversations/:id/messages` | Tin nhắn trong cuộc hội thoại |
| GET | `/conversations/:id/evaluations` | Kết quả đánh giá |
| GET | `/conversations/export` | Xuất dữ liệu |

### Công việc
| Method | Path | Mô tả |
|--------|------|-------|
| GET | `/jobs` | Danh sách công việc |
| POST | `/jobs` | Tạo công việc mới |
| GET | `/jobs/:id` | Chi tiết công việc |
| PUT | `/jobs/:id` | Cập nhật công việc |
| DELETE | `/jobs/:id` | Xóa công việc |
| POST | `/jobs/:id/trigger` | Chạy ngay |
| POST | `/jobs/:id/test-run` | Chạy thử |
| GET | `/jobs/:id/results` | Kết quả đánh giá |

### Dashboard
| Method | Path | Mô tả |
|--------|------|-------|
| GET | `/dashboard` | Thống kê tổng quan |

### Cài đặt
| Method | Path | Mô tả |
|--------|------|-------|
| GET | `/settings` | Xem cài đặt |
| PUT | `/settings/ai` | Cấu hình AI |
| PUT | `/settings/general` | Cài đặt chung |
| POST | `/settings/ai/test` | Test kết nối AI |

### Người dùng
| Method | Path | Mô tả |
|--------|------|-------|
| GET | `/users` | Danh sách thành viên |
| POST | `/users/invite` | Mời thành viên |
| PUT | `/users/:id/role` | Thay đổi role |
| PUT | `/users/:id/reset-password` | Đặt lại mật khẩu (owner/admin) |
| DELETE | `/users/:id` | Xóa thành viên |
