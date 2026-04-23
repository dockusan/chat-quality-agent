# Deploy trên Railway

Hướng dẫn này dành cho mô hình đơn giản nhất trên Railway: `1 app service + 1 MySQL service + 1 volume`.

## 1. Tạo project

1. Import repository này vào Railway.
2. Railway sẽ build bằng `Dockerfile` ở root.
3. Vào `Settings > Networking` và bấm **Generate Domain** để lấy domain public.

Ứng dụng đã hỗ trợ sẵn `PORT` do Railway inject, nên không cần hard-code cổng chạy.

## 2. Tạo MySQL

1. Trong cùng project, tạo thêm service **MySQL**.
2. CQA sẽ tự nhận các biến Railway cung cấp như `MYSQLHOST`, `MYSQLPORT`, `MYSQLUSER`, `MYSQLPASSWORD`, `MYSQLDATABASE`.

Không cần map thủ công các biến `DB_*` nếu dùng MySQL do Railway quản lý.

## 3. Gắn volume cho file đính kèm

Tạo volume và mount vào app service, ví dụ mount path `/data`.

App sẽ tự dùng `RAILWAY_VOLUME_MOUNT_PATH` nếu có. File upload và attachment chat sẽ được lưu persistent thay vì mất sau mỗi lần redeploy.

## 4. Khai báo biến môi trường bắt buộc

Thêm các biến sau trong app service:

- `APP_ENV=production`
- `JWT_SECRET=<chuỗi ngẫu nhiên ít nhất 32 ký tự>`
- `ENCRYPTION_KEY=<đúng 32 ký tự>`
- `TZ=Asia/Ho_Chi_Minh`

Khuyến nghị thêm:

- `APP_URL=https://<domain-public-hoặc-custom-domain>`
- `RATE_LIMIT_PER_IP=500`
- `RATE_LIMIT_PER_USER=1000`

Nếu không set `APP_URL`, app sẽ tự fallback sang `RAILWAY_PUBLIC_DOMAIN`.

## 5. Healthcheck và callback

Set healthcheck path là `/health`.

Sau khi deploy xong:

- mở app lần đầu để chạy Setup
- vào cấu hình Zalo/Facebook
- dùng domain Railway hoặc custom domain làm callback URL, ví dụ `https://your-app.up.railway.app/api/v1/channels/facebook/callback`

## 6. Lưu ý vận hành

- Railway cấp HTTPS tự động cho public domain.
- Service có volume sẽ có một khoảng downtime ngắn khi redeploy.
- Nếu dùng notification hoặc OAuth callback, nên set `APP_URL` cố định theo custom domain để tránh đổi URL khi chuyển môi trường.
