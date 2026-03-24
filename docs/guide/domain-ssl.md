# Tên miền & SSL

CQA hỗ trợ SSL tự động qua Let's Encrypt. Certificate được tạo và gia hạn hoàn toàn tự động.

## Trỏ tên miền

1. Đăng nhập vào nhà cung cấp domain (123HOST, GoDaddy, Cloudflare...)
2. Tạo bản ghi DNS:
   - **Loại**: A
   - **Tên**: `cqa` (hoặc tên bạn muốn, ví dụ `chat`)
   - **Giá trị**: IP VPS của bạn
   - **TTL**: 300 (hoặc Auto)

3. Chờ DNS cập nhật (thường 5-15 phút). Kiểm tra:

```bash
ping cqa.yourdomain.com
```

Nếu trả về đúng IP VPS là DNS đã cập nhật.

## Bật SSL

Mở file `.env` trên VPS:

```bash
nano /opt/cqa/.env
```

Thêm hoặc sửa 2 dòng:

```env
LEGO_DOMAIN=cqa.yourdomain.com
LEGO_EMAIL=admin@yourdomain.com
```

Khởi động lại:

```bash
cd /opt/cqa
docker compose down
docker compose up -d
```

CQA sẽ tự động:
- Tạo SSL certificate từ Let's Encrypt
- Chuyển hướng HTTP → HTTPS
- Kiểm tra và gia hạn certificate mỗi 7 ngày

Truy cập: `https://cqa.yourdomain.com`

## Chạy không cần SSL (HTTP only)

Nếu không cần SSL (ví dụ test local hoặc mạng nội bộ), **không cần** điền `LEGO_DOMAIN`. CQA sẽ tự chạy ở chế độ HTTP trên port 80.

## Kiểm tra SSL

```bash
docker compose logs nginx --tail=20
```

Bạn sẽ thấy:

```
Certificate obtained for cqa.yourdomain.com
Starting nginx with SSL...
```

Hoặc kiểm tra trên trình duyệt — bấm vào icon khóa bên cạnh URL để xem thông tin certificate.

## Xử lý lỗi SSL

| Lỗi | Nguyên nhân | Cách sửa |
|-----|-------------|----------|
| `Could not obtain certificate` | DNS chưa trỏ đúng | Kiểm tra DNS A record |
| `Too many requests` | Đã request quá 5 lần/tuần | Chờ 1 tuần hoặc dùng staging |
| `Port 80 already in use` | Có service khác dùng port 80 | Dừng service đó (Apache, nginx cũ...) |

## Bước tiếp theo

- [Thiết lập ban đầu](/guide/initial-setup) — Tạo admin, cấu hình AI
